package routes

import (
    "delivery-service/src/infrastructure/controller"
    "github.com/gin-gonic/gin"
)

func DeliveryRoutes(router *gin.Engine, deliveryAlertController *controllers.DeliveryAlertController) {
    router.POST("/api/delivery", deliveryAlertController.PublishDeliveryAlert)
    router.GET("/api/delivery", deliveryAlertController.GetDeliveryAlerts)
}