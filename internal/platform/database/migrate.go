package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

func MigrateAndSeed(ctx context.Context, handle *Handle) error {
	db := handle.Gorm.WithContext(ctx)

	if err := db.AutoMigrate(
		&User{},
		&UserProfile{},
		&TestCategory{},
		&Test{},
		&Question{},
		&QuestionOption{},
		&TestResultProfile{},
		&TestSession{},
		&TestAnswer{},
		&TestResult{},
		&UserTrait{},
		&UserEvent{},
	); err != nil {
		return fmt.Errorf("auto migrate database: %w", err)
	}

	if err := seedCategories(ctx, handle.Gorm); err != nil {
		return fmt.Errorf("seed test categories: %w", err)
	}

	if err := seedInitialTests(ctx, handle.Gorm); err != nil {
		return fmt.Errorf("seed initial tests: %w", err)
	}

	return nil
}

func seedCategories(ctx context.Context, db *gorm.DB) error {
	categories := []TestCategory{
		{Code: "love", Title: "عشق و رابطه", SortOrder: 10, IsActive: true},
		{Code: "challenge", Title: "چالشی و باحال", SortOrder: 20, IsActive: true},
		{Code: "self_knowledge", Title: "خودشناسی عمیق‌تر", SortOrder: 30, IsActive: true},
		{Code: "adult", Title: "۱۸+", SortOrder: 40, IsActive: true},
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, category := range categories {
			var row TestCategory
			if err := tx.Where("Code = ?", category.Code).
				FirstOrCreate(&row, category).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

