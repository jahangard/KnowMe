package roadmap

import (
	"context"
	"database/sql"
	"fmt"
)

type Recommendation struct {
	TestID int64
	Title  string
	Reason string
}

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) NextTest(ctx context.Context, userID int64) (*Recommendation, error) {
	const query = "SELECT TOP (1) t.Id, t.Title " +
		"FROM dbo.Tests t " +
		"WHERE t.IsActive = 1 " +
		"AND NOT EXISTS (" +
		"SELECT 1 FROM dbo.TestSessions ts " +
		"WHERE ts.UserId = @UserId AND ts.TestId = t.Id AND ts.Status = 'completed') " +
		"ORDER BY t.SortOrder, t.Id;"

	var rec Recommendation
	err := s.db.QueryRowContext(ctx, query, sql.Named("UserId", userID)).Scan(&rec.TestID, &rec.Title)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("select next roadmap test: %w", err)
	}

	rec.Reason = "next_not_completed"
	return &rec, nil
}
