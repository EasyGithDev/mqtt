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
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/easygithdev/mqtt/client"
	"github.com/easygithdev/mqtt/client/conn"
)

const (
	connHost = "test.mosquitto.org"
	// connHost = "mqtt.eclipseprojects.io"

	// connPort = "1883"
	connPort = "1884"

	DEBUG = true
)

func init() {

	if !DEBUG {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

}

var onTcpConnect = func(conn net.Conn) {
	fmt.Println("Connecting to server " + conn.RemoteAddr().String())
}

var onTcpDisconnect = func() {
	fmt.Println("Disconnect from server")
}

var onConnect = func() {
	fmt.Println("Connecting to MQTT server ...")
}

var onPing = func() {
	fmt.Println("Ping server ...")
}

var onPublish = func(topic string, msg string, qos int) {
	fmt.Printf("Publish on %s:%s with QoS:%d\n", topic, msg, qos)
}

var onSubscribe = func(topic string, qos int) {
	fmt.Printf("Subscribe to %s with QoS:%d\n", topic, qos)
}

var OnMessage = func(msg string) {
	fmt.Println("msg: " + msg)
}

func idGenerator(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {

	pub := flag.Bool("pub", false, "publish")
	sub := flag.Bool("sub", false, "subscribe")
	host := flag.String("h", "", "the hostname")
	port := flag.String("p", "1883", "the port")
	id := flag.String("i", idGenerator(10), "the client Id")
	topic := flag.String("t", "", "the topic name")
	msg := flag.String("m", "", "the message to send")
	qos := flag.Int("qos", 0, "quality of service")
	login := flag.String("log", "", "the login")
	pwd := flag.String("pwd", "", "the password")

	flag.Parse()

	if !*pub && !*sub {
		fmt.Fprintf(os.Stderr, "missing required -pub/-sub flag\n")
		os.Exit(2)
	}

	if strings.TrimSpace(*host) == "" {
		fmt.Fprintf(os.Stderr, "missing required -h argument/flag\n")
		os.Exit(2)
	}

	if strings.TrimSpace(*topic) == "" {
		fmt.Fprintf(os.Stderr, "missing required -t argument/flag\n")
		os.Exit(2)
	}

	if *pub && strings.TrimSpace(*msg) == "" {
		fmt.Fprintf(os.Stderr, "missing required -m flag\n")
		os.Exit(2)
	}

	var credentialsOption client.Option = nil
	if *login != "" && *pwd != "" {
		credentialsOption = client.WithCredentials(*login, *pwd)
	}

	///////////////////////////////////////////////////////////
	// Init
	///////////////////////////////////////////////////////////
	mc := client.New(
		// client Id
		*id,
		// credentials
		credentialsOption,
		// connection infos
		client.WithConnInfos(conn.New(*host, conn.WithPort(*port))),
	)

	// handlers

	mc.OnTcpConnect = onTcpConnect
	mc.OnTcpDisconnect = onTcpDisconnect
	mc.OnConnect = onConnect
	mc.OnPing = onPing
	mc.OnPublish = onPublish
	mc.OnSubscribe = onSubscribe
	mc.OnMessage = OnMessage

	// Connection

	_, connErr := mc.Connect()
	if connErr != nil {
		log.Print("Error connecting:", connErr.Error())
		os.Exit(1)
	}
	defer mc.Close()

	if *pub {
		mc.Publish(*topic, *msg, *qos)
	} else if *sub {

		respSub, errSub := mc.Subscribe(*topic)
		if errSub != nil {
			log.Printf("Subscribe Error: %s\n", errSub)
		}

		if respSub {
			mc.LoopForever()
		}
	}

}
