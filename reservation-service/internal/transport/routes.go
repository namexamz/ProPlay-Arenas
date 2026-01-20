package transport

import (
	"reservation/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	reservationServ service.BookingService,
	jwtSecret string,
){
	reservationHandler := NewBookingHandler( reservationServ)

	reservationHandler.Register(r, jwtSecret)
}