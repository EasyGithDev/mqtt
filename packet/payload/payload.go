package payload

import "github.com/easygithdev/mqtt/packet/util"

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
