package mqttout

import (
	"github.com/elastic/beats/libbeat/outputs/codec"
)

type config struct {
	Host     string       `config:"host"`
	Port     int          `config:"port"` 
	User     string       `config:"user"`
	Password string       `config:"password"`
	Topic    string       `config:"topic"`
	Codec    codec.Config `config:"codec"`
}

var (
	defaultConfig = config{
		Port: 1883,
	}
)

func (c *config) Validate() error {
	return nil
}
