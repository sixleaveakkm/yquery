package yquery

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

var yq YQuery

func TestMain(m *testing.M) {
	_, err := yq.New([]byte(data))
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
		fmt.Printf("Testing: %s\n", c.Parser)
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
