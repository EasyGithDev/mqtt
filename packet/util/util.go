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

func StringEncode(str string) []byte {
	// size := mp.LengthEncode(len(str))
	size := Uint162bytes(uint16(len(str)))
	var buffer bytes.Buffer
	buffer.Write(size)
	buffer.WriteString(str)

	return buffer.Bytes()
}

// func (mp *MqttPacket) computeLength(buffer []byte) uint16 {
// 	return uint16(len(buffer))
// }
