package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"wallet/internal/database"
	"wallet/internal/event"
	"wallet/internal/event/handler"
	createaccount "wallet/internal/usecase/create_account"
	createclient "wallet/internal/usecase/create_client"
	createtransaction "wallet/internal/usecase/create_transaction"
	"wallet/internal/web"
	"wallet/internal/web/webserver"
	"wallet/pkg/events"
	"wallet/pkg/kafka"
	"wallet/pkg/uow"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		"root", "root", "mysql", 3306, "wallet", "utf8", true, "Local")
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)
	uow.Register("AccountRepository", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})
	uow.Register("TransactionRepository", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := createclient.NewCreateClientUseCase(clientDb)
	createAccountUseCase := createaccount.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := createtransaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)

	webserver := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)
	webserver.AddHandler("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	webserver.Start()
}
