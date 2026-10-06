package eventtracking

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Log(ctx context.Context, userID int64, eventType string, testID, questionID *int64, metadata any) error {
	var metadataJSON any
	if metadata != nil {
		raw, err := json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("marshal event metadata: %w", err)
		}
		metadataJSON = string(raw)
	}

	_, err := s.db.ExecContext(ctx,
		"INSERT INTO dbo.UserEvents (UserId, EventType, TestId, QuestionId, MetadataJson, CreatedAt) "+
			"VALUES (@UserId, @EventType, @TestId, @QuestionId, @MetadataJson, SYSUTCDATETIME())",
		sql.Named("UserId", userID),
		sql.Named("EventType", eventType),
		sql.Named("TestId", nullableInt64(testID)),
		sql.Named("QuestionId", nullableInt64(questionID)),
		sql.Named("MetadataJson", metadataJSON),
	)
	if err != nil {
		return fmt.Errorf("insert user event: %w", err)
	}
	return nil
}

func nullableInt64(value *int64) any {
	if value == nil {
		return nil
	}
	return *value
}
