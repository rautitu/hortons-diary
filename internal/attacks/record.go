package attacks

import (
	"time"

	"github.com/google/uuid"
)

type PainRecord struct {
	ID              uuid.UUID
	StartTime       time.Time
	Severity        int
	DurationMinutes int
}
