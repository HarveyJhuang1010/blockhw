package transaction

import (
	"github.com/HarveyJhuang1010/blockhw/internal/appcontext"
	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type transactionController struct {
	in digIn
}

func newTransactionController(in digIn) bo.TransactionController {
	return &transactionController{
		in: in,
	}
}

func (t *transactionController) GetTransactionDetail(ginCtx *gin.Context) {
	txHash := ginCtx.Param("txHash")

	if txHash == "" {
		ginCtx.JSON(400, gin.H{
			"error": "invalid transaction id",
		})
		return
	}

	res, err := t.in.UseCase.GetTransactionDetail(ginCtx, txHash)
	if err != nil {
		appcontext.GetLogger().Error("GetTransactionDetail Failed", zap.Error(err))
		ginCtx.JSON(500, gin.H{
			"error": "internal error",
		})
		return
	}

	ginCtx.JSON(200, res)
}
