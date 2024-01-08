package po

type Block struct {
	Base
	Number       uint64         `gorm:"column:number;type:uint;primaryKey" json:"block_number"`
	Hash         string         `gorm:"column:hash;type:varchar(67);uniqueIndex" json:"block_hash"`
	ParentHash   string         `gorm:"column:parent;type:varchar(67)" json:"parent_hash"`
	Time         uint64         `gorm:"column:time;type:uint" json:"block_time"`
	Transactions []*Transaction `gorm:"foreignKey:BlockNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type BlockSyncRecord struct {
	Base
	Number uint64 `gorm:"column:number;type:uint;primaryKey"`
	// created, synced, confirmed
	Status string `gorm:"column:status;type:varchar(10)"`
}
