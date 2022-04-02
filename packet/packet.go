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
package packet

import (
	"bytes"
)

/////////////////////////////////////////////////
// Interface PacketContent
/////////////////////////////////////////////////

type PaquetContent interface {
	Encode() []byte
	Len() int
}

func Encode(header PaquetContent) []byte {
	return header.Encode()
}

func Len(header PaquetContent) int {
	return header.Len()
}

/////////////////////////////////////////////////
// MqttPacket
/////////////////////////////////////////////////

type MqttPacket struct {
	Header         PaquetContent
	VariableHeader PaquetContent
	Payload        PaquetContent
}

func NewMqttPacket() *MqttPacket {
	return &MqttPacket{}
}

func (mp *MqttPacket) Encode() []byte {

	// compute the fields
	var mqttBuffer bytes.Buffer

	mqttBuffer.Write(PaquetContent.Encode(mp.Header))

	if mp.VariableHeader != nil {
		if PaquetContent.Len(mp.VariableHeader) > 0 {
			mqttBuffer.Write(PaquetContent.Encode(mp.VariableHeader))
		}
	}

	if mp.Payload != nil {
		if PaquetContent.Len(mp.Payload) > 0 {
			mqttBuffer.Write(PaquetContent.Encode(mp.Payload))
		}
	}

	return mqttBuffer.Bytes()
}

func (mp *MqttPacket) Decode(data []byte) {

	// mp.Header.Decode(data)

	// // check the first byte
	// switch data[0] {
	// case header.CONNACK:
	// case header.PUBACK:
	// case header.SUBACK:

	// }

}
