package message

// TODO: document

type Communicator struct {
	SendToEventManager    func(*Message)
	UnsubFromEventManager func()
}

func NewCommunicator() *Communicator {
	return &Communicator{}
}

func (c *Communicator) SetEventManagerSendCallback(func_ func(*Message)) {
	c.SendToEventManager = func_
}

func (c *Communicator) SetEventManagerUnsubCallback(func_ func()) {
	c.UnsubFromEventManager = func_
}
