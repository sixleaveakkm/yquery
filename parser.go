package yquery

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// functions handle parse string
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
