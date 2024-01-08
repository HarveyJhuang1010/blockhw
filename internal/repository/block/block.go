package block

import (
	"context"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/HarveyJhuang1010/blockhw/internal/models/po"
	"github.com/pkg/errors"
	"gorm.io/gorm"
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

func (b *blockRepo) SyncBlock(ctx context.Context, block *po.Block) error {
	if err := b.in.RDB.Transaction(func(tx *gorm.DB) error {
		var blockSyncRecord po.BlockSyncRecord
		if err := tx.Where("number = ?", block.Number).First(&blockSyncRecord).Error; err != nil {
			return errors.WithStack(err)
		}

		if blockSyncRecord.Status == "synced" || blockSyncRecord.Status == "confirmed" {
			return nil
		}

		if err := tx.Model(&blockSyncRecord).Update("status", "synced").Error; err != nil {
			return errors.WithStack(err)
		}

		if err := tx.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(block).Error; err != nil {
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (b *blockRepo) SaveBlockSyncRecord(ctx context.Context, record *po.BlockSyncRecord) error {
	if err := b.in.RDB.Save(record).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (b *blockRepo) GetMinUnSyncRecord(ctx context.Context) (*po.BlockSyncRecord, error) {
	var res po.BlockSyncRecord

	if err := b.in.RDB.Where("status = ?", "created").First(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	return &res, nil
}

func (b *blockRepo) GetUnConfirmedRecord(ctx context.Context, blockNum uint64) ([]*po.BlockSyncRecord, error) {
	var res []*po.BlockSyncRecord

	maxConfirmNumber := blockNum - 20
	if err := b.in.RDB.Where("number <= ? AND status = ?", maxConfirmNumber, "synced").Find(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	return res, nil
}
