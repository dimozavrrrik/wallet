package service

import (
	"context"

	"github.com/dimozavrrrik/wallet/internal/model"
)

type WalletRepository interface {
	Create(ctx context.Context) (model.Wallet, error)
}
