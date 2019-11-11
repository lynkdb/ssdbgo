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
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type Client struct {
	sock     *net.TCPConn
	recv_buf bytes.Buffer
}

func Connect(ip string, port int) (*Client, error) {

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, err
	}

	sock, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &Client{sock: sock}, nil
}

func (c *Client) Cmd(args ...interface{}) *Result {

	r := &Result{
		Status: ResultError,
	}

	if err := c.send(args); err != nil {
		r.Status = ResultFail
		r.Items = []ResultBytes{[]byte(err.Error())}
		return r
	}

	resp, err := c.recv()
	if err != nil || len(resp) < 1 {
		r.Status = ResultFail
		if err != nil {
			r.Items = []ResultBytes{[]byte(err.Error())}
		} else {
			r.Items = []ResultBytes{[]byte("network error")}
		}
		return r
	}

	switch resp[0].String() {

	case ResultOK, ResultNotFound, ResultError, ResultFail, ResultClientError:
		r.Status = resp[0].String()

	default:
		r.Status = ResultFail
		r.Items = []ResultBytes{}
	}

	if r.Status == ResultOK && len(resp) > 1 {
		r.Items = append(r.Items, resp[1:]...)
	}

	return r
}

func send_buf(args []interface{}) ([]byte, error) {

	var buf bytes.Buffer

	for _, arg := range args {

		var s string

		switch argt := arg.(type) {

		case string:
			s = argt

		case []byte:
			send_buf_bs(&buf, argt)
			continue

		case [][]byte:
			for _, bs := range argt {
				send_buf_bs(&buf, bs)
			}
			continue

		case []string:
			for _, s := range argt {
				send_buf_ss(&buf, &s)
			}
			continue

		case int:
			s = strconv.FormatInt(int64(argt), 10)

		case int8:
			s = strconv.FormatInt(int64(argt), 10)

		case int16:
			s = strconv.FormatInt(int64(argt), 10)

		case int32:
			s = strconv.FormatInt(int64(argt), 10)

		case int64:
			s = strconv.FormatInt(argt, 10)

		case uint:
			s = strconv.FormatUint(uint64(argt), 10)

		case uint8:
			s = strconv.FormatUint(uint64(argt), 10)

		case uint16:
			s = strconv.FormatUint(uint64(argt), 10)

		case uint32:
			s = strconv.FormatUint(uint64(argt), 10)

		case uint64:
			s = strconv.FormatUint(argt, 10)

		case float32:
			s = strconv.FormatFloat(float64(argt), 'f', -1, 32)

		case float64:
			s = strconv.FormatFloat(argt, 'f', -1, 64)

		case bool:
			if argt {
				s = "1"
			} else {
				s = "0"
			}

		case nil:
			s = ""

		default:
			return []byte{}, errors.New("bad arguments")
		}

		send_buf_ss(&buf, &s)
	}

	buf.WriteByte('\n')

	return buf.Bytes(), nil
}

func send_buf_bs(buf *bytes.Buffer, data []byte) {
	buf.WriteString(strconv.FormatInt(int64(len(data)), 10))
	buf.WriteByte('\n')
	buf.Write(data)
	buf.WriteByte('\n')
}

func send_buf_ss(buf *bytes.Buffer, data *string) {
	buf.WriteString(strconv.FormatInt(int64(len(*data)), 10))
	buf.WriteByte('\n')
	buf.WriteString(*data)
	buf.WriteByte('\n')
}

func (c *Client) send(args []interface{}) error {

	buf, err := send_buf(args)
	if err == nil {
		_, err = c.sock.Write(buf)
	}

	return err
}

func (c *Client) recv() ([]ResultBytes, error) {

	var buf [8192]byte

	for {

		if resp := c.parse(); resp == nil || len(resp) > 0 {
			return resp, nil
		}

		n, err := c.sock.Read(buf[0:])
		if err != nil {
			return nil, err
		}

		c.recv_buf.Write(buf[0:n])
	}
}

func (c *Client) parse() []ResultBytes {

	var (
		resp   = []ResultBytes{}
		buf    = c.recv_buf.Bytes()
		idx    = 0
		offset = 0
	)

	for {

		if idx = bytes.IndexByte(buf[offset:], '\n'); idx == -1 {
			break
		}

		p := buf[offset : offset+idx]
		offset += idx + 1

		if len(p) == 0 || (len(p) == 1 && p[0] == '\r') {

			if len(resp) == 0 {
				continue
			} else {
				c.recv_buf.Next(offset)
				return resp
			}
		}

		size, err := strconv.Atoi(string(p))
		if err != nil || size < 0 {
			return nil
		}
		if offset+size >= c.recv_buf.Len() {
			break
		}

		resp = append(resp, ResultBytes(bytesClone(buf[offset:offset+size])))
		offset += size + 1
	}

	return []ResultBytes{}
}

// Close The Client Connection
func (c *Client) Close() error {

	if c.sock == nil {
		return nil
	}

	return c.sock.Close()
}

func bytesClone(src []byte) []byte {

	dst := make([]byte, len(src))
	copy(dst, src)

	return dst
}
