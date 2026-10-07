package database

import "time"

type User struct {
	ID             int64     `gorm:"column:Id;primaryKey;autoIncrement"`
	TelegramUserID int64     `gorm:"column:TelegramUserId;uniqueIndex;not null"`
	Username       *string   `gorm:"column:Username;size:255"`
	FirstSeenAt    time.Time `gorm:"column:FirstSeenAt;autoCreateTime"`
	LastSeenAt     time.Time `gorm:"column:LastSeenAt"`
}

func (User) TableName() string { return "Users" }

type UserProfile struct {
	UserID                 int64     `gorm:"column:UserId;primaryKey"`
	Name                   *string   `gorm:"column:Name;size:200"`
	Age                    *int      `gorm:"column:Age"`
	Gender                 *string   `gorm:"column:Gender;size:50"`
	Mobile                 *string   `gorm:"column:Mobile;size:50"`
	ProfileCompletionLevel int       `gorm:"column:ProfileCompletionLevel;not null;default:0"`
	PendingField           *string   `gorm:"column:PendingField;size:50"`
	NamePrompted           bool      `gorm:"column:NamePrompted;not null;default:false"`
	MobilePrompted         bool      `gorm:"column:MobilePrompted;not null;default:false"`
	UpdatedAt              time.Time `gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (UserProfile) TableName() string { return "UserProfiles" }

type TestCategory struct {
	ID        int64  `gorm:"column:Id;primaryKey;autoIncrement"`
	Code      string `gorm:"column:Code;size:100;uniqueIndex;not null"`
	Title     string `gorm:"column:Title;size:200;not null"`
	SortOrder int    `gorm:"column:SortOrder;not null;default:0"`
	IsActive  bool   `gorm:"column:IsActive;not null;default:true"`
}

func (TestCategory) TableName() string { return "TestCategories" }

type Test struct {
	ID          int64     `gorm:"column:Id;primaryKey;autoIncrement"`
	CategoryID  int64     `gorm:"column:CategoryId;not null;index"`
	Code        string    `gorm:"column:Code;size:100;uniqueIndex;not null"`
	Title       string    `gorm:"column:Title;size:250;not null"`
	Description *string   `gorm:"column:Description;size:1000"`
	SortOrder   int       `gorm:"column:SortOrder;not null;default:0"`
	IsActive    bool      `gorm:"column:IsActive;not null;default:true"`
	CreatedAt   time.Time `gorm:"column:CreatedAt;autoCreateTime"`
}

func (Test) TableName() string { return "Tests" }

type Question struct {
	ID           int64  `gorm:"column:Id;primaryKey;autoIncrement"`
	TestID       int64  `gorm:"column:TestId;not null;index:UX_Questions_Test_Order,unique"`
	Text         string `gorm:"column:Text;size:1000;not null"`
	Order        int    `gorm:"column:Order;not null;index:UX_Questions_Test_Order,unique"`
	QuestionType string `gorm:"column:QuestionType;size:50;not null"`
	IsActive     bool   `gorm:"column:IsActive;not null;default:true"`
}

func (Question) TableName() string { return "Questions" }

type QuestionOption struct {
	ID         int64   `gorm:"column:Id;primaryKey;autoIncrement"`
	QuestionID int64   `gorm:"column:QuestionId;not null;index"`
	Text       string  `gorm:"column:Text;size:500;not null"`
	Order      int     `gorm:"column:Order;not null"`
	Score      float64 `gorm:"column:Score"`
	TraitKey   *string `gorm:"column:TraitKey;size:100"`
}

func (QuestionOption) TableName() string { return "QuestionOptions" }

type TestResultProfile struct {
	ID          int64  `gorm:"column:Id;primaryKey;autoIncrement"`
	TestID      int64  `gorm:"column:TestId;not null;index:UX_TestResultProfiles_Test_Trait,unique"`
	TraitKey    string `gorm:"column:TraitKey;size:100;not null;index:UX_TestResultProfiles_Test_Trait,unique"`
	Label       string `gorm:"column:Label;size:200;not null"`
	Title       string `gorm:"column:Title;size:250;not null"`
	Subtitle    string `gorm:"column:Subtitle;size:500"`
	Description string `gorm:"column:Description;size:2000"`
	SortOrder   int    `gorm:"column:SortOrder;not null;default:0"`
	IsActive    bool   `gorm:"column:IsActive;not null;default:true"`
}

func (TestResultProfile) TableName() string { return "TestResultProfiles" }

type TestSession struct {
	ID                int64      `gorm:"column:Id;primaryKey;autoIncrement"`
	UserID            int64      `gorm:"column:UserId;not null;index:IX_TestSessions_User_Status"`
	TestID            int64      `gorm:"column:TestId;not null;index"`
	CurrentQuestionID *int64     `gorm:"column:CurrentQuestionId"`
	Status            string     `gorm:"column:Status;size:30;not null;index:IX_TestSessions_User_Status"`
	StartedAt         time.Time  `gorm:"column:StartedAt;autoCreateTime"`
	CompletedAt       *time.Time `gorm:"column:CompletedAt"`
}

func (TestSession) TableName() string { return "TestSessions" }

type TestAnswer struct {
	ID             int64     `gorm:"column:Id;primaryKey;autoIncrement"`
	SessionID      int64     `gorm:"column:SessionId;not null;index:UX_TestAnswers_Session_Question,unique"`
	QuestionID     int64     `gorm:"column:QuestionId;not null;index:UX_TestAnswers_Session_Question,unique"`
	OptionID       *int64    `gorm:"column:OptionId"`
	AnswerValue    *string   `gorm:"column:AnswerValue;size:1000"`
	ResponseTimeMS *int64    `gorm:"column:ResponseTimeMs"`
	AnsweredAt     time.Time `gorm:"column:AnsweredAt;autoCreateTime"`
}

func (TestAnswer) TableName() string { return "TestAnswers" }

type TestResult struct {
	ID            int64     `gorm:"column:Id;primaryKey;autoIncrement"`
	SessionID     int64     `gorm:"column:SessionId;not null;uniqueIndex"`
	ResultType    string    `gorm:"column:ResultType;size:100;not null"`
	ScoreJSON     *string   `gorm:"column:ScoreJson"`
	Summary       *string   `gorm:"column:Summary"`
	ShareImageURL *string   `gorm:"column:ShareImageUrl;size:1000"`
	CreatedAt     time.Time `gorm:"column:CreatedAt;autoCreateTime"`
}

func (TestResult) TableName() string { return "TestResults" }

type UserTrait struct {
	UserID     int64     `gorm:"column:UserId;primaryKey"`
	TraitKey   string    `gorm:"column:TraitKey;size:100;primaryKey"`
	Score      float64   `gorm:"column:Score;not null"`
	Confidence float64   `gorm:"column:Confidence;not null"`
	UpdatedAt  time.Time `gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (UserTrait) TableName() string { return "UserTraits" }

type UserEvent struct {
	ID           int64     `gorm:"column:Id;primaryKey;autoIncrement"`
	UserID       int64     `gorm:"column:UserId;not null;index:IX_UserEvents_User_CreatedAt"`
	EventType    string    `gorm:"column:EventType;size:100;not null"`
	TestID       *int64    `gorm:"column:TestId"`
	QuestionID   *int64    `gorm:"column:QuestionId"`
	MetadataJSON *string   `gorm:"column:MetadataJson"`
	CreatedAt    time.Time `gorm:"column:CreatedAt;autoCreateTime;index:IX_UserEvents_User_CreatedAt,sort:desc"`
}

func (UserEvent) TableName() string { return "UserEvents" }
