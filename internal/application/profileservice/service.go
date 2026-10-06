package profileservice

import (
	"context"
	"database/sql"
	"fmt"
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
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Ensure(ctx context.Context, userID int64) error {
	_, err := s.db.ExecContext(ctx,
		"IF NOT EXISTS (SELECT 1 FROM dbo.UserProfiles WHERE UserId=@UserId) "+
			"INSERT INTO dbo.UserProfiles (UserId, ProfileCompletionLevel, UpdatedAt) VALUES (@UserId, 0, SYSUTCDATETIME())",
		sql.Named("UserId", userID),
	)
	if err != nil {
		return fmt.Errorf("ensure user profile: %w", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, userID int64) (*Profile, error) {
	if err := s.Ensure(ctx, userID); err != nil {
		return nil, err
	}

	var p Profile
	err := s.db.QueryRowContext(ctx,
		"SELECT UserId, Name, Age, Gender, Mobile, ProfileCompletionLevel, PendingField, NamePrompted, MobilePrompted "+
			"FROM dbo.UserProfiles WHERE UserId=@UserId",
		sql.Named("UserId", userID),
	).Scan(
		&p.UserID,
		&p.Name,
		&p.Age,
		&p.Gender,
		&p.Mobile,
		&p.ProfileCompletionLevel,
		&p.PendingField,
		&p.NamePrompted,
		&p.MobilePrompted,
	)
	if err != nil {
		return nil, fmt.Errorf("load user profile: %w", err)
	}
	return &p, nil
}

func (s *Service) SetGender(ctx context.Context, userID int64, gender string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE dbo.UserProfiles SET Gender=@Gender, PendingField='age', ProfileCompletionLevel=1, UpdatedAt=SYSUTCDATETIME() WHERE UserId=@UserId",
		sql.Named("Gender", gender),
		sql.Named("UserId", userID),
	)
	if err != nil {
		return fmt.Errorf("set gender: %w", err)
	}
	return nil
}

func (s *Service) BeginAgeCapture(ctx context.Context, userID int64) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE dbo.UserProfiles SET PendingField='age', UpdatedAt=SYSUTCDATETIME() WHERE UserId=@UserId",
		sql.Named("UserId", userID),
	)
	if err != nil {
		return fmt.Errorf("begin age capture: %w", err)
	}
	return nil
}

func (s *Service) SetAge(ctx context.Context, userID int64, age int) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE dbo.UserProfiles SET Age=@Age, PendingField=NULL, ProfileCompletionLevel=2, UpdatedAt=SYSUTCDATETIME() WHERE UserId=@UserId",
		sql.Named("Age", age),
		sql.Named("UserId", userID),
	)
	if err != nil {
		return fmt.Errorf("set age: %w", err)
	}
	return nil
}

func (s *Service) BeginNameCapture(ctx context.Context, userID int64) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE dbo.UserProfiles SET PendingField='name', NamePrompted=1, UpdatedAt=SYSUTCDATETIME() WHERE UserId=@UserId",
		sql.Named("UserId", userID),
	)
	if err != nil {
		return fmt.Errorf("begin name capture: %w", err)
	}
	return nil
}

func (s *Service) SetName(ctx context.Context, userID int64, name string) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE dbo.UserProfiles SET Name=@Name, PendingField=NULL, NamePrompted=1, "+
			"ProfileCompletionLevel=CASE WHEN ProfileCompletionLevel < 3 THEN 3 ELSE ProfileCompletionLevel END, "+
			"UpdatedAt=SYSUTCDATETIME() WHERE UserId=@UserId",
		sql.Named("Name", name),
		sql.Named("UserId", userID),
	)
	if err != nil {
		return fmt.Errorf("set name: %w", err)
	}
	return nil
}

func (s *Service) ClearPending(ctx context.Context, userID int64) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE dbo.UserProfiles SET PendingField=NULL, UpdatedAt=SYSUTCDATETIME() WHERE UserId=@UserId",
		sql.Named("UserId", userID),
	)
	if err != nil {
		return fmt.Errorf("clear pending profile field: %w", err)
	}
	return nil
}

func (s *Service) CompletedTestCount(ctx context.Context, userID int64) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(1) FROM dbo.TestSessions WHERE UserId=@UserId AND [Status]='completed'",
		sql.Named("UserId", userID),
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count completed tests: %w", err)
	}
	return count, nil
}
