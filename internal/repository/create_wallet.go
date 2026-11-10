package repository

import (
	"context"
	"fmt"

	"github.com/dimozavrrrik/wallet/internal/model"
)

const queryCreateWallet = `
INSERT INTO wallets (id, balance) 
VALUES (gen_random_uuid(), 0)
RETURNING id, balance;
`

func (r *WalletRepository) Create(ctx context.Context) (model.Wallet, error) {

	var wallet model.Wallet

	err := r.pool.QueryRow(ctx, queryCreateWallet).Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		return model.Wallet{}, fmt.Errorf("create new wallet: %w", err)
	}

	return wallet, nil

}
