package block

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/HarveyJhuang1010/blockhw/internal/models/dto"
	"github.com/HarveyJhuang1010/blockhw/internal/models/po"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
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

	res.Transactions = []string{}
	for _, tx := range block.Transactions {
		res.Transactions = append(res.Transactions, tx.Hash)
	}

	return &res, nil
}

func (b *blockUseCase) SyncBlockByNum(ctx context.Context, blockNum uint64) error {
	bn := big.NewInt(int64(blockNum))
	block, err := b.in.EvmClient.BlockByNumber(ctx, bn)
	if err != nil {
		return errors.Wrap(err, "get block by number")
	}

	poBlock := &po.Block{
		Number:     block.Number().Uint64(),
		Hash:       block.Hash().String(),
		ParentHash: block.ParentHash().String(),
		Time:       block.Time(),
	}

	chainID, err := b.in.EvmClient.NetworkID(ctx)
	if err != nil {
		return errors.Wrap(err, "get network id")
	}

	var txns []*po.Transaction
	for _, txn := range block.Transactions() {
		msg, err := core.TransactionToMessage(txn, types.LatestSignerForChainID(chainID), nil)
		if err != nil {
			return errors.Wrap(err, "get message")
		}

		tr, err := b.in.EvmClient.TransactionReceipt(ctx, txn.Hash())
		if err != nil {
			return errors.Wrap(err, "get transaction receipt")
		}

		var logs []*po.TransactionLog
		for _, l := range tr.Logs {
			logs = append(logs, &po.TransactionLog{
				TransactionHash: l.TxHash.String(),
				Index:           l.Index,
				Data:            hex.EncodeToString(l.Data),
			})
		}

		var to, value string
		if txn.To() != nil {
			to = txn.To().String()
		}
		if txn.Value() != nil {
			value = txn.Value().String()
		}
		poTxn := &po.Transaction{
			Hash:        txn.Hash().String(),
			From:        msg.From.String(),
			To:          to,
			Nonce:       txn.Nonce(),
			Data:        hex.EncodeToString(txn.Data()),
			Value:       value,
			BlockNumber: blockNum,
			Logs:        logs,
		}

		txns = append(txns, poTxn)
	}

	poBlock.Transactions = txns

	if err := b.in.Repo.SyncBlock(ctx, poBlock); err != nil {
		return errors.Wrap(err, "create block")
	}

	return nil
}
