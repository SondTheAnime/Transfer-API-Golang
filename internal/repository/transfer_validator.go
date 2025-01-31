package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

// TransferValidatorImpl implementa a interface PendingTransferValidator
type TransferValidatorImpl struct {
	db               *sql.DB
	pendingTransfers map[int]float64
	mu               sync.RWMutex
}

// NewTransferValidator cria uma nova instância de PendingTransferValidator
func NewTransferValidator(db *sql.DB) PendingTransferValidator {
	return &TransferValidatorImpl{
		db:               db,
		pendingTransfers: make(map[int]float64),
	}
}

// ValidateTransfer valida se uma transferência pode ser realizada
func (v *TransferValidatorImpl) ValidateTransfer(ctx context.Context, userID int, amount float64) error {
	v.mu.RLock()
	pendingAmount := v.pendingTransfers[userID]
	v.mu.RUnlock()

	var currentBalance float64
	err := v.db.QueryRowContext(ctx, "SELECT balance FROM users WHERE id = $1", userID).Scan(&currentBalance)
	if err == sql.ErrNoRows {
		return ErrUserNotFound
	}
	if err != nil {
		return fmt.Errorf("erro ao verificar saldo: %w", err)
	}

	if currentBalance-(pendingAmount+amount) < 0 {
		return fmt.Errorf("saldo insuficiente considerando transferências pendentes (saldo: %.2f, pendente: %.2f, transferência: %.2f)",
			currentBalance, pendingAmount, amount)
	}

	return nil
}

// RegisterPendingTransfer registra uma transferência pendente e retorna uma função de cleanup
func (v *TransferValidatorImpl) RegisterPendingTransfer(ctx context.Context, userID int, amount float64) func() {
	v.mu.Lock()
	v.pendingTransfers[userID] += amount
	v.mu.Unlock()

	return func() {
		v.mu.Lock()
		v.pendingTransfers[userID] -= amount
		v.mu.Unlock()
	}
}
