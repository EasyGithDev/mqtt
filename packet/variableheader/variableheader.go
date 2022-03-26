package variableheader

import (
	"bytes"

	"github.com/easygithdev/mqtt/packet/util"
)

type MqttVariableHeader struct {
	// Length of protocol name (expl MQTT4 -> 4)
	ProtocolLength uint16

	// Protocol + version (expl MQTT4)
	ProtocolName    string
	ProtocolVersion byte

	// Connect flag (expl clean session)
	ConnectFlag byte

	// Keep alive (2 bytes)
	KeepAlive uint16
}

func NewMqttVariableHeader() *MqttVariableHeader {
	return &MqttVariableHeader{}
}

// 10 bytes = 2 (protocol length) + 4 (protocol name) + 1 (protocol version) + 1 (connect flag) + 2 (keep alive)
func (mvh *MqttVariableHeader) Encode() []byte {
	var buffer bytes.Buffer

	buf := util.Uint162bytes(mvh.ProtocolLength)
	buffer.Write(buf)

	buffer.WriteString(mvh.ProtocolName)

	buffer.WriteByte(mvh.ProtocolVersion)

	buffer.WriteByte(mvh.ConnectFlag)

	buf = util.Uint162bytes(mvh.KeepAlive)
	buffer.Write(buf)

	return buffer.Bytes()
}

func (mvh *MqttVariableHeader) ComputeProtocolLength() {
	mvh.ProtocolLength = uint16(len(mvh.ProtocolName))
}

func (mvh *MqttVariableHeader) Len() int {
	if mvh.ProtocolLength != 0 {
		return 10
	}
	return 0
}
