package message

import (
	"main/utils"
)

// Messenger struct represents the connection between a Message emitter/sink and an Event Manager
type Messenger struct {
	utils.ID
	communicators           []*Communicator
	ReceiveFromEventManager func(*Message)
}

// Return a new Messenger object
func NewMessenger(id utils.ID) *Messenger {
	return &Messenger{
		ID: id,
	}
}

func (m *Messenger) PushCommunicator(Communicator *Communicator) {
	m.communicators = append(m.communicators, Communicator)
}

func (m *Messenger) GetActiveCommunicator() *Communicator {
	if len(m.communicators) < 1 {
		return nil
	}
	return m.communicators[len(m.communicators)-1]
}

func (m *Messenger) PopCommunicator() {
	communicator_ := m.GetActiveCommunicator()
	if communicator_ == nil {
		return
	}

	m.communicators = m.communicators[:len(m.communicators)-1]
	communicator_.UnsubFromEventManager()
}

// Pop all communicators
func (m *Messenger) UnsubFromAllEventManagers() {
	for len(m.communicators) > 0 {
		m.PopCommunicator()
	}
}

// todo document
func (m *Messenger) SetEventManagerReceiveCallback(func_ func(*Message)) {
	m.ReceiveFromEventManager = func_
}
