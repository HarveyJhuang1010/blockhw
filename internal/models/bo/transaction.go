package bo

import (
	"context"

	"github.com/HarveyJhuang1010/blockhw/internal/models/dto"
	"github.com/HarveyJhuang1010/blockhw/internal/models/po"
	"github.com/gin-gonic/gin"
)

type TransactionRepo interface {
	GetTransactions(ctx context.Context, blockNum uint64) ([]*po.Transaction, error)
	GetTransactionDetail(ctx context.Context, txHash string) (*po.Transaction, error)
}

type TransactionUseCase interface {
	GetTransactionDetail(ctx context.Context, txHash string) (*dto.TransactionDetail, error)
}

type TransactionController interface {
	GetTransactionDetail(ginCtx *gin.Context)
}
