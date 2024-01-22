package api

import (
	"github.com/CasperHollemans/PepperChain/internal/config"
	"github.com/CasperHollemans/PepperChain/internal/crypto"
	"github.com/CasperHollemans/PepperChain/internal/logging"
	"github.com/CasperHollemans/PepperChain/internal/network"
	"github.com/CasperHollemans/PepperChain/internal/transaction"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ServerCtx struct {
	server             *http.Server
	transactionService transaction.TransactionService
	memoryPoolService  transaction.MemoryPoolService
	cfg                *config.Config
	logger             logging.Logger
	cryptoService      crypto.CryptoService
	test               *bool
	networkService     network.NetworkService
}

func SetupServerContext() *ServerCtx {
	cfg := config.NewConfig()
	logger := logging.NewZapLogger(true)
	memoryPoolService := transaction.NewMemoryPoolService(logger)
	cryptoService := crypto.NewCryptoService(logger)
	networkService := network.NewNetworkService(logger)
	transactionService := transaction.NewTransactionService(memoryPoolService, logger, cryptoService, networkService)

	ctx := &ServerCtx{
		transactionService: transactionService,
		cfg:                cfg,
		logger:             logger,
		memoryPoolService:  memoryPoolService,
		cryptoService:      cryptoService,
		networkService:     networkService,
	}

	r := gin.Default()
	r.POST("/transactions", addTransaction(ctx))

	ctx.server = &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}
	return ctx
}

func StartServer() {
	ctx := SetupServerContext()
	err := ctx.server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func addTransaction(ctx *ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tx transaction.Transaction
		if err := c.ShouldBindJSON(&tx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.transactionService.AddTransaction(tx)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
