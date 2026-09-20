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

func validatePainRecordInput(
	startTime time.Time,
	severity int,
	durationMinutes int,
) error {
	if severity < 1 || severity > 10 {
		return fmt.Errorf("severity must be between 1 and 10, got %d", severity)
	}
	if durationMinutes <= 0 {
		return fmt.Errorf("duration in minutes must be positive, got %d", durationMinutes)
	}
	if startTime.IsZero() {
		return fmt.Errorf("start time is required")
	}
	return nil
}

func NewPainRecord(
	startTime time.Time,
	severity int,
	durationMinutes int,
) (PainRecord, error) {
	if err := validatePainRecordInput(startTime, severity, durationMinutes); err != nil {
		return PainRecord{}, err
	}
	return PainRecord{
		ID:              uuid.New(),
		StartTime:       startTime,
		Severity:        severity,
		DurationMinutes: durationMinutes,
	}, nil
}

func (r *PainRecord) EditPainRecord(
	startTime time.Time,
	severity int,
	durationMinutes int,
) error {
	if err := validatePainRecordInput(
		startTime, severity, durationMinutes,
	); err != nil {
		return err
	}

	r.StartTime = startTime
	r.Severity = severity
	r.DurationMinutes = durationMinutes
	return nil
}
