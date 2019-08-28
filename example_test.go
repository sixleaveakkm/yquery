package yquery

import (
	"fmt"
)

// language=yaml
var exampleData = `
intA: 111
stringB: this is a string
mapC:
  intD: 222
  stringE: string e
  listF:
    - list item 1
    - list item 2 
A: &anchorA
  B: string b
C: *anchorA

D: "*anchorA"
`

func ExampleGetInt() {

	var yq YQuery
	_, err := yq.New([]byte(exampleData))
	if err != nil {
		// failed to unmarshal data
	}

	dataA, err := yq.Get("intA")
	if err != nil {
		// failed get
	}
	fmt.Println(dataA)
	// Output: 111
}

func ExampleGetString() {
	var yq YQuery
	_, err := yq.New([]byte(exampleData))
	if err != nil {
		// failed to unmarshal data
	}
	dataB, _ := yq.Get("stringB")
	fmt.Println(dataB)
	// Output: this is a string
}

func ExampleGetMapItem() {
	var yq YQuery
	_, err := yq.New([]byte(exampleData))
	if err != nil {
		// failed to unmarshal data
	}
	dataD, _ := yq.Get("mapC.intD")
	fmt.Println(dataD)
	// Output: 222
}

func ExampleGetList() {
	var yq YQuery
	_, err := yq.New([]byte(exampleData))
	if err != nil {
		// failed to unmarshal data
	}
	// list index starts from 0
	dataF2, _ := yq.Get("mapC.listF[1]")
	fmt.Println(dataF2)
	// Output: list item 2
}

func ExampleGetWithDelimiter() {
	data := `
example.com:
  admin: admin@example.com
`
	yq := &YQuery{}
	_, _ = yq.New([]byte(data))
	admin, _ := yq.Get("example.com;admin", ";")
	fmt.Println(admin)
	// Output: admin@example.com
}

func ExampleGetAnchor() {
	var yq YQuery
	_, _ = yq.New([]byte(exampleData))
	// skip error handle

	dataBinC, _ := yq.Get("C")
	fmt.Println(dataBinC)
	// Output: B: string b
}

func ExampleGetAnchorOrigin() {
	var yq YQuery
	_, _ = yq.New([]byte(exampleData))
	// skip error handle

	dataA, _ := yq.Get("A")
	fmt.Println(dataA)
	// Output: B: string b
}

func ExampleGetAnchorRaw() {
	var yq YQuery
	_, _ = yq.New([]byte(exampleData))
	// skip error handle
	rawC, _ := yq.GetRaw("C")
	fmt.Println(rawC)
	// Output: *anchorA
}

func ExampleGetAnchorRawOrigin() {
	var yq YQuery
	_, _ = yq.New([]byte(exampleData))
	// skip error handle

	dataA, _ := yq.GetRaw("A")
	fmt.Println(dataA)
	// Output: &anchorA
	// B: string b
}

func ExampleGetAstString() {
	var yq YQuery
	_, _ = yq.New([]byte(exampleData))
	// skip error handle
	dataD, _ := yq.Get("D")
	fmt.Println(dataD)
	// Output: *anchorA
}
