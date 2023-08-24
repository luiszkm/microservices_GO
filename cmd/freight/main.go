package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	routeusecase "github.com/luiszkm/microservices_GO/internal/Aplicatioon/usecase/route_usecase"
	"github.com/luiszkm/microservices_GO/internal/Domain/freight/entity"
	"github.com/luiszkm/microservices_GO/internal/Infra/repository"
	"github.com/luiszkm/microservices_GO/pkg/kafka"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306)/routes?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	msgChan := make(chan *ckafka.Message)
	topics := []string{"routes"}
	servers := "host.docker.internal:9094"
	go kafka.Consumer(topics, servers, msgChan)

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
