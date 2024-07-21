package models

type WalletUpdate struct {
	Total    int    `json:"total" database:"total"`
	Username string `json:"username" database:"username"`
	Action   string `json:"action"`
}

const (
	WALLET_TOTAL_ADD       = "ADD"
	WALLET_TOTAL_SUBSTRUCT = "SUBSTRUCT"
)
