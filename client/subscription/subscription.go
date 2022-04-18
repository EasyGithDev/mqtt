package subscription

type Subscriptions map[string]Subscription

type Subscription struct {
	Topic string
	Qos   byte
}

func New(topic string, qos byte) *Subscription {
	return &Subscription{topic, qos}
}
