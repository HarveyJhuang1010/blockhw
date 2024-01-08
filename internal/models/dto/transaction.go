package dto

type TransactionDetail struct {
	Hash  string                `json:"tx_hash"`
	From  string                `json:"from"`
	To    string                `json:"to"`
	Nonce uint64                `json:"nonce"`
	Data  string                `json:"data"`
	Value int64                 `json:"value"`
	Logs  []*TransactionLogList `json:"logs"`
}

type TransactionLogList struct {
	Index int64  `json:"index"`
	Data  string `json:"data"`
}
