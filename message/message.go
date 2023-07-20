package message

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Message struct {
	Id      string                 `json:"id,omitempty"`
	Type    MessageType            `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

// Return a new Message object
func NewMessage(type_ MessageType, payload map[string]interface{}) *Message {
	return &Message{
		Id:      uuid.New().String(),
		Type:    type_,
		Payload: payload,
	}
}

// Unmarshal given JSON encoded bytes to this Message object
func (m *Message) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, m)
}

// Marshal this object to JSON encoded bytes
func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

type Marshallable interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}
