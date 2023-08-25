package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	routeusecase "github.com/luiszkm/microservices_GO/internal/Aplicatioon/usecase/route_usecase"
	"github.com/luiszkm/microservices_GO/internal/Domain/freight/entity"
	"github.com/luiszkm/microservices_GO/internal/Infra/repository"
	"github.com/luiszkm/microservices_GO/pkg/kafka"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	routesCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "routes_created_total",
			Help: "Total number of created routes",
		},
	)

	routesStarted = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "routes_started_total",
			Help: "Total number of started routes",
		},
		[]string{"status"},
	)

	errorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
	)
)

func init() {
	prometheus.MustRegister(routesStarted)
	prometheus.MustRegister(errorsTotal)
	prometheus.MustRegister(routesCreated)
}

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
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":8080", nil)
	}()
	go http.ListenAndServe(":8081", nil)
	repository := repository.NewRouterRepositoryMysql(db)
	freight := entity.NewFreight(10)

	createRouteUseCase := routeusecase.NewCreateRouteUseCase(repository, freight)
	changeRouteStatusUseCase := routeusecase.NewChangeRouteStatusUseCase(repository)

	for msg := range msgChan {
		input := routeusecase.CreateRouteInput{}
		json.Unmarshal(msg.Value, &input)

		switch input.Event {
		case "RouteCreated":
			_, err := createRouteUseCase.Execute(input)
			if err != nil {
				fmt.Println(err)
				errorsTotal.Inc()
			} else {
				routesCreated.Inc()
			}
		case "RouteStarted":

			input := routeusecase.ChangeStatusInput{}
			json.Unmarshal(msg.Value, &input)
			_, err := changeRouteStatusUseCase.Execute(input)
			if err != nil {
				fmt.Println(err)
				errorsTotal.Inc()
			} else {
				routesStarted.WithLabelValues("started").Inc()
			}
		case "RouteFinished":
			input := routeusecase.ChangeStatusInput{}
			json.Unmarshal(msg.Value, &input)
			_, err := changeRouteStatusUseCase.Execute(input)
			if err != nil {
				fmt.Println(err)
				errorsTotal.Inc()
			} else {
				routesStarted.WithLabelValues("finished").Inc()
			}
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
