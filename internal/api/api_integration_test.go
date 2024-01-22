package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/CasperHollemans/PepperChain/internal/config"
	"github.com/CasperHollemans/PepperChain/internal/transaction"
	"net/http"
	"testing"
	"time"
)

func setup() *ServerCtx {
	ctx := SetupServerContext()
	go ctx.server.ListenAndServe()
	return ctx
}

func TestSetupServer(t *testing.T) {
	server := SetupServerContext()
	if server == nil {
		t.Errorf("Expected server to be setup, got nil")
	}
}

/*
When: new transaction is created
Expect: 200 status code
*/
func TestCreateTransactionExpect200(t *testing.T) {
	// Setup
	ctx := setup()
	defer ctx.server.Close()
	tx, _ := ctx.transactionService.CreateAndSignTransaction(10, "recipient")
	// Execute
	post, _ := executeTransactionsPost(tx)
	// Verify
	if post.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %v", post.StatusCode)
	}
}

/*
When: new transaction is created
Expect: no error
*/
func TestCreateTransactionExpectNoError(t *testing.T) {
	// Setup
	ctx := setup()
	defer ctx.server.Close()
	tx, _ := ctx.transactionService.CreateAndSignTransaction(10, "recipient")
	// Execute
	_, err := executeTransactionsPost(tx)
	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

/*
When: new transaction is created
Expect: transaction is added to the memory pool
*/
func TestCreateTransactionExpectTransactionInMemoryPool(t *testing.T) {
	// Setup
	ctx := setup()
	defer ctx.server.Close()
	tx, _ := ctx.transactionService.CreateAndSignTransaction(10, "recipient")
	// Execute
	executeTransactionsPost(tx)
	// Verify
	if len(ctx.memoryPoolService.GetTransactions()) != 1 {
		t.Errorf("Expected 1 transaction in memory pool, got %v", len(ctx.memoryPoolService.GetTransactions()))
	}
}

/*
When: new transaction is created with invalid signature
Expect: transaction is not added to the memory pool
*/
func TestCreateTransactionExpectNoTransactionInMemoryPool(t *testing.T) {
	// Setup
	ctx := setup()
	defer ctx.server.Close()
	tx, _ := ctx.transactionService.CreateAndSignTransaction(10, "recipient")
	tx.Signature = []byte("invalid")
	// Execute
	executeTransactionsPost(tx)
	// Verify
	if len(ctx.memoryPoolService.GetTransactions()) != 0 {
		t.Errorf("Expected 0 transaction in memory pool, got %v", len(ctx.memoryPoolService.GetTransactions()))
	}
}

/*
When: new transaction is created with invalid amount
Expect: transaction is not added to the memory pool
*/
func TestCreateTransactionExpectNoTransactionInMemoryPoolInvalidAmount(t *testing.T) {
	// Setup
	ctx := setup()
	defer ctx.server.Close()
	tx, _ := ctx.transactionService.CreateAndSignTransaction(-10, "recipient")
	// Execute
	executeTransactionsPost(tx)
	// Verify
	if len(ctx.memoryPoolService.GetTransactions()) != 0 {
		t.Errorf("Expected 0 transaction in memory pool, got %v", len(ctx.memoryPoolService.GetTransactions()))
	}
}

/*
When: new transaction is created with empty recipient
Expect: transaction is not added to the memory pool
*/
func TestCreateTransactionExpectNoTransactionInMemoryPoolInvalidRecipient(t *testing.T) {
	// Setup
	ctx := setup()
	defer ctx.server.Close()
	tx, _ := ctx.transactionService.CreateAndSignTransaction(10, "")
	// Execute
	executeTransactionsPost(tx)
	// Verify
	if len(ctx.memoryPoolService.GetTransactions()) != 0 {
		t.Errorf("Expected 0 transaction in memory pool, got %v", len(ctx.memoryPoolService.GetTransactions()))
	}
}

/*
When: new transaction is created with invalid public key
Expect: transaction is not added to the memory pool
*/
func TestCreateTransactionExpectNoTransactionInMemoryPoolInvalidPublicKey(t *testing.T) {
	// Setup
	ctx := setup()
	defer ctx.server.Close()
	tx, _ := ctx.transactionService.CreateAndSignTransaction(10, "recipient")
	tx.PublicKey = nil
	// Execute
	executeTransactionsPost(tx)
	// Verify
	if len(ctx.memoryPoolService.GetTransactions()) != 0 {
		t.Errorf("Expected 0 transaction in memory pool, got %v", len(ctx.memoryPoolService.GetTransactions()))
	}
}

/*
When: new transaction is created with invalid timestamp
Expect: transaction is not added to the memory pool
*/
func TestCreateTransactionExpectNoTransactionInMemoryPoolInvalidTimestamp(t *testing.T) {
	// Setup
	ctx := setup()
	defer ctx.server.Close()
	tx, _ := ctx.transactionService.CreateAndSignTransaction(10, "recipient")
	tx.TimeStamp = time.Now().Add(time.Hour).Unix()
	// Execute
	executeTransactionsPost(tx)
	// Verify
	if len(ctx.memoryPoolService.GetTransactions()) != 0 {
		t.Errorf("Expected 0 transaction in memory pool, got %v", len(ctx.memoryPoolService.GetTransactions()))
	}
}

/*
When: new transaction is created
Expect: transaction is broadcasted
*/
func TestCreateTransactionExpectTransactionBroadcasted(t *testing.T) {
	// Setup
	receivedChan := make(chan bool, 1)
	timeout := time.After(2 * time.Second)
	server := &http.Server{
		Addr: ":8081",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedChan <- true
			w.WriteHeader(200)
			w.Write([]byte("OK"))
			return
		}),
	}
	go server.ListenAndServe()
	defer server.Close()
	time.Sleep(500 * time.Millisecond)
	ctx := setup()
	defer ctx.server.Close()
	ctx.networkService.RegisterNode("http://localhost:8081")
	tx, _ := ctx.transactionService.CreateAndSignTransaction(10, "recipient")

	// Execute
	executeTransactionsPost(tx)

	// Verify
	select {
	case <-receivedChan:
		fmt.Println("Transaction received")
		break
	case <-timeout:
		t.Errorf("Expected transaction to be broadcasted")
	}
}

func getUrl(path string) string {
	var cfg = config.NewConfig()
	return fmt.Sprintf("%s/%s", cfg.BaseUrl, path)
}

func executeTransactionsPost(tx transaction.Transaction) (*http.Response, error) {
	encoded, _ := json.Marshal(tx)
	body := bytes.NewBuffer(encoded)
	post, err := http.Post(getUrl("transactions"), "application/json", body)
	return post, err
}
