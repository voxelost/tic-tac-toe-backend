package message

type Messenger struct {
	SendToEventManager              func(*Message)
	ReceiveFromEventManagerCallback func(*Message)
	UnsubFromEventManager           func()
}

func NewMessenger(receiveFromEventManagerCallback func(*Message)) *Messenger {
	return &Messenger{
		ReceiveFromEventManagerCallback: receiveFromEventManagerCallback,
	}
}
