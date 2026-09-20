package attacks

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/tonnomolt/hortons-diary/internal/auth"
)

type Repository interface {
	Create(
		ctx context.Context,
		userID auth.UserID,
		record PainRecord,
	) (PainRecord, error)

	FindRecord(
		ctx context.Context,
		userID auth.UserID,
		recordID uuid.UUID,
	) (*PainRecord, error)

	Update(
		ctx context.Context,
		userID auth.UserID,
		record PainRecord,
	) error

	Delete(
		ctx context.Context,
		userID auth.UserID,
		recordID uuid.UUID,
	) error
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

func (s *Service) Update(
	ctx context.Context,
	userID auth.UserID,
	recordID uuid.UUID,
	input CreateInput,
) error {
	r, errFind := s.repo.FindRecord(ctx, userID, recordID)
	if errFind != nil {
		return errFind
	}
	errEdit := r.EditPainRecord(
		input.StartTime,
		input.Severity,
		input.DurationMinutes,
	)
	if errEdit != nil {
		return errEdit
	}
	return s.repo.Update(ctx, userID, *r)
}

func (s *Service) Delete(
	ctx context.Context,
	userID auth.UserID,
	recordID uuid.UUID,
) error {
	return s.repo.Delete(ctx, userID, recordID)
}
