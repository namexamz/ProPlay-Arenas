package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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
	if _, err := uuid.Parse(bookingID); err != nil {
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

	type fetchResult struct {
		body   []byte
		status int
		err    error
	}

	var wg sync.WaitGroup
	var venueRes fetchResult
	var paymentRes fetchResult

	wg.Add(2)
	go func() {
		defer wg.Done()
		venueRes.body, venueRes.status, venueRes.err = a.fetchJSON(c, fmt.Sprintf("%s/venues/%d", a.venueURL, bookingLookup.VenueID))
	}()
	go func() {
		defer wg.Done()
		paymentRes.body, paymentRes.status, paymentRes.err = a.fetchJSON(c, fmt.Sprintf("%s/bookings/%s/payment", a.paymentURL, bookingID))
	}()
	wg.Wait()

	resp := models.BookingSummaryResponse{
		Booking: bookingBody,
	}

	if venueRes.err != nil {
		msg := "venue service unavailable"
		resp.VenueError = &msg
	} else if venueRes.status != http.StatusOK {
		msg := fmt.Sprintf("venue service status %d", venueRes.status)
		resp.VenueError = &msg
	} else {
		resp.Venue = venueRes.body
	}

	if paymentRes.err != nil {
		msg := "payment service unavailable"
		resp.PaymentError = &msg
	} else if paymentRes.status != http.StatusOK {
		msg := fmt.Sprintf("payment service status %d", paymentRes.status)
		resp.PaymentError = &msg
	} else {
		resp.Payment = paymentRes.body
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