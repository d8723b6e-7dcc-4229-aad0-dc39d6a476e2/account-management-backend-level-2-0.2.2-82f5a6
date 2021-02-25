package account

// Account defines account structure. It has balance and id
type Account struct {
	ID      string `json:"account_id"`
	Balance int    `json:"balance"`
}
