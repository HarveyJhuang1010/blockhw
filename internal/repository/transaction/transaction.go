package transaction

import (
	"context"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/HarveyJhuang1010/blockhw/internal/models/po"
	"github.com/pkg/errors"
)

type transactionRepo struct {
	in digIn
}

func newTransactionRepo(in digIn) bo.TransactionRepo {
	return &transactionRepo{
		in: in,
	}
}
func (t *transactionRepo) GetTransactions(ctx context.Context, blockNum uint64) ([]*po.Transaction, error) {
	var res []*po.Transaction
	if err := t.in.RDB.Where("block_num = ?", blockNum).Preload("Logs").Find(&res).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (t *transactionRepo) GetTransactionDetail(ctx context.Context, txHash string) (*po.Transaction, error) {
	var res po.Transaction

	if err := t.in.RDB.Where("tx_hash = ?", txHash).Preload("Logs").First(&res).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return &res, nil
}
