package client

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/easygithdev/mqtt/packet"
)

type MqttConnectOptions struct {
	Login    string
	Password string
}

type MqttClient struct {
	conn *net.Conn
}

func NewMqttClient() *MqttClient {
	return &MqttClient{}
}

func (mc *MqttClient) SetConn(conn *net.Conn) {
	mc.conn = conn
}

func (mc *MqttClient) Connect(idClient string, options *MqttConnectOptions) (bool, error) {

	mp := packet.NewMqttPacket()
	mp.Control = packet.CONNECT
	mp.ProtocolName = "MQTT"
	mp.ProtocolVersion = 4
	mp.ConnectFlag = 0x2
	mp.KeepAlive = 60
	mp.Payload = idClient

	if options != nil {
		mp.Control = mp.Control & (0x01 << 7) & (0x01 << 6)
		mp.Payload += options.Login + options.Password
	}

	log.Printf("Sending command: 0x%x \n", mp.Control)

	buffer := mp.Encode()

	log.Printf("Packet: %s\n", mc.ShowPacket(buffer))

	n, err := (*mc.conn).Write(buffer)
	if err != nil {
		log.Printf("Sender: Write Error: %s\n", err)
		return false, err
	}

	log.Printf("Sender: Wrote %d byte(s)\n", n)

	response, err := mc.Read()
	if err != nil {
		log.Printf("Sender: Read Error: %s\n", err)
		return false, err
	}

	myPacket := mp.GetPacket(response)

	log.Printf("Packet: %s\n", mc.ShowPacket(myPacket))

	if response[0] == packet.CONNACK {

		switch response[3] {
		case packet.CONNECT_REFUSED_1:
			return false, errors.New("connection Refused, unacceptable protocol version")
		case packet.CONNECT_REFUSED_2:
			return false, errors.New("connection Refused, identifier rejected")
		case packet.CONNECT_REFUSED_3:
			return false, errors.New("connection Refused, Server unavailable")
		case packet.CONNECT_REFUSED_4:
			return false, errors.New("connection Refused, bad user name or password")
		case packet.CONNECT_REFUSED_5:
			return false, errors.New("connection Refused, not authorized")
		default:
			return false, nil
		}
	}

	return true, nil
}

func (mc *MqttClient) Disconnect() (bool, error) {
	mp := packet.NewMqttPacket()
	mp.Control = packet.DISCONNECT

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

func (mp *MqttClient) Read() ([]byte, error) {

	buffer := make([]byte, 100)
	_, err := (*mp.conn).Read(buffer)
	if err != nil {
		return buffer, err
	}

	return buffer, nil
}

func (mp *MqttClient) ShowPacket(buffer []byte) string {
	str := ""
	for i := 0; i < len(buffer); i++ {
		str += fmt.Sprintf("0x%X ", buffer[i])
	}
	str += "\n"

	return str
}
