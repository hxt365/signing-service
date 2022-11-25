package usecase

import (
	"coordinator/reposistory"
	"coordinator/utils"
	"database/sql"
)

type ProgressUseCase struct {
	pr         *reposistory.ProgressRepo
	batchSize  int
	numRecords int
	numKeys    int
}

func NewProgressUseCase(pr *reposistory.ProgressRepo, batchSize, numRecords, numKeys int) *ProgressUseCase {
	return &ProgressUseCase{
		pr:         pr,
		batchSize:  batchSize,
		numRecords: numRecords,
		numKeys:    numKeys,
	}
}

func (pu *ProgressUseCase) GetNextRange(keyID int) (int, int, error) {
	if keyID < 1 || keyID > pu.numKeys {
		return 0, 0, utils.ValidationErr{Err: "invalid key ID"}
	}

	var (
		rangeStart     = (keyID-1)*pu.numRecords/pu.numKeys + 1
		rangeEnd       = keyID * pu.numRecords / pu.numKeys
		nextRangeStart int
		nextRangeEnd   int
	)

	if keyID == pu.numKeys {
		rangeEnd = pu.numRecords
	}

	curProgress, err := pu.pr.GetCommittedOffset(keyID)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0, 0, err
		}
		nextRangeStart = rangeStart
	} else {
		nextRangeStart = curProgress.Committed + 1
	}
	nextRangeEnd = nextRangeStart + pu.batchSize - 1
	if nextRangeEnd > rangeEnd {
		nextRangeEnd = rangeEnd
	}

	return nextRangeStart, nextRangeEnd, nil
}

func (pu *ProgressUseCase) SetCommittedOffset(keyID int, offset int) error {
	if keyID < 1 || keyID > pu.numKeys {
		return utils.ValidationErr{Err: "invalid key ID"}
	}
	if offset < 1 {
		return utils.ValidationErr{Err: "invalid committed offset"}
	}

	err := pu.pr.SetCommittedOffset(keyID, offset)
	return err
}
