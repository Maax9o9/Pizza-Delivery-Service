package main

import (
    "log"
    infra "delivery-service/src/infrastructure"
    routes "delivery-service/src/infrastructure/routes"
    "delivery-service/src/infrastructure/controller"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
    r := gin.Default()
    r.Use(cors.Default())

    deliveryAlertService, rabbitMQ := infra.Init()
    deliveryAlertController := controllers.NewDeliveryAlertController(deliveryAlertService, rabbitMQ)
    routes.DeliveryRoutes(r, deliveryAlertController)

    if err := r.Run(":7070"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}