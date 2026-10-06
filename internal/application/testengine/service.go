package testengine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	store "github.com/jahangard/KnowMe/internal/platform/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Option struct {
	ID   int64
	Text string
}

type QuestionView struct {
	SessionID  int64
	TestID     int64
	TestTitle  string
	QuestionID int64
	Text       string
	Current    int
	Total      int
	Options    []Option
}

type Result struct {
	SessionID int64
	TestID    int64
	TestTitle string
	TraitKey  string
	Score     float64
	Scores    map[string]float64
}

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Start(ctx context.Context, userID, testID int64) (*QuestionView, error) {
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("begin start test: %w", tx.Error)
	}
	defer tx.Rollback()

	var session store.TestSession
	err := tx.
		Where("UserId = ? AND TestId = ? AND Status = ?", userID, testID, "active").
		Order("Id DESC").
		First(&session).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		session = store.TestSession{
			UserID: userID,
			TestID: testID,
			Status: "active",
		}
		if err := tx.Create(&session).Error; err != nil {
			return nil, fmt.Errorf("create test session: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("load active test session: %w", err)
	}

	view, err := s.nextQuestionTx(tx, session.ID, testID)
	if err != nil {
		return nil, err
	}
	if view == nil {
		return nil, fmt.Errorf("test has no active questions")
	}

	if err := tx.Model(&store.TestSession{}).
		Where("Id = ?", session.ID).
		Update("CurrentQuestionId", view.QuestionID).Error; err != nil {
		return nil, fmt.Errorf("set current question: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit start test: %w", err)
	}
	return view, nil
}

func (s *Service) Answer(ctx context.Context, userID, sessionID, questionID, optionID int64) (*QuestionView, *Result, error) {
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, nil, fmt.Errorf("begin answer: %w", tx.Error)
	}
	defer tx.Rollback()

	var session store.TestSession
	if err := tx.
		Where("Id = ? AND UserId = ?", sessionID, userID).
		First(&session).Error; err != nil {
		return nil, nil, fmt.Errorf("load test session: %w", err)
	}
	if session.Status != "active" {
		return nil, nil, fmt.Errorf("test session is not active")
	}

	var option store.QuestionOption
	if err := tx.
		Where("Id = ? AND QuestionId = ?", optionID, questionID).
		First(&option).Error; err != nil {
		return nil, nil, fmt.Errorf("validate question option: %w", err)
	}

	var question store.Question
	if err := tx.
		Where("Id = ? AND TestId = ?", questionID, session.TestID).
		First(&question).Error; err != nil {
		return nil, nil, fmt.Errorf("validate question: %w", err)
	}

	var existing int64
	if err := tx.Model(&store.TestAnswer{}).
		Where("SessionId = ? AND QuestionId = ?", sessionID, questionID).
		Count(&existing).Error; err != nil {
		return nil, nil, fmt.Errorf("check previous answer: %w", err)
	}

	if existing == 0 {
		answer := store.TestAnswer{
			SessionID:  sessionID,
			QuestionID: questionID,
			OptionID:   &optionID,
		}
		if err := tx.Create(&answer).Error; err != nil {
			return nil, nil, fmt.Errorf("save answer: %w", err)
		}
	}

	next, err := s.nextQuestionTx(tx, sessionID, session.TestID)
	if err != nil {
		return nil, nil, err
	}

	if next != nil {
		if err := tx.Model(&store.TestSession{}).
			Where("Id = ?", sessionID).
			Update("CurrentQuestionId", next.QuestionID).Error; err != nil {
			return nil, nil, fmt.Errorf("advance session: %w", err)
		}
		if err := tx.Commit().Error; err != nil {
			return nil, nil, fmt.Errorf("commit answer: %w", err)
		}
		return next, nil, nil
	}

	result, err := s.finishTx(tx, userID, sessionID, session.TestID)
	if err != nil {
		return nil, nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, nil, fmt.Errorf("commit finish test: %w", err)
	}
	return nil, result, nil
}

func (s *Service) nextQuestionTx(tx *gorm.DB, sessionID, testID int64) (*QuestionView, error) {
	var answeredQuestionIDs []int64
	if err := tx.Model(&store.TestAnswer{}).
		Where("SessionId = ?", sessionID).
		Pluck("QuestionId", &answeredQuestionIDs).Error; err != nil {
		return nil, fmt.Errorf("load answered questions: %w", err)
	}

	query := tx.Model(&store.Question{}).
		Where("TestId = ? AND IsActive = ?", testID, true)
	if len(answeredQuestionIDs) > 0 {
		query = query.Not("Id IN ?", answeredQuestionIDs)
	}

	var question store.Question
	if err := query.
		Order(clause.OrderByColumn{Column: clause.Column{Name: "Order"}}).
		Order("Id ASC").
		First(&question).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("load next question: %w", err)
	}

	var test store.Test
	if err := tx.Where("Id = ?", testID).First(&test).Error; err != nil {
		return nil, fmt.Errorf("load test title: %w", err)
	}

	var total int64
	if err := tx.Model(&store.Question{}).
		Where("TestId = ? AND IsActive = ?", testID, true).
		Count(&total).Error; err != nil {
		return nil, fmt.Errorf("count test questions: %w", err)
	}

	var rows []store.QuestionOption
	if err := tx.
		Where("QuestionId = ?", question.ID).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "Order"}}).
		Order("Id ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("load question options: %w", err)
	}

	options := make([]Option, 0, len(rows))
	for _, row := range rows {
		options = append(options, Option{ID: row.ID, Text: row.Text})
	}

	return &QuestionView{
		SessionID:  sessionID,
		TestID:     testID,
		TestTitle:  test.Title,
		QuestionID: question.ID,
		Text:       question.Text,
		Current:    question.Order,
		Total:      int(total),
		Options:    options,
	}, nil
}

func (s *Service) finishTx(tx *gorm.DB, userID, sessionID, testID int64) (*Result, error) {
	type scoreRow struct {
		TraitKey string
		Score    float64
	}

	var scored []scoreRow
	if err := tx.
		Table("TestAnswers AS a").
		Select("qo.TraitKey AS TraitKey, SUM(COALESCE(qo.Score, 0)) AS Score").
		Joins("JOIN QuestionOptions AS qo ON qo.Id = a.OptionId").
		Where("a.SessionId = ? AND qo.TraitKey IS NOT NULL", sessionID).
		Group("qo.TraitKey").
		Scan(&scored).Error; err != nil {
		return nil, fmt.Errorf("score test: %w", err)
	}
	if len(scored) == 0 {
		return nil, fmt.Errorf("test produced no scores")
	}

	scores := make(map[string]float64, len(scored))
	for _, row := range scored {
		scores[row.TraitKey] = row.Score
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].Score == scored[j].Score {
			return scored[i].TraitKey < scored[j].TraitKey
		}
		return scored[i].Score > scored[j].Score
	})
	winner := scored[0]

	var test store.Test
	if err := tx.Where("Id = ?", testID).First(&test).Error; err != nil {
		return nil, fmt.Errorf("load test title: %w", err)
	}

	rawScores, err := json.Marshal(scores)
	if err != nil {
		return nil, fmt.Errorf("marshal test scores: %w", err)
	}
	scoreJSON := string(rawScores)
	now := time.Now().UTC()

	if err := tx.Model(&store.TestSession{}).
		Where("Id = ?", sessionID).
		Updates(map[string]any{
			"Status":            "completed",
			"CompletedAt":       now,
			"CurrentQuestionId": nil,
		}).Error; err != nil {
		return nil, fmt.Errorf("complete session: %w", err)
	}

	testResult := store.TestResult{
		SessionID:  sessionID,
		ResultType: winner.TraitKey,
		ScoreJSON:  &scoreJSON,
	}
	if err := tx.Create(&testResult).Error; err != nil {
		return nil, fmt.Errorf("save test result: %w", err)
	}

	for key, score := range scores {
		trait := store.UserTrait{
			UserID:     userID,
			TraitKey:   key,
			Score:      score,
			Confidence: 1,
			UpdatedAt:  now,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "UserId"},
				{Name: "TraitKey"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"Score":      score,
				"Confidence": 1,
				"UpdatedAt":  now,
			}),
		}).Create(&trait).Error; err != nil {
			return nil, fmt.Errorf("update user trait: %w", err)
		}
	}

	return &Result{
		SessionID: sessionID,
		TestID:    testID,
		TestTitle: test.Title,
		TraitKey:  winner.TraitKey,
		Score:     winner.Score,
		Scores:    scores,
	}, nil
}
