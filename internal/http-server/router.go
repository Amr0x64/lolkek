package router

import (
	"net/http"
	"wb-l0/internal/http-server/handlers"

	"github.com/gin-gonic/gin"
)

func Router(h *handlers.Handlers) {
	r := gin.Default()
	r.LoadHTMLGlob("D:/VsCodeProjects/wbl_l0/templates/index.html")

	r.GET("/getOrder", h.GetOrder)
	r.GET("/order", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.Run("localhost:6002")

}
