package testcatalog

import (
	"context"
	"errors"
	"fmt"

	store "github.com/jahangard/KnowMe/internal/platform/database"
	"gorm.io/gorm"
)

type Category struct {
	ID    int64
	Code  string
	Title string
}

type Test struct {
	ID          int64
	Title       string
	Description string
}

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Categories(ctx context.Context) ([]Category, error) {
	var rows []store.TestCategory
	if err := s.db.WithContext(ctx).
		Where("IsActive = ?", true).
		Order("SortOrder ASC").
		Order("Id ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("load test categories: %w", err)
	}

	items := make([]Category, 0, len(rows))
	for _, row := range rows {
		items = append(items, Category{
			ID:    row.ID,
			Code:  row.Code,
			Title: row.Title,
		})
	}
	return items, nil
}

func (s *Service) Category(ctx context.Context, categoryID int64) (*Category, error) {
	var row store.TestCategory
	if err := s.db.WithContext(ctx).
		Where("Id = ? AND IsActive = ?", categoryID, true).
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("load test category: %w", err)
	}

	return &Category{ID: row.ID, Code: row.Code, Title: row.Title}, nil
}

func (s *Service) Tests(ctx context.Context, categoryID int64) ([]Test, error) {
	var rows []store.Test
	if err := s.db.WithContext(ctx).
		Where("CategoryId = ? AND IsActive = ?", categoryID, true).
		Order("SortOrder ASC").
		Order("Id ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("load category tests: %w", err)
	}

	items := make([]Test, 0, len(rows))
	for _, row := range rows {
		description := ""
		if row.Description != nil {
			description = *row.Description
		}
		items = append(items, Test{
			ID:          row.ID,
			Title:       row.Title,
			Description: description,
		})
	}
	return items, nil
}
