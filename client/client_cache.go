package client

import (
	"main/utils"
)

// TODO: DOCUMENT
type ClientCache struct {
	cache map[utils.ID]*Client
}

func NewClientCache() *ClientCache {
	return &ClientCache{
		cache: make(map[utils.ID]*Client),
	}
}

func (cc *ClientCache) Register(c *Client) {
	cc.cache[c.GetId()] = c
}

func (cc *ClientCache) Unregister(c *Client) {
	delete(cc.cache, c.GetId())
}

func (cc *ClientCache) GetClientForId(Id utils.ID) *Client {
	return cc.cache[Id]
}

func (cc *ClientCache) Count() int {
	return len(cc.cache)
}
