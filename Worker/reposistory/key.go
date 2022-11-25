package reposistory

import (
	"context"
	"database/sql"
)

type KeyRepo struct {
	db *sql.DB
}

func NewKeyRepo(db *sql.DB) *KeyRepo {
	return &KeyRepo{db: db}
}

type KeyWithLock struct {
	ID     int
	KeyID  int
	KeyVal string
	tx     *sql.Tx
}

func (k *KeyWithLock) Done() {
	k.tx.Commit()
}

func (k *KeyWithLock) Abort() {
	k.tx.Rollback()
}

func (k *KeyRepo) LockKey(ctx context.Context) (*KeyWithLock, error) {
	tx, err := k.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	kl := &KeyWithLock{tx: tx}
	q := "SELECT id, identifier, value FROM secret_keys ORDER BY id LIMIT  1 FOR UPDATE SKIP LOCKED"
	if err = tx.QueryRowContext(ctx, q).Scan(&kl.ID, &kl.KeyID, &kl.KeyVal); err != nil {
		tx.Rollback()
		return kl, err
	}

	q = "DELETE FROM secret_keys WHERE id = ?"
	if _, err = tx.ExecContext(ctx, q, kl.ID); err != nil {
		tx.Rollback()
		return kl, err
	}

	q = "INSERT INTO secret_keys (identifier, value) VALUES (?, ?)"
	if _, err = tx.ExecContext(ctx, q, kl.KeyID, kl.KeyVal); err != nil {
		tx.Rollback()
		return kl, err
	}

	return kl, nil
}
