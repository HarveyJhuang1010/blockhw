package bo

import (
	"context"

	"github.com/HarveyJhuang1010/blockhw/internal/models/dto"
	"github.com/HarveyJhuang1010/blockhw/internal/models/po"
	"github.com/gin-gonic/gin"
)

type BlockRepo interface {
	GetLatestBlocks(ctx context.Context, limit int) ([]*po.Block, error)
	GetBlockDetail(ctx context.Context, blockNum uint64) (*po.Block, error)
}

type BlockUseCase interface {
	GetLatestBlocks(ctx context.Context, limit int) (*dto.BlockListResp, error)
	GetBlockDetail(ctx context.Context, blockNum uint64) (*dto.BlockDetail, error)
}

type BlockController interface {
	GetLatestBlocks(ginCtx *gin.Context)
	GetBlockDetail(ginCtx *gin.Context)
}
