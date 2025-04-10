package handler

import (
	"balance/internal/usecase/update_account_balance"
	"balance/pkg/events"
	"log"
	"sync"
)

type BalanceUpdatedPayload struct {
	AccountIdFrom        string  `json:"account_id_from"`
	AccountIdTo          string  `json:"account_id_to"`
	BalanceAccountIdFrom float64 `json:"balance_account_id_from"`
	BalanceAccountIdTo   float64 `json:"balance_account_id_to"`
}

type BalanceUpdatedKafkaHandler struct {
	UpdateBalanceUseCase *update_account_balance.UpdateAccountBalanceUseCase
}

func NewBalanceUpdatedKafkaHandler(
	updateBalanceUseCase *update_account_balance.UpdateAccountBalanceUseCase,
) *BalanceUpdatedKafkaHandler {
	return &BalanceUpdatedKafkaHandler{
		UpdateBalanceUseCase: updateBalanceUseCase,
	}
}

func (h *BalanceUpdatedKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()

	if message.GetName() != "BalanceUpdated" {
		log.Print("Received message with wrong event name")
		return
	}

	msgPayload, ok := message.GetPayload().(map[string]interface{})
	if !ok {
		log.Print("Failed to cast message payload to map[string]interface{}")
		return
	}

	payload := &BalanceUpdatedPayload{
		AccountIdFrom:        msgPayload["account_id_from"].(string),
		AccountIdTo:          msgPayload["account_id_to"].(string),
		BalanceAccountIdFrom: msgPayload["balance_account_id_from"].(float64),
		BalanceAccountIdTo:   msgPayload["balance_account_id_to"].(float64),
	}

	input := update_account_balance.UpdateAccountBalanceInputDTO{
		AccountID: payload.AccountIdFrom,
		Balance:   payload.BalanceAccountIdFrom,
	}
	output, err := h.UpdateBalanceUseCase.Execute(input)
	if err != nil {
		log.Printf("Failed to update balance for account %s: %v", payload.AccountIdFrom, err)
		return
	}
	log.Printf("Updated balance for account %s: %f\n", output.AccountID, output.Balance)

	input2 := update_account_balance.UpdateAccountBalanceInputDTO{
		AccountID: payload.AccountIdTo,
		Balance:   payload.BalanceAccountIdTo,
	}
	output2, err := h.UpdateBalanceUseCase.Execute(input2)
	if err != nil {
		log.Printf("Failed to update balance for account %s: %v", payload.AccountIdTo, err)
		return
	}
	log.Printf("Updated balance for account %s: %f\n", output2.AccountID, output2.Balance)
}
