package yquery_test

import (
	"fmt"

	"github.com/sixleaveakkm/yquery"
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
	yq, _ := yquery.Unmarshal([]byte(exampleData))

	dataA, err := yq.Get("intA")
	if err != nil {
		// failed get
	}
	fmt.Println(dataA)
	// Output: 111
}

func ExampleGetString() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	dataB, _ := yq.Get("stringB")
	fmt.Println(dataB)
	// Output: this is a string
}

func ExampleGetMapItem() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	dataD, _ := yq.Get("mapC.intD")
	fmt.Println(dataD)
	// Output: 222
}

func ExampleGetList() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
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
	yq, _ := yquery.Unmarshal([]byte(data))
	admin, _ := yq.Get("example.com;admin", ";")
	fmt.Println(admin)
	// Output: admin@example.com
}

func ExampleGetAnchor() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	// skip error handle

	dataBinC, _ := yq.Get("C")
	fmt.Println(dataBinC)
	// Output: B: string b
}

func ExampleGetAnchorOrigin() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	// skip error handle

	dataA, _ := yq.Get("A")
	fmt.Println(dataA)
	// Output: B: string b
}

func ExampleGetAnchorRaw() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	// skip error handle
	rawC, _ := yq.GetRaw("C")
	fmt.Println(rawC)
	// Output: *anchorA
}

func ExampleGetAnchorRawOrigin() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	// skip error handle

	dataA, _ := yq.GetRaw("A")
	fmt.Println(dataA)
	// Output: &anchorA
	// B: string b
}

func ExampleGetAstString() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	// skip error handle
	dataD, _ := yq.Get("D")
	fmt.Println(dataD)
	// Output: *anchorA
}
