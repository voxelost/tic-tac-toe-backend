package message

import (
	"encoding/json"
	"main/utils"
)

type Message struct {
	utils.ID           `json:"id,omitempty"`
	RecipientId        utils.ID    `json:"recipient_id,omitempty"`
	Type               MessageType `json:"type"`
	Origin             *Origin     `json:"origin,omitempty"`
	EventManagerOrigin *Origin     `json:"event_manager_origin,omitempty"`
	Payload            interface{} `json:"payload,omitempty"`
}

// Return a new Message object
func NewMessage(type_ MessageType, payload interface{}) *Message {
	return &Message{
		ID:      *utils.NewId(),
		Type:    type_,
		Payload: payload,
	}
}

// Return a new Message object from given JSON encoded bytes
func NewMessageFromBytes(bytes []byte) *Message {
	message := new(Message)
	err := message.Unmarshal(bytes)
	if err != nil {
		return nil
	}
	message.ID = *utils.NewId()
	return message
}

// Set message Origin
func (m *Message) SetOrigin(origin *Origin) {
	m.Origin = origin
}

// Set message EventManagerOrigin
func (m *Message) SetEventManagerOrigin(origin *Origin) {
	m.EventManagerOrigin = origin
}

// Unmarshal given JSON encoded bytes to this Message object
func (m *Message) Unmarshal(bytes []byte) error {
	defer recover()

	return json.Unmarshal(bytes, m)
}

// Marshal this object to JSON encoded bytes
func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
