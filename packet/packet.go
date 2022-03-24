package packet

import (
	"bytes"
	"encoding/binary"
	"log"
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

//  MQTT Connection responses

var CONNECT_ACCEPTED byte = 0x00
var CONNECT_REFUSED_1 byte = 0x01
var CONNECT_REFUSED_2 byte = 0x02
var CONNECT_REFUSED_3 byte = 0x03
var CONNECT_REFUSED_4 byte = 0x04
var CONNECT_REFUSED_5 byte = 0x05

// control + length + protocol name + Protocol Level +Connect Flags + keep alive +Payload

type MqttPacket struct {

	// Fixed header
	// 1 byte  = Packet type  (4bits) + flags (4bits)
	Control byte

	// Length of the paquet in bytes
	// Optionnal
	PacketLength byte

	// Remaining length 1-4 bytes
	// This is the total length without fixed header
	RemainingLength []byte

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
	mp.RemainingLength = mp.RemaingLengthEncode(bodyLength + payloadLength)

	log.Printf("Protocol length Hexadecimal: Ox%X Dec:%d\n", mp.ProtocolLength, mp.ProtocolLength)
	log.Printf("Payload length Hexadecimal: Ox%X Dec:%d\n", mp.Length, mp.Length)
	log.Printf("RemainingLength Dec: %d\n", mp.RemainingLength)

	// compute the fields

	var mqttBuffer bytes.Buffer
	// 1 byte
	mqttBuffer.WriteByte(mp.Control)

	// 1-4 bytes
	mqttBuffer.Write(mp.RemainingLength)

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

func (mp *MqttPacket) RemaingLengthEncode(x int) []byte {

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

func (mp *MqttPacket) RemaingLengthDecode(x []byte) int {

	var multiplier int = 1

	var value int = 0

	var encodedByte byte = 0

	for i := 0; i < len(x); i++ {

		encodedByte = x[i]

		value += int(encodedByte&byte(127)) * multiplier

		multiplier *= 128

		// if (multiplier > 128*128*128)

		//    throw Error(Malformed Remaining Length)

		//  if  (encodedByte & 128) == 0 {
		// 	 break
		//  }

	}

	return value
}

func (mp *MqttPacket) computeLength(buffer []byte) uint16 {
	return uint16(len(buffer))
}

func uint162bytes(val uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, val)
	return buf
}
