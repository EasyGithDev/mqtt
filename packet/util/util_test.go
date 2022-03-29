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
