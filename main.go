package main

import (
    infra "delivery-service/src/infrastructure"
    routes "delivery-service/src/infrastructure/routes"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
    r := gin.Default()
    r.Use(cors.Default())

    deliveryAlertController := infra.Init()
    routes.DeliveryRoutes(r, deliveryAlertController)

    r.Run(":7070")
}