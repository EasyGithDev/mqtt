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
