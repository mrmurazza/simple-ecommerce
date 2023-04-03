package handler

import (
	"ecommerce/domain/order"
	"ecommerce/dto/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderApiHandler struct {
	orderSvc order.Service
}

func NewOrderApiHandler(orderSvc order.Service) *OrderApiHandler {
	return &OrderApiHandler{
		orderSvc: orderSvc,
	}
}

func (h *OrderApiHandler) Order(c *gin.Context) {
	userId, err := fetchUserId(c)
	if userId == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req request.CreateOrderRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.UserId = *userId
	err = req.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.orderSvc.Order(req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
