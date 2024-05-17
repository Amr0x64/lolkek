package handlers

import "wb-l0/internal/cache"

type Handlers struct {
	c cache.Cache
}

func New(c cache.Cache) *Handlers {
	return &Handlers{
		c: c,
	}
}
