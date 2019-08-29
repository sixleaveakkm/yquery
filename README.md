# yquery : a yq style parser for your golang  project

## Overview [![GoDoc](https://godoc.org/github.com/sixleaveakkm/yquery?status.svg)](https://godoc.org/github.com/sixleaveakkm/yquery) [![Build Status](https://travis-ci.org/sixleaveakkm/yquery.svg?branch=master)](https://travis-ci.org/sixleaveakkm/yquery)

yquery is a yq style parse to let you handle yaml file without provide data struct
You get get string item by provide string (e.g., "a.b[0]") in your golang project
This package use [go-yaml v3](https://github.com/go-yaml/yaml/tree/v3) to do the base parse work, thanks for their great job

## Install

```
go get github.com/sixleaveakkm/yquery
```

## Example

```
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
```

## Initialize
```go
yq, _ := yquery.Unmarshal([]byte(exampleData))
```

## Get Data
```
dataA, err := yq.Get("intA")
fmt.Println(dataA)
// Output: 111

dataB, _ := yq.Get("stringB")
fmt.Println(dataB)
// Output: this is a string
```

## Get Data From map
```go
dataD, _ := yq.Get("mapC.intD")
fmt.Println(dataD)
// Output: 222
```

## Get Data from list
```
yq, _ := yquery.Unmarshal([]byte(exampleData))
// list index starts from 0
dataF2, _ := yq.Get("mapC.listF[1]")
fmt.Println(dataF2)
// Output: list item 2
```

## Use Self Defined Delimiter
It use go-yaml v3 to do the base parse job, which not support
key like: `[example.com]` yet.
```go
	data := `
example.com:
  admin: admin@example.com
`
yq, _ := yquery.Unmarshal([]byte(data))
admin, _ := yq.Get("example.com;admin", ";")
fmt.Println(admin)
// Output: admin@example.com
```

## Get Object
```go
dataBinC, _ := yq.Get("C")
fmt.Println(dataBinC)
// Output: B: string b
```

## Get Anchor
```go
dataA, _ := yq.Get("A")
fmt.Println(dataA)
// Output: B: string b
```

## Go Raw Data
```go
yq, _ := yquery.Unmarshal([]byte(exampleData))
// skip error handle
rawC, _ := yq.GetRaw("C")
fmt.Println(rawC)
// Output: *anchorA

dataA, _ := yq.GetRaw("A")
fmt.Println(dataA)
// Output: &anchorA
// B: string b
```


## Author

sixleveakkm@gmail.com

## License

Apache 2.0.
