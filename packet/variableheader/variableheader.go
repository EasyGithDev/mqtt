package variableheader

import (
	"bytes"

	"github.com/easygithdev/mqtt/packet/util"
)

var CONNECT_FLAG_CLEAN_SESSION byte = 0x02
var CONNECT_FLAG_WILL_FLAG byte = 0x04
var PUBLCONNECT_FLAG_WILL_QOS_1 byte = 0x08
var PUBLCONNECT_FLAG_WILL_QOS_2 byte = 0x10
var CONNECT_FLAG_WILL_RETAIN byte = 0x20
var CONNECT_FLAG_PASSWORD byte = 0x40
var CONNECT_FLAG_USERNAME byte = 0x80

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
