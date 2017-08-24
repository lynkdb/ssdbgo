// Copyright 2014 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"code.hooto.com/lynkdb/ssdbgo"
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
		fmt.Println("Connect Error:", err)
		return
	}
	defer conn.Close()

	// API::Bool() bool
	conn.Cmd("set", "true", "True")
	if conn.Cmd("get", "true").Bool() {
		fmt.Println("set bool OK")
	}

	conn.Cmd("set", "aa", "val-aaaaaaaaaaaaaaaaaa")
	conn.Cmd("multi_set", []string{
		"bb", "val-bbbbbbbbbbbbbbbbbb",
		"cc", "val-cccccccccccccccccc",
	})
	// API::String() string
	if rs := conn.Cmd("get", "aa"); rs.OK() {
		fmt.Println("get OK\n\t", rs.String())
	}
	if rs := conn.Cmd("get", []byte("aa")); rs.OK() {
		fmt.Println("get bytes OK\n\t", rs.String())
	}

	// API::KvEach()
	if rs := conn.Cmd("multi_get", "aa", "bb"); rs.OK() {
		fmt.Println("multi_get OK")
		rs.KvEach(func(key, value ssdbgo.ResultBytes) {
			fmt.Println("\t", key.String(), value.String())
		})
	}
	// API::KvEach() bytes
	bkeys := [][]byte{[]byte("aa"), []byte("bb")}
	if rs := conn.Cmd("multi_get", bkeys); rs.OK() {
		fmt.Println("multi_get bytes OK")
		rs.KvEach(func(key, value ssdbgo.ResultBytes) {
			fmt.Println("\t", key, value)
		})
	}

	if rs := conn.Cmd("scan", "aa", "cc", 10); rs.OK() {
		fmt.Println("scan OK")
		n := rs.KvEach(func(key, value ssdbgo.ResultBytes) {
			fmt.Println("\t", key.String(), value.String())
		})
		fmt.Println("\t got", n)
	}

	if rs := conn.Cmd("scan", []byte("aa"), []byte("cc"), 10); rs.OK() {
		fmt.Println("scan bytes OK")
		n := rs.KvEach(func(key, value ssdbgo.ResultBytes) {
			fmt.Println("\t", key.String(), value.String())
		})
		fmt.Println("\t got", n)
	}

	conn.Cmd("zset", "z", "a", 3)
	conn.Cmd("multi_zset", "z", "b", -2, "c", 5, "d", 3)
	if rs := conn.Cmd("zrscan", "z", "", "", "", 10); rs.OK() {
		fmt.Println("zrscan OK")
		rs.KvEach(func(key, value ssdbgo.ResultBytes) {
			fmt.Println("\t", key.String(), value.Int())
		})
	}

	conn.Cmd("set", "key", 10)
	if rs := conn.Cmd("incr", "key", 1).Int(); rs > 0 {
		fmt.Println("incr OK\n\t", rs)
	}

	// API::Int() int
	// API::Int64() int64
	conn.Cmd("setx", "key", 123456, 300)
	if rs := conn.Cmd("ttl", "key").Int(); rs > 0 {
		fmt.Println("ttl OK\n\t", rs)
	}

	if rs := conn.Cmd("multi_hset", "zone", "c1", "v-01", "c2", "v-02"); rs.OK() {
		fmt.Println("multi_hset OK")
	}
	if rs := conn.Cmd("multi_hget", "zone", "c1", "c2"); rs.OK() {
		fmt.Println("multi_hget OK")
		rs.KvEach(func(key, value ssdbgo.ResultBytes) {
			fmt.Println("\t", key.String(), value.String())
		})
	}

	// API::Float64() float64
	conn.Cmd("set", "float", 123.456)
	if rs := conn.Cmd("get", "float").Float64(); rs > 0 {
		fmt.Println("float OK\n\t", rs)
	}

	// API::List()
	conn.Cmd("qpush", "queue", "q-1111111111111")
	conn.Cmd("qpush", "queue", "q-2222222222222")
	if rs := conn.Cmd("qpop", "queue", 10); rs.OK() {
		fmt.Println("qpop list OK")
		for i, value := range rs.List() {
			fmt.Println("\t", i, value.String())
		}
	}

	// ssdbgo.Result.JsonDecode(obj interface{}) error
	conn.Cmd("set", "json_key", "{\"name\": \"test obj.name\", \"value\": \"test obj.value\"}")
	if rs := conn.Cmd("get", "json_key"); rs.OK() {
		var rs_obj struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}
		if err := rs.JsonDecode(&rs_obj); err == nil {
			fmt.Println("JsonDecode OK")
			fmt.Println("\tname :", rs_obj.Name)
			fmt.Println("\tvalue:", rs_obj.Value)
		} else {
			fmt.Println("json_key ERR", err)
		}
	}

	// Batch()
	batch := conn.Batch()
	batch.Cmd("set", "batch.aaa", "batch.value.aaa")
	batch.Cmd("set", "batch.bbb", "batch.value.bbb")
	if rs, err := batch.Exec(); err == nil {
		fmt.Println("Batch OK")
		for _, bv := range rs {
			fmt.Println("\t", bv.String())
		}

		fmt.Println("\t", conn.Cmd("get", "batch.aaa").String())
		fmt.Println("\t", conn.Cmd("get", "batch.bbb").String())
	}
}
