package web

import (
	"balance/internal/usecase/get_account_balance"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type BalanceHandler struct {
	GetAccountBalanceUseCase *get_account_balance.GetAccountBalanceUseCase
}

func NewBalanceHandler(
	getAccountBalanceUseCase *get_account_balance.GetAccountBalanceUseCase,
) *BalanceHandler {
	return &BalanceHandler{
		GetAccountBalanceUseCase: getAccountBalanceUseCase,
	}
}

func (h *BalanceHandler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "account_id")

	input := get_account_balance.GetAccountBalanceInputDTO{
		AccountID: accountID,
	}

	output, err := h.GetAccountBalanceUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
