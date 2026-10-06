package roadmap

import (
	"context"
	"errors"
	"fmt"

	store "github.com/jahangard/KnowMe/internal/platform/database"
	"gorm.io/gorm"
)

type Recommendation struct {
	TestID int64
	Title  string
	Reason string
}

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) NextTest(ctx context.Context, userID int64) (*Recommendation, error) {
	var completedTestIDs []int64
	if err := s.db.WithContext(ctx).
		Model(&store.TestSession{}).
		Where("UserId = ? AND Status = ?", userID, "completed").
		Distinct().
		Pluck("TestId", &completedTestIDs).Error; err != nil {
		return nil, fmt.Errorf("load completed roadmap tests: %w", err)
	}

	query := s.db.WithContext(ctx).
		Model(&store.Test{}).
		Where("IsActive = ?", true).
		Order("SortOrder ASC").
		Order("Id ASC")
	if len(completedTestIDs) > 0 {
		query = query.Not("Id IN ?", completedTestIDs)
	}

	var test store.Test
	if err := query.First(&test).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("select next roadmap test: %w", err)
	}

	return &Recommendation{
		TestID: test.ID,
		Title:  test.Title,
		Reason: "next_not_completed",
	}, nil
}
