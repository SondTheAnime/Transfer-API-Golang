package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"transfer-api/internal/models"
)

// ErrInsufficientBalance representa um erro de saldo insuficiente
var ErrInsufficientBalance = errors.New("saldo insuficiente")

// ErrUserNotFound representa um erro de usuário não encontrado
var ErrUserNotFound = errors.New("usuário não encontrado")

// UserReader define operações de leitura para usuários
type UserReader interface {
	GetUser(ctx context.Context, id int) (*models.User, error)
}

// UserWriter define operações de escrita para usuários
type UserWriter interface {
	UpdateBalance(ctx context.Context, id int, amount float64) error
}

// TransferManager define operações relacionadas a transferências
type TransferManager interface {
	Transfer(ctx context.Context, fromID, toID int, amount float64) error
}

// TransactionManager define operações de transação
type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(*sql.Tx) error) error
}

// PendingTransferValidator define operações para validação de transferências pendentes
type PendingTransferValidator interface {
	ValidateTransfer(ctx context.Context, userID int, amount float64) error
	RegisterPendingTransfer(ctx context.Context, userID int, amount float64) func()
}

// Repository combina todas as interfaces necessárias
type Repository interface {
	UserReader
	UserWriter
	TransferManager
}

type UserRepository struct {
	db               *sql.DB
	mu               sync.RWMutex
	txManager        TransactionManager
	validator        PendingTransferValidator
	pendingTransfers map[int]float64
	pendingMutex     sync.Mutex
}

func NewUserRepository(db *sql.DB, validator PendingTransferValidator) *UserRepository {
	return &UserRepository{
		db:               db,
		txManager:        NewTransactionManager(db),
		validator:        validator,
		pendingTransfers: make(map[int]float64),
	}
}

func (r *UserRepository) GetUser(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, "SELECT id, balance FROM users WHERE id = $1", id).Scan(&user.ID, &user.Balance)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}
	return user, nil
}

func (r *UserRepository) UpdateBalance(ctx context.Context, id int, amount float64) error {
	result, err := r.db.ExecContext(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", amount, id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar saldo: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar atualização: %w", err)
	}
	if rows != 1 {
		return ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) Transfer(ctx context.Context, fromID, toID int, amount float64) error {
	// Valida a transferência considerando transferências pendentes
	if err := r.validator.ValidateTransfer(ctx, fromID, amount); err != nil {
		return err
	}

	// Registra a transferência pendente e obtém função de cleanup
	cleanup := r.validator.RegisterPendingTransfer(ctx, fromID, amount)
	defer cleanup()

	// Executa a transferência dentro de uma transação
	err := r.txManager.WithinTransaction(ctx, func(tx *sql.Tx) error {
		// Verifica e bloqueia o usuário de origem
		var fromBalance float64
		err := tx.QueryRowContext(ctx, "SELECT balance FROM users WHERE id = $1 FOR UPDATE", fromID).Scan(&fromBalance)
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		if err != nil {
			return fmt.Errorf("erro ao obter usuário de origem: %w", err)
		}

		if fromBalance < amount {
			return ErrInsufficientBalance
		}

		// Verifica e bloqueia o usuário de destino
		var toBalance float64
		err = tx.QueryRowContext(ctx, "SELECT balance FROM users WHERE id = $1 FOR UPDATE", toID).Scan(&toBalance)
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		if err != nil {
			return fmt.Errorf("erro ao obter usuário de destino: %w", err)
		}

		// Atualiza os saldos
		if err := updateBalanceInTx(ctx, tx, fromID, -amount); err != nil {
			return err
		}

		if err := updateBalanceInTx(ctx, tx, toID, amount); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("Erro na transferência: %v", err)
		return err
	}

	return nil
}

// updateBalanceInTx é uma função auxiliar para atualizar o saldo dentro de uma transação
func updateBalanceInTx(ctx context.Context, tx *sql.Tx, userID int, amount float64) error {
	result, err := tx.ExecContext(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", amount, userID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar saldo: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar atualização: %w", err)
	}
	if rows != 1 {
		return ErrUserNotFound
	}
	return nil
}
