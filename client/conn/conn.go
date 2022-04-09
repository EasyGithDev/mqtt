package conn

const (
	DEFAULT_TRANSPORT  = "tcp"
	DEFAULT_PORT       = "1883"
	DEFAULT_KEEP_ALIVE = 60
)

// Store mqtt conn in struct
type MqttConn struct {
	Host        string
	Port        string
	Transport   string
	KeepAlive   uint16
	BindAddress string
}

func New(host string, opts ...ConnOption) *MqttConn {
	mc := &MqttConn{
		Host:        host,
		Port:        DEFAULT_PORT,
		Transport:   DEFAULT_TRANSPORT,
		KeepAlive:   DEFAULT_KEEP_ALIVE,
		BindAddress: ""}

	// Apply options
	for _, applyOpt := range opts {
		applyOpt(mc)
	}

	return mc
}

type ConnOption func(f *MqttConn)

func WithPort(port string) ConnOption {
	return func(mc *MqttConn) {
		mc.Port = port
	}
}

func WithKeepAlive(keepAlive uint16) ConnOption {
	return func(mc *MqttConn) {
		mc.KeepAlive = keepAlive
	}
}

func WithTransport(transport string) ConnOption {
	return func(mc *MqttConn) {
		mc.Transport = transport
	}
}
