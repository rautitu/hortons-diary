package attacks

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tonnomolt/hortons-diary/internal/auth"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(
	ctx context.Context,
	userID auth.UserID,
	record PainRecord,
) (PainRecord, error) {
	const query = `
        INSERT INTO pain_record (
            user_id, pain_record_id, start_time,
            duration_minutes, severity
        ) VALUES ($1, $2, $3, $4, $5)
    `

	_, err := r.db.Exec(ctx, query,
		uuid.UUID(userID),
		record.ID,
		record.StartTime,
		record.DurationMinutes,
		record.Severity,
	)
	if err != nil {
		return PainRecord{}, fmt.Errorf("create pain record: %w", err)
	}
	return record, nil
}

func (r *PostgresRepository) FindRecord(
	ctx context.Context,
	userID auth.UserID,
	recordID uuid.UUID,
) (*PainRecord, error) {
	const query = `
        SELECT pain_record_id, start_time,
               duration_minutes, severity
        FROM pain_record
        WHERE user_id = $1 AND pain_record_id = $2
    `

	var record PainRecord
	err := r.db.QueryRow(ctx, query,
		uuid.UUID(userID),
		recordID,
	).Scan(
		&record.ID,
		&record.StartTime,
		&record.DurationMinutes,
		&record.Severity,
	)
	if err != nil {
		return nil, fmt.Errorf("find pain record: %w", err)
	}

	return &record, nil
}
