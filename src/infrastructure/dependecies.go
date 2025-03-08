package infrastructure

import (
    "log"
    "delivery-service/src/application/services"
    "delivery-service/src/application/repositorys"
    "delivery-service/src/infrastructure/adapters"
    "delivery-service/src/infrastructure/controller"
)

func Init() *controllers.DeliveryAlertController {
    rabbitMQ, err := adapters.NewRabbitMQ("amqp://max:123@44.213.165.25:5672/", "Pizzas")
    if err != nil {
        log.Fatalf("Failed to initialize RabbitMQ: %v", err)
    }

    rabbitRepo := repositorys.NewRabbitRepository(rabbitMQ)
    repo := NewMySQL()

    deliveryAlertUseCase := services.NewDeliveryAlertUseCase(repo, rabbitRepo)
    deliveryAlertService := services.NewDeliveryAlertService(rabbitMQ, deliveryAlertUseCase)
    deliveryAlertController := controllers.NewDeliveryAlertController(deliveryAlertService)

    go deliveryAlertService.ProcessDeliveryAlerts()

    return deliveryAlertController
}