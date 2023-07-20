package message

import (
	"main/utils"
	"sync"
)

// Struct EventManager represents an Event Manager responsible for receiving Events (Messages) and broadcasting them to
// all subscribers
type EventManager struct {
	sync.Mutex

	Subscribers map[utils.ID]*Messenger
	Router      *EventManagerRouter
	Origin      *Origin
}

// return new EventManager object
func NewEventManager(origin *Origin) *EventManager {
	return &EventManager{
		Subscribers: make(map[utils.ID]*Messenger),
		Router:      NewEventManagerRouter(),
		Origin:      origin,
	}
}

// Subscribe a Messenger object to this EventManager
func (em *EventManager) SubscribeMessenger(messenger *Messenger) {
	communicator := NewCommunicator()
	communicator.SetEventManagerSendCallback(em.Receive)
	communicator.SetEventManagerUnsubCallback(func() { em.UnsubscribeMessenger(messenger) })

	messenger.PushCommunicator(communicator)

	em.Lock()
	defer em.Unlock()
	em.Subscribers[messenger.GetId()] = messenger
}

// Unsubscribe a Messenger object from this EventManager
func (em *EventManager) UnsubscribeMessenger(messenger *Messenger) {
	if messenger != nil {
		em.Lock()
		defer em.Unlock()
		delete(em.Subscribers, messenger.GetId())
	}
}

// Receive a message from Messenger
func (em *EventManager) Receive(message *Message) {
	if matchedRoute, ok := em.Router.Match(message); ok {
		if matchedRoute(message) {
			em.Publish(message)
		}
	}
}

// Broadcast message to all Subscribers
func (em *EventManager) Publish(message *Message) {
	if message.Origin == nil {
		return
	}
	em.Lock()
	defer em.Unlock()

	message.EventManagerOrigin = em.Origin // messages will always be sent from a single Event Manager only
	msgCopy := *message

	for _, subscriber := range em.Subscribers {
		if subscriber.ReceiveFromEventManager != nil {
			subscriber.ReceiveFromEventManager(&msgCopy)
		}
	}
}
