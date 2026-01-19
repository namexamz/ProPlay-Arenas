package transport

import (
	"reservation/internal/dto"
	"reservation/internal/middleware"
	"reservation/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(c *gin.Engine, bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

func (r *BookingHandler) Register(c *gin.Engine, jwtSecret string) {
	c.POST("/booking", r.CreateReservation)
	c.POST("/bookings/:id/cancel", r.CancelReservation)
	c.GET("/bookings/:id", r.GetByID)
	c.GET("/bookings", middleware.AuthMiddleware(jwtSecret), r.GetUserReservations)
	c.PUT("/bookings/:id", r.UpdateReservation)
}

func (r *BookingHandler) CreateReservation(c *gin.Context) {
	var dto dto.ReservationCreate

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	reservation, err := r.bookingService.CreateReservation(&dto)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, reservation)
}

func (r *BookingHandler) CancelReservation(c *gin.Context) {
	var dto dto.ReservationCancel
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid reservation ID"})
		return
	}

	reservation, err := r.bookingService.ReservationCancel(uint(id), dto.Reason)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, reservation)
}

func (r *BookingHandler) GetByID(c *gin.Context) {
	idstr := c.Param("id")

	id, err := strconv.Atoi(idstr)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	reservation, err := r.bookingService.GetByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, reservation)

}

func (r *BookingHandler) GetUserReservations(c *gin.Context) {

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	clientID, ok := userID.(uint)
	if !ok {
		c.JSON(401, gin.H{"error": "invalid user ID type"})
		return
	}

	reservation, err := r.bookingService.GetUserReservations(clientID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, reservation)
}

func (r *BookingHandler) UpdateReservation(c *gin.Context) {
	var dto dto.ReservationUpdate

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid reservation ID"})
		return
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	reservation, err := r.bookingService.ReservationUpdate(uint(id), &dto)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, reservation)

}
