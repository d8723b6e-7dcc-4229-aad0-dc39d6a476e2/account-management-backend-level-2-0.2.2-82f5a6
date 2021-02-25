package account

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetAccountById(ctx context.Context, accountID string) (*Account, error)
	SaveNewAccount(ctx context.Context, account *Account) error
	SetAccountBalance(ctx context.Context, accountID string, newBalance int) error
}

type sqlLiteRepository struct {
	db *sql.DB
}

func NewSqlLiteRepository(db *sql.DB) *sqlLiteRepository {
	return &sqlLiteRepository{db}
}

// GetAccountById returns pointer to account if its found and error if error is recieved while
// trying to fetch account from database. If account with the provided id does not exist. nil for account
// and error is returned
func (r *sqlLiteRepository) GetAccountById(ctx context.Context, accountID string) (*Account, error) {
	row := r.db.QueryRow("SELECT account_id, balance FROM account where account_id=?", accountID)
	var account Account
	err := row.Scan(&account.ID, &account.Balance)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &account, nil
	default:
		return nil, err
	}
}

func (r *sqlLiteRepository) SaveNewAccount(ctx context.Context, account *Account) error {
	stmt, err := r.db.Prepare("INSERT INTO account (account_id, balance) values(?,?)")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(account.ID, account.Balance)
	if err != nil {
		return err
	}
	return nil
}

func (r *sqlLiteRepository) SetAccountBalance(ctx context.Context, accountID string, newBalance int) error {
	stmt, err := r.db.Prepare("update account set balance=? where account_id=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(newBalance, accountID)
	if err != nil {
		return err
	}
	return nil
}
