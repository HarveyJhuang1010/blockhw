package block

import (
	"strconv"

	"github.com/HarveyJhuang1010/blockhw/internal/appcontext"
	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type blockController struct {
	in digIn
}

func newBlockController(in digIn) bo.BlockController {
	return &blockController{
		in: in,
	}
}

func (b *blockController) GetLatestBlocks(ginCtx *gin.Context) {
	// get query param
	q := ginCtx.Query("limit")
	limit, err := strconv.Atoi(q)
	if err != nil {
		appcontext.GetLogger().Error("GetBlockDetail Failed", zap.Error(err))
		ginCtx.JSON(400, gin.H{
			"error": "invalid limit",
		})
		return
	}
	if limit == 0 {
		limit = 10
	}

	res, err := b.in.UseCase.GetLatestBlocks(ginCtx, limit)
	if err != nil {
		appcontext.GetLogger().Error("GetBlockDetail Failed", zap.Error(err))
		ginCtx.JSON(500, gin.H{
			"error": "internal error",
		})
		return
	}

	ginCtx.JSON(200, res)
}

func (b *blockController) GetBlockDetail(ginCtx *gin.Context) {
	id := ginCtx.Param("id")

	if id == "" {
		ginCtx.JSON(400, gin.H{
			"error": "invalid block id",
		})
		return
	}

	blockNum, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		ginCtx.JSON(400, gin.H{
			"error": "invalid block id",
		})
		return
	}

	res, err := b.in.UseCase.GetBlockDetail(ginCtx, blockNum)
	if err != nil {
		appcontext.GetLogger().Error("GetBlockDetail Failed", zap.Error(err))
		ginCtx.JSON(500, gin.H{
			"error": "internal error",
		})
		return
	}

	ginCtx.JSON(200, res)
}
