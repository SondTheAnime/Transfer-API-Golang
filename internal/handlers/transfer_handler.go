package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"transfer-api/internal/repository"
)

type TransferHandler struct {
	repo *repository.UserRepository
}

type TransferRequest struct {
	FromID int     `json:"from_id"`
	ToID   int     `json:"to_id"`
	Amount float64 `json:"amount"`
}

func NewTransferHandler(repo *repository.UserRepository) *TransferHandler {
	return &TransferHandler{repo: repo}
}

func (h *TransferHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID de usuário inválido", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"balance": user.Balance})
}

func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao decodificar requisição", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Valor da transferência deve ser positivo", http.StatusBadRequest)
		return
	}

	if req.FromID == req.ToID {
		http.Error(w, "Não é permitido realizar transferência para si mesmo", http.StatusBadRequest)
		return
	}

	if err := h.repo.Transfer(req.FromID, req.ToID, req.Amount); err != nil {
		switch err.Error() {
		case "usuário de origem não encontrado", "usuário de destino não encontrado":
			http.Error(w, err.Error(), http.StatusNotFound)
		case "saldo insuficiente":
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Erro interno ao processar transferência", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transferência realizada com sucesso"})
}
