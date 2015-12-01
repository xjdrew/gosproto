# gosproto
[sproto](https://github.com/cloudwu/sproto)'s encoder and decoder in golang.

# type map
sproto type      | golang type
---------------- | -------------------------------------------------
string           | \*string, []byte
integer          | \*int8, \*uint8, \*int16, \*uint16, \*int32, \*uint32, \*int64, \*uint64, \*int, \*uint
boolean          | \*bool
object           | \*struct
array of string  | []string
array of integer | []int8, []uint8, []int16, []uint16, []int32, []uint32, []int64, []uint64, []int, []uint
array of boolean | []bool
array of object  | []\*struct

# test
```
go test github.com/xjdrew/gosproto
```

# benchmark
```
$ go test -bench . github.com/xjdrew/gosproto
PASS
BenchmarkEncode-4         300000      4122 ns/op
BenchmarkDecode-4         300000      5417 ns/op
BenchmarkEncodePacked-4122  300000      4712 ns/op
BenchmarkDecodePacked-4712  200000      6267 ns/op
ok      github.com/xjdrew/gosproto5.752s
```
