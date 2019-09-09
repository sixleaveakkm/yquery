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

func ExampleGetValueInAnchor() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	dataAB, _ := yq.Get("A.B")
	fmt.Println(dataAB)
	// Output: string b
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

func ExampleSetInt() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))

	_ = yq.Set("intA", "333")
	newA, _ := yq.Get("intA")
	fmt.Println(newA)
	// Output: 333
}

func ExampleSetString() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	_ = yq.Set("stringB", "string modified")
	dataB, _ := yq.Get("stringB")
	fmt.Println(dataB)
	// Output: string modified
}

func ExampleSetMapItem() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	_ = yq.Set("mapC.intD", "555")
	dataD, _ := yq.Get("mapC.intD")
	fmt.Println(dataD)
	// Output: 555
}

func ExampleSetList() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	_ = yq.Set("mapC.listF[0]", "item to be 0")
	dataF2, _ := yq.Get("mapC.listF[0]")
	fmt.Println(dataF2)
	// Output: item to be 0
}

func ExampleSetAnchor() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	// skip error handle
	_ = yq.Set("A.B", "new b")
	dataBinA, _ := yq.Get("A.B")
	dataBinC, _ := yq.Get("C.B")
	fmt.Println(dataBinA)
	fmt.Println(dataBinC)
	// Output: new b
	// new b
}

func ExampleSetAnchorReferenceError() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	// skip error handle
	err := yq.Set("C.B", "new b")
	fmt.Printf("%s\n", err)
	// Output: the item 'C.B' is unable to write because it is in an anchor reference or is an item of the merged item
}
