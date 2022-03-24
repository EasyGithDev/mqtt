package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/easygithdev/mqtt/client"
	"github.com/easygithdev/mqtt/packet"
)

const (
	// connHost = "localhost"
	// connPort = "8080"
	// connType = "tcp"

	connHost = "test.mosquitto.org"
	connPort = "1883"
	connType = "tcp"

	DELIMITER byte = '\n'
	QUIT_SIGN      = "quit!"
)

// func Read(conn net.Conn, delim byte) (string, error) {
// 	reader := bufio.NewReader(conn)
// 	var buffer bytes.Buffer
// 	for {
// 		ba, isPrefix, err := reader.ReadLine()
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			return "", err
// 		}
// 		buffer.Write(ba)
// 		if !isPrefix {
// 			break
// 		}
// 	}
// 	return buffer.String(), nil
// }

func Write(conn net.Conn, content string) (int, error) {
	writer := bufio.NewWriter(conn)
	number, err := writer.WriteString(content)
	if err == nil {
		err = writer.Flush()
	}
	return number, err
}

func handleConnection(conn net.Conn) {
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		fmt.Println("Client left.")
		conn.Close()
		return
	}

	log.Println("Client message:", string(buffer[:len(buffer)-1]))

	conn.Write(buffer)

	handleConnection(conn)
}

func startServer() {
	fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {

		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client connected.")
		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		handleConnection(c)
	}
}

func myRead(client *client.MqttClient, conn net.Conn) {

	message, err := client.Read(conn)

	if err != nil {
		log.Printf("Sender: Read error: %s", err)
	}

	log.Print("Server relay:", message)

}

func main() {

	res := packet.NewMqttPacket().RemaingLengthEncode(321)

	println(res)

	os.Exit(0)

	//go startServer()

	var dialer net.Dialer
	dialer.Timeout = time.Minute * 10
	dialer.KeepAlive = time.Minute * 10

	fmt.Println("Connecting to " + connType + " server " + connHost + ":" + connPort)
	// conn, err := net.Dial(connType, connHost+":"+connPort)

	conn, err := dialer.Dial(connType, connHost+":"+connPort)

	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}

	// mp := connect.MQTTPaquet{}
	// mp.Control = connect.CONNECT
	// // mp.PacketLength = 23
	// mp.RemainingLength = 23

	// mp.ProtocolLength = 4
	// mp.ProtocolName = "MQTT"
	// mp.ProtocolVersion = 4
	// mp.ConnectFlag = 0x2
	// mp.KeepAlive = 60
	// mp.Length = 7
	// mp.Payload = "python_test"

	// // start to calcule the strings length

	// var bufProtocolLen bytes.Buffer
	// bufProtocolLen.WriteString(mp.ProtocolName)
	// println("len ProtocoleLen:", len(bufProtocolLen.Bytes()))
	// mp.ProtocolLength = uint16(len(bufProtocolLen.Bytes()))

	// var bufPayload bytes.Buffer
	// bufPayload.WriteString(mp.Payload)
	// println("len Payload:", len(bufPayload.Bytes()))
	// mp.Length = uint16(len(bufPayload.Bytes()))

	// fmt.Printf("Hexadecimal: Ox%X Dec:%d\n", mp.Control, mp.Control)
	// fmt.Printf("Hexadecimal: Ox%X Dec:%d\n", mp.Length, mp.Length)

	// var network bytes.Buffer // Stand-in for a network connection
	// network.WriteByte(mp.Control)

	// network.WriteByte(mp.RemainingLength)

	// buf := make([]byte, 2)
	// binary.BigEndian.PutUint16(buf, mp.ProtocolLength)
	// network.Write(buf)

	// network.WriteString(mp.ProtocolName)

	// network.WriteByte(mp.ProtocolVersion)

	// network.WriteByte(mp.ConnectFlag)

	// buf = make([]byte, 2)
	// binary.BigEndian.PutUint16(buf, mp.KeepAlive)
	// network.Write(buf)

	// buf = make([]byte, 2)
	// binary.BigEndian.PutUint16(buf, mp.Length)
	// network.Write(buf)

	// network.WriteString(mp.Payload)

	// fmt.Println("len", len(network.Bytes()))
	// fmt.Println("cap", cap(network.Bytes()))

	// fmt.Println(network.Bytes())

	// for _, v := range network.Bytes() {
	// 	fmt.Printf("0x%x ", v)
	// }

	// s := fmt.Sprintf("%v", mp)
	// b := []byte(s)
	// for _, v := range b {
	// 	fmt.Printf("0x%x ", v)
	// }

	//fmt.Println(str)

	// this is the []byte

	//	os.Exit(0)

	// num, err := conn.Write(network.Bytes())

	// Connect
	client := client.NewMqttClient()

	// Read continue

	clientBuffer := client.Connect("python_test")

	log.Printf("Client Buffer:  %x \n", clientBuffer)

	num, err := conn.Write(clientBuffer)
	if err != nil {
		log.Printf("Sender: Write Error: %s\n", err)
	}

	log.Printf("Sender: Wrote %d byte(s)\n", num)

	myRead(client, conn)

	// Disconnect
	clientBuffer = client.Disconnect()
	log.Printf("Client Buffer:  %x \n", clientBuffer)

	n, err := conn.Write(clientBuffer)
	if err != nil {
		log.Printf("Sender: Write Error: %s\n", err)
	}

	log.Printf("Sender: Wrote %d byte(s)\n", n)

	myRead(client, conn)
}
