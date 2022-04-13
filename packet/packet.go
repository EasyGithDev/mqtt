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
	"fmt"

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
	Hexa() string
}

func Encode(pContent PaquetContent) []byte {
	return pContent.Encode()
}

func Len(pContent PaquetContent) int {
	return pContent.Len()
}

func Hexa(pContent PaquetContent) string {
	return pContent.Hexa()
}

/////////////////////////////////////////////////
// MqttPacket
/////////////////////////////////////////////////

type MqttPacket struct {
	Header         *header.MqttHeader
	VariableHeader PaquetContent
	Payload        PaquetContent
}

type OptionPacket func(mp *MqttPacket)

func NewMqttPacket(header *header.MqttHeader, opts ...OptionPacket) *MqttPacket {
	mp := &MqttPacket{Header: header}

	for _, applyOpt := range opts {
		if applyOpt != nil {
			applyOpt(mp)
		}
	}

	return mp
}

func WithVariableHeader(variableHeader PaquetContent) OptionPacket {
	return func(mp *MqttPacket) {
		mp.VariableHeader = variableHeader
	}
}

func WithPayload(payload PaquetContent) OptionPacket {
	return func(mp *MqttPacket) {
		mp.Payload = payload
	}
}

func (mp *MqttPacket) Encode() []byte {

	// compute the fields
	var mqttBuffer bytes.Buffer

	var vhLen, pLen int = 0, 0

	if mp.VariableHeader != nil {
		vhLen = PaquetContent.Len(mp.VariableHeader)
	}

	if mp.Payload != nil {
		pLen = PaquetContent.Len(mp.Payload)
	}

	mp.Header.RemainingLength = header.RemainingLengthEncode(vhLen + pLen)

	mqttBuffer.Write(PaquetContent.Encode(mp.Header))

	if vhLen > 0 {
		mqttBuffer.Write(PaquetContent.Encode(mp.VariableHeader))
	}

	if pLen > 0 {
		mqttBuffer.Write(PaquetContent.Encode(mp.Payload))
	}

	return mqttBuffer.Bytes()
}

func Decode(data []byte) *MqttPacket {

	var mp *MqttPacket = nil

	bb := bytes.NewBuffer(data)
	control, _ := bb.ReadByte()

	// check the first byte
	switch data[0] {
	case header.CONNECT:
	case header.CONNACK:

		header := header.New(header.WithControl(control), header.WithRemainingLength(2))
		header.Control = control
		rl, _ := bb.ReadByte()
		header.RemainingLength = []byte{rl}
		vHeader := vheader.NewGenericHeader(bb.Bytes())
		mp = NewMqttPacket(header, WithVariableHeader(vHeader))

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

	return mp
}

func (mp *MqttPacket) String() string {

	strHeader := "****************\tHeader\t****************\n" +
		mp.Header.String() +
		fmt.Sprintf("\nLen:%d bytes", mp.Header.Len()) +
		fmt.Sprintf("\nHexa:%s", mp.Header.Hexa())

	strVheader := "\n****************\tvHeader\t****************\n"
	if mp.VariableHeader != nil {
		strVheader += mp.VariableHeader.String() +
			fmt.Sprintf("\nLen:%d bytes", mp.VariableHeader.Len()) +
			fmt.Sprintf("\nHexa:%s", mp.VariableHeader.Hexa())

	} else {
		strVheader += "No variable header"
	}

	strPayload := "\n****************\tPayload\t****************\n"
	if mp.Payload != nil {
		strPayload += mp.Payload.String() +
			fmt.Sprintf("\nLen:%d bytes", mp.Payload.Len()) +
			fmt.Sprintf("\nHexa:%s", mp.Payload.Hexa())

	} else {
		strPayload += "No payload"
	}

	return "\n****************\t" + header.ControlToString(mp.Header.Control) + "\t****************\n" +
		strHeader +
		strVheader +
		strPayload

}
