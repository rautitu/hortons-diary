package attacks

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PainRecord struct {
	ID              uuid.UUID
	StartTime       time.Time
	Severity        int
	DurationMinutes int
}

func NewPainRecord(
	startTime time.Time,
	severity int,
	durationMinutes int,
) (PainRecord, error) {
	if severity < 1 || severity > 10 {
		return PainRecord{}, fmt.Errorf("severity must be between 1 and 10, got %d", severity)
	}
	if durationMinutes <= 0 {
		return PainRecord{}, fmt.Errorf("duration in minutes must be positive, got %d", durationMinutes)
	}
	if startTime.IsZero() {
		return PainRecord{}, fmt.Errorf("start time is required")
	}
	return PainRecord{
		ID:              uuid.New(),
		StartTime:       startTime,
		Severity:        severity,
		DurationMinutes: durationMinutes,
	}, nil
}
