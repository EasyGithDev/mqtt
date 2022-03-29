package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/easygithdev/mqtt/packet"
	"github.com/easygithdev/mqtt/packet/header"
	"github.com/easygithdev/mqtt/packet/payload"
	"github.com/easygithdev/mqtt/packet/util"
	"github.com/easygithdev/mqtt/packet/variableheader"
)

// Fixed size for the read buffer
const READ_BUFFER_SISE = 1024

type onConnect func(string, string, byte)

// type onDisConnect func(userdata, flags, rc)

// type onMessage func(userdata, message)

type MqttConnectOptions struct {
	Login    string
	Password string
	QoS      byte
}

type MqttClient struct {
	clientId string
	conn     net.Conn
	options  *MqttConnectOptions
}

func NewMqttClient(clientId string) *MqttClient {
	return &MqttClient{clientId: clientId}
}

func (mc *MqttClient) SetConn(conn net.Conn) {
	mc.conn = conn
}

func (mc *MqttClient) SetOptions(options *MqttConnectOptions) {
	mc.options = options
}

func (mc *MqttClient) Connect(host string, port string) (bool, error) {
	conn, err := net.Dial("tcp", host+":"+port)

	if err != nil {
		log.Print("Error connecting:", err.Error())
		return false, err
	}
	mc.conn = conn
	return true, nil
}

func (mc *MqttClient) Disconnect() {
	mc.conn.Close()
}

func (mc *MqttClient) MqttConnect() (bool, error) {

	mh := header.NewMqttHeader()
	mh.Control = header.CONNECT

	mvh := variableheader.NewMqttVariableHeader()
	mvh.BuildConnect("MQTT", 4, variableheader.CONNECT_FLAG_CLEAN_SESSION, 60)

	mpl := payload.NewMqttPayload()
	mpl.AddString(mc.clientId)

	mp := packet.NewMqttPacket()
	mp.Header = mh
	mp.VariableHeader = mvh
	mp.Payload = mpl

	// log.Printf("Payload: %s .................................\n", string(mp.Payload))

	if mc.options != nil {

		//mvh.ConnectFlag = variableheader.CONNECT_FLAG_CLEAN_SESSION | variableheader.CONNECT_FLAG_USERNAME | variableheader.CONNECT_FLAG_PASSWORD

		mp.Header.Control = mp.Header.Control | (0x01 << 7) | (0x01 << 6)
		mp.Payload.AddString(mc.options.Login)
		mp.Payload.AddString(mc.options.Password)
	}

	// log.Printf("Payload: %s .................................\n", string(mp.Payload))

	// log.Printf("Sending command: 0x%x \n", mp.Control)

	buffer := mp.Encode()

	log.Printf("Packet: %s\n", mc.ShowPacket(buffer))

	n, err := mc.conn.Write(buffer)
	if err != nil {
		log.Printf("Sender: Write Error: %s\n", err)
		return false, err
	}

	log.Printf("Wrote %d byte(s)\n", n)

	response, err := mc.Read()
	if err != nil {
		log.Printf("Read Error: %s\n", err)
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

func (mc *MqttClient) Subscribe(topic string) (bool, error) {

	// Adding connection to mc
	_, err := mc.MqttConnect()

	if err != nil {
		return false, nil
	}

	mh := header.NewMqttHeader()
	mh.Control = header.SUBSCRIBE | 1<<1

	var packetId uint16 = 33

	mvh := variableheader.NewMqttVariableHeader()
	mvh.BuildSubscribe(packetId, topic)

	mpl := payload.NewMqttPayload()
	mpl.AddQos(0x00)

	mp := packet.NewMqttPacket()
	mp.Header = mh
	mp.VariableHeader = mvh
	mp.Payload = mpl

	buffer := mp.Encode()

	log.Printf("Packet: %s\n", mc.ShowPacket(buffer))

	n, err := mc.conn.Write(buffer)
	if err != nil {
		log.Printf("Sender: Write Error: %s\n", err)
		return false, err
	}

	log.Printf("Wrote %d byte(s)\n", n)

	response, err := mc.Read()
	if err != nil {
		log.Printf("Read Error: %s\n", err)
		return false, err
	}

	fmt.Println(response)

	log.Printf("Packet: %s\n", mc.ShowPacket(response))

	if response[0] == header.SUBACK {
		return true, nil
	}

	return false, nil
}

// QoS 0: There won’t be any response
// QoS 1: PUBACK – Publish acknowledgement response
// QoS 2 :
// wait for PUBREC – Publish received.
// send back PUBREL – Publish release.
// wait for PUBCOMP – Publish complete.
func (mc *MqttClient) Publish(topic string, message string) (bool, error) {

	// Adding connection to mc
	_, err := mc.MqttConnect()

	if err != nil {
		return false, nil
	}

	// if resp {
	//log.Printf("Connection established: \n")
	// Perform the on connect
	// }

	mh := header.NewMqttHeader()
	mh.Control = header.PUBLISH

	mvh := variableheader.NewMqttVariableHeader()
	mvh.BuildPublish(topic)

	mpl := payload.NewMqttPayload()
	mpl.AddString(message)

	mp := packet.NewMqttPacket()
	mp.Header = mh
	mp.VariableHeader = mvh
	mp.Payload = mpl

	buffer := mp.Encode()

	log.Printf("Packet: %s\n", mc.ShowPacket(buffer))

	n, err := mc.conn.Write(buffer)
	if err != nil {
		log.Printf("Sender: Write Error: %s\n", err)
		return false, err
	}

	log.Printf("Wrote %d byte(s)\n", n)

	// performing section for qos1 & qos 2

	return true, nil
}

/**
The PINGREQ Packet is sent from a Client to the Server. It can be used to:

Indicate to the Server that the Client is alive in the absence of any other Control Packets being sent from the Client to the Server.
Request that the Server responds to confirm that it is alive.
Exercise the network to indicate that the Network Connection is active.

This Packet is used in Keep Alive processing, see Section 3.1.2.10 for more details.

The PINGREQ Packet has no variable header.
The PINGREQ Packet has no payload.

*/
func (mc *MqttClient) Ping() (bool, error) {

	// Adding connection to mc
	_, err := mc.MqttConnect()

	if err != nil {
		return false, nil
	}

	mh := header.NewMqttHeader()
	mh.Control = header.PINGREQ

	mvh := variableheader.NewMqttVariableHeader()

	mpl := payload.NewMqttPayload()

	mp := packet.NewMqttPacket()
	mp.Header = mh
	mp.VariableHeader = mvh
	mp.Payload = mpl

	writeBuffer := mp.Encode()

	log.Printf("Packet: %s\n", mc.ShowPacket(writeBuffer))

	n, err := mc.conn.Write(writeBuffer)
	if err != nil {
		log.Printf("Write Error: %s\n", err)
		return false, err
	}

	log.Printf("Wrote %d byte(s)\n", n)

	// Read PINGRESP

	readBuffer := make([]byte, READ_BUFFER_SISE)
	n, err = mc.conn.Read(readBuffer)
	if err != nil {
		log.Printf("Read Error: %s\n", err)
		return false, err
	}

	bb := bytes.NewBuffer(readBuffer[:n])
	control, _ := bb.ReadByte()

	if control == header.PINGRESP {
		return true, nil
	}

	return false, nil
}

// func (mc *MqttClient) Disconnect() (bool, error) {
// 	mp := packet.NewMqttPacket()
// 	mp.Control = packet.DISCONNECT

// 	log.Printf("Sending command: 0x%x \n", mp.Control)

// 	buffer := mp.Encode()

// 	n, err := mc.conn.Write(buffer)
// 	if err != nil {
// 		log.Printf("Sender: Write Error: %s\n", err)
// 		return false, err
// 	}

// 	log.Printf("Wrote %d byte(s)\n", n)

// 	_, err = mc.Read()
// 	if err != nil {
// 		log.Printf("Sender: Read Error: %s\n", err)
// 		return false, err
// 	}

// 	return true, nil

// }

func (mp *MqttClient) Read() ([]byte, error) {

	buffer := make([]byte, 100)
	n, err := mp.conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	log.Printf("Read: %d byte(s)\n", n)

	return buffer[:n+1], nil
}

func (mp *MqttClient) ReadLoop() {

	for {
		buffer := make([]byte, 1024)
		n, err := mp.conn.Read(buffer)
		if err != nil {
			log.Printf("Error: %s\n", err)
			if err == io.EOF {
				break
			}
		}

		//	res := append([]byte(nil), buffer[:n]...)
		// b1 := bytes.NewBuffer(buffer[:n])
		// fmt.Println(b1)

		log.Printf("Read: %d byte(s)\n", n)

		b1 := bytes.NewBuffer(buffer[:n])
		log.Printf("Len of buffer: %d byte(s)\n", b1.Len())

		header := header.NewMqttHeader()
		header.Control, _ = b1.ReadByte()

		log.Printf("Len of buffer: %d byte(s)\n", b1.Len())

		nb, rLength := header.RemaingLengthDecode(b1.Bytes())
		log.Printf("Read length: %d %d \n", nb, rLength)
		header.RemainingLength = b1.Next(nb)

		log.Printf("Len of buffer: %d byte(s)\n", b1.Len())

		// topicLen, topicMsg := util.StringDecode(b1.Bytes())

		topicLen := util.Bytes2uint16(b1.Next(2))
		topicMsg := string(b1.Next(int(topicLen)))
		log.Printf("Read Topic: %d  [%s] \n", topicLen, topicMsg)

		log.Printf("Len of buffer: %d byte(s)\n", b1.Len())

		msgMsg := string(b1.Bytes())
		log.Printf("Read msg: [%s]\n", msgMsg)

	}

}

func (mp *MqttClient) ShowPacket(buffer []byte) string {
	str := ""
	for i := 0; i < len(buffer); i++ {
		str += fmt.Sprintf("0x%X ", buffer[i])
	}
	str += "\n"

	return str
}
