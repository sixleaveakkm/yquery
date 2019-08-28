// YQuery is a yq style parse to let you handle yaml file without provide data struct
// You get get string item by provide string (e.g., "a.b[0]")
// This package use [go-yaml v3](https://github.com/go-yaml/yaml/tree/v3) to do the base parse work
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

// RootNode is the root node holds unmarshal struct
// It has type of Node from gopkg.in/yaml.v3 , you can operate it directly if you want
var RootNode *yaml.Node

var MaxMergeInOneLayer int = 3

// New unmarshal bytes data into a struct (Node) inside this package, return error if meets problem
// Optionally you could override maximum number of merge struct directly in one node
// the default value will be 3. With should cover most case.
// e.g.
// ```
// ...
// a:
//   <<: *mergeFromSomewhere
//   <<: *mergeFromSomewhereElse
// ...
// ```
// In the example above, the number of merge in one node is 2.
// And in the example below, the number is 1
// ```
// a:
//   <<: &a
//      a1: foo
// b:
//   <<: &b
//     <<: *a
//     b1: bar
// c:
//   <<: &b
// ```
// It use RootNode to store data, which type is *yaml.Node, comes from go-yaml
// And since it holds data inside, personal I prefer `New` rather then `Unmarshal`
func New(in []byte, maxMerge ...int) error {
	if len(maxMerge) > 0 {
		MaxMergeInOneLayer = maxMerge[0]
	}

	node := yaml.Node{}
	err := yaml.Unmarshal(in, &node)
	if err != nil {
		return err
	}
	RootNode = node.Content[0]
	return nil
}

// Marshal struct, return bytes data if no error
// Wrapper of gopkg.in/yaml.v3
func Marshal() ([]byte, error) {
	return yaml.Marshal(RootNode)
}

// Get return the parsed data string of the parser if no error
// Receives a parser string (e.g. "a.b") with optional delimiter character.
// Example, if if the data is like following:
// ```yaml
// a:
//   b: data of b
// ```
// Get("a.b") should return "data of b"
//
// Optional parameter "customDelimiter"
// For a struct like following,
// ```
// example.com:
//   admin: admin@example.com
// ```
// there is no way to know "example.com.admin" means "admin" in "example.com" or "admin" in "com" in "example"
// Currently go yaml v3 seems don't support key string with bracket, e.g. "[example.com]"
// therefore you could provide a custom delimiter, e.g. `Get("example.com;admin",";")`
func Get(parser string, customDelimiter ...string) (string, error) {
	return getNodeString(parser, false, customDelimiter...)

}

// GetRaw return the raw string if there is no error
// GetRaw is similar to Get, but keep the anchor and merge item
// ```yaml
// a: &anchorA
//   b: data of b
// c: *anchorA
// d:
//   <<: &mergeC
//      e: 1
// f:
//   <<: *mergeC
// ```
// For the data above
// Using `Get("c")`, it should return "b: data of b"
// Using `GetRaw("c")`, you can get `*anchorA`
func GetRaw(parser string, customDelimiter ...string) (string, error) {
	return getNodeString(parser, true, customDelimiter...)
}

func getNodeString(parser string, raw bool, customDelimiter ...string) (string, error) {
	node, err := GetNode(parser, raw, customDelimiter...)
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
func GetNode(parser string, raw bool, customDelimiter ...string) (*yaml.Node, error) {
	delimiter, err := getDelimiter(customDelimiter)
	if err != nil {
		return RootNode, err
	}
	slices := getParserSlice(parser, delimiter)
	result := parseNode(slices, delimiter, RootNode, 0, raw)
	return result.Node, result.Err
}

type parseResult struct {
	Node *yaml.Node
	Err  error
}

func parseNode(slices []string, delimiter string, currentNode *yaml.Node, i int, isRaw bool) parseResult {
	var e error
	ch := make(chan parseResult, MaxMergeInOneLayer)

	// anchor use
	if currentNode.Alias != nil {
		if isRaw {
			currentNode.Value = "*" + currentNode.Value
		} else {
			currentNode = currentNode.Alias
		}
	}

	// anchor definition
	if currentNode.Anchor != "" {
		if !isRaw {
			currentNode.Anchor = ""
		}
	}
	if i >= len(slices) {
		return parseResult{currentNode, nil}
	}
	nodeTag := currentNode.Tag
	switch nodeTag {
	case seqTag:
		index, err := getSequenceNum(slices[i])
		if err != nil {
			return parseResult{currentNode, err}
		}
		return parseNode(slices, delimiter, currentNode.Content[index], i+1, isRaw)
	case mapTag:
		for index, content := range currentNode.Content {
			if index%2 == 1 {
				continue
			}

			if content.Tag == mergeTag {
				ch <- parseNode(slices, delimiter, currentNode.Content[index+1], i, isRaw)
			}
			if content.Value == slices[i] {
				return parseNode(slices, delimiter, currentNode.Content[index+1], i+1, isRaw)
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
		for i := 0; i < MaxMergeInOneLayer; i++ {
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
}
