package cache

import (
	"errors"
	"sync"
	"wb-l0/internal/models"
)

type Cache struct {
	sync.RWMutex
	OrderCache map[string]models.Order
}

func New() Cache {
	return Cache{
		OrderCache: make(map[string]models.Order),
	}
}

func (c *Cache) Get(id string) (models.Order, error) {
	c.RLock()
	defer c.RUnlock()
	if o, ok := c.OrderCache[id]; ok {
		return o, nil
	}

	return models.Order{}, errors.New("order not found")
}

func (c *Cache) Set(order models.Order) {
	c.Lock()
	defer c.Unlock()
	c.OrderCache[order.OrderUid] = order
}
