package transport

import (
	"net/http"
	"strings"

	"gateway/internal/middleware"
	"gateway/internal/proxy"

	"github.com/gin-gonic/gin"
)

type Config struct {
	JWTSecret             string
	UserServiceURL        string
	VenueServiceURL       string
	ReservationServiceURL string
	PaymentServiceURL     string
}

func RegisterRoutes(r *gin.Engine, cfg Config) error {
	userUpstream, err := proxy.NewUpstream(cfg.UserServiceURL, "/api", nil)
	if err != nil {
		return err
	}
	venueUpstream, err := proxy.NewUpstream(cfg.VenueServiceURL, "/api", nil)
	if err != nil {
		return err
	}
	reservationUpstream, err := proxy.NewUpstream(cfg.ReservationServiceURL, "/api", nil)
	if err != nil {
		return err
	}
	paymentUpstream, err := proxy.NewUpstream(cfg.PaymentServiceURL, "/api", nil)
	if err != nil {
		return err
	}

	api := r.Group("/api")
	api.Use(middleware.AuthUnless(cfg.JWTSecret, isPublicRequest))

	api.Any("/auth/*path", gin.WrapH(http.HandlerFunc(userUpstream.ServeHTTP)))
	api.Any("/users/*path", gin.WrapH(http.HandlerFunc(userUpstream.ServeHTTP)))

	api.Any("/venues/:id/availability", gin.WrapH(http.HandlerFunc(reservationUpstream.ServeHTTP)))
	api.Any("/venues/:id/bookings", gin.WrapH(http.HandlerFunc(reservationUpstream.ServeHTTP)))

	api.Any("/venues", gin.WrapH(http.HandlerFunc(venueUpstream.ServeHTTP)))
	api.Any("/venues/*path", gin.WrapH(http.HandlerFunc(venueUpstream.ServeHTTP)))
	api.Any("/venue-types", gin.WrapH(http.HandlerFunc(venueUpstream.ServeHTTP)))
	api.Any("/venue-types/*path", gin.WrapH(http.HandlerFunc(venueUpstream.ServeHTTP)))

	api.Any("/bookings/:id/payment", gin.WrapH(http.HandlerFunc(paymentUpstream.ServeHTTP)))
	api.Any("/payments/*path", gin.WrapH(http.HandlerFunc(paymentUpstream.ServeHTTP)))
	api.Any("/payments", gin.WrapH(http.HandlerFunc(paymentUpstream.ServeHTTP)))

	aggregator := NewAggregator(cfg)
	aggregator.Register(api)

	api.Any("/bookings", gin.WrapH(http.HandlerFunc(reservationUpstream.ServeHTTP)))
	api.Any("/bookings/*path", gin.WrapH(http.HandlerFunc(reservationUpstream.ServeHTTP)))

	return nil
}

func isPublicRequest(r *http.Request) bool {
	path := r.URL.Path
	method := r.Method

	if strings.HasPrefix(path, "/api/auth/") {
		return true
	}

	if method != http.MethodGet {
		return false
	}

	if path == "/api/venue-types" || strings.HasPrefix(path, "/api/venue-types/") {
		return true
	}

	if path == "/api/venues" || strings.HasPrefix(path, "/api/venues/") {
		if strings.HasSuffix(path, "/bookings") {
			return false
		}
		return true
	}

	return false
}

 
