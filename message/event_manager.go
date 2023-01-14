package message

// Struct EventManager represents an Event Manager responsible for receiving Events (Messages) and broadcasting them to
// all subscribers
type EventManager struct {
	Subscribers map[string]*EventManagerSubscriber
}

// return new EventManager object
func NewEventManager() *EventManager {
	em := &EventManager{
		Subscribers: make(map[string]*EventManagerSubscriber),
	}
	return em
}

// Subscribe a Messenger object to this EventManager
func (em *EventManager) SubscribeMessenger(messenger *Messenger) {
	sub := NewEventManagerSubscriber(messenger.ReceiveFromEventManagerCallback)
	em.Subscribers[sub.Id] = sub
	messenger.UnsubFromEventManager = func() {
		delete(em.Subscribers, sub.Id)
	}

	messenger.SendToEventManager = em.Publish
}

// Unsubscribe a Messenger object from this EventManager
func (em *EventManager) UnsubscribeMessenger(messenger *Messenger) {
	messenger.UnsubFromEventManager()
}

// Broadcast message to all Subscribers
func (em *EventManager) Publish(message *Message) {
	for _, v := range em.Subscribers {
		v.Callback(message)
	}
}
