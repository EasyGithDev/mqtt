package credentials

// Store login/pass in struct
type MqttCredentials struct {
	Login    string
	Password string
}

func New(login string, password string) *MqttCredentials {
	return &MqttCredentials{Login: login, Password: password}
}
