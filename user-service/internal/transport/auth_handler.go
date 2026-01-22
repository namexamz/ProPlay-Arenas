package transport

import (
	"log/slog"
	"net/http"

	"user-service/internal/models"
	service "user-service/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	logger      *slog.Logger
	authService *service.AuthService
}

func NewAuthHandler(
	logger *slog.Logger,
	authService *service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		logger: logger.With(
			slog.String("layer", "transport"),
			slog.String("handler", "auth"),
		),
		authService: authService,
	}
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
}

func (h *AuthHandler) Register(c *gin.Context) {
	h.logger.Info("register request started", slog.String("ip", c.ClientIP()))

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("register validation failed", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": err.Error()})
		return
	}

	token, err := h.authService.RegisterUser(req)
	if err != nil {
		h.logger.Error("register failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": map[string]string{"access_token": token}, "error": nil})
}

func (h *AuthHandler) Login(c *gin.Context) {
	h.logger.Info("login request started", slog.String("ip", c.ClientIP()))

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("login validation failed", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": err.Error()})
		return
	}

	token, err := h.authService.LoginUser(req)
	if err != nil {
		h.logger.Warn("login failed", slog.Any("error", err))
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": map[string]string{"access_token": token}, "error": nil})
}
