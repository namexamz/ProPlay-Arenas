package transport

import (
	"reservation/internal/dto"
	
	"reservation/internal/service"
	"strconv"

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
	c.POST("/bookings/:id/cancel", r.CancelReservation)
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

func (r *BookingHandler) GetByID(c *gin.Context){
	idstr := c.Param("id")

	id, err := strconv.Atoi(idstr)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	reservation, error :=  r.bookingService.GetByID(uint(id))
	if error != nil {
		c.JSON(404, gin.H{
			"error":err.Error(),
		})
		return
	}

	c.JSON(200, reservation)
	
}
