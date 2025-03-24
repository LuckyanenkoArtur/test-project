package handlers

import (
	"net/http"
	"strings"

	"github.com/LuckyanenkoArtur/go-wallet-test-task/internal/models"
	"github.com/LuckyanenkoArtur/go-wallet-test-task/internal/services/db"
	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	db          *db.PostgresDB
}

func NewWalletHandler(db *db.PostgresDB) *WalletHandler {
	return &WalletHandler{
		db:          db,
	}
}

func (h *WalletHandler) ListWalletsHandler(ctx *gin.Context) {
	query := `SELECT * FROM wallets`
	rows, err := h.db.DB.Query(query)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	var wallets []models.Wallet
	for rows.Next() {
		var wallet models.Wallet
		err := rows.Scan(
			&wallet.ID,
			&wallet.UserID,
			&wallet.Balance,
		)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		wallets = append(wallets, wallet)
	}

	if len(wallets) == 0 {
		ctx.JSON(404, gin.H{
			"error": "No wallets found",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"data": wallets,
	})
}
func (h *WalletHandler) GetWallet(ctx *gin.Context) {
	id := ctx.Params.ByName("wallet_uuid")
	query := `SELECT * FROM wallets WHERE id = $1`

	rows, err := h.db.DB.Query(query, id)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	var wallet models.Wallet
	if rows.Next() {
		err := rows.Scan(
			&wallet.ID,
			&wallet.UserID,
			&wallet.Balance,
		)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		ctx.JSON(404, gin.H{
			"error": "Wallet not found",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"data": wallet,
	})
}

func (h *WalletHandler) UpdateWallet(ctx *gin.Context) {
	var request models.WalletRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	if request.Amount <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}

	operationType := strings.ToUpper(request.OperationType)
	if operationType != "DEPOSIT" && operationType != "WITHDRAW" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation type"})
		return
	}

	tx, err := h.db.DB.Begin()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	var currentBalance float64
	query := `SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(query, request.WalletID).Scan(&currentBalance)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Wallet not found"})
		return
	}

	var newBalance float64
	if operationType == "DEPOSIT" {
		newBalance = currentBalance + request.Amount
	} else if operationType == "WITHDRAW" {
		if currentBalance < request.Amount {
			tx.Rollback()
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
			return
		}
		newBalance = currentBalance - request.Amount
	}

	updateQuery := `UPDATE wallets SET balance = $1 WHERE id = $2`
	_, err = tx.Exec(updateQuery, newBalance, request.WalletID)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet balance"})
		return
	}

	var operationTypeID int
	opQuery := `SELECT id FROM operation_type WHERE name = $1`
	err = tx.QueryRow(opQuery, operationType).Scan(&operationTypeID)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Operation type not found"})
		return
	}

	transactionQuery := `INSERT INTO transaction_logs (wallet_id, operation_type_id, amount) VALUES ($1, $2, $3)`
	_, err = tx.Exec(transactionQuery, request.WalletID, operationTypeID, request.Amount)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log transaction"})
		return
	}

	err = tx.Commit()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Wallet updated successfully",
		"status":     200,
		"newBalance": newBalance,
	})
}

