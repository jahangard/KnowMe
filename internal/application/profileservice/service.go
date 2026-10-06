package profileservice

import (
	"context"
	"fmt"

	store "github.com/jahangard/KnowMe/internal/platform/database"
	"gorm.io/gorm"
)

type Profile struct {
	UserID                 int64
	Name                   *string
	Age                    *int
	Gender                 *string
	Mobile                 *string
	ProfileCompletionLevel int
	PendingField            *string
	NamePrompted            bool
	MobilePrompted          bool
}

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Ensure(ctx context.Context, userID int64) error {
	row := store.UserProfile{UserID: userID}
	if err := s.db.WithContext(ctx).
		Where("UserId = ?", userID).
		FirstOrCreate(&row).Error; err != nil {
		return fmt.Errorf("ensure user profile: %w", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, userID int64) (*Profile, error) {
	if err := s.Ensure(ctx, userID); err != nil {
		return nil, err
	}

	var row store.UserProfile
	if err := s.db.WithContext(ctx).
		Where("UserId = ?", userID).
		First(&row).Error; err != nil {
		return nil, fmt.Errorf("load user profile: %w", err)
	}

	return &Profile{
		UserID:                 row.UserID,
		Name:                   row.Name,
		Age:                    row.Age,
		Gender:                 row.Gender,
		Mobile:                 row.Mobile,
		ProfileCompletionLevel: row.ProfileCompletionLevel,
		PendingField:            row.PendingField,
		NamePrompted:            row.NamePrompted,
		MobilePrompted:          row.MobilePrompted,
	}, nil
}

func (s *Service) SetGender(ctx context.Context, userID int64, gender string) error {
	return s.update(ctx, userID, map[string]any{
		"Gender":                 gender,
		"PendingField":           "age",
		"ProfileCompletionLevel": 1,
	}, "set gender")
}

func (s *Service) BeginAgeCapture(ctx context.Context, userID int64) error {
	return s.update(ctx, userID, map[string]any{
		"PendingField": "age",
	}, "begin age capture")
}

func (s *Service) SetAge(ctx context.Context, userID int64, age int) error {
	return s.update(ctx, userID, map[string]any{
		"Age":                    age,
		"PendingField":           nil,
		"ProfileCompletionLevel": 2,
	}, "set age")
}

func (s *Service) BeginNameCapture(ctx context.Context, userID int64) error {
	return s.update(ctx, userID, map[string]any{
		"PendingField": "name",
		"NamePrompted": true,
	}, "begin name capture")
}

func (s *Service) SetName(ctx context.Context, userID int64, name string) error {
	return s.update(ctx, userID, map[string]any{
		"Name":                   name,
		"PendingField":           nil,
		"NamePrompted":           true,
		"ProfileCompletionLevel": gorm.Expr(
			"CASE WHEN ProfileCompletionLevel < ? THEN ? ELSE ProfileCompletionLevel END",
			3, 3,
		),
	}, "set name")
}

func (s *Service) BeginMobileCapture(ctx context.Context, userID int64) error {
	return s.update(ctx, userID, map[string]any{
		"PendingField":   "mobile",
		"MobilePrompted": true,
	}, "begin mobile capture")
}

func (s *Service) SetMobile(ctx context.Context, userID int64, mobile string) error {
	return s.update(ctx, userID, map[string]any{
		"Mobile":                 mobile,
		"PendingField":           nil,
		"MobilePrompted":         true,
		"ProfileCompletionLevel": gorm.Expr(
			"CASE WHEN ProfileCompletionLevel < ? THEN ? ELSE ProfileCompletionLevel END",
			4, 4,
		),
	}, "set mobile")
}

func (s *Service) SkipMobile(ctx context.Context, userID int64) error {
	return s.update(ctx, userID, map[string]any{
		"PendingField":   nil,
		"MobilePrompted": true,
	}, "skip mobile")
}

func (s *Service) ClearPending(ctx context.Context, userID int64) error {
	return s.update(ctx, userID, map[string]any{
		"PendingField": nil,
	}, "clear pending profile field")
}

func (s *Service) CompletedTestCount(ctx context.Context, userID int64) (int, error) {
	var count int64
	if err := s.db.WithContext(ctx).
		Model(&store.TestSession{}).
		Where("UserId = ? AND Status = ?", userID, "completed").
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count completed tests: %w", err)
	}
	return int(count), nil
}

func (s *Service) update(ctx context.Context, userID int64, values map[string]any, action string) error {
	result := s.db.WithContext(ctx).
		Model(&store.UserProfile{}).
		Where("UserId = ?", userID).
		Updates(values)
	if result.Error != nil {
		return fmt.Errorf("%s: %w", action, result.Error)
	}
	if result.RowsAffected == 0 {
		if err := s.Ensure(ctx, userID); err != nil {
			return err
		}
		if err := s.db.WithContext(ctx).
			Model(&store.UserProfile{}).
			Where("UserId = ?", userID).
			Updates(values).Error; err != nil {
			return fmt.Errorf("%s after ensure: %w", action, err)
		}
	}
	return nil
}
