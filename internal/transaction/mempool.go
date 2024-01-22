package transaction

import (
	"github.com/CasperHollemans/PepperChain/internal/logging"
	"sync"
)

type MemoryPool struct {
	Transactions []Transaction
}

type MemoryPoolService interface {
	AddTransaction(tx Transaction)
	GetTransactions() []Transaction
}

type MemoryPoolServiceImpl struct {
	memoryPool *MemoryPool
	logger     logging.Logger
	mu         sync.Mutex
}

func NewMemoryPoolService(logger logging.Logger) *MemoryPoolServiceImpl {
	return &MemoryPoolServiceImpl{
		memoryPool: &MemoryPool{
			Transactions: []Transaction{},
		},
		logger: logger,
	}
}

func (m *MemoryPoolServiceImpl) AddTransaction(tx Transaction) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Debug("Adding transaction to memory pool", tx)
	m.memoryPool.Transactions = append(m.memoryPool.Transactions, tx)
}

func (m *MemoryPoolServiceImpl) GetTransactions() []Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Debug("Getting transactions from memory pool")
	return m.memoryPool.Transactions
}
