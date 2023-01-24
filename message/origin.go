package message

import "main/utils"

type OriginType string

// possible Message Origins
var (
	Client = OriginType("client")
	Game   = OriginType("game")
	Server = OriginType("server")
)

// Origin represents the entity that emmited a message
type Origin struct {
	*utils.ID `json:"id"`
	Type      OriginType `json:"type"`
}

// Return new Origin object
func NewOrigin(type_ OriginType, id *utils.ID) *Origin {
	return &Origin{
		Type: type_,
		ID:   id,
	}
}
