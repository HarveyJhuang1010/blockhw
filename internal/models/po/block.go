package po

type Block struct {
	Base
	Number       uint64         `gorm:"column:number;type:numeric(20,0);primaryKey" json:"number"`
	Hash         string         `gorm:"column:hash;type:varchar(66);uniqueIndex" json:"hash"`
	Parent       string         `gorm:"column:parent;type:varchar(66)" json:"parent"`
	Time         uint64         `gorm:"column:time;type:numeric(20,0)" json:"time"`
	Transactions []*Transaction `gorm:"foreignKey:BlockNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"transactions"`
}
