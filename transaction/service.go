package transaction

import (
	"context"
)

// Service defines methods for account service as specified by API
type Service interface {
	CreateNewTransaction(ctx context.Context, transactionID, accountID string, amount int) error
	UpdateTransaction(ctx context.Context, transactionID, accountID string, amount int) error
	GetTransactionByID(ctx context.Context, transactionID string) (*Transaction, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) UpdateTransaction(ctx context.Context, transactionID, accountID string, amount int) error {
	err := s.repository.UpdateTransaction(ctx, &Transaction{
		ID:        transactionID,
		AccountID: accountID,
		Amount:    amount,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CreateNewTransaction(ctx context.Context, transactionID, accountID string, amount int) error {
	err := s.repository.NewTransaction(ctx, &Transaction{
		ID:        transactionID,
		AccountID: accountID,
		Amount:    amount,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetTransactionByID(ctx context.Context, transactionID string) (*Transaction, error) {
	transaction, err := s.repository.GetTransactionById(ctx, transactionID)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
