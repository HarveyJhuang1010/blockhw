package cron

import (
	"fmt"

	"github.com/HarveyJhuang1010/blockhw/internal/appcontext"
	"github.com/HarveyJhuang1010/blockhw/internal/config"
	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/HarveyJhuang1010/blockhw/internal/models/po"
	"go.uber.org/zap"
)

type blockSyncTask struct {
	in digIn
}

var (
	_            bo.CronTask = (*blockSyncTask)(nil)
	currentBlock uint64
)

func newBlockSyncTask(in digIn) bo.CronTask {
	return &blockSyncTask{
		in: in,
	}
}

func (t *blockSyncTask) Schedule() string {
	//return "0 */3 * * * *" // 3 minutes for test
	return "*/10 * * * * *"
}

func (t *blockSyncTask) Run() {
	ctx := appcontext.GetContext()
	if currentBlock == 0 {
		latestBlockInDB, err := t.in.BlockRepo.GetLatestBlocks(ctx, 1)
		if err != nil {
			appcontext.GetLogger().Error("get latest block failed", zap.Error(err))
			return
		}
		if len(latestBlockInDB) == 0 {
			currentBlock = config.GetConfig().Worker.StartNumber
		} else {
			currentBlock = latestBlockInDB[0].Number
		}
	}

	maxNumber, err := t.in.EvmClient.BlockNumber(ctx)
	if err != nil {
		appcontext.GetLogger().Error("get block number failed", zap.Error(err))
		return
	}

	for currentBlock <= maxNumber {
		rc := &po.BlockSyncRecord{Number: currentBlock, Status: "created"}
		if err := t.in.BlockRepo.SaveBlockSyncRecord(ctx, rc); err != nil {
			appcontext.GetLogger().Error("save block sync record failed", zap.Error(err), zap.Uint64("block", currentBlock))
			return
		}
		if err := t.in.NatClient.Publish(ctx, fmt.Sprintf("block.%d", currentBlock), currentBlock); err != nil {
			appcontext.GetLogger().Error("publish block failed", zap.Error(err), zap.Uint64("block", currentBlock))
			return
		}
		currentBlock++
	}
}
