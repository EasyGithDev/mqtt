package net

import (
	"io"
	"net"
)

type MqttConn struct {
	conn      net.Conn
	connected bool
}

func NewMqttConn(conn net.Conn) *MqttConn {
	return &MqttConn{conn: conn, connected: false}
}

func (mc *MqttConn) Read(buffer []byte) (int, error) {
	n, err := mc.conn.Read(buffer)
	if err != nil {
		if err == io.EOF {
			mc.connected = false
		}
	}
	return n, err
}

func (mc *MqttConn) Write(buffer []byte) (int, error) {
	n, err := mc.conn.Write(buffer)
	if err != nil {
		if err == io.EOF {
			mc.connected = false
		}
	}
	return n, err
}

func (mc *MqttConn) IsConnected() bool {
	return mc.connected
}

func (mc *MqttConn) SetConnected(connected bool) {
	mc.connected = connected
}
