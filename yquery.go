package yquery

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// RootNode is the root node holds unmarshal struct
// It has type of Node from gopkg.in/yaml.v3 , you can operate it directly if you want
var RootNode *yaml.Node

// MaxMergeInOneLayer is the number that
var MaxMergeInOneLayer = 3

const (
	seqTag   = "!!seq"
	mapTag   = "!!map"
	mergeTag = "!!merge"
)

// New loads bytes data and unmarshal to node struct
// It use RootNode to store data, which type is *yaml.Node, comes from gopkg.in/yaml.v3
func New(in []byte) error {
	node := yaml.Node{}
	err := yaml.Unmarshal(in, &node)
	if err != nil {
		return err
	}
	RootNode = node.Content[0]
	return nil
}

// Marshal the node struct to bytes
// Wrapper of gopkg.in/yaml.v3
func Marshal() ([]byte, error) {
	return yaml.Marshal(RootNode)
}

// Get corresponding item string.
// Receives a parser string (e.g. "a.b") with optional delimiter character.
// For a struct like following,
// ```
// example.com:
//   admin: admin@example.com
// ```
// there is no way to know "example.com.admin" means "admin" in "example.com" or "admin" in "com" in "example"
// Currently go yaml v3 seems don't support key string with bracket, e.g. "[example.com]"
// therefore you could provide a custom delimiter, e.g. `Get("example.com;admin",";")`
func Get(parser string, customDelimiter ...string) (string, error) {
	node, err := GetNode(parser, customDelimiter...)
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
func GetNode(parser string, customDelimiter ...string) (*yaml.Node, error) {
	delimiter, err := getDelimiter(customDelimiter)
	if err != nil {
		return RootNode, err
	}
	slices := getParserSlice(parser, delimiter)
	result := parseNode(slices, delimiter, RootNode, 0)
	return result.Node, result.Err
}

type parseResult struct {
	Node *yaml.Node
	Err  error
}

func parseNode(slices []string, delimiter string, currentNode *yaml.Node, i int) parseResult {
	var e error
	ch := make(chan parseResult, MaxMergeInOneLayer)
	if i >= len(slices) {
		return parseResult{currentNode, nil}
	}

	if currentNode.Alias != nil {
		currentNode = currentNode.Alias
	}
	nodeTag := currentNode.Tag
	switch nodeTag {
	case seqTag:
		index, err := getSequenceNum(slices[i])
		if err != nil {
			return parseResult{currentNode, err}
		}
		return parseNode(slices, delimiter, currentNode.Content[index], i+1)
	case mapTag:
		for index, content := range currentNode.Content {
			if index%2 == 1 {
				continue
			}

			if content.Tag == mergeTag {
				ch <- parseNode(slices, delimiter, currentNode.Content[index+1], i)
			}
			if content.Value == slices[i] {
				return parseNode(slices, delimiter, currentNode.Content[index+1], i+1)
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

func getSequenceNum(slice string) (int, error) {
	regex := regexp.MustCompile(`\[(\d+)]`)
	result := regex.FindStringSubmatch(slice)
	if len(result) == 2 {
		return strconv.Atoi(result[1])
	}
	return 0, fmt.Errorf("cannot match %s to index", slice)
}

func getDelimiter(option []string) (string, error) {
	if len(option) > 1 {
		return "", fmt.Errorf("get could only get 0 or 1 string for delimiter, got %s", option)
	}
	if len(option) > 0 {
		deli := option[0]
		if deli == "[" || deli == "]" {
			return "", fmt.Errorf("custom delimiter cannot be '[' or ']'")
		}
		return deli, nil
	}
	return ".", nil
}

func getParserSlice(parser string, delimiter string) []string {
	parser = strings.Replace(parser, "[", delimiter+"[", -1)
	return strings.Split(parser, delimiter)
}
