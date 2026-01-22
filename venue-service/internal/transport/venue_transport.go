package transport

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"venue-service/internal/models"
	"venue-service/internal/services"

	"github.com/gin-gonic/gin"
)

type VenueHandler struct {
	service services.VenueService
	logger  *slog.Logger
}

type GetVenuesQuery struct {
	District  string `form:"district"`
	VenueType string `form:"venue_type"`
	HourPrice int    `form:"hour_price"`
	IsActive  *bool  `form:"is_active"`
	OwnerID   uint   `form:"owner_id"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Limit     int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

func NewVenueHandler(service services.VenueService, logger *slog.Logger) *VenueHandler {
	return &VenueHandler{
		service: service,
		logger:  logger.With("layer", "transport"),
	}
}

func (h *VenueHandler) RegisterRoutes(r *gin.Engine) {
	venues := r.Group("/venues")
	{
		venues.GET("", h.GetList)
		venues.POST("", h.Create)
		venues.GET("/:id/schedule", h.GetSchedule)
		venues.PUT("/:id/schedule", h.UpdateSchedule)
		venues.GET("/:id", h.GetByID)
		venues.PUT("/:id", h.Update)
		venues.DELETE("/:id", h.Delete)
	}

	venueTypes := r.Group("/venue-types")
	{
		venueTypes.GET("", h.GetVenueTypes)
	}

	users := r.Group("/users")
	{
		users.GET("/:id/venues", h.GetByOwnerID)
	}
}

func (h *VenueHandler) GetList(c *gin.Context) {
	var query GetVenuesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		h.logger.Error("Ошибка парсинга query параметров", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Установка дефолтных значений для пагинации
	if query.Limit == 0 {
		query.Limit = 10
	}
	if query.Page == 0 {
		query.Page = 1
	}

	filter := services.VenueFilter{
		District:  query.District,
		HourPrice: query.HourPrice,
		IsActive:  query.IsActive,
		OwnerID:   query.OwnerID,
		Page:      query.Page,
		Limit:     query.Limit,
	}

	// Валидация VenueType
	if query.VenueType != "" {
		venueType := models.VenueType(query.VenueType)
		if !venueType.IsValid() {
			h.logger.Error("Неверный тип площадки", "venue_type", query.VenueType)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("неверный тип площадки: %s", query.VenueType),
			})
			return
		}
		filter.VenueType = venueType
	}

	venues, err := h.service.GetList(filter)
	if err != nil {
		h.logger.Error("Ошибка получения списка площадок", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Конвертируем модели в DTO
	dtoList := ToVenueDTOList(venues)
	c.JSON(http.StatusOK, dtoList)
}

func (h *VenueHandler) GetByID(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		return
	}

	venue, err := h.service.GetByID(id)
	if err != nil {
		if err == services.ErrVenueNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Error("Ошибка получения площадки по ID", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Площадка успешно получена", "id", id)
	// Конвертируем модель в DTO
	venueDTO := ToVenueDTO(venue)
	c.JSON(http.StatusOK, venueDTO)
}

func (h *VenueHandler) Create(c *gin.Context) {
	var dto VenueDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		h.logger.Error("Ошибка парсинга JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Валидация VenueType
	venueType := models.VenueType(dto.VenueType)
	if !venueType.IsValid() {
		h.logger.Error("Неверный тип площадки", "venue_type", dto.VenueType)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("неверный тип площадки: %s", dto.VenueType),
		})
		return
	}

	// Конвертируем DTO в модель
	venue, err := FromVenueDTO(&dto)
	if err != nil {
		h.logger.Error("Ошибка конвертации DTO", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.Create(venue); err != nil {
		h.logger.Error("Ошибка создания площадки", "venue_type", venue.VenueType, "owner_id", venue.OwnerID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Площадка успешно создана", "id", venue.ID, "venue_type", venue.VenueType, "owner_id", venue.OwnerID)
	// Конвертируем модель обратно в DTO для ответа
	venueDTO := ToVenueDTO(venue)
	c.JSON(http.StatusCreated, venueDTO)
}

func (h *VenueHandler) Update(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		return
	}

	// PUT-семантика: требуем все обязательные поля для полного обновления записи
	var dto VenueDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		h.logger.Error("Ошибка парсинга JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("PUT требует все обязательные поля: %v", err),
		})
		return
	}

	// Валидация всех обязательных полей для PUT
	// binding:"required" работает для строк, но для int/uint нужна ручная проверка

	// Проверка VenueType (не пустой и валидный)
	if dto.VenueType == "" {
		h.logger.Error("Отсутствует обязательное поле venue_type")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "PUT требует все поля: venue_type обязателен",
		})
		return
	}
	venueType := models.VenueType(dto.VenueType)
	if !venueType.IsValid() {
		h.logger.Error("Неверный тип площадки", "venue_type", dto.VenueType)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("неверный тип площадки: %s", dto.VenueType),
		})
		return
	}

	// Проверка OwnerID (не может быть 0)
	if dto.OwnerID == 0 {
		h.logger.Error("Отсутствует обязательное поле owner_id")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "PUT требует все поля: owner_id обязателен и не может быть 0",
		})
		return
	}

	// Проверка HourPrice (может быть 0, но не отрицательным)
	if dto.HourPrice < 0 {
		h.logger.Error("Некорректное значение hour_price", "hour_price", dto.HourPrice)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "hour_price не может быть отрицательным",
		})
		return
	}

	// Проверка District (binding:"required" работает, но проверяем явно для ясности)
	if dto.District == "" {
		h.logger.Error("Отсутствует обязательное поле district")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "PUT требует все поля: district обязателен",
		})
		return
	}

	// Конвертируем DTO в модель
	venue, err := FromVenueDTO(&dto)
	if err != nil {
		h.logger.Error("Ошибка конвертации DTO", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.Update(id, venue); err != nil {
		if err == services.ErrVenueNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Error("Ошибка обновления площадки", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Площадка успешно обновлена", "id", id)
	// Получаем обновленную площадку для ответа
	updatedVenue, err := h.service.GetByID(id)
	if err != nil {
		h.logger.Error(
			"КРИТИЧЕСКАЯ ОШИБКА: обновление площадки успешно, но запись недоступна при повторном чтении",
			"id", id,
			"error", err,
			"severity", "critical",
			"anomaly", true,
		)
		// Возвращаем 204 No Content, т.к. обновление прошло успешно, но вернуть данные не можем
		c.Status(http.StatusNoContent)
		return
	}
	venueDTO := ToVenueDTO(updatedVenue)
	c.JSON(http.StatusOK, venueDTO)
}

func (h *VenueHandler) Delete(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		return
	}

	if err := h.service.Delete(id); err != nil {
		if err == services.ErrVenueNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Error("Ошибка деактивации площадки", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Площадка успешно деактивирована", "id", id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Площадка деактивирована",
	})
}

func (h *VenueHandler) GetSchedule(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		return
	}

	venue, err := h.service.GetSchedule(id)
	if err != nil {
		if err == services.ErrVenueNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Error("Ошибка получения расписания", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Конвертируем модель в ScheduleDTO
	scheduleDTO := ToScheduleDTO(venue)
	c.JSON(http.StatusOK, scheduleDTO)
}

func (h *VenueHandler) UpdateSchedule(c *gin.Context) {
	id, err := h.parseID(c)
	if err != nil {
		return
	}

	var dto ScheduleUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		h.logger.Error("Ошибка парсинга JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Конвертируем DTO в models.Weekdays
	weekdays, err := FromScheduleUpdateDTO(&dto)
	if err != nil {
		h.logger.Error("Ошибка конвертации DTO расписания", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.UpdateSchedule(id, services.ScheduleUpdate{Weekdays: *weekdays}); err != nil {
		if err == services.ErrVenueNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		h.logger.Error("Ошибка обновления расписания", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("Расписание успешно обновлено", "id", id)
	// Возвращаем обновленное расписание
	updatedVenue, err := h.service.GetSchedule(id)
	if err != nil {
		h.logger.Error(
			"КРИТИЧЕСКАЯ ОШИБКА: обновление расписания успешно, но запись недоступна при повторном чтении",
			"id", id,
			"error", err,
			"severity", "critical",
			"anomaly", true,
		)
		c.Status(http.StatusNoContent)
		return
	}
	scheduleDTO := ToScheduleDTO(updatedVenue)
	c.JSON(http.StatusOK, scheduleDTO)
}

func (h *VenueHandler) GetVenueTypes(c *gin.Context) {
	types := []gin.H{
		{"value": string(models.VenueFootball), "label": "Футбол"},
		{"value": string(models.VenueBasketball), "label": "Баскетбол"},
		{"value": string(models.VenueTennis), "label": "Теннис"},
		{"value": string(models.VenueGym), "label": "Тренажерный зал"},
		{"value": string(models.VenueSwimming), "label": "Плавание"},
	}
	c.JSON(http.StatusOK, types)
}

func (h *VenueHandler) GetByOwnerID(c *gin.Context) {
	ownerID, err := h.parseID(c)
	if err != nil {
		return
	}

	venues, err := h.service.GetByOwnerID(ownerID)
	if err != nil {
		h.logger.Error("Ошибка получения площадок владельца", "owner_id", ownerID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Конвертируем модели в DTO
	dtoList := ToVenueDTOList(venues)
	c.JSON(http.StatusOK, dtoList)
}

// parseID вспомогательная функция для парсинга ID из параметра
func (h *VenueHandler) parseID(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error("Неверный формат ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный формат ID",
		})
		return 0, err
	}
	if id == 0 {
		h.logger.Error("ID не может быть 0")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID не может быть 0",
		})
		return 0, strconv.ErrRange
	}
	return uint(id), nil
}
