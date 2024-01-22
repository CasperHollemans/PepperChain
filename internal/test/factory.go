package test

import (
	"github.com/CasperHollemans/PepperChain/internal/transaction"
	"time"
)

func CreateTransaction() transaction.Transaction {
	return transaction.Transaction{
		TimeStamp: time.Now().Unix(),
		Sender:    "",
		Recipient: "",
		Amount:    0,
		Signature: "",
	}
}
