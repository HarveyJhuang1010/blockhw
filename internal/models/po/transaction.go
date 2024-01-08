package po

type Transaction struct {
	Base
	Hash        string           `gorm:"column:hash;type:varchar(66);primaryKey" json:"hash"`
	From        string           `gorm:"column:from;type:varchar(42)" json:"from"`
	To          string           `gorm:"column:to;type:varchar(42)" json:"to"`
	Nonce       uint64           `gorm:"column:nonce;type:numeric(20,0)" json:"nonce"`
	Data        string           `gorm:"column:data;type:varchar(66)" json:"data"`
	Value       string           `gorm:"column:value;type:varchar(66)" json:"value"`
	BlockNumber uint64           `gorm:"column:block_number;type:numeric(20,0);index" json:"block_number"`
	Logs        []TransactionLog `gorm:"foreignKey:TransactionHash;references:Hash;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"logs"`
}

type TransactionLog struct {
	Base
	ID              uint64 `gorm:"column:id;type:numeric(20,0);primaryKey" json:"id"`
	TransactionHash string `gorm:"column:transaction_hash;type:varchar(66);index" json:"transaction_hash"`
	Index           uint64 `gorm:"column:index;type:numeric(20,0)" json:"index"`
	Data            string `gorm:"column:data;type:varchar(66)" json:"data"`
}
