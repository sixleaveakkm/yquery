// Package yquery is a yq style parse to let you handle yaml file without provide data struct
// You get get string item by provide string (e.g., "a.b[0]") in your golang project
// This package use [go-yaml v3](https://github.com/go-yaml/yaml/tree/v3) to do the base parse work.
package yquery

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	seqTag   = "!!seq"
	mapTag   = "!!map"
	mergeTag = "!!merge"
)

// YQuery is the data struct hold necessary unmarshal data
// RootNode is the root node holds unmarshal struct.
// It has type of Node from gopkg.in/yaml.v3 , you can operate it directly if you want.
type YQuery struct {
	RootNode *yaml.Node

	maxMergeInOneLayer int
}

// Unmarshal bytes data into a struct (Node) inside this package, return error if meets problem
// Optionally you could override maximum number of merge struct directly in one node.
// The default value will be 3. With should cover most case.
// e.g.
//     ...
//      a:
//        <<: *mergeFromSomewhere
//        <<: *mergeFromSomewhereElse
//     ...
//
// In the example above, the number of merge in one node is 2.
// And in the example below, the number is 1.
//  	a:
//  	  <<: &a
//  	     a1: foo
//  	b:
//  	  <<: &b
//  	    <<: *a
//  	    b1: bar
//  	c:
//  	  <<: &b
// It use RootNode to store data, which type is *yaml.Node, comes from go-yaml.
func Unmarshal(in []byte, maxMerge ...int) (*YQuery, error) {
	y := YQuery{}
	if len(maxMerge) > 0 {
		y.maxMergeInOneLayer = maxMerge[0]
	} else {
		y.maxMergeInOneLayer = 3
	}

	node := yaml.Node{}
	err := yaml.Unmarshal(in, &node)
	if err != nil {
		return nil, err
	}
	y.RootNode = node.Content[0]
	return &y, nil
}

// Marshal struct, return bytes data if no error
// Wrapper of gopkg.in/yaml.v3.
func (y *YQuery) Marshal() ([]byte, error) {
	return yaml.Marshal(y.RootNode)
}

// Get return the parsed data string of the parser if no error
// Receives a parser string (e.g. "a.b") with optional delimiter character.
// Example, if if the data is like following:
//      a:
//         b: data of b
//
// Get("a.b") should return "data of b".
// Optional parameter "customDelimiter".
// For a struct like following,
//
//    example.com:
//        admin: admin@example.com
//
// there is no way to know "example.com.admin" means "admin" in "example.com" or "admin" in "com" in "example".
// Currently go yaml v3 seems don't support key string with bracket, e.g. "[example.com]".
// therefore you could provide a custom delimiter, e.g. `Get("example.com;admin",";")`.
func (y *YQuery) Get(parser string, customDelimiter ...string) (string, error) {
	return y.getNodeString(parser, false, customDelimiter...)

}

// GetRaw return the raw string if there is no error
// GetRaw is similar to Get, but keep the anchor and merge item.
//
//     a: &anchorA
//       b: data of b
//     c: *anchorA
//     d:
//       <<: &mergeC
//           e: 1
//     f:
//       <<: *mergeC
//
// For the data above,
// using `Get("c")`, it should return "b: data of b",
// using `GetRaw("c")`, you can get `*anchorA`.
func (y *YQuery) GetRaw(parser string, customDelimiter ...string) (string, error) {
	return y.getNodeString(parser, true, customDelimiter...)
}

func (y *YQuery) getNodeString(parser string, raw bool, customDelimiter ...string) (string, error) {
	if len(customDelimiter) > 1 {
		return "", fmt.Errorf("get could only get 0 or 1 string for delimiter, got %s", customDelimiter)
	}
	delimiter := ""
	if len(customDelimiter) > 0 {
		delimiter = customDelimiter[0]
	}
	config := Config{Delimiter: delimiter}
	node, err := y.GetNode(parser, raw, config)
	if err != nil {
		return "", err
	}

	if node.Value != "" {
		return node.Value, nil
	}
	str, err := yaml.Marshal(&node)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

// GetNode corresponding node
// Receives a parser string (e.g. "a.b") with optional delimiter character.
func (y *YQuery) GetNode(parser string, raw bool, config ...Config) (*yaml.Node, error) {
	delimiter, err := getDelimiter(config)
	if err != nil {
		return y.RootNode, err
	}
	slices := getParserSlice(parser, delimiter)
	// make a copy of root. Prevent modify
	node := y.RootNode
	result := y.parseNode(slices, delimiter, node, 0, raw)
	return result.Node, result.Err
}

// Set the value of responding node
// Cannot set value inside anchor reference's, and not able to override sub item of a merge item.
func (y *YQuery) Set(parser string, value string, config ...Config) error {
	delimiter, err := getDelimiter(config)
	if err != nil {
		return err
	}
	slices := getParserSlice(parser, delimiter)
	result := y.parseNode(slices, delimiter, y.RootNode, 0, true, parseParameter{
		setParameter: setParameter{
			Config:  config[0],
			Value:   value,
			InMerge: false,
		},
	})

	return fmt.Errorf("the item '%s' is unable to write because it is in an anchor reference or is an item of the merged item", parser)

	result.Node.Value = value
	return nil

}

// Config is the optional parameter for set
// All elements are optional.
type Config struct {
	// Delimiter is the custom delimiter
	Delimiter string
	// ForceInMerge should be set to true if you want to set a element inside a merge node and this element is not a root element of the merge node
	// It is used to prevent unexpected data loss because all other data related to this path will be loss.
	// e.g. there is a merge
	//     o:
	//       <<: *mergeA
	//          a:
	//             b: 333
	//             c: testValue
	//          d: otherValue
	//
	// If you set "o.a.b" to 444, due to the yaml design, a new node "a" will be added directly to "o" to override the merge.
	// The result will be
	//     o:
	//       <<: *mergeA
	//          a:
	//             b: 333
	//             c: testValue
	//          d: otherValue
	//       a:
	//         b: 444
	//
	// Which means "o.a.c" is a unset value, and this may be not what you want. Therefore, yquery will return an error tell you it is illegal.
	// If you trust yourself knowing there is no other value in the structure,
	// you could use this parameter to force set the value.
	ForceInMerge bool
	// Recursive could be set to true to create missing node in the middle of the path of your parser string
	// Or yquery will return an error when it could not found the element.
	// You don't need to set this to true if only the last element of your parser string is not exist.
	Recursive bool
}

type setParameter struct {
	Config
	Value   string
	InMerge bool
}

type getParameter struct {
	IsRaw bool
}

type parseParameter struct {
	setParameter
	getParameter
}

type parseResult struct {
	Node *yaml.Node
	Err  error
}

func (y *YQuery) parseNode(slices []string, delimiter string, currentNode *yaml.Node, i int, isSet bool, parameter parseParameter) parseResult {
	var e error
	ch := make(chan parseResult, y.maxMergeInOneLayer)

	if i < len(slices) {
		// anchor use, return error
		if currentNode.Alias != nil {
			if isSet {
				return parseResult{currentNode, fmt.Errorf("the item '%s' reaches an anchor reference. You cannot modify value from anchor reference", strings.Join(slices[:i], delimiter))}
			}
			currentNode = currentNode.Alias
		}

		nodeTag := currentNode.Tag
		switch nodeTag {
		case seqTag:
			index, err := getSequenceNum(slices[i])
			if err != nil {
				return parseResult{currentNode, err}
			}
			return y.parseNode(slices, delimiter, currentNode.Content[index], i+1, isSet, parameter)
		case mapTag:
			for index, content := range currentNode.Content {
				if index%2 == 1 {
					continue
				}

				if content.Tag == mergeTag {
					parameter.InMerge = true
					ch <- y.parseNode(slices, delimiter, currentNode.Content[index+1], i, isSet, parameter)
				}
				if content.Value == slices[i] {
					return y.parseNode(slices, delimiter, currentNode.Content[index+1], i+1, isSet, parameter)
				}
			}
			e = fmt.Errorf("cannot find item %s", strings.Join(slices[:i], delimiter))

		default:
			e = fmt.Errorf("unable continue to parse item %s, get value: %s",
				strings.Join(slices[:i], delimiter), currentNode.Value)
		}
		close(ch)
		if e != nil {
			var resArr []parseResult
			for i := 0; i < y.maxMergeInOneLayer; i++ {
				res, ok := <-ch
				if !ok {
					continue
				}
				if res.Err == nil {
					resArr = append(resArr, res)
				}
			}
			if len(resArr) > 0 {
				return resArr[len(resArr)-1]
			}
		}
		return parseResult{currentNode, e}
	} else {
		if !isSet {
			// anchor use, change anchor value to *value when get raw
			if currentNode.Alias != nil {
				if parameter.IsRaw {
					currentNode.Value = "*" + currentNode.Value
				}
			}
			// anchor definition, remove data when it is raw
			if currentNode.Anchor != "" {
				if !parameter.IsRaw {
					currentNode.Anchor = ""
				}
			}
			return parseResult{currentNode, nil}
		} else {
			return parseResult{currentNode, nil}
		}
	}
}
