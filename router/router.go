package router

import (
	"net/http"
	"strconv"

	"github.com/around-project/handler"
	"github.com/gin-gonic/gin"
)

type Router struct {
	handler *handler.Handler
}

func NewRouter(h *handler.Handler) *Router {
	return &Router{handler: h}
}

func (r *Router) MakeOrder(c *gin.Context) {

	orderNumber := c.Param("order_id")
	if orderNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order number"})
		return
	}

	var req handler.OrderDetails

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	r.handler.MakeOrder(orderNumber, req.Latitude, req.Longitude)
	c.Status(http.StatusOK)
}

func (r *Router) GetOrder(c *gin.Context) {

	var err error
	orderNumber := c.Param("order_id")
	if orderNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order number"})
		return
	}

	maxInt := 0
	max, ok := c.GetQuery("max")
	if ok {
		maxInt, err = strconv.Atoi(max)
		if err != nil {
			c.JSON(http.StatusPreconditionFailed, gin.H{"error": err})
			return
		}
	}

	resp, err := r.handler.GetOrder(orderNumber, maxInt)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (r *Router) DeleteOrder(c *gin.Context) {

	orderNumber := c.Param("order_id")
	if orderNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order number"})
		return
	}

	err := r.handler.DeleteOrder(orderNumber)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.Status(http.StatusOK)
}
