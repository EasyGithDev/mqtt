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
package payload

import (
	"github.com/easygithdev/mqtt/packet/util"
)

type MqttPayload struct {
	// Length	(2 bytes)
	// length of payload
	// Length uint16

	// Payload
	Payload []byte
	// Payload string
}

func NewMqttPayload() *MqttPayload {
	return &MqttPayload{}
}

func (mp *MqttPayload) Encode() []byte {
	return mp.Payload
}

func (mp *MqttPayload) Len() int {
	return len(mp.Payload)
}

func (mp *MqttPayload) AddString(str string) {
	mp.Payload = append(mp.Payload, util.StringEncode(str)...)
}

func (mp *MqttPayload) AddQos(qos byte) {
	mp.Payload = append(mp.Payload, []byte{0}...)
}

func (mp *MqttPayload) ShowMessage(start int, end int) string {
	str := ""
	buffer := mp.Payload[start:end]
	len := len(buffer)
	for len > 0 {
		n, s := util.StringDecode(buffer[:len])
		str += s
		len -= n
	}

	return str
}
