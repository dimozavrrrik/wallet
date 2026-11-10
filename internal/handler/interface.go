package handler

import (
	"context"

	"github.com/dimozavrrrik/wallet/internal/model"
)

type WalletService interface {
	Create(ctx context.Context) (model.Wallet, error)
}
