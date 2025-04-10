package main

import (
	"balance/internal/database"
	"balance/internal/event"
	"balance/internal/event/handler"
	"balance/internal/usecase/get_account_balance"
	"balance/internal/usecase/update_account_balance"
	"balance/internal/web"
	"balance/internal/web/webserver"
	"balance/pkg/events"
	"balance/pkg/kafka"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

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
		"group.id":          "balance-service",
		"auto.offset.reset": "earliest",
	}

	// Create balance database gateway
	balanceDb := database.NewBalanceDB(db)

	// Create use cases
	getAccountBalanceUseCase := get_account_balance.NewGetAccountBalanceUseCase(balanceDb)
	updateAccountBalanceUseCase := update_account_balance.NewUpdateAccountBalanceUseCase(balanceDb)

	// Create the Kafka consumer
	consumer := kafka.NewConsumer(&configMap, []string{"balances"})

	// Create the balance updated handler
	balanceUpdatedHandler := handler.NewBalanceUpdatedKafkaHandler(updateAccountBalanceUseCase)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("BalanceUpdated", balanceUpdatedHandler)

	msgChan := make(chan *ckafka.Message)
	go consumer.Consume(msgChan)

	var wg sync.WaitGroup
	go func() {
		for msg := range msgChan {
			wg.Add(1)
			var balanceUpdatedEvent event.BalanceUpdated
			if err := json.Unmarshal(msg.Value, &balanceUpdatedEvent); err != nil {
				fmt.Println("Error parsing JSON:", err)
				return
			}
			balanceUpdatedHandler.Handle(&balanceUpdatedEvent, &wg)
		}
	}()

	// Setup web server
	webserver := webserver.NewWebServer(":3003")
	balanceHandler := web.NewBalanceHandler(getAccountBalanceUseCase)
	webserver.AddGetHandler("/balances/{account_id}", balanceHandler.GetAccountBalance)
	webserver.AddGetHandler("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Start the web server
	fmt.Println("Balance service started on :3003")
	webserver.Start()
}
