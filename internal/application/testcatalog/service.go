package testcatalog

import (
	"context"
	"errors"
	"fmt"

	store "github.com/jahangard/KnowMe/internal/platform/database"
	"gorm.io/gorm"
)

type Category struct {
	ID          int64
	Code        string
	Title       string
	ParentID    *int64
	CatalogType string
}

type Test struct {
	ID          int64
	CategoryID  int64
	Title       string
	Description string
	IsReady     bool
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
		Where("IsActive = ? AND CatalogType = ? AND ParentId IS NULL", true, "tests").
		Order("SortOrder ASC").
		Order("Id ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("load test categories: %w", err)
	}

	items := make([]Category, 0, len(rows))
	for _, row := range rows {
		items = append(items, Category{
			ID:          row.ID,
			Code:        row.Code,
			Title:       row.Title,
			ParentID:    row.ParentID,
			CatalogType: row.CatalogType,
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

	return &Category{ID: row.ID, Code: row.Code, Title: row.Title, ParentID: row.ParentID, CatalogType: row.CatalogType}, nil
}

func (s *Service) Children(ctx context.Context, parentID int64, catalogType string) ([]Category, error) {
	var rows []store.TestCategory
	if err := s.db.WithContext(ctx).
		Where("ParentId = ? AND CatalogType = ? AND IsActive = ?", parentID, catalogType, true).
		Order("SortOrder ASC").Order("Id ASC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("load child topics: %w", err)
	}
	items := make([]Category, 0, len(rows))
	for _, row := range rows {
		items = append(items, Category{ID: row.ID, Code: row.Code, Title: row.Title, ParentID: row.ParentID, CatalogType: row.CatalogType})
	}
	return items, nil
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
			CategoryID:  row.CategoryID,
			Title:       row.Title,
			Description: description,
			IsReady:     row.IsReady,
		})
	}
	return items, nil
}

func (s *Service) Test(ctx context.Context, testID int64) (*Test, error) {
	var row store.Test
	if err := s.db.WithContext(ctx).
		Where("Id = ? AND IsActive = ?", testID, true).
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("load test: %w", err)
	}

	description := ""
	if row.Description != nil {
		description = *row.Description
	}
	return &Test{ID: row.ID, CategoryID: row.CategoryID, Title: row.Title, Description: description, IsReady: row.IsReady}, nil
}
