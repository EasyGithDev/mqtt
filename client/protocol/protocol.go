package protocol

// Default name for connect in variable Header
const PROTOCOL_NAME string = "MQTT"

// Default level for connect in variable Header
const PROTOCOL_LEVEL byte = 4

// Store mqtt name/level in struct
type MqttProtocol struct {
	Name  string
	Level byte
}

func New(name string, level byte) *MqttProtocol {
	return &MqttProtocol{Name: name, Level: level}
}
