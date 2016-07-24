# Bloblog

[![GoDoc](https://godoc.org/github.com/zhuharev/bloblog?status.svg)](https://godoc.org/github.com/zhuharev/bloblog)

Bloblog ― append key-value database with automatic autoincremental numeric keys

## Usage

```go
import "gopkg.in/zhuharev/bloblog.v2"
...

bl, e := bloblog.Open("db.bl")
if e!=nil {
	panic(e)
}

id, e := bl.Insert([]byte("Hello"))
if e!=nil {
	panic(e)
}

// will print 1, if "db.bl" new database
println(id)

data, e := bl.Get(1)

// will print "Hello"
fmt.Println(string(data))
```