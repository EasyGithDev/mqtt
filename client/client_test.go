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
	"log"
	"os"
	"testing"
)

const (
	connHost = "mqtt.eclipseprojects.io"
	connPort = "1883"
)

func TestConnectDisconnect(t *testing.T) {

	mc := NewMqttClient("test-golang-mqtt")
	_, connErr := mc.TcpConnect(connHost, connPort)
	if connErr != nil {
		log.Print("Error connecting:", connErr.Error())
		os.Exit(1)
	}
	defer mc.TcpDisconnect()

	response, err := mc.MqttConnect()

	if err != nil {
		t.Errorf("Mqtt connection fail %s", err)
	} else if !response {
		t.Errorf("Mqtt connection fail")
	} else {

		response, err := mc.MqttDisconnect()

		if err != nil {
			t.Errorf("Mqtt disconnect fail %s", err)
		} else if !response {
			t.Errorf("Mqtt disconnect fail")
		}
	}
}
