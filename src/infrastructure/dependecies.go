package infrastructure

import (
	"delivery-service/src/application"
	"delivery-service/src/application/repositorys"
	"delivery-service/src/application/services"
	"delivery-service/src/infrastructure/adapters"
	"log"
)

func Init() (*services.DeliveryAlertService, *adapters.RabbitMQ) {
    rabbitMQ, err := adapters.NewRabbitMQ("amqp://max:123@44.213.165.25:5672/", "Pizzas")
    if err != nil {
        log.Fatalf("Failed to initialize RabbitMQ: %v", err)
    }

    rabbitRepo := repositorys.NewRabbitRepository(rabbitMQ)
    repo := NewMySQL()

    deliveryAlertUseCase := application.NewDeliveryAlertUseCase(repo, rabbitRepo)
    deliveryAlertService := services.NewDeliveryAlertService(rabbitMQ, deliveryAlertUseCase)

    go deliveryAlertService.ProcessDeliveryAlerts()

    return deliveryAlertService, rabbitMQ
}