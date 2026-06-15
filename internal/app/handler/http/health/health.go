package rhealth

import (
	"net/http"

	rhandler "github.com/Lagwick/order-service/internal/app/handler/http"
	"github.com/gin-gonic/gin"
)

type handler struct{}

func NewHandler() rhandler.Health {
	return &handler{}
}

func (h *handler) LastCheck(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
