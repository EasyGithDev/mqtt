package util

import (
	"bytes"
	"encoding/binary"
	"log"
)

func Uint162bytes(val uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, val)
	return buf
}

func Bytes2uint16(val []byte) uint16 {
	return binary.LittleEndian.Uint16(val)
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

	buffSize := make([]byte, 2)
	buffer.Read(buffSize)

	size := Bytes2uint16(buffSize)

	buffStr := make([]byte, size)
	n, err := buffer.Read(buffStr)

	if err != nil {
		log.Printf("String decode: %s %s \n", buffStr, err)
	}

	return n, string(buffStr[:n])
}

// func (mp *MqttPacket) computeLength(buffer []byte) uint16 {
// 	return uint16(len(buffer))
// }
