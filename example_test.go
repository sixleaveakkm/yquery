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

// language=yaml
var dataToSet = `newData: this is a new string`

func ExampleGetInt() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))

	dataA, err := yq.Get("intA")
	rawA, _ := yq.GetRaw("intA")
	if err != nil {
		// failed get
	}
	fmt.Println(dataA)
	fmt.Println(rawA)
	// Output: 111
	// 111
}

func ExampleGetString() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	dataB, _ := yq.Get("stringB")
	rawB, _ := yq.GetRaw("stringB")
	fmt.Println(dataB)
	fmt.Println(rawB)
	// Output: this is a string
	// this is a string
}

func ExampleGetMapItem() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	dataD, _ := yq.Get("mapC.intD")
	rawDataD, _ := yq.GetRaw("mapC.intD")
	fmt.Println(dataD)
	fmt.Println(rawDataD)
	// Output: 222
	// 222
}

func ExampleGetList() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	// list index starts from 0
	dataF2, _ := yq.Get("mapC.listF[1]")
	rawF2, _ := yq.GetRaw("mapC.listF[1]")
	fmt.Println(dataF2)
	fmt.Println(rawF2)
	// Output: list item 2
	// list item 2
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

func ExampleGetAnchorReference() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	dataBinC, _ := yq.Get("C")
	fmt.Println(dataBinC)
	rawC, _ := yq.GetRaw("C")
	fmt.Println(rawC)
	// Output: B: string b
	// *anchorA
}

func ExampleGetAnchorDefine() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	dataA, _ := yq.Get("A")
	rawA, _ := yq.GetRaw("A")
	fmt.Println(dataA)
	fmt.Println("---")
	fmt.Println(rawA)
	// Output: B: string b
	// ---
	// &anchorA
	// B: string b
}

func ExampleGetValueInAnchor() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	dataAB, _ := yq.Get("A.B")
	dataCB, _ := yq.Get("C.B")
	fmt.Println(dataAB)
	fmt.Println(dataCB)
	// Output: string b
	// string b
}

func ExampleGetAstString() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	// skip error handle
	dataD, _ := yq.Get("D")
	rawD, _ := yq.GetRaw("D")
	fmt.Println(dataD)
	fmt.Println(rawD)
	// Output: *anchorA
	// *anchorA
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

func ExampleSetAddItem() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	_ = yq.Set("notExist", "new value")
	newItem, _ := yq.Get("notExist")
	fmt.Println(newItem)
	// Output: new value
}

func ExampleSetMapItem() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	_ = yq.Set("mapC.intD", "555")
	dataD, _ := yq.Get("mapC.intD")
	fmt.Println(dataD)
	// Output: 555
}

func ExampleSetMapNewItem() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	_ = yq.Set("mapC.newItem", "555")
	newItem, _ := yq.Get("mapC.newItem")
	fmt.Println(newItem)
	// Output: 555
}

func ExampleSetList() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	_ = yq.Set("mapC.listF[0]", "item to be 0")
	dataF1, _ := yq.Get("mapC.listF[0]")
	fmt.Println(dataF1)
	// Output: item to be 0
}

func ExampleSetAddStruct() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	_ = yq.Set("G", dataToSet)
	GNewData, _ := yq.Get("G.newData")
	fmt.Println(GNewData)
	// Output: this is a new string
}

func ExampleSetListNewItem() {
	yq, _ := yquery.Unmarshal([]byte(exampleData))
	_ = yq.Set("mapC.listF[2]", "new item 3")
	dataF3, _ := yq.Get("mapC.listF[2]")
	fmt.Println(dataF3)
	// Output: new item 3
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
	// Output: the item 'C' reaches an anchor reference. You can not modify value from anchor reference
}
