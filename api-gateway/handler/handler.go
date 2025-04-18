package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Assignment/api-gateway/service"
)

func RegisterAuthRoutes(router *gin.RouterGroup, authService service.AuthService) {
	handler := &AuthHandler{authService: authService}
	router.POST("/auth/login", handler.Login)
	router.POST("/auth/register", handler.Register)
}

type AuthHandler struct {
	authService service.AuthService
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.GenerateToken(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.GenerateToken(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func RegisterProxyRoutes(router *gin.RouterGroup, basePath string, proxyService service.ProxyService) {
	router.Any(basePath, proxyHandler(proxyService))
	router.Any(basePath+"/*path", proxyHandler(proxyService))
}

func proxyHandler(proxyService service.ProxyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")

		headers := make(map[string]string)
		headers["X-Request-ID"] = c.GetHeader("X-Request-ID")

		if exists {
			headers["X-User-ID"] = userID.(string)
		}

		resp, err := proxyService.ProxyRequest(
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Body,
			headers,
		)

		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Gateway error: " + err.Error()})
			return
		}

		c.Status(resp.StatusCode)

		for k, v := range resp.Headers {
			c.Header(k, v)
		}

		c.Writer.Write(resp.Body)
	}
}
