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
package vheader

import (
	"fmt"

	"github.com/easygithdev/mqtt/packet/util"
)

var CONNECT_FLAG_CLEAN_SESSION byte = 0x02
var CONNECT_FLAG_WILL_FLAG byte = 0x04
var PUBLCONNECT_FLAG_WILL_QOS_1 byte = 0x08
var PUBLCONNECT_FLAG_WILL_QOS_2 byte = 0x10
var CONNECT_FLAG_WILL_RETAIN byte = 0x20
var CONNECT_FLAG_PASSWORD byte = 0x40
var CONNECT_FLAG_USERNAME byte = 0x80

/////////////////////////////////////////////////
// Empty header
/////////////////////////////////////////////////

type EmptyHeader struct {
}

func NewEmptyHeader() *EmptyHeader {
	return &EmptyHeader{}
}

func (eh *EmptyHeader) Encode() []byte {
	return nil
}

func (eh *EmptyHeader) Len() int {
	return 0
}

/////////////////////////////////////////////////
// Connect header
/////////////////////////////////////////////////

type ConnectHeader struct {

	// Packet Identifier field
	PacketIdentifier uint16

	// Protocol (expl MQTT)
	ProtocolName string

	// Protocol level (expl 4)
	ProtocolVersion byte

	// Connect flag (expl clean session)
	Flag byte

	// Keep alive (2 bytes)
	KeepAlive uint16
}

func NewConnectHeader(protocolName string, protocolVersion byte, flag byte, keepAlive uint16) *ConnectHeader {
	return &ConnectHeader{ProtocolName: protocolName, ProtocolVersion: protocolVersion, Flag: flag, KeepAlive: keepAlive}
}

func (ch *ConnectHeader) Encode() []byte {
	var content []byte

	content = append(content, util.StringEncode(ch.ProtocolName)...)
	content = append(content, []byte{ch.ProtocolVersion}...)
	content = append(content, []byte{ch.Flag}...)
	content = append(content, util.Uint162bytes(ch.KeepAlive)...)

	return content
}

func (ch *ConnectHeader) Len() int {
	return len(ch.Encode())
}

func (ch *ConnectHeader) String() string {
	return fmt.Sprintf("protocol: %s\nversion: %d\nflag: %b\nkeepalive: %d", ch.ProtocolName, ch.ProtocolVersion, ch.Flag, ch.KeepAlive)
}

/////////////////////////////////////////////////
// Subscribe header
/////////////////////////////////////////////////

type SubscribeHeader struct {
	TopicName string
	PacketId  uint16
}

func NewSubscribeHeader(packetId uint16, topicName string) *SubscribeHeader {
	return &SubscribeHeader{PacketId: packetId, TopicName: topicName}
}

func (sh *SubscribeHeader) Encode() []byte {
	var content []byte

	content = append(content, util.Uint162bytes(sh.PacketId)...)
	content = append(content, util.StringEncode(sh.TopicName)...)

	return content
}

func (sh *SubscribeHeader) Len() int {
	return len(sh.Encode())
}

func (sh *SubscribeHeader) String() string {
	return fmt.Sprintf("packetId: %d\ntopicName: %s", sh.PacketId, sh.TopicName)
}

/////////////////////////////////////////////////
// Publish header
/////////////////////////////////////////////////

type PublishHeader struct {
	TopicName string
}

func NewPublishHeader(topicName string) *PublishHeader {
	return &PublishHeader{TopicName: topicName}
}

func (ph *PublishHeader) Encode() []byte {

	var content []byte

	content = append(content, util.StringEncode(ph.TopicName)...)

	return content
}

func (ph *PublishHeader) Len() int {
	return len(ph.Encode())
}

func (ph *PublishHeader) String() string {
	return fmt.Sprintf("topicName: %s", ph.TopicName)
}
