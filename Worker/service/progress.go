package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"Worker/constant"
)

var ProgressServiceURL = fmt.Sprintf("%s/progresses", constant.CoordinatorServiceURL)

type ProgressService struct {
	client *http.Client
}

func NewProgressService(c *http.Client) *ProgressService {
	return &ProgressService{client: c}
}

type GetNextBatchRangeResponse struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func (ps *ProgressService) GetNextBatchRange(ctx context.Context, keyID int) (int, int, error) {
	req, err := http.NewRequest(http.MethodGet, ProgressServiceURL, nil)
	if err != nil {
		return 0, 0, err
	}
	req = req.WithContext(ctx)

	q := req.URL.Query()
	q.Add("key_id", strconv.Itoa(keyID))
	req.URL.RawQuery = q.Encode()

	resp, err := ps.client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf(resp.Status)
	}

	var batchRange GetNextBatchRangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&batchRange); err != nil {
		return 0, 0, err
	}

	return batchRange.Start, batchRange.End, nil
}

func (ps *ProgressService) UpdateCommittedOffset(ctx context.Context, keyID, offset int) error {
	body := map[string]int{
		"key_id": keyID,
		"offset": offset,
	}
	jsonValue, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", ProgressServiceURL, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	resp, err := ps.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(resp.Status)
	}

	return nil
}
