package block

import (
	"context"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/HarveyJhuang1010/blockhw/internal/models/dto"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

type blockUseCase struct {
	in digIn
}

func newBlockUseCase(in digIn) bo.BlockUseCase {
	return &blockUseCase{
		in: in,
	}
}

func (b *blockUseCase) GetLatestBlocks(ctx context.Context, limit int) (*dto.BlockListResp, error) {
	var resp dto.BlockListResp

	blocks, err := b.in.Repo.GetLatestBlocks(ctx, limit)
	if err != nil {
		return nil, errors.Wrap(err, "get latest blocks")
	}

	var blockList []*dto.BlockList

	if err := copier.Copy(&blockList, &blocks); err != nil {
		return nil, errors.Wrap(err, "copy block list")
	}

	resp.Blocks = blockList

	return &resp, nil
}

func (b *blockUseCase) GetBlockDetail(ctx context.Context, blockNum uint64) (*dto.BlockDetail, error) {
	var res dto.BlockDetail

	block, err := b.in.Repo.GetBlockDetail(ctx, blockNum)
	if err != nil {
		return nil, errors.Wrap(err, "get block detail")
	}

	if err := copier.Copy(&res, &block); err != nil {
		return nil, errors.Wrap(err, "copy block detail")
	}

	for _, tx := range block.Transactions {
		res.Transactions = append(res.Transactions, tx.Hash)
	}

	return &res, nil
}
