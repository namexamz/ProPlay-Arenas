package transport

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"user-service/internal/dto"
	"user-service/internal/models"
	service "user-service/internal/services"

	"github.com/gin-gonic/gin"
)

func getClaims(c *gin.Context) (*models.Claims, error) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, errors.New("unauthorized")
	}
	userClaims, ok := claims.(*models.Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return userClaims, nil
}



type UserHandler struct {
	service service.UserService
	logger  *slog.Logger
}

func NewUserHandler(service service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger: logger.With(
			slog.String("layer", "transport"),
			slog.String("handler", "user"),
		),
	}
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.GET("/me", h.GetMe)
		users.PUT("/me", h.UpdateMe)
		users.POST("/me/become-owner", h.BecomeOwner)
		users.GET("/:id", h.GetPublicProfile)
	}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	h.logger.Debug("get me request started", slog.String("ip", c.ClientIP()))

	userClaims, err := getClaims(c)
	if err != nil {
		h.logger.Warn("get me unauthorized", slog.Any("error", err))
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": err.Error()})
		return
	}

	

	user, err := h.service.GetMe(userClaims.UserID)
	if err != nil {
		h.logger.Error("get me failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dto.UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     string(user.Role),
	}, "error": nil})
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	h.logger.Info("update me request started", slog.String("ip", c.ClientIP()))

	userClaims, err := getClaims(c)
	if err != nil {
		h.logger.Warn("update me unauthorized", slog.Any("error", err))
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": err.Error()})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("update me validation failed", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": err.Error()})
		return
	}

	user, err := h.service.UpdateMe(userClaims.UserID, req)
	if err != nil {
		h.logger.Error("update me failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dto.UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     string(user.Role),
	}, "error": nil})
}

func (h *UserHandler) BecomeOwner(c *gin.Context) {
	h.logger.Info("become owner request started", slog.String("ip", c.ClientIP()))

	userClaims, err := getClaims(c)
	if err != nil {
		h.logger.Warn("become owner unauthorized", slog.Any("error", err))
		c.JSON(http.StatusUnauthorized, gin.H{"data": nil, "error": err.Error()})
		return
	}
	

	user, err := h.service.BecomeOwner(userClaims.UserID)
	if err != nil {
		h.logger.Warn("become owner failed", slog.Any("error", err))
		c.JSON(http.StatusForbidden, gin.H{"data": nil, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dto.UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     string(user.Role),
	}, "error": nil})
}

func (h *UserHandler) GetPublicProfile(c *gin.Context) {
	h.logger.Debug("get public profile request started", slog.String("ip", c.ClientIP()))

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.logger.Warn("get public profile validation failed", slog.String("id", idParam))
		c.JSON(http.StatusBadRequest, gin.H{"data": nil, "error": "invalid user id"})
		return
	}

	profile, err := h.service.GetPublicProfile(uint(id))
	if err != nil {
		h.logger.Error("get public profile failed", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"data": nil, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": profile, "error": nil})
}
