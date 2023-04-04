package transaction

import (
	"context"
	"database/sql"
)

// Repository defines methods for repository
type Repository interface {
	GetTransactionById(ctx context.Context, transactionID string) (*Transaction, error)
	NewTransaction(ctx context.Context, transaction *Transaction) error
	UpdateTransaction(ctx context.Context, transaction *Transaction) error
}

type sqlLiteRepository struct {
	db *sql.DB
}

func NewSqlLiteRepository(db *sql.DB) *sqlLiteRepository {
	return &sqlLiteRepository{db}
}

func (r *sqlLiteRepository) GetTransactionById(ctx context.Context, transactionID string) (*Transaction, error) {
	row := r.db.QueryRow("SELECT transaction_id, account_id, amount FROM 'transaction' where transaction_id=?", transactionID)
	var transaction Transaction
	err := row.Scan(&transaction.ID, &transaction.AccountID, &transaction.Amount)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &transaction, nil
	default:
		return nil, err
	}
}

func (r *sqlLiteRepository) NewTransaction(ctx context.Context, transaction *Transaction) error {
	stmt, err := r.db.Prepare("INSERT INTO 'transaction' (transaction_id, account_id, amount) VALUES (?,?,?)")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(transaction.ID, transaction.AccountID, transaction.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (r *sqlLiteRepository) UpdateTransaction(ctx context.Context, transaction *Transaction) error {
	stmt, err := r.db.Prepare("update 'transaction' set amount=?, account_id=?  where transaction_id=?")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(transaction.Amount, transaction.AccountID, transaction.ID)
	if err != nil {
		return err
	}
	return nil
}
