package po

type Block struct {
	Base
	Number       uint64         `gorm:"column:number;type:numeric(20,0);primaryKey" json:"block_number"`
	Hash         string         `gorm:"column:hash;type:varchar(66);uniqueIndex" json:"block_hash"`
	Parent       string         `gorm:"column:parent;type:varchar(66)" json:"parent_hash"`
	Time         uint64         `gorm:"column:time;type:numeric(20,0)" json:"block_time"`
	Transactions []*Transaction `gorm:"foreignKey:BlockNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
