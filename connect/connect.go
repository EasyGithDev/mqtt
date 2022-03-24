package connect

import (
	"net"
	"strconv"
)

// Control message

var CONNECT byte = 0x10
var CONNACK byte = 0x20
var PUBLISH byte = 3

// control + length + protocol name + Protocol Level +Connect Flags + keep alive +Payload

type MQTTPaquet struct {

	// Fixed header
	// packet type (4bits) + flags (4bits)
	Control byte //`json:"Control"`

	// Length of the paquet in bytes
	// OPtionnal
	PacketLength byte //`json:"PacketLength"`

	// Remaining length
	// This is the total length
	RemainingLength byte //`json:"RemainingLength"`

	// Length of protocol name (expl MQTT4 -> 4)
	ProtocolLength uint16 //`json:"ProtocolLength"`

	// Protocol + version (expl MQTT4)
	ProtocolName    string //`json:"ProtocolName"`
	ProtocolVersion byte   //`json:"ProtocolVersion"`

	// Connect flag (expl clean session)
	ConnectFlag byte //`json:"ConnectFlag"`

	// Keep alive (2 bytes)
	KeepAlive uint16 //`json:"KeepAlive"`

	// Length	(2 bytes)
	// length of payload
	Length uint16 //`json:"Length"`

	// Payload
	Payload string //`json:"Payload"`

	// clientId        string
	// cleanSession    bool
	// username        string
	// password        string
	// lastWilltopic   string
	// lastWillQos     string
	// lastWillMessage string
	// lastWillRetain  bool
	// keepAlive       int
}

func Read(conn net.Conn) (string, error) {

	buffer := make([]byte, 100)
	_, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}

	switch buffer[0] {
	case CONNACK:
		return "CONNACK", nil

	}

	return "", nil
}

// func NewMQTTPaquet(clientId string) *MQTTPaquet {
// 	return &MQTTPaquet{clientId: clientId, cleanSession: true, lastWillMessage: "unexpected exit", lastWillRetain: false, keepAlive: 60}
// }

// func (cm *MQTTPaquet) String() string {
// 	return cm.clientId + "\n" +
// 		strconv.FormatBool(cm.cleanSession) + "\n" +
// 		cm.username + "\n" +
// 		cm.password + "\n" +
// 		cm.lastWilltopic + "\n" +
// 		cm.lastWillQos + "\n" +
// 		cm.lastWillMessage + "\n" +
// 		strconv.FormatBool(cm.lastWillRetain) + "\n" +
// 		strconv.Itoa(cm.keepAlive)
// }

type Connack struct {
	sessionPresent bool
	returnCode     int
}

func NewConnack() *Connack {
	return &Connack{}
}

func (ck *Connack) String() string {
	return strconv.FormatBool(ck.sessionPresent) + "\n" + strconv.Itoa(ck.returnCode)
}
