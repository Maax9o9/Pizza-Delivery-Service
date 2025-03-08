package controllers

import (
	"delivery-service/src/application/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeliveryAlertController struct {
	deliveryAlertService *services.DeliveryAlertService
}

func NewDeliveryAlertController(deliveryAlertService *services.DeliveryAlertService) *DeliveryAlertController {
	return &DeliveryAlertController{
		deliveryAlertService: deliveryAlertService,
	}
}

func (dac *DeliveryAlertController) ProcessDeliveryAlerts(c *gin.Context) {
	go dac.deliveryAlertService.ProcessDeliveryAlerts()
	c.JSON(http.StatusOK, gin.H{"status": "Processing delivery alerts"})
}

func (dac *DeliveryAlertController) GetDeliveryAlerts(c *gin.Context) {
	alerts, err := dac.deliveryAlertService.GetAllDeliveryAlerts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, alerts)
}
