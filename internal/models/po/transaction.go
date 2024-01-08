package po

type Transaction struct {
	Base
	Hash        string            `gorm:"column:hash;type:varchar(67);primaryKey" json:"tx_hash"`
	From        string            `gorm:"column:from;type:varchar(43)" json:"from"`
	To          string            `gorm:"column:to;type:varchar(43)" json:"to"`
	Nonce       uint64            `gorm:"column:nonce;type:uint" json:"nonce"`
	Data        string            `gorm:"column:data;type:text" json:"data"`
	Value       string            `gorm:"column:value;type:varchar(67)" json:"value"`
	BlockNumber uint64            `gorm:"column:block_number;type:uint"`
	Logs        []*TransactionLog `gorm:"foreignKey:TransactionHash;references:Hash;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"logs"`
}

type TransactionLog struct {
	Base
	TransactionHash string `gorm:"column:transaction_hash;type:varchar(67);primaryKey"`
	Index           uint   `gorm:"column:index;type:uint;primaryKey" json:"index"`
	Data            string `gorm:"column:data;type:text" json:"data"`
}
