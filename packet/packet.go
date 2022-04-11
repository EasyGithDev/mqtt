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
	"github.com/easygithdev/mqtt/packet/vheader"
)

/////////////////////////////////////////////////
// Interface PacketContent
/////////////////////////////////////////////////

type PaquetContent interface {
	Encode() []byte
	Len() int
	String() string
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

func NewMqttPacket(header PaquetContent, variableHeader PaquetContent, payload PaquetContent) *MqttPacket {
	return &MqttPacket{Header: header, VariableHeader: variableHeader, Payload: payload}
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

	bb := bytes.NewBuffer(data)
	control, _ := bb.ReadByte()

	// check the first byte
	switch data[0] {
	case header.CONNECT:
	case header.CONNACK:

		header := header.NewMqttHeader(2)
		header.Control = control
		header.RemainingLength = make([]byte, 1)
		header.RemainingLength[0], _ = bb.ReadByte()
		vHeader := vheader.NewGenericHeader(bb.Bytes())

		mp.Header = header
		mp.VariableHeader = vHeader

	case header.PUBLISH:
	case header.PUBACK:
	case header.PUBREC:
	case header.PUBREL:
	case header.PUBCOMP:
	case header.SUBSCRIBE:
	case header.SUBACK:
	case header.UNSUBSCRIBE:
	case header.UNSUBACK:
	case header.PINGREQ:
	case header.PINGRESP:
	case header.DISCONNECT:
	case header.AUTH:
	}
}

func (mp *MqttPacket) String() string {
	return "****************\tHeader\t****************\n" +
		mp.Header.String() +
		"\n****************\tvHeader\t****************\n" +
		mp.VariableHeader.String() +
		"\n****************\tPayload\t****************\n" +
		mp.Payload.String()
}
