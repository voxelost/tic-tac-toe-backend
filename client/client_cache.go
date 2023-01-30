package client

import (
	"main/utils"
	"sync"
)

// TODO: DOCUMENT
type ClientCache struct {
	sync.Mutex
	cache map[utils.ID]*Client
}

func NewClientCache() *ClientCache {
	return &ClientCache{
		cache: make(map[utils.ID]*Client),
	}
}

func (cc *ClientCache) Register(c *Client) {
	cc.Lock()
	defer cc.Unlock()

	cc.cache[c.GetId()] = c
}

func (cc *ClientCache) Unregister(c *Client) {
	cc.Lock()
	defer cc.Unlock()

	delete(cc.cache, c.GetId())
}

func (cc *ClientCache) GetClientForId(Id utils.ID) *Client {
	cc.Lock()
	defer cc.Unlock()

	return cc.cache[Id]
}

func (cc *ClientCache) Count() int {
	return len(cc.cache)
}
