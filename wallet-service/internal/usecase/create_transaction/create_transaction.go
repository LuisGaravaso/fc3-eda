package createtransaction

import (
	"context"
	"wallet/internal/entity"
	"wallet/internal/gateway"
	"wallet/pkg/events"
	"wallet/pkg/uow"
)

type CreateTransactionInputDTO struct {
	AccountIdFrom string  `json:"account_id_from"`
	AccountIdTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionOutputDTO struct {
	Id            string  `json:"id"`
	AccountIdFrom string  `json:"account_id_from"`
	AccountIdTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type BalanceUpdatedOutputDTO struct {
	AccountIdFrom        string  `json:"account_id_from"`
	AccountIdTo          string  `json:"account_id_to"`
	BalanceAccountIdFrom float64 `json:"balance_account_id_from"`
	BalanceAccountIdTo   float64 `json:"balance_account_id_to"`
}

type CreateTransactionUseCase struct {
	Uow                     uow.UowInterface
	EventDispatcher         events.EventDispatcherInterface
	TransactionCreatedEvent events.EventInterface
	BalanceUpdatedEvent     events.EventInterface
}

func NewCreateTransactionUseCase(
	uow uow.UowInterface,
	eventsDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
	balanceUpdated events.EventInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		Uow:                     uow,
		EventDispatcher:         eventsDispatcher,
		TransactionCreatedEvent: transactionCreated,
		BalanceUpdatedEvent:     balanceUpdated,
	}
}

func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	var transactionOutput *CreateTransactionOutputDTO
	var balanceOutput *BalanceUpdatedOutputDTO

	err := uc.Uow.Do(ctx, func(uow *uow.Uow) error {
		// Get repositories
		accountGateway, err := uc.getAccountRepository(ctx)
		if err != nil {
			return err
		}

		transactionGateway, err := uc.getTransactionRepository(ctx)
		if err != nil {
			return err
		}

		// Find accounts
		accountFrom, err := accountGateway.FindById(input.AccountIdFrom)
		if err != nil {
			return err
		}

		accountTo, err := accountGateway.FindById(input.AccountIdTo)
		if err != nil {
			return err
		}

		// Create transaction
		transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
		if err != nil {
			return err
		}

		// Update balances
		err = accountGateway.UpdateBalance(accountFrom)
		if err != nil {
			return err
		}

		err = accountGateway.UpdateBalance(accountTo)
		if err != nil {
			return err
		}

		// Save transaction
		err = transactionGateway.Create(transaction)
		if err != nil {
			return err
		}

		balanceOutput = &BalanceUpdatedOutputDTO{
			AccountIdFrom:        accountFrom.Id,
			AccountIdTo:          accountTo.Id,
			BalanceAccountIdFrom: accountFrom.Balance,
			BalanceAccountIdTo:   accountTo.Balance,
		}

		transactionOutput = &CreateTransactionOutputDTO{
			Id:            transaction.Id,
			AccountIdFrom: transaction.AccountFrom.Id,
			AccountIdTo:   transaction.AccountTo.Id,
			Amount:        transaction.Amount,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Dispatch the event after the transaction is successfully committed
	uc.TransactionCreatedEvent.SetPayload(transactionOutput)
	uc.EventDispatcher.Dispatch(uc.TransactionCreatedEvent)

	uc.BalanceUpdatedEvent.SetPayload(balanceOutput)
	uc.EventDispatcher.Dispatch(uc.BalanceUpdatedEvent)
	return transactionOutput, nil
}

func (uc *CreateTransactionUseCase) getAccountRepository(ctx context.Context) (gateway.AccountGateway, error) {
	accountRepository, err := uc.Uow.GetRepository(ctx, "AccountRepository")
	if err != nil {
		return nil, err
	}
	return accountRepository.(gateway.AccountGateway), nil
}

func (uc *CreateTransactionUseCase) getTransactionRepository(ctx context.Context) (gateway.TransactionGateway, error) {
	transactionRepository, err := uc.Uow.GetRepository(ctx, "TransactionRepository")
	if err != nil {
		return nil, err
	}
	return transactionRepository.(gateway.TransactionGateway), nil
}
