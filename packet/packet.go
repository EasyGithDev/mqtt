package packet

import (
	"bytes"

	"github.com/easygithdev/mqtt/packet/header"
	"github.com/easygithdev/mqtt/packet/payload"
	"github.com/easygithdev/mqtt/packet/variableheader"
)

type MqttPacket struct {
	Header         *header.MqttHeader
	VariableHeader *variableheader.MqttVariableHeader
	Payload        *payload.MqttPayload
}

func NewMqttPacket() *MqttPacket {
	return &MqttPacket{}
}

func (mp *MqttPacket) Encode() []byte {

	// compute the fields

	var mqttBuffer bytes.Buffer

	// mp.VariableHeader.ComputeProtocolLength()
	mp.Header.ComputeRemainingLength(mp.VariableHeader.Len() + mp.Payload.Len())

	mqttBuffer.Write(mp.Header.Encode())

	if mp.VariableHeader.Len() > 0 {
		mqttBuffer.Write(mp.VariableHeader.Encode())
	}

	if mp.Payload.Len() > 0 {
		mqttBuffer.Write(mp.Payload.Encode())
	}

	return mqttBuffer.Bytes()
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
