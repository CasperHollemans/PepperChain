package transaction

import (
	"bytes"
	"encoding/json"
	"github.com/CasperHollemans/PepperChain/internal/crypto"
	"github.com/CasperHollemans/PepperChain/internal/logging"
	"github.com/CasperHollemans/PepperChain/internal/network"
	"net/http"
	"time"
)

type Transaction struct {
	TimeStamp int64
	Sender    string
	Recipient string
	Amount    int64
	Signature []byte
	PublicKey []byte
}

type TransactionService interface {
	AddTransaction(tx Transaction)
	CreateAndSignTransaction(amount int64, recipient string) (Transaction, error)
}

type TransactionServiceImpl struct {
	memoryPoolService MemoryPoolService
	logger            logging.Logger
	cryptoService     crypto.CryptoService
	networkService    network.NetworkService
}

func NewTransactionService(
	memoryPoolService MemoryPoolService,
	logger logging.Logger,
	cryptoService crypto.CryptoService,
	networkService network.NetworkService,
) *TransactionServiceImpl {
	return &TransactionServiceImpl{
		memoryPoolService: memoryPoolService,
		logger:            logger,
		cryptoService:     cryptoService,
		networkService:    networkService,
	}
}

func (t *TransactionServiceImpl) AddTransaction(tx Transaction) {
	t.logger.Info("Received new transaction")
	isValid := t.validateTransaction(tx)
	if !isValid {
		t.logger.Error("Transaction is not valid", tx)
		return
	}
	t.logger.Debug("Adding transaction to memory pool", tx)
	t.memoryPoolService.AddTransaction(tx)
	t.BroadcastTransaction(tx)
}

func (t *TransactionServiceImpl) validateTransaction(tx Transaction) bool {
	if tx.Amount <= 0 {
		t.logger.Error("Transaction amount is not valid", tx)
		return false
	}
	if tx.Recipient == "" {
		t.logger.Error("Transaction recipient is not valid", tx)
		return false
	}
	if tx.Signature == nil {
		t.logger.Error("Transaction signature is not valid", tx)
		return false
	}
	if tx.PublicKey == nil {
		t.logger.Error("Transaction public key is not valid", tx)
		return false
	}
	if tx.TimeStamp > time.Now().Unix() {
		t.logger.Error("Transaction timestamp is not valid", tx)
		return false
	}

	isValid := t.verifySignature(tx)
	if !isValid {
		t.logger.Error("Transaction signature is not valid", tx)
		return false
	}

	return true
}

func (t *TransactionServiceImpl) verifySignature(tx Transaction) bool {
	// Transactions are signed without signature and public key
	txCopy := Transaction{
		TimeStamp: tx.TimeStamp,
		Sender:    tx.Sender,
		Recipient: tx.Recipient,
		Amount:    tx.Amount,
	}
	txJSON, err := json.Marshal(txCopy)
	if err != nil {
		t.logger.Error("Error marshalling transaction", err, tx)
		return false
	}

	hash := t.cryptoService.Hash(txJSON)
	return t.cryptoService.VerifySignature(tx.PublicKey, hash, tx.Signature)
}

func (t *TransactionServiceImpl) CreateAndSignTransaction(amount int64, recipient string) (Transaction, error) {
	privateKey, _, err := t.cryptoService.GenerateKeyPair()
	if err != nil {
		return Transaction{}, err
	}

	tx := Transaction{
		TimeStamp: time.Now().Unix(),
		Sender:    "",
		Recipient: recipient,
		Amount:    amount,
	}

	txJSON, err := json.Marshal(tx)
	if err != nil {
		panic(err)
	}

	hash := t.cryptoService.Hash(txJSON)
	signature := t.cryptoService.Sign(privateKey, hash)
	tx.Signature = signature
	tx.PublicKey = append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return tx, nil
}

func (t *TransactionServiceImpl) BroadcastTransaction(tx Transaction) error {
	t.logger.Info("Broadcasting transaction", tx)
	for _, node := range t.networkService.GetNodes() {
		jsonTx, err := json.Marshal(tx)
		if err != nil {
			t.logger.Error("Error marshalling transaction", err)
			return err
		}
		url := node + "/transactions"
		t.logger.Debug("Broadcasting transaction to node", url)
		http.Post(url, "application/json", bytes.NewBuffer(jsonTx))
	}
	return nil
}
