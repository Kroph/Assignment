package handler

import (
	"net/http"
	"strconv"
	"time"

	"order-service/internal/domain"
	"order-service/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order domain.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	order.UserID = userID.(string)

	createdOrder, err := h.orderService.CreateOrder(c.Request.Context(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdOrder)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderService.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists || order.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var request struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := domain.OrderStatus(request.Status)
	validStatuses := map[domain.OrderStatus]bool{
		domain.OrderStatusPending:   true,
		domain.OrderStatusPaid:      true,
		domain.OrderStatusShipped:   true,
		domain.OrderStatusDelivered: true,
		domain.OrderStatusCancelled: true,
	}

	if !validStatuses[status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	if err := h.orderService.UpdateOrderStatus(c.Request.Context(), id, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var filter domain.OrderFilter
	filter.UserID = userID.(string)

	if status := c.Query("status"); status != "" {
		filter.Status = domain.OrderStatus(status)
	}

	if fromDateStr := c.Query("from_date"); fromDateStr != "" {
		fromDate, err := time.Parse(time.RFC3339, fromDateStr)
		if err == nil {
			filter.FromDate = &fromDate
		}
	}

	if toDateStr := c.Query("to_date"); toDateStr != "" {
		toDate, err := time.Parse(time.RFC3339, toDateStr)
		if err == nil {
			filter.ToDate = &toDate
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filter.Page = page
	filter.PageSize = pageSize

	orders, total, err := h.orderService.ListOrders(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":    orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func RegisterOrderRoutes(router *gin.RouterGroup, orderService service.OrderService) {
	handler := NewOrderHandler(orderService)

	router.POST("/orders", handler.CreateOrder)
	router.GET("/orders/:id", handler.GetOrderByID)
	router.PATCH("/orders/:id", handler.UpdateOrderStatus)
	router.GET("/orders", handler.ListOrders)
}
