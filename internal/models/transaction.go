package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionLog struct {
	ID              uuid.UUID `json:"id"`
	WalletID        uuid.UUID `json:"wallet_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	CreatedAt       time.Time `json:"created_at"`
}