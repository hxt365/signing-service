package reposistory

import (
	"database/sql"
	"fmt"
	"strings"
)

type SignatureRepo struct {
	db *sql.DB
}

func NewSignatureRepo(db *sql.DB) *SignatureRepo {
	return &SignatureRepo{db: db}
}

type Signature struct {
	ID       int
	KeyID    int
	RecordID int
	Value    string
}

func (sr *SignatureRepo) BulkInsert(signatures []*Signature) error {
	valueStrings := make([]string, 0, len(signatures))
	valueArgs := make([]interface{}, 0, len(signatures)*3)
	for _, s := range signatures {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, s.RecordID)
		valueArgs = append(valueArgs, s.KeyID)
		valueArgs = append(valueArgs, s.Value)
	}

	stmt := fmt.Sprintf("INSERT INTO signatures (record_id, key_id, value) VALUES %s", strings.Join(valueStrings, ","))
	if _, err := sr.db.Exec(stmt, valueArgs...); err != nil {
		return err
	}
	return nil
}

