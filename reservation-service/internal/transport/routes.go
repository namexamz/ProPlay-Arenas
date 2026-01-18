package transport

import (
	"reservation/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	reservationServ service.BookingService,
){
	reservationHandler := NewBookingHandler(r, reservationServ)

	reservationHandler.Register(r)
}