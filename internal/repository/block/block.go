package block

import (
	"context"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/HarveyJhuang1010/blockhw/internal/models/po"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

type blockRepo struct {
	in digIn
}

func newBlockRepo(in digIn) bo.BlockRepo {
	return &blockRepo{
		in: in,
	}
}

func (b *blockRepo) GetLatestBlocks(ctx context.Context, limit int) ([]*po.Block, error) {
	var res []*po.Block
	if err := b.in.RDB.Preload("Transactions").Order("number desc").Limit(limit).Find(&res).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (b *blockRepo) GetBlockDetail(ctx context.Context, blockNum uint64) (*po.Block, error) {
	var res po.Block

	if err := b.in.RDB.Where("number = ?", blockNum).Preload("Transactions").First(&res).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return &res, nil
}

func (b *blockRepo) CreateBlock(ctx context.Context, block *po.Block) error {
	if err := b.in.RDB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(block).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
