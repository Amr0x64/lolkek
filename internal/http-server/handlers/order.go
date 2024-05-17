package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetOrder(c *gin.Context) {
	uid, ok := c.GetQuery("order_uid")
	log.Println(uid)
	if !ok {
		c.JSON(400, "Get params is empty")
	}
	
	order, err := h.c.Get(uid)
	if err != nil {
		c.JSON(400, err.Error())
		return
	}
	c.JSON(200, order)
}
