package client

import (
	"log"
	"net"

	"github.com/easygithdev/mqtt/packet"
)

type MqttClient struct {
	conn *net.Conn
}

func NewMqttClient() *MqttClient {
	return &MqttClient{}
}

func (mc *MqttClient) SetConn(conn *net.Conn) {
	mc.conn = conn
}

func (mc *MqttClient) Connect(idClient string) (bool, error) {

	mp := packet.NewMqttPacket()
	mp.Control = packet.CONNECT
	mp.ProtocolName = "MQTT"
	mp.ProtocolVersion = 4
	mp.ConnectFlag = 0x2
	mp.KeepAlive = 60
	mp.Payload = idClient

	log.Printf("Sending command: 0x%x \n", mp.Control)

	buffer := mp.Encode()

	n, err := (*mc.conn).Write(buffer)
	if err != nil {
		log.Printf("Sender: Write Error: %s\n", err)
		return false, err
	}

	log.Printf("Sender: Wrote %d byte(s)\n", n)

	_, err = mc.Read()
	if err != nil {
		log.Printf("Sender: Read Error: %s\n", err)
		return false, err
	}

	return true, nil
}

func (mc *MqttClient) Disconnect() []byte {
	mp := packet.NewMqttPacket()
	mp.Control = packet.DISCONNECT

	log.Printf("Sending command: 0x%x \n", mp.Control)

	return mp.Encode()
}

func (mp *MqttClient) Read() (string, error) {

	buffer := make([]byte, 100)
	_, err := (*mp.conn).Read(buffer)
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
