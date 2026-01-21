package seed

import (
	"fmt"
	"log/slog"
	"time"
	"venue-service/internal/models"

	"gorm.io/gorm"
)

// SeedVenues заполняет базу данных тестовыми площадками
func SeedVenues(db *gorm.DB, logger *slog.Logger) error {
	logger = logger.With("layer", "seed")

	// Проверяем, есть ли уже данные
	var count int64
	if err := db.Model(&models.Venue{}).Count(&count).Error; err != nil {
		return fmt.Errorf("ошибка проверки данных: %w", err)
	}

	if count > 0 {
		logger.Info("База данных уже содержит данные, пропускаем сиды", "count", count)
		return nil
	}

	now := time.Now()

	venues := []models.Venue{
		// Футбольные площадки
		{
			VenueType: models.VenueFootball,
			OwnerID:   1,
			IsActive:  true,
			HourPrice: 1500,
			District:  "Центральный",
			StartTime: time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.UTC),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.UTC),
			Weekdays: models.Weekdays{
				Monday:    true,
				Tuesday:   true,
				Wednesday: true,
				Thursday:  true,
				Friday:    true,
				Saturday:  true,
				Sunday:    true,
			},
		},
		{
			VenueType: models.VenueFootball,
			OwnerID:   1,
			IsActive:  true,
			HourPrice: 2000,
			District:  "Северный",
			StartTime: time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.UTC),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day(), 21, 0, 0, 0, time.UTC),
			Weekdays: models.Weekdays{
				Monday:    true,
				Tuesday:   true,
				Wednesday: true,
				Thursday:  true,
				Friday:    true,
				Saturday:  true,
				Sunday:    false, // Воскресенье выходной
			},
		},
		{
			VenueType: models.VenueFootball,
			OwnerID:   2,
			IsActive:  true,
			HourPrice: 1200,
			District:  "Южный",
			StartTime: time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, time.UTC),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, time.UTC),
			Weekdays: models.Weekdays{
				Monday:    true,
				Tuesday:   true,
				Wednesday: true,
				Thursday:  true,
				Friday:    true,
				Saturday:  true,
				Sunday:    true,
			},
		},
		// Баскетбольные площадки
		{
			VenueType: models.VenueBasketball,
			OwnerID:   2,
			IsActive:  true,
			HourPrice: 1800,
			District:  "Центральный",
			StartTime: time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, time.UTC),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.UTC),
			Weekdays: models.Weekdays{
				Monday:    true,
				Tuesday:   true,
				Wednesday: true,
				Thursday:  true,
				Friday:    true,
				Saturday:  true,
				Sunday:    true,
			},
		},
		{
			VenueType: models.VenueBasketball,
			OwnerID:   3,
			IsActive:  true,
			HourPrice: 2200,
			District:  "Западный",
			StartTime: time.Date(now.Year(), now.Month(), now.Day(), 8, 30, 0, 0, time.UTC),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day(), 20, 30, 0, 0, time.UTC),
			Weekdays: models.Weekdays{
				Monday:    true,
				Tuesday:   true,
				Wednesday: true,
				Thursday:  true,
				Friday:    true,
				Saturday:  false, // Суббота выходной
				Sunday:    false, // Воскресенье выходной
			},
		},
		// Теннисные корты
		{
			VenueType: models.VenueTennis,
			OwnerID:   3,
			IsActive:  true,
			HourPrice: 2500,
			District:  "Восточный",
			StartTime: time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.UTC),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day(), 21, 0, 0, 0, time.UTC),
			Weekdays: models.Weekdays{
				Monday:    true,
				Tuesday:   true,
				Wednesday: true,
				Thursday:  true,
				Friday:    true,
				Saturday:  true,
				Sunday:    true,
			},
		},
		{
			VenueType: models.VenueTennis,
			OwnerID:   1,
			IsActive:  true,
			HourPrice: 3000,
			District:  "Центральный",
			StartTime: time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.UTC),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, time.UTC),
			Weekdays: models.Weekdays{
				Monday:    true,
				Tuesday:   true,
				Wednesday: true,
				Thursday:  true,
				Friday:    true,
				Saturday:  true,
				Sunday:    true,
			},
		},
		// Тренажерные залы
		{
			VenueType: models.VenueGym,
			OwnerID:   4,
			IsActive:  true,
			HourPrice: 500,
			District:  "Северный",
			StartTime: time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, time.UTC),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC),
			Weekdays: models.Weekdays{
				Monday:    true,
				Tuesday:   true,
				Wednesday: true,
				Thursday:  true,
				Friday:    true,
				Saturday:  true,
				Sunday:    true,
			},
		},
		{
			VenueType: models.VenueGym,
			OwnerID:   4,
			IsActive:  true,
			HourPrice: 800,
			District:  "Центральный",
			StartTime: time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, time.UTC),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, time.UTC),
			Weekdays: models.Weekdays{
				Monday:    true,
				Tuesday:   true,
				Wednesday: true,
				Thursday:  true,
				Friday:    true,
				Saturday:  true,
				Sunday:    false, // Воскресенье выходной
			},
		},
		// Бассейны
		{
			VenueType: models.VenueSwimming,
			OwnerID:   5,
			IsActive:  true,
			HourPrice: 1000,
			District:  "Южный",
			StartTime: time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, time.UTC),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, time.UTC),
			Weekdays: models.Weekdays{
				Monday:    true,
				Tuesday:   true,
				Wednesday: true,
				Thursday:  true,
				Friday:    true,
				Saturday:  true,
				Sunday:    true,
			},
		},
		{
			VenueType: models.VenueSwimming,
			OwnerID:   5,
			IsActive:  false, // Деактивированная площадка для тестирования
			HourPrice: 1500,
			District:  "Восточный",
			StartTime: time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.UTC),
			EndTime:   time.Date(now.Year(), now.Month(), now.Day(), 19, 0, 0, 0, time.UTC),
			Weekdays: models.Weekdays{
				Monday:    true,
				Tuesday:   true,
				Wednesday: true,
				Thursday:  true,
				Friday:    true,
				Saturday:  false,
				Sunday:    false,
			},
		},
	}

	// Создаем площадки пакетами для лучшей производительности
	batchSize := 5
	for i := 0; i < len(venues); i += batchSize {
		end := i + batchSize
		if end > len(venues) {
			end = len(venues)
		}

		batch := venues[i:end]
		if err := db.Create(&batch).Error; err != nil {
			logger.Error("Ошибка создания сидов", "batch", i, "error", err)
			return fmt.Errorf("ошибка создания сидов (batch %d): %w", i, err)
		}

		logger.Info("Создан batch площадок", "batch", i, "count", len(batch))
	}

	logger.Info("Сиды успешно созданы", "total", len(venues))
	return nil
}

// SeedVenuesForce принудительно перезаписывает все данные (удаляет существующие)
func SeedVenuesForce(db *gorm.DB, logger *slog.Logger) error {
	logger = logger.With("layer", "seed")

	// Удаляем все существующие данные
	if err := db.Exec("DELETE FROM venues").Error; err != nil {
		return fmt.Errorf("ошибка удаления данных: %w", err)
	}

	logger.Info("Существующие данные удалены")
	return SeedVenues(db, logger)
}
