package transaction

type Transaction struct {
	ID        string `json:"-"`
	AccountID string `json:"account_id"`
	Amount    int    `json:"amount"`
}
