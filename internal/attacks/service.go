package attacks

import (
	"context"
	"time"

	"github.com/tonnomolt/hortons-diary/internal/auth"
)

type Repository interface {
	Create(
		ctx context.Context,
		userID auth.UserID,
		record PainRecord,
	) (PainRecord, error)
}

type CreateInput struct {
	StartTime       time.Time
	DurationMinutes int
	Severity        int
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(
	ctx context.Context,
	userID auth.UserID,
	input CreateInput,
) (PainRecord, error) {
	record, err := NewPainRecord(
		input.StartTime,
		input.Severity,
		input.DurationMinutes,
	)
	if err != nil {
		return PainRecord{}, err
	}

	return s.repo.Create(ctx, userID, record)
}
