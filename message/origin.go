package message

// possible Origins
var (
	Client = "client"
	Game   = "game"
	Server = "server"
)

// Origin represents the entity that emmited a message
type Origin struct {
	Type string
	Id   *string
}

// Return new Origin object
func NewOrigin(type_ string, id *string) *Origin {
	return &Origin{
		Type: type_,
		Id:   id,
	}
}
