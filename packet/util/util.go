// MIT License

// Copyright (c) 2022 Florent Brusciano

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package util

import (
	"bytes"
	"encoding/binary"
)

func Uint162bytes(val uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, val)
	return buf
}

func Bytes2uint16(val []byte) uint16 {
	return binary.BigEndian.Uint16(val)
}

func StringEncode(str string) []byte {

	size := Uint162bytes(uint16(len(str)))
	var buffer bytes.Buffer
	buffer.Write(size)
	buffer.WriteString(str)

	return buffer.Bytes()
}

func StringDecode(b []byte) (int, string) {

	buffer := bytes.NewBuffer(b)

	buffSize := buffer.Next(2)

	size := Bytes2uint16(buffSize)

	buffStr := buffer.Next(int(size))

	return 2 + len(buffStr), string(buffStr)
}
