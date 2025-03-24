package main

import (
	"log"
	"time"

	handlers "github.com/LuckyanenkoArtur/go-wallet-test-task/internal/api/handlers"
	"github.com/LuckyanenkoArtur/go-wallet-test-task/internal/services/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	postgresDB, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect to PostgreeDB: %v", err)
	}
	defer postgresDB.Close()

	indexHandler := &handlers.IndexHandler{}
	walletHandlers := handlers.NewWalletHandler(postgresDB)
	router := gin.Default()

	// Setup CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "OPTIONS", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/", indexHandler.RootHandler)
	router.GET("/api/v1/wallets/", walletHandlers.ListWalletsHandler) // Get all wallets
	router.GET("/api/v1/wallets/:wallet_uuid", walletHandlers.GetWallet) // Get Wallet by uuid
	router.POST("/api/v1/wallets/", walletHandlers.UpdateWallet) // Manage money of wallet

	log.Fatal(router.Run(":5000"))
}