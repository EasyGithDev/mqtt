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
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"time"

	"github.com/easygithdev/mqtt/client/conn"
	"github.com/easygithdev/mqtt/client/credentials"
	"github.com/easygithdev/mqtt/client/protocol"
	"github.com/easygithdev/mqtt/packet"
	"github.com/easygithdev/mqtt/packet/header"
	"github.com/easygithdev/mqtt/packet/payload"
	"github.com/easygithdev/mqtt/packet/util"
	"github.com/easygithdev/mqtt/packet/vheader"
)

// Fixed size for the read buffer
const READ_BUFFER_SISE = 1024

// Default clean session for connect in variable Header
var CLEAN_SESSION bool = true

// Define the Mqtt client
type MqttClient struct {
	// Connection
	conn      *net.Conn
	connInfos *conn.MqttConn

	// Credentials
	credentials *credentials.MqttCredentials

	// parameters
	clientId     string
	cleanSession bool
	userData     interface{}
	protocol     *protocol.MqttProtocol

	// behaviours
	mqttConnected   bool
	topicSubscribed string

	// callbacks
	OnConnect     func(mc MqttClient, userData interface{}, rc net.Conn)
	OnDisconnect  func(mc MqttClient, userData interface{}, rc net.Conn)
	OnPublish     func(mc MqttClient, userData interface{}, mid int)
	OnSubscribe   func(mc MqttClient, userData interface{}, mid int)
	OnUnsubscribe func(mc MqttClient, userData interface{}, mid int)
	OnMessage     func(mc MqttClient, userData interface{}, message string)
}

// client_id=””, clean_session=True, userdata=None, protocol=MQTTv311)
func New(clientId string, opts ...Option) *MqttClient {
	mc := &MqttClient{
		conn:         nil,
		clientId:     clientId,
		cleanSession: CLEAN_SESSION,
		userData:     nil,
		protocol:     protocol.New(protocol.PROTOCOL_NAME, protocol.PROTOCOL_LEVEL),
	}

	for _, applyOpt := range opts {
		if applyOpt != nil {
			applyOpt(mc)
		}
	}

	return mc
}

type Option func(f *MqttClient)

func WithCleanSession(cleanSession bool) Option {
	return func(mc *MqttClient) {
		mc.cleanSession = cleanSession
	}
}

func WithUserData(userData interface{}) Option {
	return func(mc *MqttClient) {
		mc.userData = userData
	}
}

func WithProtocol(name string, level byte) Option {
	return func(mc *MqttClient) {
		mc.protocol = protocol.New(name, level)
	}
}

func WithCredentials(login string, password string) Option {
	return func(mc *MqttClient) {
		mc.credentials = credentials.New(login, password)
	}
}

func WithConnInfos(connInfos *conn.MqttConn) Option {
	return func(mc *MqttClient) {
		mc.connInfos = connInfos
	}
}

func (mc *MqttClient) Connect() (bool, error) {

	conn, err := net.Dial(mc.connInfos.Transport, mc.connInfos.Host+":"+mc.connInfos.Port)

	if err != nil {
		log.Println("Error connecting:", err.Error())
		return false, err
	}
	mc.conn = &conn

	return true, nil
}

func (mc *MqttClient) Close() {
	(*mc.conn).Close()
}

func (mc *MqttClient) Read() (*bytes.Buffer, error) {
	readBuffer := make([]byte, READ_BUFFER_SISE)
	n, readErr := (*mc.conn).Read(readBuffer)
	if readErr != nil {
		return nil, readErr
	}
	return bytes.NewBuffer(readBuffer[:n]), nil
}

func (mc *MqttClient) Write(buffer []byte) (int, error) {
	return (*mc.conn).Write(buffer)
}

// connect(host, port=1883, keepalive=60, bind_address="")
func (mc *MqttClient) MqttConnect() (bool, error) {

	if mc.conn == nil {
		return false, nil
	}

	if mc.mqttConnected {
		return true, nil
	}

	var connectFlag byte = 0
	if mc.cleanSession {
		connectFlag |= vheader.CONNECT_FLAG_CLEAN_SESSION
	}

	if mc.credentials != nil {
		connectFlag |= vheader.CONNECT_FLAG_USERNAME | vheader.CONNECT_FLAG_PASSWORD
	}

	mvh := vheader.NewConnectHeader(mc.protocol.Name, mc.protocol.Level, connectFlag, mc.connInfos.KeepAlive)

	mpl := payload.NewMqttPayload()
	mpl.AddString(mc.clientId)

	if mc.credentials != nil {
		// mp.Header.Control = mp.Header.Control | (0x01 << 7) | (0x01 << 6)
		mpl.AddString(mc.credentials.Login)
		mpl.AddString(mc.credentials.Password)
	}

	mh := header.NewMqttHeader(mvh.Len() + mpl.Len())
	mh.UseConnect()

	mp := packet.NewMqttPacket(mh, mvh, mpl)
	writeBuffer := mp.Encode()

	log.Printf("\n%s\n\n", mp)
	// log.Printf("Packet: %s\n", util.ShowHexa(writeBuffer))

	_, err := (*mc.conn).Write(writeBuffer)
	if err != nil {
		log.Printf("Write Error: %s\n", err)
		return false, err
	}

	// log.Printf("Wrote %d byte(s)\n", n)

	// Read CONNHACK

	bb, readErr := mc.Read()
	if err != nil {
		log.Printf("Read Error: %s\n", readErr)
		return false, readErr
	}

	control, _ := bb.ReadByte()

	if control == header.CONNACK {

		bb.Next(2)

		accepted, _ := bb.ReadByte()

		switch accepted {
		case header.CONNECT_ACCEPTED:
			if mc.OnConnect != nil {
				mc.OnConnect(*mc, nil, *mc.conn)
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
	mh.UseDisconnect()

	mp := packet.NewMqttPacket(mh, nil, nil)

	// log.Printf("Sending command: 0x%x \n", mh.Control)

	writeBuffer := mp.Encode()

	n, err := (*mc.conn).Write(writeBuffer)
	if err != nil {
		log.Printf("Write Error: %s\n", err)
		return false, err
	}

	log.Printf("Wrote %d byte(s)\n", n)

	// mc.TcpDisconnect()
	mc.mqttConnected = false

	return true, nil

}

// The SUBSCRIBE Packet is sent from the Client to the Server to create one or more Subscriptions.
// Each Subscription registers a Client’s interest in one or more Topics. The Server sends PUBLISH Packets to the Client in order to forward Application Messages that were published to Topics that match these Subscriptions. The SUBSCRIBE Packet also specifies (for each Subscription) the maximum QoS with which the Server can send Application Messages to the Client.
func (mc *MqttClient) Subscribe(topic string) (bool, error) {

	mc.topicSubscribed = topic

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
	mh.UseSubscribe()

	mp := packet.NewMqttPacket(mh, mvh, mpl)
	writeBuffer := mp.Encode()

	log.Printf("\n%s\n\n", mp)

	// log.Printf("Packet: %s\n", util.ShowHexa(buffer))

	n, err := (*mc.conn).Write(writeBuffer)
	if err != nil {
		log.Printf("Write Error: %s\n", err)
		return false, err
	}

	log.Printf("Wrote %d byte(s)\n", n)

	// Read SUBACK

	bb, readErr := mc.Read()
	if err != nil {
		log.Printf("Read Error: %s\n", readErr)
		return false, readErr
	}
	control, _ := bb.ReadByte()

	if control == header.SUBACK {
		if mc.OnSubscribe != nil {
			mc.OnSubscribe(*mc, nil, 0)
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
	mh.UsePublish()

	if qos == 1 {
		mh.UseQos1()
	}

	if qos == 2 {
		mh.UseQos2()
	}

	// retain
	//mh.Control |= 1 << 3

	mp := packet.NewMqttPacket(mh, mvh, mpl)
	writeBuffer := mp.Encode()

	log.Printf("\n%s\n\n", mp)

	n, err := (*mc.conn).Write(writeBuffer)
	if err != nil {
		log.Printf("Write Error: %s\n", err)
		return false, err
	}

	log.Printf("Publish wrote %d byte(s)\n", n)

	// Nothing to read for Qos 0
	if qos == 0 {
		if mc.OnPublish != nil {
			mc.OnPublish(*mc, nil, 0)
		}
	} else if qos == 1 {

		// Read PUBACK

		bb, err := mc.Read()
		if err != nil {
			log.Printf("Read Error: %s\n", err)
			return false, err
		}

		control, _ := bb.ReadByte()

		if control == header.PUBACK {
			if mc.OnPublish != nil {
				mc.OnPublish(*mc, nil, 0)
			}
			return true, nil
		}

	} else if qos == 2 {
		// Read PUBREC

		bb, err := mc.Read()
		if err != nil {
			log.Printf("Read Error: %s\n", err)
			return false, err
		}

		control, _ := bb.ReadByte()
		bb.Read(make([]byte, 1))
		vh := make([]byte, 2)
		bb.Read(vh)

		if control == header.PUBREC {

			// Send a PUBREL
			// The variable header contains the same Packet Identifier as the PUBREC Packet that is being acknowledged
			mvh := vheader.NewGenericHeader(vh)
			mh := header.NewMqttHeader(mvh.Len())
			mh.UsePubrel()
			mp := packet.NewMqttPacket(mh, mvh, nil)
			writeBuffer := mp.Encode()
			_, err := (*mc.conn).Write(writeBuffer)
			if err != nil {
				log.Printf("Write Error: %s\n", err)
				return false, err
			}

			// Read PUBCOMB
			bb, err := mc.Read()
			if err != nil {
				log.Printf("Read Error: %s\n", err)
				return false, err
			}

			control, _ := bb.ReadByte()

			if control == header.PUBCOMP {
				if mc.OnPublish != nil {
					mc.OnPublish(*mc, nil, 0)
				}
			}

		} else {
			return false, nil
		}

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
	mh.UsePingreq()

	mp := packet.NewMqttPacket(mh, nil, nil)
	writeBuffer := mp.Encode()

	log.Printf("\n%s\n\n", mp)

	// log.Printf("Packet: %s\n", util.ShowHexa(writeBuffer))

	n, err := (*mc.conn).Write(writeBuffer)
	if err != nil {
		log.Printf("Write Error: %s\n", err)
		return false, err
	}

	log.Printf("Wrote %d byte(s)\n", n)

	// Read PINGRESP

	bb, err := mc.Read()
	if err != nil {
		log.Printf("Read Error: %s\n", err)
		return false, err
	}

	control, _ := bb.ReadByte()

	if control == header.PINGRESP {
		return true, nil
	}

	return false, nil
}

func (mc *MqttClient) LoopStart() {
	start := time.Now()
	for {
		elapsed := time.Since(start).Seconds()
		if int(math.Ceil(elapsed)) >= int(mc.connInfos.KeepAlive/3) {
			mc.Ping()
			start = time.Now()
		}
	}
}

func (mc *MqttClient) LoopForever() {

	for {
		if !mc.mqttConnected {
			_, connErr := mc.Connect()
			if connErr != nil {
				log.Println("Trying reset the connection...:", connErr.Error())
				time.Sleep(time.Second * 1)
				continue
			}

			// Try subscribe
			mc.Subscribe(mc.topicSubscribed)
		}

		for {

			b1, err := mc.Read()
			if err != nil {
				log.Printf("Error: %s\n", err)
				if err == io.EOF {
					// Do something smart
				}
				mc.mqttConnected = false
				mc.Close()
				break
			}

			// log.Printf("Read: %d byte(s)\n", n)

			// b1 := bytes.NewBuffer(buffer[:n])
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

			msg := string(b1.Bytes())
			// log.Printf("Read msg: [%s]\n", msgMsg)

			if mc.OnMessage != nil {
				mc.OnMessage(*mc, nil, msg)
			}

		}
	}

}
