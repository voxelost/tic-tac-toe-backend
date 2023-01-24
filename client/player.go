package client

type Player struct {
	*Client
	Payload map[string]string
}

func NewPlayer(c *Client) *Player {
	return &Player{
		Client:  c,
		Payload: make(map[string]string),
	}
}
