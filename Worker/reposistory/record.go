package reposistory

import (
	"context"
	"database/sql"
)

type RecordRepo struct {
	db *sql.DB
}

func NewRecordRepo(db *sql.DB) *RecordRepo {
	return &RecordRepo{db: db}
}

type Record struct {
	ID    int
	Value string
}

func (rr *RecordRepo) GetRecordsByRange(ctx context.Context, start, end int) ([]*Record, error) {
	q := "SELECT id, value FROM records WHERE id BETWEEN ? AND ?"
	rows, err := rr.db.QueryContext(ctx, q, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*Record

	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.ID, &r.Value); err != nil {
			return records, err
		}
		records = append(records, &r)
	}

	if err = rows.Err(); err != nil {
		return records, err
	}
	return records, nil
}
