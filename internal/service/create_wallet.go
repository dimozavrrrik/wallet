package service

import (
	"context"

	"github.com/dimozavrrrik/wallet/internal/model"
)

func (s *WalletService) Create(ctx context.Context) (model.Wallet, error) {
	return s.repo.Create(ctx)
}
