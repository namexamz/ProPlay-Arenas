package transport

import (
	"reservation/internal/dto"
	"reservation/internal/service"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

func (r *BookingHandler) Register(c *gin.Engine) {
	c.POST("/booking", r.CreateReservation)
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
