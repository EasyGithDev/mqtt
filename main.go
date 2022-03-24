package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/easygithdev/mqtt/client"
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

	DEBUG = true
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

func init() {

	if !DEBUG {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

}

func main() {

	conn, err := net.Dial(connType, connHost+":"+connPort)

	if err != nil {
		log.Print("Error connecting:", err.Error())
		os.Exit(1)
	}
	log.Print("Connecting to " + connType + " server " + connHost + ":" + connPort)

	// Connect
	client := client.NewMqttClient()

	// Adding connection to client
	client.SetConn(&conn)

	// Read continue

	resp, err := client.Connect("golang-1")

	if err != nil {
		log.Printf("Connect Error: %s\n", err)
	}

	if resp {
		log.Printf("Connection established: \n")
	}

	resp, err = client.Disconnect()

	if err != nil {
		log.Printf("Connect Error: %s\n", err)
	}

	if resp {
		log.Printf("Disconnected: \n")
	}

	// Disconnect
	// clientBuffer = client.Disconnect()
	// log.Printf("Client Buffer:  %x \n", clientBuffer)

	// n, err := conn.Write(clientBuffer)
	// if err != nil {
	// 	log.Printf("Sender: Write Error: %s\n", err)
	// }

	// log.Printf("Sender: Wrote %d byte(s)\n", n)

	// myRead(client, conn)
}
