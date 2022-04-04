// MIT License

// Copyright (c) 2022 Florent Brusciano

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"time"

	"github.com/easygithdev/mqtt/packet"
	"github.com/easygithdev/mqtt/packet/header"
	"github.com/easygithdev/mqtt/packet/payload"
	"github.com/easygithdev/mqtt/packet/util"
	"github.com/easygithdev/mqtt/packet/vheader"
)

// Fixed connection type
const CONN_TYPE = "tcp"

// Fixed size for the read buffer
const READ_BUFFER_SISE = 1024

// Default name for connect in variable Header
var PROTOCOL_NAME string = "MQTT"

// Default level for connect in variable Header
var PROTOCOL_LEVEL byte = 4

// Default timeout for connect in variable Header
var TIME_OUT uint16 = 60

type MqttCredentials struct {
	login    string
	password string
}

func NewMqttCredentials(login string, password string) *MqttCredentials {
	return &MqttCredentials{login: login, password: password}
}

type MqttClient struct {
	// Tcp connection
	conn *net.Conn

	// Credentials
	credentials *MqttCredentials

	// parameters
	clientId     string
	Timeout      uint16
	CleanSession bool

	// behaviours
	mqttConnected bool

	// handlers
	OnTcpConnect    func(conn net.Conn)
	OnTcpDisconnect func()
	OnConnect       func()
	OnPing          func()
	OnPublish       func(topic string, message string, qos int)
	OnSubscribe     func(topic string, qos int)
	OnMessage       func(messgage string)
}

func NewMqttClient(clientId string) *MqttClient {
	return &MqttClient{clientId: clientId, Timeout: TIME_OUT, CleanSession: true}
}

func (mc *MqttClient) SetConn(conn *net.Conn) {
	mc.conn = conn
}

func (mc *MqttClient) SetCredentials(credentials *MqttCredentials) {
	mc.credentials = credentials
}

func (mc *MqttClient) TcpConnect(host string, port string) (bool, error) {
	conn, err := net.Dial(CONN_TYPE, host+":"+port)

	if err != nil {
		log.Print("Error connecting:", err.Error())
		return false, err
	}
	mc.conn = &conn

	if mc.OnTcpConnect != nil {
		mc.OnTcpConnect(conn)
	}

	return true, nil
}

func (mc *MqttClient) TcpDisconnect() {
	(*mc.conn).Close()

	if mc.OnTcpConnect != nil {
		mc.OnTcpDisconnect()
	}
}

func (mc *MqttClient) MqttConnect() (bool, error) {

	if mc.mqttConnected {
		return true, nil
	}

	var connectFlag byte = 0
	if mc.CleanSession {
		connectFlag |= vheader.CONNECT_FLAG_CLEAN_SESSION
	}

	if mc.credentials != nil {
		connectFlag |= vheader.CONNECT_FLAG_USERNAME | vheader.CONNECT_FLAG_PASSWORD
	}

	mvh := vheader.NewConnectHeader(PROTOCOL_NAME, PROTOCOL_LEVEL, connectFlag, mc.Timeout)

	mpl := payload.NewMqttPayload()
	mpl.AddString(mc.clientId)

	if mc.credentials != nil {
		// mp.Header.Control = mp.Header.Control | (0x01 << 7) | (0x01 << 6)
		mpl.AddString(mc.credentials.login)
		mpl.AddString(mc.credentials.password)
	}

	mh := header.NewMqttHeader(mvh.Len() + mpl.Len())
	mh.Control = header.CONNECT

	mp := packet.NewMqttPacket(mh, mvh, mpl)
	writeBuffer := mp.Encode()

	log.Printf("\n%s\n", mp)
	// log.Printf("Packet: %s\n", util.ShowHexa(writeBuffer))

	n, err := (*mc.conn).Write(writeBuffer)
	if err != nil {
		log.Printf("Write Error: %s\n", err)
		return false, err
	}

	// log.Printf("Wrote %d byte(s)\n", n)

	// Read CONNHACK

	readBuffer := make([]byte, READ_BUFFER_SISE)
	n, readErr := (*mc.conn).Read(readBuffer)
	if readErr != nil {
		log.Printf("Read Error: %s\n", readErr)
		return false, readErr
	}

	bb := bytes.NewBuffer(readBuffer[:n])
	control, _ := bb.ReadByte()

	if control == header.CONNACK {

		bb.Next(2)

		accepted, _ := bb.ReadByte()

		switch accepted {
		case header.CONNECT_ACCEPTED:
			if mc.OnConnect != nil {
				mc.OnConnect()
			}
			mc.mqttConnected = true
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

func (mc *MqttClient) MqttDisconnect() (bool, error) {

	if !mc.mqttConnected {
		return true, nil
	}

	mh := header.NewMqttHeader(0)
	mh.Control = header.DISCONNECT

	mp := packet.NewMqttPacket(mh, nil, nil)

	// log.Printf("Sending command: 0x%x \n", mh.Control)

	writeBuffer := mp.Encode()

	n, err := (*mc.conn).Write(writeBuffer)
	if err != nil {
		log.Printf("Write Error: %s\n", err)
		return false, err
	}

	log.Printf("Wrote %d byte(s)\n", n)

	mc.TcpDisconnect()
	mc.mqttConnected = false

	return true, nil

}

// The SUBSCRIBE Packet is sent from the Client to the Server to create one or more Subscriptions.
// Each Subscription registers a Client’s interest in one or more Topics. The Server sends PUBLISH Packets to the Client in order to forward Application Messages that were published to Topics that match these Subscriptions. The SUBSCRIBE Packet also specifies (for each Subscription) the maximum QoS with which the Server can send Application Messages to the Client.
func (mc *MqttClient) Subscribe(topic string) (bool, error) {

	// Adding connection to mc
	_, err := mc.MqttConnect()

	if err != nil {
		return false, nil
	}

	//The variable header component of many of the Control Packet types includes a 2 byte Packet Identifier field.
	//These Control Packets are PUBLISH (where QoS > 0), PUBACK, PUBREC, PUBREL, PUBCOMP, SUBSCRIBE, SUBACK, UNSUBSCRIBE, UNSUBACK.
	var packetId uint16 = uint16(rand.Intn(math.MaxInt16))

	mvh := vheader.NewSubscribeHeader(packetId, topic)
	// mvh.BuildSubscribe(packetId, topic)

	mpl := payload.NewMqttPayload()
	mpl.AddQos(0x00)

	mh := header.NewMqttHeader(mvh.Len() + mpl.Len())
	// Bits 3,2,1 and 0 of the fixed header of the SUBSCRIBE Control Packet are reserved and MUST be set to 0,0,1 and 0 respectively.
	mh.Control = header.SUBSCRIBE | 1<<1

	mp := packet.NewMqttPacket(mh, mvh, mpl)
	writeBuffer := mp.Encode()

	log.Printf("\n%s\n", mp)

	// log.Printf("Packet: %s\n", util.ShowHexa(buffer))

	n, err := (*mc.conn).Write(writeBuffer)
	if err != nil {
		log.Printf("Write Error: %s\n", err)
		return false, err
	}

	// log.Printf("Wrote %d byte(s)\n", n)

	// Read SUBACK

	readBuffer := make([]byte, READ_BUFFER_SISE)
	n, readErr := (*mc.conn).Read(readBuffer)
	if readErr != nil {
		log.Printf("Read Error: %s\n", readErr)
		return false, readErr
	}

	bb := bytes.NewBuffer(readBuffer[:n])
	control, _ := bb.ReadByte()

	if control == header.SUBACK {
		if mc.OnSubscribe != nil {
			mc.OnSubscribe(topic, 0)
		}
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
func (mc *MqttClient) PublishQos0(topic string, message string) (bool, error) {
	return mc.Publish(topic, message, 0)
}

func (mc *MqttClient) PublishQos1(topic string, message string) (bool, error) {
	return mc.Publish(topic, message, 1)
}

func (mc *MqttClient) PublishQos2(topic string, message string) (bool, error) {
	return mc.Publish(topic, message, 2)
}

func (mc *MqttClient) Publish(topic string, message string, qos int) (bool, error) {

	// Adding connection to mc
	_, err := mc.MqttConnect()

	if err != nil {
		return false, nil
	}

	mvh := vheader.NewPublishHeader(topic)

	mpl := payload.NewMqttPayload()
	mpl.AddString(message)

	mh := header.NewMqttHeader(mvh.Len() + mpl.Len())
	mh.Control = header.PUBLISH

	if qos == 1 {
		mh.Control |= 1 << 1

	}

	if qos == 2 {
		mh.Control |= 1 << 2
	}

	// retain
	//mh.Control |= 1 << 3

	mp := packet.NewMqttPacket(mh, mvh, mpl)
	writeBuffer := mp.Encode()

	log.Printf("\n%s\n", mp)

	n, err := (*mc.conn).Write(writeBuffer)
	if err != nil {
		log.Printf("Write Error: %s\n", err)
		return false, err
	}

	// log.Printf("Publish wrote %d byte(s)\n", n)

	// Nothing to read for Qos 0
	if qos == 0 {
		if mc.OnPublish != nil {
			mc.OnPublish(topic, message, qos)
		}
	} else if qos == 1 {

		// Read PUBACK

		readBuffer := make([]byte, READ_BUFFER_SISE)
		n, err = (*mc.conn).Read(readBuffer)
		if err != nil {
			log.Printf("Read Error: %s\n", err)
			return false, err
		}

		bb := bytes.NewBuffer(readBuffer[:n])
		control, _ := bb.ReadByte()

		if control == header.PUBACK {
			if mc.OnPublish != nil {
				mc.OnPublish(topic, message, qos)
			}
			return true, nil
		}

	} else if qos == 2 {
		fmt.Println("2")
	}

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

	mh := header.NewMqttHeader(0)
	mh.Control = header.PINGREQ

	mp := packet.NewMqttPacket(mh, nil, nil)
	writeBuffer := mp.Encode()

	log.Printf("\n%s\n", mp)

	// log.Printf("Packet: %s\n", util.ShowHexa(writeBuffer))

	n, err := (*mc.conn).Write(writeBuffer)
	if err != nil {
		log.Printf("Write Error: %s\n", err)
		return false, err
	}

	log.Printf("Wrote %d byte(s)\n", n)

	// Read PINGRESP

	readBuffer := make([]byte, READ_BUFFER_SISE)
	n, err = (*mc.conn).Read(readBuffer)
	if err != nil {
		log.Printf("Read Error: %s\n", err)
		return false, err
	}

	bb := bytes.NewBuffer(readBuffer[:n])
	control, _ := bb.ReadByte()

	if control == header.PINGRESP {
		if mc.OnPing != nil {
			mc.OnPing()
		}
		return true, nil
	}

	return false, nil
}

func (mc *MqttClient) LoopStart() {
	start := time.Now()
	for {
		elapsed := time.Since(start).Seconds()
		if int(math.Ceil(elapsed)) >= int(TIME_OUT/3) {
			mc.Ping()
			start = time.Now()
		}
	}
}

func (mc *MqttClient) LoopForever() {

	for {

		// Adding connection to mc
		_, err := mc.MqttConnect()

		if err != nil {
			log.Printf("LoopForever Error: %s\n", err)
		}

		buffer := make([]byte, READ_BUFFER_SISE)
		n, err := (*mc.conn).Read(buffer)
		if err != nil {
			log.Printf("Error: %s\n", err)
			if err == io.EOF {
				mc.mqttConnected = false
				continue
			}
		}

		// log.Printf("Read: %d byte(s)\n", n)

		b1 := bytes.NewBuffer(buffer[:n])
		// log.Printf("Len of buffer: %d byte(s)\n", b1.Len())

		control, _ := b1.ReadByte()

		log.Printf("Header control: %b\n", control)

		// log.Printf("Len of buffer: %d byte(s)\n", b1.Len())

		nb, _ := header.RemaingLengthDecode(b1.Bytes())
		// log.Printf("Read length: %d %d \n", nb, rLength)
		remainingLength := b1.Next(nb)

		log.Printf("Header remainingLength: %b\n", remainingLength)

		// log.Printf("Len of buffer: %d byte(s)\n", b1.Len())

		// topicLen, topicMsg := util.StringDecode(b1.Bytes())

		topicLen := util.Bytes2uint16(b1.Next(2))
		b1.Next(int(topicLen))

		// topicMsg := string(b1.Next(int(topicLen)))
		// log.Printf("Read Topic: %d  [%s] \n", topicLen, topicMsg)

		// log.Printf("Len of buffer: %d byte(s)\n", b1.Len())

		msgMsg := string(b1.Bytes())
		// log.Printf("Read msg: [%s]\n", msgMsg)

		if mc.OnMessage != nil {
			mc.OnMessage(msgMsg)
		}

	}

}
