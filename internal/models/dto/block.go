package dto

type BlockListResp struct {
	Blocks []*BlockList `json:"blocks"`
}

type BlockList struct {
	Number     uint64 `json:"block_number"`
	Hash       string `json:"block_hash"`
	Time       int64  `json:"block_time"`
	ParentHash string `json:"parent_hash"`
}

type BlockDetail struct {
	Number       uint64   `json:"block_number"`
	Hash         string   `json:"block_hash"`
	Time         int64    `json:"block_time"`
	ParentHash   string   `json:"parent_hash"`
	Transactions []string `json:"transactions"`
}
