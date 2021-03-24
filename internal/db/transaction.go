package db

import (
	"database/sql"
)

const (
	eventTypeCreate = "create"
	eventTypeUpdate = "update"
	eventTypeDelete = "delete"
)

type Transaction interface {
	UserTransaction
}

type txImpl struct {
	tx *sql.Tx
	userTransactionImpl
}

func newTransaction(tx *sql.Tx, redisClient RedisClient) Transaction {
	return &txImpl{
		tx:                  tx,
		userTransactionImpl: userTransactionImpl{tx: tx, redisClient: redisClient},
	}
}
