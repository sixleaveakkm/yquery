package yquery_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sixleaveakkm/yquery"
)

// language=yaml
var data = `
a: title
# comment with empty line following

b: 112
c: &cPtr
  # comment in c
  d: "d in c"
  e: "e in c"
f: *cPtr

g:
  <<: &gAnchor
      h:
        - li1
        - li2
      h2:
        - 1
        - 2
      i: other item
  g2: other in g

g3:
  <<: &g3Anchor
    h:
      - ui1
      - ui2
j:
    !!merge <<: *gAnchor
    <<: *g3Anchor
    i: override i
    k: other item in j
l: |
  this is a code
  with second line

m: >
  this is a code
  in one line

n:
- 1
- 2
- - 3.1
  - 3.2

o: 2010-10-10T12:34:56Z
`

type casePair struct {
	Parser string
	Value  string
}

var yq *yquery.YQuery

func TestMain(m *testing.M) {
	var err error
	yq, err = yquery.Unmarshal([]byte(data))
	if err != nil {
		log.Fatalf("Failed to parse data to node")
	}
	fmt.Println("Start Testing...")
	os.Exit(m.Run())
}

func TestGetScalars(t *testing.T) {
	asserts := assert.New(t)
	testCases := []casePair{
		{"a", "title"},
		{"c.d", "d in c"},
		{"l", "this is a code\nwith second line\n"},
		{"m", "this is a code in one line\n"},
		{"o", "2010-10-10T12:34:56Z"},
	}
	for _, c := range testCases {
		res, err := yq.Get(c.Parser)
		asserts.NoError(err)
		asserts.Equal(c.Value, res)
	}
}

func TestGetList(t *testing.T) {
	asserts := assert.New(t)
	testCases := []casePair{
		{"n[0]", "1"},
		{"n[1]", "2"},
		{"n[2][0]", "3.1"},
		{"g.h[0]", "li1"},
	}
	for _, c := range testCases {
		res, err := yq.Get(c.Parser)
		asserts.NoError(err)
		asserts.Equal(c.Value, res)
	}
}

func TestGetAnchor(t *testing.T) {
	asserts := assert.New(t)
	testCases := []casePair{
		{"f.d", "d in c"},
	}
	for _, c := range testCases {
		res, err := yq.Get(c.Parser)
		asserts.NoError(err)
		asserts.Equal(c.Value, res)
	}
}

func TestGetMerge(t *testing.T) {
	asserts := assert.New(t)
	testCases := []casePair{
		{"j.h[0]", "ui1"},
		{"j.h2[0]", "1"},
		{"j.i", "override i"},
	}
	for _, c := range testCases {
		res, err := yq.Get(c.Parser)
		asserts.NoError(err)
		asserts.Equal(c.Value, res)
	}
}

func TestGetMergeOverride(t *testing.T) {
	asserts := assert.New(t)
	testCases := []casePair{
		{"j.h", "- ui1\n- ui2"},
		{"j.h[1]", "ui2"},
	}
	for _, c := range testCases {
		res, err := yq.Get(c.Parser)
		asserts.NoError(err)
		asserts.Equal(c.Value, res)
	}
}

func TestBlank(t *testing.T) {
	asserts := assert.New(t)
	_, err := yq.Get("")
	asserts.Error(err)
}

func TestNotExist(t *testing.T) {
	asserts := assert.New(t)
	testCases := []casePair{
		{"", ""},
		{"notExist", ""},
		{"j..i", ""},
		{"n[999]", ""},
		{"n[3]", ""},
		{"n[-1]", ""},
		{"n[2][999]", ""},
		{"g.h.notExist", ""},
		{"c.d.notExist", ""},
	}
	for _, c := range testCases {
		_, err := yq.Get(c.Parser)
		asserts.Error(err)
	}
}

func TestSetStraightShouldError(t *testing.T) {
	asserts := assert.New(t)
	testCases := []casePair{
		{"", "not valid"},
		{"f..d", "not valid"},
		{"n[999]", "index out range"},
		{"n.notIndex", "string in list"},
		{"n[-1]", "index not correct"},
		{"f.d", "setting a anchor reference"},
	}
	for _, c := range testCases {
		err := yq.Set(c.Parser, c.Value)
		asserts.Error(err)
	}
}

func TestSet(t *testing.T) {
	asserts := assert.New(t)
	testCases := []casePair{
		{"g.h", "override map to string"},
		{"j.k.newItem", "override string to map"},
	}
	for _, c := range testCases {
		err := yq.Set(c.Parser, c.Value)
		res, _ := yq.Get(c.Parser)
		asserts.NoError(err)
		asserts.Equal(c.Value, res)
	}
}

func TestStrangeDelimiterError(t *testing.T) {
	asserts := assert.New(t)
	_, err := yq.Get("g;h", "too many args", ";", "-")
	asserts.Error(err)
	_, err = yq.Get("g[h", "forbidden delimiter", "[")
	asserts.Error(err)
	_, err = yq.Get("g]h", "forbidden delimiter", "]")
	asserts.Error(err)
	err = yq.Set("g]h", "forbidden delimiter", yquery.Config{
		Delimiter: "]",
	})
	asserts.Error(err)
	err = yq.Set("g;h", "too many args", yquery.Config{
		Delimiter: ";",
	}, yquery.Config{
		Delimiter: ";",
	})
	asserts.Error(err)
}
