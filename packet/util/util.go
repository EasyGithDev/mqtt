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
