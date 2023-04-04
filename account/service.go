package account

import (
	"context"
	"errors"
	"log"

	"github.com/37d7fcd0-f13b-4571-8fd3-fc12d70c7b7d/account-management-backend-level-2-0.1.5-79653d/transaction"
)

// Service defines methods for account service as specified by API
type Service interface {
	UpdateAccountBalance(ctx context.Context, accountID string, transactionID string, amount int) error
	GetAccountBalance(ctx context.Context, accountID string) (int, error)
}

type service struct {
	repository         Repository
	transactionService transaction.Service
	// Basic implementation of caching account balances
	accountsBalances map[string]int
}

func NewService(repository Repository, transactionService transaction.Service) Service {
	return &service{repository: repository, transactionService: transactionService, accountsBalances: map[string]int{}}

}

// GetAccountBalance first checks service accountsBalances cache to know if balance is already known.
// If balance is not known it fetches it from repository.
func (s *service) GetAccountBalance(ctx context.Context, accountID string) (int, error) {
	v, ok := s.accountsBalances[accountID]
	if ok {
		return v, nil
	}
	account, err := s.repository.GetAccountById(ctx, accountID)
	if err != nil {
		return 0, err
	}
	if account == nil {
		return 0, errors.New("not found")
	}
	s.accountsBalances[accountID] = account.Balance
	return account.Balance, nil
}

// UpdateAccountBalance updates account balance and sets the new balance for accountsBalances in service
func (s *service) UpdateAccountBalance(ctx context.Context, accountID, transactionID string, amount int) error {
	var balance int
	account, err := s.repository.GetAccountById(ctx, accountID)
	if err != nil {
		return err
	}
	if account == nil {
		err = s.repository.SaveNewAccount(ctx, &Account{
			ID:      accountID,
			Balance: 0,
		})
		if err != nil {
			return err
		}
		balance = 0
	} else {
		balance = account.Balance
	}

	//check if transaction does not exist already. If it exists return.
	t, err := s.transactionService.GetTransactionByID(ctx, transactionID)
	if err != nil {
		return err
	}
	if t != nil {
		return nil
	}

	err = s.transactionService.CreateNewTransaction(ctx, transactionID, accountID, amount)
	if err != nil {
		return err
	}

	newBalance := balance + amount
	err = s.repository.SetAccountBalance(ctx, accountID, newBalance)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Update accounts balances cache.
	s.accountsBalances[accountID] = newBalance
	return nil
}
