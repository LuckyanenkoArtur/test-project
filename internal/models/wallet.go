package models

import (
	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Balance   float64   `json:"balance"`
}

type WalletRequest struct {
	WalletID     uuid.UUID `json:"walletId"`
	OperationType string    `json:"operationType"`
	Amount       float64   `json:"amount"`
}