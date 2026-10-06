package testengine

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
)

type Option struct {
	ID   int64
	Text string
}

type QuestionView struct {
	SessionID int64
	TestID    int64
	TestTitle string
	QuestionID int64
	Text      string
	Current   int
	Total     int
	Options   []Option
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
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Start(ctx context.Context, userID, testID int64) (*QuestionView, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin start test: %w", err)
	}
	defer tx.Rollback()

	var sessionID int64
	err = tx.QueryRowContext(ctx,
		"SELECT TOP (1) Id FROM dbo.TestSessions WHERE UserId=@UserId AND TestId=@TestId AND [Status]='active' ORDER BY Id DESC",
		sql.Named("UserId", userID),
		sql.Named("TestId", testID),
	).Scan(&sessionID)

	if err == sql.ErrNoRows {
		err = tx.QueryRowContext(ctx,
			"INSERT INTO dbo.TestSessions (UserId, TestId, [Status], StartedAt) OUTPUT INSERTED.Id VALUES (@UserId, @TestId, 'active', SYSUTCDATETIME())",
			sql.Named("UserId", userID),
			sql.Named("TestId", testID),
		).Scan(&sessionID)
	}
	if err != nil {
		return nil, fmt.Errorf("get or create test session: %w", err)
	}

	view, err := s.nextQuestionTx(ctx, tx, sessionID, testID)
	if err != nil {
		return nil, err
	}
	if view == nil {
		return nil, fmt.Errorf("test has no active questions")
	}

	if _, err := tx.ExecContext(ctx,
		"UPDATE dbo.TestSessions SET CurrentQuestionId=@QuestionId WHERE Id=@SessionId",
		sql.Named("QuestionId", view.QuestionID),
		sql.Named("SessionId", sessionID),
	); err != nil {
		return nil, fmt.Errorf("set current question: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit start test: %w", err)
	}
	return view, nil
}

func (s *Service) Answer(ctx context.Context, userID, sessionID, questionID, optionID int64) (*QuestionView, *Result, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("begin answer: %w", err)
	}
	defer tx.Rollback()

	var testID int64
	var status string
	err = tx.QueryRowContext(ctx,
		"SELECT TestId, [Status] FROM dbo.TestSessions WHERE Id=@SessionId AND UserId=@UserId",
		sql.Named("SessionId", sessionID),
		sql.Named("UserId", userID),
	).Scan(&testID, &status)
	if err != nil {
		return nil, nil, fmt.Errorf("load test session: %w", err)
	}
	if status != "active" {
		return nil, nil, fmt.Errorf("test session is not active")
	}

	var valid int
	err = tx.QueryRowContext(ctx,
		"SELECT COUNT(1) FROM dbo.QuestionOptions qo INNER JOIN dbo.Questions q ON q.Id=qo.QuestionId WHERE qo.Id=@OptionId AND q.Id=@QuestionId AND q.TestId=@TestId",
		sql.Named("OptionId", optionID),
		sql.Named("QuestionId", questionID),
		sql.Named("TestId", testID),
	).Scan(&valid)
	if err != nil {
		return nil, nil, fmt.Errorf("validate option: %w", err)
	}
	if valid != 1 {
		return nil, nil, fmt.Errorf("invalid question option")
	}

	var already int
	err = tx.QueryRowContext(ctx,
		"SELECT COUNT(1) FROM dbo.TestAnswers WHERE SessionId=@SessionId AND QuestionId=@QuestionId",
		sql.Named("SessionId", sessionID),
		sql.Named("QuestionId", questionID),
	).Scan(&already)
	if err != nil {
		return nil, nil, fmt.Errorf("check previous answer: %w", err)
	}

	if already == 0 {
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO dbo.TestAnswers (SessionId, QuestionId, OptionId, AnsweredAt) VALUES (@SessionId, @QuestionId, @OptionId, SYSUTCDATETIME())",
			sql.Named("SessionId", sessionID),
			sql.Named("QuestionId", questionID),
			sql.Named("OptionId", optionID),
		); err != nil {
			return nil, nil, fmt.Errorf("save answer: %w", err)
		}
	}

	next, err := s.nextQuestionTx(ctx, tx, sessionID, testID)
	if err != nil {
		return nil, nil, err
	}

	if next != nil {
		if _, err := tx.ExecContext(ctx,
			"UPDATE dbo.TestSessions SET CurrentQuestionId=@QuestionId WHERE Id=@SessionId",
			sql.Named("QuestionId", next.QuestionID),
			sql.Named("SessionId", sessionID),
		); err != nil {
			return nil, nil, fmt.Errorf("advance session: %w", err)
		}
		if err := tx.Commit(); err != nil {
			return nil, nil, fmt.Errorf("commit answer: %w", err)
		}
		return next, nil, nil
	}

	result, err := s.finishTx(ctx, tx, userID, sessionID, testID)
	if err != nil {
		return nil, nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("commit finish test: %w", err)
	}
	return nil, result, nil
}

func (s *Service) nextQuestionTx(ctx context.Context, tx *sql.Tx, sessionID, testID int64) (*QuestionView, error) {
	const query = "SELECT TOP (1) q.Id, q.[Text], q.[Order], t.Title, " +
		"(SELECT COUNT(1) FROM dbo.Questions q2 WHERE q2.TestId=t.Id AND q2.IsActive=1) AS Total " +
		"FROM dbo.Questions q INNER JOIN dbo.Tests t ON t.Id=q.TestId " +
		"WHERE q.TestId=@TestId AND q.IsActive=1 AND NOT EXISTS (" +
		"SELECT 1 FROM dbo.TestAnswers a WHERE a.SessionId=@SessionId AND a.QuestionId=q.Id) " +
		"ORDER BY q.[Order], q.Id"

	var view QuestionView
	view.SessionID = sessionID
	view.TestID = testID
	err := tx.QueryRowContext(ctx, query,
		sql.Named("TestId", testID),
		sql.Named("SessionId", sessionID),
	).Scan(&view.QuestionID, &view.Text, &view.Current, &view.TestTitle, &view.Total)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load next question: %w", err)
	}

	rows, err := tx.QueryContext(ctx,
		"SELECT Id, [Text] FROM dbo.QuestionOptions WHERE QuestionId=@QuestionId ORDER BY [Order], Id",
		sql.Named("QuestionId", view.QuestionID),
	)
	if err != nil {
		return nil, fmt.Errorf("load question options: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var option Option
		if err := rows.Scan(&option.ID, &option.Text); err != nil {
			return nil, fmt.Errorf("scan question option: %w", err)
		}
		view.Options = append(view.Options, option)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate question options: %w", err)
	}

	return &view, nil
}

func (s *Service) finishTx(ctx context.Context, tx *sql.Tx, userID, sessionID, testID int64) (*Result, error) {
	rows, err := tx.QueryContext(ctx,
		"SELECT qo.TraitKey, SUM(COALESCE(qo.Score,0)) Score " +
			"FROM dbo.TestAnswers a INNER JOIN dbo.QuestionOptions qo ON qo.Id=a.OptionId " +
			"WHERE a.SessionId=@SessionId AND qo.TraitKey IS NOT NULL " +
			"GROUP BY qo.TraitKey",
		sql.Named("SessionId", sessionID),
	)
	if err != nil {
		return nil, fmt.Errorf("score test: %w", err)
	}
	defer rows.Close()

	scores := map[string]float64{}
	for rows.Next() {
		var key string
		var score float64
		if err := rows.Scan(&key, &score); err != nil {
			return nil, fmt.Errorf("scan test score: %w", err)
		}
		scores[key] = score
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate scores: %w", err)
	}
	if len(scores) == 0 {
		return nil, fmt.Errorf("test produced no scores")
	}

	type pair struct {
		key   string
		score float64
	}
	ranked := make([]pair, 0, len(scores))
	for key, score := range scores {
		ranked = append(ranked, pair{key: key, score: score})
	}
	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].score == ranked[j].score {
			return ranked[i].key < ranked[j].key
		}
		return ranked[i].score > ranked[j].score
	})
	winner := ranked[0]

	var title string
	if err := tx.QueryRowContext(ctx,
		"SELECT Title FROM dbo.Tests WHERE Id=@TestId",
		sql.Named("TestId", testID),
	).Scan(&title); err != nil {
		return nil, fmt.Errorf("load test title: %w", err)
	}

	scoreJSON, err := json.Marshal(scores)
	if err != nil {
		return nil, fmt.Errorf("marshal test scores: %w", err)
	}

	if _, err := tx.ExecContext(ctx,
		"UPDATE dbo.TestSessions SET [Status]='completed', CompletedAt=SYSUTCDATETIME(), CurrentQuestionId=NULL WHERE Id=@SessionId",
		sql.Named("SessionId", sessionID),
	); err != nil {
		return nil, fmt.Errorf("complete session: %w", err)
	}

	if _, err := tx.ExecContext(ctx,
		"INSERT INTO dbo.TestResults (SessionId, ResultType, ScoreJson, Summary, CreatedAt) VALUES (@SessionId, @ResultType, @ScoreJson, NULL, SYSUTCDATETIME())",
		sql.Named("SessionId", sessionID),
		sql.Named("ResultType", winner.key),
		sql.Named("ScoreJson", string(scoreJSON)),
	); err != nil {
		return nil, fmt.Errorf("save test result: %w", err)
	}

	for key, score := range scores {
		if _, err := tx.ExecContext(ctx,
			"MERGE dbo.UserTraits AS target USING (SELECT @UserId UserId, @TraitKey TraitKey) AS src " +
				"ON target.UserId=src.UserId AND target.TraitKey=src.TraitKey " +
				"WHEN MATCHED THEN UPDATE SET Score=@Score, Confidence=1, UpdatedAt=SYSUTCDATETIME() " +
				"WHEN NOT MATCHED THEN INSERT (UserId,TraitKey,Score,Confidence,UpdatedAt) VALUES (@UserId,@TraitKey,@Score,1,SYSUTCDATETIME());",
			sql.Named("UserId", userID),
			sql.Named("TraitKey", key),
			sql.Named("Score", score),
		); err != nil {
			return nil, fmt.Errorf("update user trait: %w", err)
		}
	}

	return &Result{
		SessionID: sessionID,
		TestID:    testID,
		TestTitle: title,
		TraitKey:  winner.key,
		Score:     winner.score,
		Scores:    scores,
	}, nil
}
