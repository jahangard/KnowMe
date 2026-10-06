package profile

import "time"

type User struct {
	ID             int64
	TelegramUserID int64
	Username       *string
	FirstSeenAt    time.Time
	LastSeenAt     time.Time
}

type UserProfile struct {
	UserID                 int64
	Name                   *string
	Age                    *int
	Gender                 *string
	Mobile                 *string
	ProfileCompletionLevel int
	UpdatedAt              time.Time
}

type Trait struct {
	UserID     int64
	TraitKey   string
	Score      float64
	Confidence float64
	UpdatedAt  time.Time
}
