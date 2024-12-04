package propagation

import (
	"github.com/eclipse/paho.golang/packets"
	"github.com/eclipse/paho.golang/paho"
)

// MQTTCarrier adapts http.Header to satisfy the TextMapCarrier interface.
type MQTTCarrier struct {
	*paho.Publish
}

func NewMQTTCarrier(p *paho.Publish) *MQTTCarrier {
	if p.Properties == nil {
		p.InitProperties(&packets.Properties{
			User: make([]packets.User, 0),
		})
	}
	return &MQTTCarrier{p}
}

// Get returns the value associated with the passed key.
func (mc MQTTCarrier) Get(key string) string {
	return mc.Properties.User.Get(key)
}

// Set stores the key-value pair.
func (mc MQTTCarrier) Set(key string, value string) {
	mc.Properties.User.Add(key, value)
}

// Keys lists the keys stored in this carrier.
func (mc MQTTCarrier) Keys() []string {
	keys := make([]string, 0, len(mc.Properties.User))
	for _, p := range mc.Properties.User {
		keys = append(keys, p.Key)
	}
	return keys
}
