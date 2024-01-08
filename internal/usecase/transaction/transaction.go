package transaction

import (
	"context"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/HarveyJhuang1010/blockhw/internal/models/dto"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

type transactionUseCase struct {
	in digIn
}

func newTransactionUseCase(in digIn) bo.TransactionUseCase {
	return &transactionUseCase{
		in: in,
	}
}

func (t *transactionUseCase) GetTransactionDetail(ctx context.Context, txHash string) (*dto.TransactionDetail, error) {
	var res dto.TransactionDetail

	tx, err := t.in.Repo.GetTransactionDetail(ctx, txHash)
	if err != nil {
		return nil, errors.Wrap(err, "get transaction detail")
	}

	if err := copier.Copy(&res, &tx); err != nil {
		return nil, errors.Wrap(err, "copy transaction detail")
	}

	return &res, nil
}
