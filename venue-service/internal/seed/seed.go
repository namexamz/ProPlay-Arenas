package seed

import (
	"fmt"
	"log/slog"
	"time"
	"venue-service/internal/models"

	"gorm.io/gorm"
)

// createDaySchedule создает DaySchedule для одного дня
func createDaySchedule(enabled bool, startHour, startMin, endHour, endMin int) models.DaySchedule {
	schedule := models.DaySchedule{
		Enabled: enabled,
	}
	if enabled {
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day(), startHour, startMin, 0, 0, time.UTC)
		endTime := time.Date(now.Year(), now.Month(), now.Day(), endHour, endMin, 0, 0, time.UTC)
		schedule.StartTime = &startTime
		schedule.EndTime = &endTime
	}
	return schedule
}

// createWeekdays создает Weekdays с одинаковым расписанием для всех дней
func createWeekdays(startHour, startMin, endHour, endMin int, enabledDays ...bool) models.Weekdays {
	// По умолчанию все дни включены
	mondayEnabled := true
	tuesdayEnabled := true
	wednesdayEnabled := true
	thursdayEnabled := true
	fridayEnabled := true
	saturdayEnabled := true
	sundayEnabled := true

	// Если переданы значения, используем их
	if len(enabledDays) > 0 {
		mondayEnabled = enabledDays[0]
	}
	if len(enabledDays) > 1 {
		tuesdayEnabled = enabledDays[1]
	}
	if len(enabledDays) > 2 {
		wednesdayEnabled = enabledDays[2]
	}
	if len(enabledDays) > 3 {
		thursdayEnabled = enabledDays[3]
	}
	if len(enabledDays) > 4 {
		fridayEnabled = enabledDays[4]
	}
	if len(enabledDays) > 5 {
		saturdayEnabled = enabledDays[5]
	}
	if len(enabledDays) > 6 {
		sundayEnabled = enabledDays[6]
	}

	return models.Weekdays{
		Monday:    createDaySchedule(mondayEnabled, startHour, startMin, endHour, endMin),
		Tuesday:   createDaySchedule(tuesdayEnabled, startHour, startMin, endHour, endMin),
		Wednesday: createDaySchedule(wednesdayEnabled, startHour, startMin, endHour, endMin),
		Thursday:  createDaySchedule(thursdayEnabled, startHour, startMin, endHour, endMin),
		Friday:    createDaySchedule(fridayEnabled, startHour, startMin, endHour, endMin),
		Saturday:  createDaySchedule(saturdayEnabled, startHour, startMin, endHour, endMin),
		Sunday:    createDaySchedule(sundayEnabled, startHour, startMin, endHour, endMin),
	}
}

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

	venues := []models.Venue{
		// Футбольные площадки
		{
			VenueType: models.VenueFootball,
			OwnerID:   1,
			IsActive:  true,
			HourPrice: 1500,
			District:  "Центральный",
			Weekdays:  createWeekdays(8, 0, 22, 0),
		},
		{
			VenueType: models.VenueFootball,
			OwnerID:   1,
			IsActive:  true,
			HourPrice: 2000,
			District:  "Северный",
			Weekdays:  createWeekdays(9, 0, 21, 0, true, true, true, true, true, true, false), // Воскресенье выходной
		},
		{
			VenueType: models.VenueFootball,
			OwnerID:   2,
			IsActive:  true,
			HourPrice: 1200,
			District:  "Южный",
			Weekdays:  createWeekdays(7, 0, 23, 0),
		},
		// Баскетбольные площадки
		{
			VenueType: models.VenueBasketball,
			OwnerID:   2,
			IsActive:  true,
			HourPrice: 1800,
			District:  "Центральный",
			Weekdays:  createWeekdays(10, 0, 22, 0),
		},
		{
			VenueType: models.VenueBasketball,
			OwnerID:   3,
			IsActive:  true,
			HourPrice: 2200,
			District:  "Западный",
			Weekdays:  createWeekdays(8, 30, 20, 30, true, true, true, true, true, false, false), // Суббота и воскресенье выходные
		},
		// Теннисные корты
		{
			VenueType: models.VenueTennis,
			OwnerID:   3,
			IsActive:  true,
			HourPrice: 2500,
			District:  "Восточный",
			Weekdays:  createWeekdays(9, 0, 21, 0),
		},
		{
			VenueType: models.VenueTennis,
			OwnerID:   1,
			IsActive:  true,
			HourPrice: 3000,
			District:  "Центральный",
			Weekdays:  createWeekdays(8, 0, 22, 0),
		},
		// Тренажерные залы
		{
			VenueType: models.VenueGym,
			OwnerID:   4,
			IsActive:  true,
			HourPrice: 500,
			District:  "Северный",
			Weekdays:  createWeekdays(6, 0, 23, 59),
		},
		{
			VenueType: models.VenueGym,
			OwnerID:   4,
			IsActive:  true,
			HourPrice: 800,
			District:  "Центральный",
			Weekdays:  createWeekdays(7, 0, 23, 0, true, true, true, true, true, true, false), // Воскресенье выходной
		},
		// Бассейны
		{
			VenueType: models.VenueSwimming,
			OwnerID:   5,
			IsActive:  true,
			HourPrice: 1000,
			District:  "Южный",
			Weekdays:  createWeekdays(8, 0, 20, 0),
		},
		{
			VenueType: models.VenueSwimming,
			OwnerID:   5,
			IsActive:  false, // Деактивированная площадка для тестирования
			HourPrice: 1500,
			District:  "Восточный",
			Weekdays:  createWeekdays(9, 0, 19, 0, true, true, true, true, true, false, false), // Суббота и воскресенье выходные
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

	// Удаляем все существующие данные включая soft-deleted записи
	// Unscoped() позволяет удалить записи, помеченные как удаленные
	if err := db.Unscoped().Delete(&models.Venue{}, "1 = 1").Error; err != nil {
		return fmt.Errorf("ошибка удаления данных: %w", err)
	}

	logger.Info("Существующие данные удалены")
	return SeedVenues(db, logger)
}
