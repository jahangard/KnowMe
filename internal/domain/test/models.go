package test

import "time"

type QuestionType string

const (
	QuestionTypeSingleChoice QuestionType = "single_choice"
	QuestionTypeScale        QuestionType = "scale"
	QuestionTypeText         QuestionType = "text"
)

type Test struct {
	ID          int64
	CategoryID  int64
	Code        string
	Title       string
	Description string
	IsActive    bool
}

type Question struct {
	ID       int64
	TestID   int64
	Text     string
	Order    int
	Type     QuestionType
	IsActive bool
}

type Option struct {
	ID         int64
	QuestionID int64
	Text       string
	Score      float64
	TraitKey   string
}

type Session struct {
	ID                int64
	UserID            int64
	TestID            int64
	CurrentQuestionID *int64
	Status            string
	StartedAt         time.Time
	CompletedAt       *time.Time
}

type Answer struct {
	ID             int64
	SessionID      int64
	QuestionID     int64
	OptionID       *int64
	AnswerValue    *string
	ResponseTimeMS *int64
	AnsweredAt     time.Time
}
