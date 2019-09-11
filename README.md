# yquery : a yq style parser for your golang  project

 [![GoDoc](https://godoc.org/github.com/sixleaveakkm/yquery?status.svg)](https://godoc.org/github.com/sixleaveakkm/yquery) [![Build Status](https://travis-ci.org/sixleaveakkm/yquery.svg?branch=master)](https://travis-ci.org/sixleaveakkm/yquery) [![Go Report Card](https://goreportcard.com/badge/github.com/sixleaveakkm/yquery)](https://goreportcard.com/report/github.com/sixleaveakkm/yquery) [![codecov](https://codecov.io/gh/sixleaveakkm/yquery/branch/master/graph/badge.svg)](https://codecov.io/gh/sixleaveakkm/yquery)

## Overview
yquery is a yq style parse to let you handle yaml file without provide data struct
You can **GET** or **SET** item by provide string (e.g., "a.b[0]") in your golang project
This package use [go-yaml v3](https://github.com/go-yaml/yaml/tree/v3) to do the base parse work, thanks for their great job

## Install

```
go get github.com/sixleaveakkm/yquery
```

## Checklist
- [x] able to get item
- [x] able to get item raw data
- [x] able to set exist item with simple struct data
- [x] able to set (add) new item with simple data
- [x] able to set (convert) literal node to map or list 

- [ ] able to set recursive path item with data
- [ ] able to set item with anchor or merge
- [ ] able to handler comment properly
- [ ] provide `Delete`

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

### Initialize
```go
yq, _ := yquery.Unmarshal([]byte(exampleData))
```

### Get Data
```
dataA, err := yq.Get("intA")
fmt.Println(dataA)
// Output: 111

dataB, _ := yq.Get("stringB")
fmt.Println(dataB)
// Output: this is a string
```

### Use Self Defined Delimiter
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

### Get Object
```go
dataBinC, _ := yq.Get("C")
fmt.Println(dataBinC)
// Output: B: string b
```

### Get Anchor
```go
dataA, _ := yq.Get("A")
fmt.Println(dataA)
// Output: B: string b
```

### Get Raw Data
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

**Check `example_test.go` for more example**

**Check `yquery_test.go` for edge condition**

## Author

sixleveakkm@gmail.com

## License

Apache 2.0.
