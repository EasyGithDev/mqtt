package client

import (
	"log"
	"net"

	"github.com/easygithdev/mqtt/packet"
)

type MqttClient struct {
}

func NewMqttClient() *MqttClient {
	return &MqttClient{}
}

func (mc *MqttClient) Connect(payload string) []byte {

	mp := packet.NewMqttPacket()
	mp.Control = packet.CONNECT
	mp.ProtocolName = "MQTT"
	mp.ProtocolVersion = 4
	mp.ConnectFlag = 0x2
	mp.KeepAlive = 60
	mp.Payload = payload

	log.Printf("Sending command: 0x%x \n", mp.Control)

	return mp.Encode()
}

func (mc *MqttClient) Disconnect() []byte {
	mp := packet.NewMqttPacket()
	mp.Control = packet.DISCONNECT

	log.Printf("Sending command: 0x%x \n", mp.Control)

	return mp.Encode()
}

func (mp *MqttClient) Read(conn net.Conn) (string, error) {

	buffer := make([]byte, 100)
	_, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}

	switch buffer[0] {
	case packet.CONNACK:
		return "CONNACK", nil
	case packet.DISCONNECT:
		return "DISCONNECT", nil
	}

	return "", nil
}
