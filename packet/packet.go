package packet

import (
	"bytes"
	"encoding/binary"
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

// control + length + protocol name + Protocol Level +Connect Flags + keep alive +Payload

type MqttPacket struct {

	// Fixed header
	// 1 byte  = Packet type  (4bits) + flags (4bits)
	Control byte

	// Length of the paquet in bytes
	// Optionnal
	PacketLength byte

	// Remaining length
	// This is the total length without fixed header
	RemainingLength byte

	// Length of protocol name (expl MQTT4 -> 4)
	ProtocolLength uint16

	// Protocol + version (expl MQTT4)
	ProtocolName    string
	ProtocolVersion byte

	// Connect flag (expl clean session)
	ConnectFlag byte

	// Keep alive (2 bytes)
	KeepAlive uint16

	// Length	(2 bytes)
	// length of payload
	Length uint16

	// Payload
	Payload string
}

func NewMqttPacket() *MqttPacket {
	return &MqttPacket{}
}

func (mp *MqttPacket) Encode() []byte {

	//mp.RemainingLength = 23

	//mp.ProtocolLength = 4
	// mp.ProtocolName = "MQTT"
	// mp.ProtocolVersion = 4
	// mp.ConnectFlag = 0x2
	// mp.KeepAlive = 60
	//mp.Length = 7
	//mp.Payload = payload

	// we start to calcule the strings length

	// the protocol len
	var bufProtocolLen bytes.Buffer
	bufProtocolLen.WriteString(mp.ProtocolName)
	// println("len ProtocoleLen:", len(bufProtocolLen.Bytes()))
	mp.ProtocolLength = mp.computeLength(bufProtocolLen.Bytes())

	// the payload len
	var bufPayload bytes.Buffer
	bufPayload.WriteString(mp.Payload)
	// println("len Payload:", len(bufPayload.Bytes()))
	mp.Length = mp.computeLength(bufPayload.Bytes())

	// Compute the remaining Length

	bodyLength := 0
	if mp.ProtocolLength != 0 {
		bodyLength = 10

	}

	payloadLength := 0
	if mp.Length != 0 {
		payloadLength = 2 + int(mp.Length)
	}
	mp.RemainingLength = byte(mp.RemaingLengthEncode(bodyLength + payloadLength))

	fmt.Printf("Protocol length Hexadecimal: Ox%X Dec:%d\n", mp.ProtocolLength, mp.ProtocolLength)
	fmt.Printf("Payload length Hexadecimal: Ox%X Dec:%d\n", mp.Length, mp.Length)
	fmt.Printf("RemainingLength Dec: %d\n", mp.RemainingLength)

	// compute the fields

	var mqttBuffer bytes.Buffer
	// 1 byte
	mqttBuffer.WriteByte(mp.Control)

	// 1-4 bytes
	mqttBuffer.WriteByte(mp.RemainingLength)

	var buf []byte

	if bodyLength != 0 {

		// 10 bytes = 2 (protocol length) + 4 (protocol name) + 1 (protocol version) + 1 (connect flag) + 2 (keep alive)

		buf = uint162bytes(mp.ProtocolLength)
		mqttBuffer.Write(buf)

		mqttBuffer.WriteString(mp.ProtocolName)

		mqttBuffer.WriteByte(mp.ProtocolVersion)

		mqttBuffer.WriteByte(mp.ConnectFlag)

		buf = uint162bytes(mp.KeepAlive)
		mqttBuffer.Write(buf)
	}

	if payloadLength != 0 {
		// 2 bytes
		buf = uint162bytes(mp.Length)
		mqttBuffer.Write(buf)

		// n bytes
		mqttBuffer.WriteString(mp.Payload)
	}

	return mqttBuffer.Bytes()
}

func (mp *MqttPacket) RemaingLengthEncode(x int) uint16 {

	var encodedByte int = 0
	for {
		encodedByte = x % 128
		fmt.Printf("byte :%d\n", byte(encodedByte))
		x = x / 128
		fmt.Printf("x :%d\n ", x)

		if x > 0 {
			encodedByte = encodedByte | 128
		}

		//output
		fmt.Printf("%d ", byte(encodedByte))
		if x <= 0 {
			break
		}
	}
	return uint16(encodedByte)
}

func (mp *MqttPacket) remaingLengthDecode(x int) uint16 {
	return 0
}

func (mp *MqttPacket) computeLength(buffer []byte) uint16 {
	return uint16(len(buffer))
}

func uint162bytes(val uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, val)
	return buf
}
