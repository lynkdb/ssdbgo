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

package ssdbgo // import "github.com/lynkdb/ssdbgo"

import (
	"encoding/json"
	"errors"
	"strconv"
)

const (
	ResultOK          = "ok"
	ResultNotFound    = "not_found"
	ResultError       = "error"
	ResultFail        = "fail"
	ResultClientError = "client_error"
)

type Result struct {
	Status string
	Items  []ResultBytes
}

func (r *Result) OK() bool {
	return r.Status == ResultOK
}

func (r *Result) NotFound() bool {
	return r.Status == ResultNotFound
}

func (r *Result) bytex() *ResultBytes {
	if len(r.Items) > 0 {
		return &r.Items[0]
	}
	return &ResultBytes{}
}

func (r *Result) Bytes() []byte {
	return r.bytex().Bytes()
}

func (r *Result) String() string {
	return r.bytex().String()
}

func (r *Result) Int() int {
	return r.bytex().Int()
}

func (r *Result) Int8() int8 {
	return r.bytex().Int8()
}

func (r *Result) Int16() int16 {
	return r.bytex().Int16()
}

func (r *Result) Int32() int32 {
	return r.bytex().Int32()
}

func (r *Result) Int64() int64 {
	return r.bytex().Int64()
}

func (r *Result) Uint() uint {
	return r.bytex().Uint()
}

func (r *Result) Uint8() uint8 {
	return r.bytex().Uint8()
}

func (r *Result) Uint16() uint16 {
	return r.bytex().Uint16()
}

func (r *Result) Uint32() uint32 {
	return r.bytex().Uint32()
}

func (r *Result) Uint64() uint64 {
	return r.bytex().Uint64()
}

func (r *Result) Float32() float32 {
	return r.bytex().Float32()
}

func (r *Result) Float64() float64 {
	return r.bytex().Float64()
}

func (r *Result) Bool() bool {
	return r.bytex().Bool()
}

func (r *Result) List() []ResultBytes {
	return r.Items
}

func (r *Result) KvLen() int {
	return len(r.Items) / 2
}

func (r *Result) KvEach(fn func(key, value ResultBytes)) int {
	for i := 1; i < len(r.Items); i += 2 {
		fn(r.Items[i-1], r.Items[i])
	}
	return r.KvLen()
}

// Json returns the map that marshals from the reply bytes as json in response .
func (r *Result) JsonDecode(v interface{}) error {
	return r.bytex().JsonDecode(v)
}

// Universal Bytes
type ResultBytes []byte

func (rd ResultBytes) Bytes() []byte {
	return rd
}

func (rd ResultBytes) String() string {
	return string(rd)
}

func (rd ResultBytes) Bool() bool {
	if len(rd) > 0 {
		if b, err := strconv.ParseBool(string(rd)); err == nil {
			return b
		}
	}
	return false
}

// int
func (rd ResultBytes) Int() int {
	return int(rd.Int64())
}

func (rd ResultBytes) Int8() int8 {
	return int8(rd.Int64())
}

func (rd ResultBytes) Int16() int16 {
	return int16(rd.Int64())
}

func (rd ResultBytes) Int32() int32 {
	return int32(rd.Int64())
}

func (rd ResultBytes) Int64() int64 {
	if len(rd) > 0 {
		if i64, err := strconv.ParseInt(string(rd), 10, 64); err == nil {
			return i64
		}
	}
	return 0
}

// unsigned int
func (rd ResultBytes) Uint() uint {
	return uint(rd.Uint64())
}

func (rd ResultBytes) Uint8() uint8 {
	return uint8(rd.Uint64())
}

func (rd ResultBytes) Uint16() uint16 {
	return uint16(rd.Uint64())
}

func (rd ResultBytes) Uint32() uint32 {
	return uint32(rd.Uint64())
}

func (rd ResultBytes) Uint64() uint64 {
	if len(rd) > 0 {
		if i64, err := strconv.ParseUint(string(rd), 10, 64); err == nil {
			return i64
		}
	}
	return 0
}

// float
func (rd ResultBytes) Float32() float32 {
	return float32(rd.Float64())
}

func (rd ResultBytes) Float64() float64 {

	if len(rd) < 1 {
		return 0
	}

	if f64, err := strconv.ParseFloat(string(rd), 64); err == nil {
		return f64
	}

	return 0
}

func (rd ResultBytes) JsonDecode(v interface{}) error {
	if len(rd) < 2 {
		return errors.New("json: invalid format")
	}
	return json.Unmarshal(rd, &v)
}
