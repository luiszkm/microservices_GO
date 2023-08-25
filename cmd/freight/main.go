package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	routeusecase "github.com/luiszkm/microservices_GO/internal/Aplicatioon/usecase/route_usecase"
	"github.com/luiszkm/microservices_GO/internal/Domain/freight/entity"
	"github.com/luiszkm/microservices_GO/internal/Infra/repository"
	"github.com/luiszkm/microservices_GO/pkg/kafka"
)

func main() {
	msgChan := make(chan *ckafka.Message)
	topics := []string{"routes"}
	servers := "host.docker.internal:9094"
	go kafka.Consume(topics, servers, msgChan)

	db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306)/routes?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	CreateTableRoutesMysql(db)
	repository := repository.NewRouterRepositoryMysql(db)
	freight := entity.NewFreight(10)

	createRouteUseCase := routeusecase.NewCreateRouteUseCase(repository, freight)
	changeRouteStatusUseCase := routeusecase.NewChangeRouteStatusUseCase(repository)

	for msg := range msgChan {
		input := routeusecase.CreateRouteInput{}
		json.Unmarshal(msg.Value, &input)

		switch input.Event {
		case "RouteCreated":
			output, err := createRouteUseCase.Execute(input)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(output)
		case "RouteStarted", "RouteFinished":
			input := routeusecase.ChangeStatusInput{}
			json.Unmarshal(msg.Value, &input)
			output, err := changeRouteStatusUseCase.Execute(input)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(output)
		}
	}

}

func CreateTableRoutesMysql(db *sql.DB) {

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS routes (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
    distance FLOAT NOT NULL,
    status VARCHAR(255) NOT NULL,
    freight_price FLOAT ,
    started_at DATETIME,
		finished_at DATETIME
	);
`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table created successfully")

}
