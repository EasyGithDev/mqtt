package client

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/easygithdev/mqtt/packet"
	"github.com/easygithdev/mqtt/packet/header"
	"github.com/easygithdev/mqtt/packet/payload"
	"github.com/easygithdev/mqtt/packet/variableheader"
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

	mh := header.NewMqttHeader()
	mh.Control = header.CONNECT

	mvh := variableheader.NewMqttVariableHeader()
	mvh.ProtocolName = "MQTT"
	mvh.ProtocolVersion = 4
	mvh.ConnectFlag = variableheader.CONNECT_FLAG_CLEAN_SESSION
	mvh.KeepAlive = 60

	mpl := payload.NewMqttPayload()
	mpl.AddString(idClient)

	mp := packet.NewMqttPacket()
	mp.Header = mh
	mp.VariableHeader = mvh
	mp.Payload = mpl

	// log.Printf("Payload: %s .................................\n", string(mp.Payload))

	if options != nil {
		mp.Header.Control = mp.Header.Control | (0x01 << 7) | (0x01 << 6)
		mp.Payload.AddString(options.Login)
		mp.Payload.AddString(options.Password)
	}

	// log.Printf("Payload: %s .................................\n", string(mp.Payload))

	// log.Printf("Sending command: 0x%x \n", mp.Control)

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

	// myPacket := mp.GetPacket(response)

	//	log.Printf("Packet: %s\n", mc.ShowPacket(myPacket))

	if response[0] == header.CONNACK {

		switch response[3] {
		case header.CONNECT_ACCEPTED:
			return true, nil
		case header.CONNECT_REFUSED_1:
			return false, errors.New("connection Refused, unacceptable protocol version")
		case header.CONNECT_REFUSED_2:
			return false, errors.New("connection Refused, identifier rejected")
		case header.CONNECT_REFUSED_3:
			return false, errors.New("connection Refused, Server unavailable")
		case header.CONNECT_REFUSED_4:
			return false, errors.New("connection Refused, bad user name or password")
		case header.CONNECT_REFUSED_5:
			return false, errors.New("connection Refused, not authorized")
		default:
			return false, nil
		}
	}

	return false, nil
}

// func (mc *MqttClient) Disconnect() (bool, error) {
// 	mp := packet.NewMqttPacket()
// 	mp.Control = packet.DISCONNECT

// 	log.Printf("Sending command: 0x%x \n", mp.Control)

// 	buffer := mp.Encode()

// 	n, err := (*mc.conn).Write(buffer)
// 	if err != nil {
// 		log.Printf("Sender: Write Error: %s\n", err)
// 		return false, err
// 	}

// 	log.Printf("Sender: Wrote %d byte(s)\n", n)

// 	_, err = mc.Read()
// 	if err != nil {
// 		log.Printf("Sender: Read Error: %s\n", err)
// 		return false, err
// 	}

// 	return true, nil

// }

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
