package message

import "github.com/google/uuid"

type EventManagerCallback func(*Message)

type EventManagerSubscriber struct {
	Id       string
	Callback EventManagerCallback
}

// return a new EventManagerSubscriber object
func NewEventManagerSubscriber(callback EventManagerCallback) *EventManagerSubscriber {
	return &EventManagerSubscriber{
		Id:       uuid.New().String(),
		Callback: callback,
	}
}
