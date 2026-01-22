package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gateway/internal/models"
)

type Aggregator struct {
	reservationURL string
	venueURL       string
	paymentURL     string
	client         *http.Client
}

func NewAggregator(cfg Config) *Aggregator {
	return &Aggregator{
		reservationURL: cfg.ReservationServiceURL,
		venueURL:       cfg.VenueServiceURL,
		paymentURL:     cfg.PaymentServiceURL,
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

func (a *Aggregator) Register(rg *gin.RouterGroup) {
	rg.GET("/bookings/:id/summary", a.GetBookingSummary)
}

func (a *Aggregator) GetBookingSummary(c *gin.Context) {
	bookingID := c.Param("id")
	if _, err := strconv.ParseUint(bookingID, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	bookingBody, status, err := a.fetchJSON(c, fmt.Sprintf("%s/bookings/%s", a.reservationURL, bookingID))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "reservation service unavailable"})
		return
	}
	if status != http.StatusOK {
		c.Data(status, "application/json", bookingBody)
		return
	}

	var bookingLookup models.BookingSummaryLookup
	if err := json.Unmarshal(bookingBody, &bookingLookup); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "invalid booking response"})
		return
	}

	venueBody, venueStatus, venueErr := a.fetchJSON(c, fmt.Sprintf("%s/venues/%d", a.venueURL, bookingLookup.VenueID))
	paymentBody, paymentStatus, paymentErr := a.fetchJSON(c, fmt.Sprintf("%s/bookings/%s/payment", a.paymentURL, bookingID))

	resp := models.BookingSummaryResponse{
		Booking: bookingBody,
	}

	if venueErr != nil {
		msg := "venue service unavailable"
		resp.VenueError = &msg
	} else if venueStatus != http.StatusOK {
		msg := fmt.Sprintf("venue service status %d", venueStatus)
		resp.VenueError = &msg
	} else {
		resp.Venue = venueBody
	}

	if paymentErr != nil {
		msg := "payment service unavailable"
		resp.PaymentError = &msg
	} else if paymentStatus != http.StatusOK {
		msg := fmt.Sprintf("payment service status %d", paymentStatus)
		resp.PaymentError = &msg
	} else {
		resp.Payment = paymentBody
	}

	c.JSON(http.StatusOK, resp)
}

func (a *Aggregator) fetchJSON(c *gin.Context, target string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, target, nil)
	if err != nil {
		return nil, 0, err
	}
	if auth := c.GetHeader("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}
