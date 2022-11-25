package reposistory

import (
	"context"
	"database/sql"

	"coordinator/constant"
)

type ProgressRepo struct {
	db *sql.DB
}

func NewProgressRepo(db *sql.DB) *ProgressRepo {
	return &ProgressRepo{db: db}
}

type Progress struct {
	ID        int
	KeyID     int
	Committed int
}

func (pr *ProgressRepo) GetCommittedOffset(keyID int) (*Progress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constant.DBTimeout)
	defer cancel()

	var p Progress
	q := "SELECT id, key_id, committed FROM progresses WHERE key_id = ?"
	if err := pr.db.QueryRowContext(ctx, q, keyID).Scan(&p.ID, &p.KeyID, &p.Committed); err != nil {
		return nil, err
	}

	return &p, nil
}

func (pr *ProgressRepo) SetCommittedOffset(keyID int, offset int) error {
	ctx, cancel := context.WithTimeout(context.Background(), constant.DBTimeout)
	defer cancel()

	q := "INSERT INTO progresses (key_id, committed) VALUES(?, ?) ON DUPLICATE KEY UPDATE committed = ?"
	stmt, err := pr.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, keyID, offset, offset)

	return err
}
