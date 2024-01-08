package po

type Transaction struct {
	Base
	Hash        string           `gorm:"column:hash;type:varchar(66);primaryKey" json:"tx_hash"`
	From        string           `gorm:"column:from;type:varchar(42)" json:"from"`
	To          string           `gorm:"column:to;type:varchar(42)" json:"to"`
	Nonce       uint64           `gorm:"column:nonce;type:uint" json:"nonce"`
	Data        string           `gorm:"column:data;type:varchar(66)" json:"data"`
	Value       string           `gorm:"column:value;type:varchar(66)" json:"value"`
	BlockNumber uint64           `gorm:"column:block_number;type:uint"`
	Logs        []TransactionLog `gorm:"foreignKey:TransactionHash;references:Hash;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"logs"`
}

type TransactionLog struct {
	Base
	ID              uint64 `gorm:"column:id;type:uint;primaryKey"`
	TransactionHash string `gorm:"column:transaction_hash;type:varchar(66);index"`
	Index           uint64 `gorm:"column:index;type:uint" json:"index"`
	Data            string `gorm:"column:data;type:varchar(66)" json:"data"`
}
