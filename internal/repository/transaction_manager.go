package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// TransactionManagerImpl implementa a interface TransactionManager
type TransactionManagerImpl struct {
	db *sql.DB
}

// NewTransactionManager cria uma nova instância de TransactionManager
func NewTransactionManager(db *sql.DB) TransactionManager {
	return &TransactionManagerImpl{db: db}
}

// WithinTransaction executa uma função dentro de uma transação
func (tm *TransactionManagerImpl) WithinTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic após garantir rollback
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("erro ao fazer rollback: %v (erro original: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao fazer commit: %w", err)
	}

	return nil
}
