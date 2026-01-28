# Gateway API Examples

Base URL: http://localhost:8085

Use Authorization header for all non-public endpoints:
`Authorization: Bearer <token>`

## Auth

### Register (POST /api/auth/register)
```bash
curl -X POST http://localhost:8085/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Ivan Petrov",
    "email": "ivan@example.com",
    "password": "password123"
  }'
```

### Login (POST /api/auth/login)
```bash
curl -X POST http://localhost:8085/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ivan@example.com",
    "password": "password123"
  }'
```

## Venues (public GET)

### List venues (GET /api/venues)
```bash
curl "http://localhost:8085/api/venues?venue_type=football&district=central&limit=10&page=1"
```

### Venue details (GET /api/venues/:id)
```bash
curl http://localhost:8085/api/venues/1
```

### Venue types (GET /api/venue-types)
```bash
curl http://localhost:8085/api/venue-types
```

## Venues (auth required)

### Create venue (POST /api/venues)
```bash
curl -X POST http://localhost:8085/api/venues \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "venue_type": "football",
    "owner_id": 1,
    "is_active": true,
    "hour_price": 1500,
    "district": "central",
    "start_time": "08:00",
    "end_time": "22:00",
    "weekdays": {
      "monday": true,
      "tuesday": true,
      "wednesday": true,
      "thursday": true,
      "friday": true,
      "saturday": true,
      "sunday": true
    }
  }'
```

### Venue schedule (GET /api/venues/:id/schedule)
```bash
curl http://localhost:8085/api/venues/1/schedule \
  -H "Authorization: Bearer <token>"
```

### Update schedule (PUT /api/venues/:id/schedule)
```bash
curl -X PUT http://localhost:8085/api/venues/1/schedule \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "start_time": "09:00",
    "end_time": "23:00"
  }'
```

## Bookings (auth required)

### Create booking (POST /api/bookings)
```bash
curl -X POST http://localhost:8085/api/bookings \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "venue_id": 1,
    "owner_id": 1,
    "start_at": "2026-01-22T10:00:00Z",
    "end_at": "2026-01-22T12:00:00Z",
    "price_cents": 50000,
    "status": "pending"
  }'
```

### My bookings (GET /api/bookings)
```bash
curl http://localhost:8085/api/bookings \
  -H "Authorization: Bearer <token>"
```

### Booking details (GET /api/bookings/:id)
```bash
curl http://localhost:8085/api/bookings/1 \
  -H "Authorization: Bearer <token>"
```

### Cancel booking (POST /api/bookings/:id/cancel)
```bash
curl -X POST http://localhost:8085/api/bookings/1/cancel \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Changed plans"
  }'
```

### Availability (GET /api/venues/:id/availability)
```bash
curl "http://localhost:8085/api/venues/1/availability?date=2026-01-22" \
  -H "Authorization: Bearer <token>"
```

### Venue bookings (GET /api/venues/:id/bookings)
```bash
curl http://localhost:8085/api/venues/1/bookings \
  -H "Authorization: Bearer <token>"
```

### Aggregation summary (GET /api/bookings/:id/summary)
```bash
curl http://localhost:8085/api/bookings/550e8400-e29b-41d4-a716-446655440000/summary \
  -H "Authorization: Bearer <token>"
```

## Payments (auth required)

### Create payment (POST /api/payments)
```bash
curl -X POST http://localhost:8085/api/payments \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "booking_id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "amount": 50000,
    "currency": "RUB",
    "method": "card"
  }'
```

### Payment status (GET /api/payments/:id)
```bash
curl http://localhost:8085/api/payments/1 \
  -H "Authorization: Bearer <token>"
```

### Payments history (GET /api/payments)
```bash
curl http://localhost:8085/api/payments \
  -H "Authorization: Bearer <token>"
```

### Refund (POST /api/payments/:id/refund)
```bash
curl -X POST http://localhost:8085/api/payments/1/refund \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 25000,
    "reason": "Cancel booking"
  }'
```

### Payment by booking (GET /api/bookings/:id/payment)
```bash
curl http://localhost:8085/api/bookings/550e8400-e29b-41d4-a716-446655440000/payment \
  -H "Authorization: Bearer <token>"
```
