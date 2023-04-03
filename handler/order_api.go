package handler

import (
	"ecommerce/domain/order"
	"ecommerce/dto/request"
	"net/http"
	"strconv"

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

func (h *OrderApiHandler) CheckoutOrder(c *gin.Context) {
	userId, err := fetchUserId(c)
	if userId == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	orderIdStr := c.Param("id")
	orderId, err := strconv.Atoi(orderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.orderSvc.CheckoutOrder(*userId, orderId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *OrderApiHandler) GetOrderHistories(c *gin.Context) {
	userId, err := fetchUserId(c)
	if userId == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	orders, err := h.orderSvc.GetOrderHistories(*userId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    orders,
	})
}

func (h *OrderApiHandler) GetAllOrders(c *gin.Context) {
	userId, err := fetchUserId(c)
	if userId == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	orders, err := h.orderSvc.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    orders,
	})
}

func (h *OrderApiHandler) GetAllProducts(c *gin.Context) {
	userId, err := fetchUserId(c)
	if userId == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	products, err := h.orderSvc.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    products,
	})
}
