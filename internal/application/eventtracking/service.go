package eventtracking

import (
	"context"
	"encoding/json"
	"fmt"

	store "github.com/jahangard/KnowMe/internal/platform/database"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Log(ctx context.Context, userID int64, eventType string, testID, questionID *int64, metadata any) error {
	var metadataJSON *string
	if metadata != nil {
		raw, err := json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("marshal event metadata: %w", err)
		}
		value := string(raw)
		metadataJSON = &value
	}

	event := store.UserEvent{
		UserID:       userID,
		EventType:    eventType,
		TestID:       testID,
		QuestionID:   questionID,
		MetadataJSON: metadataJSON,
	}
	if err := s.db.WithContext(ctx).Create(&event).Error; err != nil {
		return fmt.Errorf("insert user event: %w", err)
	}
	return nil
}
