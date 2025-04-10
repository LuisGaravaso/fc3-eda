package web

import (
	"encoding/json"
	"net/http"
	createclient "wallet/internal/usecase/create_client"
)

type WebClientHandler struct {
	CreateClientUseCase createclient.CreateClientUseCase
}

func NewWebClientHandler(createClientUseCase createclient.CreateClientUseCase) *WebClientHandler {
	return &WebClientHandler{
		CreateClientUseCase: createClientUseCase,
	}
}

func (h *WebClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var input createclient.CreateClientInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := h.CreateClientUseCase.Execute(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(output); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
