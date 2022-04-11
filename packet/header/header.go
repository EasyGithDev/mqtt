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
package header

import (
	"bytes"
	"fmt"
)

//  MQTT Control Packet type
var CONNECT byte = 0x10     // 16
var CONNACK byte = 0x20     // 32
var PUBLISH byte = 0x30     // 48
var PUBACK byte = 0x40      // 64
var PUBREC byte = 0x50      // 80
var PUBREL byte = 0x60      // 96
var PUBCOMP byte = 0x70     // 112
var SUBSCRIBE byte = 0x80   // 128
var SUBACK byte = 0x90      // 144
var UNSUBSCRIBE byte = 0xA0 // 160
var UNSUBACK byte = 0xB0    //176
var PINGREQ byte = 0xC0     //192
var PINGRESP byte = 0xD0    //208
var DISCONNECT byte = 0xE0  // 224
var AUTH byte = 0xF0        // 240

//  MQTT Flags

// PUBLISH Used in MQTT 3.1.1
// bit 3 -> DUP1
// bit 2 -> QoS2
// bit 1 -> QoS2
// bit 0 -> RETAIN3

// PUBREL -> bit 1
// SUBSCRIBE -> bit 1
// UNSUBSCRIBE -> bit 1

//  MQTT Connection responses

var CONNECT_ACCEPTED byte = 0x00
var CONNECT_REFUSED_1 byte = 0x01
var CONNECT_REFUSED_2 byte = 0x02
var CONNECT_REFUSED_3 byte = 0x03
var CONNECT_REFUSED_4 byte = 0x04
var CONNECT_REFUSED_5 byte = 0x05

// control + length + protocol name + Protocol Level +Connect Flags + keep alive +Payload

type MqttHeader struct {
	// Fixed header
	// 1 byte  = Packet type  (4bits) + flags (4bits)
	Control byte

	// Remaining length 1-4 bytes
	// This is the total length without fixed header
	RemainingLength []byte
}

type OptionHeader func(mh *MqttHeader)

func New(controle byte, opts ...OptionHeader) *MqttHeader {
	mh := &MqttHeader{Control: controle}

	for _, applyOpt := range opts {
		if applyOpt != nil {
			applyOpt(mh)
		}
	}

	return mh
}

func WithRemainingLength(rl int) OptionHeader {
	return func(mh *MqttHeader) {
		mh.RemainingLength = RemainingLengthEncode(rl)
	}
}

func (mh *MqttHeader) Encode() []byte {
	var buffer bytes.Buffer

	// 1 byte
	buffer.WriteByte(mh.Control)

	// 1-4 bytes
	buffer.Write(mh.RemainingLength)

	return buffer.Bytes()
}

func (mh *MqttHeader) Decode(buffer []byte) {
	mh.Control = buffer[0]
	mh.RemainingLength = buffer[1:]
}

func (mh *MqttHeader) Len() int {
	return 1 + len(mh.RemainingLength)
}

func (mh *MqttHeader) String() string {
	return fmt.Sprintf("control: %b \nremainingLength: %b", mh.Control, mh.RemainingLength)
}

func (mh *MqttHeader) UseConnect() {
	mh.Control = CONNECT
}

func (mh *MqttHeader) UseConnack() {
	mh.Control = CONNACK
}

func (mh *MqttHeader) UsePublish() {
	mh.Control = PUBLISH
}

func (mh *MqttHeader) UsePuback() {
	mh.Control = PUBACK
}

func (mh *MqttHeader) UsePubrec() {
	mh.Control = PUBREC
}

func (mh *MqttHeader) UsePubrel() {
	mh.Control = PUBREL | 1<<1
}

func (mh *MqttHeader) UsePubcomp() {
	mh.Control = PUBCOMP
}

func (mh *MqttHeader) UseSubscribe() {
	mh.Control = SUBSCRIBE | 1<<1
}

func (mh *MqttHeader) UseSuback() {
	mh.Control = SUBACK
}

func (mh *MqttHeader) UseUnsubscribe() {
	mh.Control = UNSUBSCRIBE | 1<<1
}

func (mh *MqttHeader) UseUnsuback() {
	mh.Control = UNSUBSCRIBE
}

func (mh *MqttHeader) UsePingreq() {
	mh.Control = PINGREQ
}

func (mh *MqttHeader) UsePingresp() {
	mh.Control = PINGRESP
}

func (mh *MqttHeader) UseDisconnect() {
	mh.Control = DISCONNECT
}

func (mh *MqttHeader) UseAuth() {
	mh.Control = AUTH
}

func (mh *MqttHeader) UseRetain() {
	if mh.Control == PUBLISH {
		mh.Control |= 1 << 1
	}
}

func (mh *MqttHeader) UseQos1() {
	if mh.Control == PUBLISH {
		mh.Control |= 1 << 1
	}
}

func (mh *MqttHeader) UseQos2() {
	if mh.Control == PUBLISH {
		mh.Control |= 1 << 2
	}
}

func (mh *MqttHeader) UseDup() {
	if mh.Control == PUBLISH {
		mh.Control |= 1 << 3
	}
}

func RemainingLengthEncode(x int) []byte {

	var buffer []byte = make([]byte, 0)
	var encodedByte int = 0
	for {
		encodedByte = x % 128

		x = x / 128
		// fmt.Printf("x :%d\n ", x)

		if x > 0 {
			encodedByte = encodedByte | 128
		}

		//output
		buffer = append(buffer, byte(encodedByte))

		// fmt.Println(buffer)

		if x <= 0 {
			break
		}
	}
	return buffer
}

func RemaingLengthDecode(x []byte) (int, int) {

	var multiplier int = 1

	var value int = 0

	var encodedByte byte = 0

	var nbBytes int = 0
	for i := 0; i < len(x); i++ {

		encodedByte = x[i]

		value += int(encodedByte&byte(127)) * multiplier

		multiplier *= 128

		// if (multiplier > 128*128*128)

		//    throw Error(Malformed Remaining Length)

		nbBytes++
		if (encodedByte & 128) == 0 {
			break
		}

	}

	return nbBytes, value
}
