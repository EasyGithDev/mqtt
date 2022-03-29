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
	"fmt"
	"reflect"
	"testing"
)

func TestUint2Bytes(t *testing.T) {

	str := "hello world$"
	data := []byte(str)

	sizeBuffer := Uint162bytes(uint16(len(data)))

	if len(sizeBuffer) != 2 {
		t.Errorf("Uint2bytes error  found %d; want 2", len(sizeBuffer))
	}

	res := fmt.Sprintf("%b", sizeBuffer)

	if res != "[0 1100]" {
		t.Errorf("Uint2bytes error  found %b; want [0 1100]", sizeBuffer)
	}

}

func TestStringEncode(t *testing.T) {

	str := "hello world$"

	expected := []byte{0, 12, 'h', 'e', 'l', 'l', 'o', 'w', 'o', 'r', 'l', 'd', '$'}

	encoded := StringEncode(str)

	if !reflect.DeepEqual(encoded, encoded) {
		t.Errorf("String encode error found [%b]; want [%b]", encoded, expected)
	}
}

func TestStringDecode(t *testing.T) {

	str := "hello world$"

	encoded := StringEncode(str)
	nb, decoded := StringDecode(encoded)

	if nb != 12 {
		t.Errorf("String decode error found [%d]; want 12", nb)
	}

	if decoded != str {
		t.Errorf("String decode error found [%b]; want [%b]", []byte(decoded), []byte(encoded))
	}

}
