package routes

import (
	controller "delivery-service/src/infrastructure/controller"
	"github.com/gin-gonic/gin"
)

func DeliveryRoutes(router *gin.Engine, controller *controller.DeliveryAlertController) {
	router.GET("/api/delivery", controller.ProcessDeliveryAlerts)
}