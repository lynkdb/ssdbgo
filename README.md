# ssdbgo

ssdbgo is a Go Client for SSDB (http://ssdb.io).

Some features include
* minimal, high performance
* support for all SSDB types and commands
* connection pooling support
* thread safe (goroutine safe)

## Start
* ssdbgo use ssdbgo.NewConnector(ssdbgo.Config{...}) to create connection with SSDB server. You can use ssdbgo.Config to set host, port, pool size, timeout, etc.

``` go
package main

import (
	"github.com/lynkdb/ssdbgo"
)

func main() {

	conn, err := ssdbgo.NewConnector(ssdbgo.Config{
		Host:    "127.0.0.1",
		Port:    6380,
		Timeout: 3, // timeout in second, default to 10
		MaxConn: 1, // max connection number, default to 1
		// Auth:    "foobared",
	})
	if err != nil {
		return
	}

	conn.Cmd("set", "key", "value")
	conn.Close()
}
```


Request: all SSDB operations go with ```ssdbgo.Connector.Cmd()```, it accepts variable arguments. The first argument of Cmd() is the SSDB command, for example "get", "set", etc. The rest arguments(maybe none) are the arguments of that command.

Examples:
``` go
conn.Cmd("set", "key", "value")
conn.Cmd("incr", "key-incr", 1)
conn.Cmd("hset", "name-hash", "key-1", "value-1")
```

## Response

the ssdbgo.Connector.Cmd() method will return an Object of ssdbgo.Result

### Response Status

The element of ssdbgo.Result.Status is the response code, ```"ok"``` means the current command are valid results. The response code may be ```"not_found"``` if you are calling "get" on an non-exist key. all of the codes include:

``` go
const (
	ResultOK          = "ok"
	ResultNotFound    = "not_found"
	ResultError       = "error"
	ResultFail        = "fail"
	ResultClientError = "client_error"
)
```

### Response Data

use the following method to get a dynamic data type what you want to need.

* ssdbgo.Result.Bytes() []byte
* ssdbgo.Result.String() string
* ssdbgo.Result.Bool() bool
* ssdbgo.Result.Int() int
* ssdbgo.Result.Int8() int8
* ssdbgo.Result.Int16() int16
* ssdbgo.Result.Int32() int32
* ssdbgo.Result.Int64() int64
* ssdbgo.Result.Uint() uint
* ssdbgo.Result.Uint8() uint8
* ssdbgo.Result.Uint16() uint16
* ssdbgo.Result.Uint32() uint32
* ssdbgo.Result.Uint64() uint64
* ssdbgo.Result.Float32() float32
* ssdbgo.Result.Float64() float64
* ssdbgo.Result.List() []ssdbgo.ResultBytes
* ssdbgo.Result.KvEach(fn func(key, value ssdbgo.ResultBytes)) int
* ssdbgo.Result.KvLen() int
* ssdbgo.Result.JsonDecode(obj interface{}) error

Examples:

``` go
// example 1
if rs := conn.Cmd("incr", "key-incr", 1); rs.OK() {
	fmt.Println("return int", rs.Int())
}

// example 2
var rsobject struct {
	Name string `json:"name"`
}
if rs := conn.Cmd("get", "key-json"); rs.OK() {
	if err := rs.JsonDecode(&rsobject); err == nil {
		fmt.Println("return json.name", rsobject.Name)
	}
}
```

the more examples of result.APIs can visit: [example/example.go](<example/example.go>)


## Refer URLs
* [Official API documentation](http://ssdb.io/docs/) to checkout a complete list of all avilable commands.

## Licensing 
Licensed under the Apache License, Version 2.0

