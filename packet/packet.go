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

	"github.com/easygithdev/mqtt/packet/header"
	"github.com/easygithdev/mqtt/packet/payload"
	"github.com/easygithdev/mqtt/packet/variableheader"
)

type MqttPacket struct {
	Header         *header.MqttHeader
	VariableHeader variableheader.VariableHeader
	Payload        *payload.MqttPayload
}

func NewMqttPacket() *MqttPacket {
	return &MqttPacket{}
}

func (mp *MqttPacket) Encode() []byte {

	// compute the fields

	var mqttBuffer bytes.Buffer

	// mp.VariableHeader.ComputeProtocolLength()
	mp.Header.ComputeRemainingLength(variableheader.Len(mp.VariableHeader) + mp.Payload.Len())

	mqttBuffer.Write(mp.Header.Encode())

	if variableheader.Len(mp.VariableHeader) > 0 {
		mqttBuffer.Write(variableheader.Encode(mp.VariableHeader))
	}

	if mp.Payload.Len() > 0 {
		mqttBuffer.Write(mp.Payload.Encode())
	}

	return mqttBuffer.Bytes()
}

func (mp *MqttPacket) Decode(data []byte) {

	mp.Header.Decode(data)

	// check the first byte
	switch data[0] {
	case header.CONNACK:
	case header.PUBACK:
	case header.SUBACK:

	}

}

// func (mp *MqttPacket) ShowBytes() string {

// 	str := ""
// 	for i := 0; i < len(buffer); i++ {
// 		str += fmt.Sprintf("0x%X ", buffer[i])
// 	}
// 	str += "\n"

// 	return str
// }

// func (mp *MqttPacket) GetPacket(buffer []byte) []byte {
// 	mp.Control = buffer[0]
// 	size := mp.RemaingLengthDecode(buffer[1:2])

// 	return buffer[:size+2]
// }
