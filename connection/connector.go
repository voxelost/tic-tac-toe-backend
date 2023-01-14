package connection

import "main/message"

type Connector struct {
	ConnIn  chan message.Message
	ConnOut chan message.Message
}

func NewConnector() *Connector {
	return &Connector{
		ConnIn:  make(chan message.Message, 1024),
		ConnOut: make(chan message.Message, 1024),
	}
}
