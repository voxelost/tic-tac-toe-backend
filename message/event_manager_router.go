package message

type (
	EventManagerRouterCallback func(*Message) (shouldPublish bool)
	RouteMap                   map[OriginType]map[MessageType]EventManagerRouterCallback
)

// EventManagerRouter holds all Routes to all registered EventManagerRouterCallbacks
type EventManagerRouter struct {
	Routes RouteMap
}

// Return a new EventManagerRouter object
func NewEventManagerRouter() *EventManagerRouter {
	return &EventManagerRouter{
		Routes: make(RouteMap),
	}
}

// Register a given Route
func (r *EventManagerRouter) Route(originType OriginType, messageType MessageType, eventManagerRouterCallback EventManagerRouterCallback) *EventManagerRouter {
	_, ok := r.Routes[originType]
	if !ok {
		r.Routes[originType] = make(map[MessageType]EventManagerRouterCallback)
	}
	r.Routes[originType][messageType] = eventManagerRouterCallback
	return r
}

// Match a message agains the EventManagerRouter. If a Route is matched, return the corresponding
// EventManagerRouterCallback, true, else return nil, false
func (r *EventManagerRouter) Match(message *Message) (EventManagerRouterCallback, bool) {
	if message == nil {
		return nil, false
	}
	if _, ok := r.Routes[message.Origin.Type]; !ok {
		return nil, false
	}

	if _, ok := r.Routes[message.Origin.Type][message.Type]; !ok {
		return nil, false
	}

	return r.Routes[message.Origin.Type][message.Type], true
}
