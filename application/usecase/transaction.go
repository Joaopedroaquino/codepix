package usecase

import (
	"errors"
	"log"

	"github.com/joaopedroaquino/codepix/codepix-go/domain/model"
)

type TransactionUseCase struct {
	TransactionRepository model.TransactionRepositoryInterface
	PixRepository         model.PixKeyRepositoryInterface
}

func (t *TransactionUseCase) Register(accountId string, amount float64, pixKeyto string, pixKeyKindTo string, description string) (*model.Transaction, error) {
	account, err := t.PixRepository.FindAccount(accountId)

	if err != nil {
		return nil, err
	}

	pixKey, err := t.PixRepository.FindKeyByKind(pixKeyto, pixKeyKindTo)

	if err != nil {
		return nil, err
	}

	transaction, err := model.NewTransaction(account, amount, pixKey, description)

	if err != nil {
		return nil, err
	}

	t.TransactionRepository.Save(transaction)

	if transaction.ID == "" {

		return nil, errors.New("unable to process this transaction")

	}

	return transaction, nil
}

func (t *TransactionUseCase) Complete(transactionId string) (*model.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(transactionId)

	if err != nil {
		log.Println("transaction not found")
		return nil, err
	}

	transaction.Status = model.TransacttionCompleted
	err = t.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, errors.New("unable to complete this transaction")
	}

	return transaction, nil
}

func (t *TransactionUseCase) Error(transactionId string, reason string) (*model.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(transactionId)

	if err != nil {
		log.Println("transaction not found")
		return nil, err
	}

	transaction.Status = model.TransactionError
	transaction.CancelDescription = reason

	err = t.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, errors.New("unable to complete this transaction")
	}

	return transaction, nil
}
