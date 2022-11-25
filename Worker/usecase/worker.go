package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"time"

	"Worker/constant"
	"Worker/reposistory"
	"Worker/service"
)

type Worker struct {
	kr *reposistory.KeyRepo
	rr *reposistory.RecordRepo
	sr *reposistory.SignatureRepo
	ps *service.ProgressService
}

func NewWorker(kr *reposistory.KeyRepo, rr *reposistory.RecordRepo, sr *reposistory.SignatureRepo, ps *service.ProgressService) *Worker {
	return &Worker{
		kr: kr,
		rr: rr,
		sr: sr,
		ps: ps,
	}
}

func (w *Worker) Start() {
	for {
		w.Sign()
		time.Sleep(time.Second)
	}
}

func (w *Worker) Sign() {
	ctx, cancel := context.WithTimeout(context.Background(), constant.ProcessTimeout)
	defer cancel()

	key, err := w.kr.LockKey(ctx)
	defer key.Done()
	if err != nil {
		log.Printf("could not acquire a key: %s\n", err)
		return
	}

	start, end, err := w.ps.GetNextBatchRange(ctx, key.KeyID)
	if err != nil {
		log.Printf("could not get next batch range for key %d: %s\n", key.KeyID, err)
		return
	}
	if start > end {
		log.Printf("key %d done its job\n", key.KeyID)
		return
	}

	records, err := w.rr.GetRecordsByRange(ctx, start, end)
	if err != nil {
		log.Printf("could not get records by range %d - %d: %s\n", start, end, err)
		return
	}

	signedRecords := signRecords(records, key.KeyVal)
	if err := w.sr.BulkInsert(signedRecordsToSignatures(signedRecords, key.KeyID)); err != nil {
		if !reposistory.IsDuplicateErr(err) {
			log.Printf("could not bulk insert signatures: %s\n", err)
			return
		}
		log.Printf("bulk insert duplicated for batch %d-%d: %s", start, end, err)
	}

	if err := w.ps.UpdateCommittedOffset(ctx, key.KeyID, end); err != nil {
		log.Printf("could not update committed offset: %s\n", err)
		return
	}

	log.Printf("Done signing batch %d-%d with key %d", start, end, key.KeyID)
}

func signRecords(records []*reposistory.Record, key string) []*reposistory.Record {
	results := make([]*reposistory.Record, 0, len(records))

	for _, r := range records {
		results = append(results, &reposistory.Record{
			ID:    r.ID,
			Value: computeHmac256(r.Value, key),
		})
	}

	return results
}

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func signedRecordsToSignatures(records []*reposistory.Record, keyID int) []*reposistory.Signature {
	results := make([]*reposistory.Signature, 0, len(records))
	for _, r := range records {
		results = append(results, &reposistory.Signature{
			KeyID:    keyID,
			RecordID: r.ID,
			Value:    r.Value,
		})
	}
	return results
}
