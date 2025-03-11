package controllers

import (
    "net/http"
    "encoding/json"
    "delivery-service/src/application/services"
    "delivery-service/src/infrastructure/adapters"
    "delivery-service/src/domain/entities"
    "github.com/gin-gonic/gin"
)

type DeliveryAlertController struct {
    deliveryAlertService *services.DeliveryAlertService
    rabbitMQ             *adapters.RabbitMQ
}

func NewDeliveryAlertController(deliveryAlertService *services.DeliveryAlertService, rabbitMQ *adapters.RabbitMQ) *DeliveryAlertController {
    return &DeliveryAlertController{
        deliveryAlertService: deliveryAlertService,
        rabbitMQ:             rabbitMQ,
    }
}

func (dac *DeliveryAlertController) GetDeliveryAlerts(c *gin.Context) {
    alert, err := dac.deliveryAlertService.GetLatestDeliveryAlert()
    if err != nil {
        c.JSON(http.StatusNoContent, gin.H{})
        return
    }
    c.JSON(http.StatusOK, alert)
}


func (dac *DeliveryAlertController) GetLatestDeliveryAlert(c *gin.Context) {
    alert, err := dac.deliveryAlertService.GetLatestDeliveryAlert()
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "No hay alertas recientes", "details": err.Error()})
        return
    }
    c.JSON(http.StatusOK, alert)
}

func (dac *DeliveryAlertController) PublishDeliveryAlert(c *gin.Context) {
    var alert entities.DeliveryAlert
    if err := c.ShouldBindJSON(&alert); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato JSON inválido", "details": err.Error()})
        return
    }

    message, err := json.Marshal(alert)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la alerta", "details": err.Error()})
        return
    }

    if err := dac.rabbitMQ.PublishToQueue("Alert", string(message)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al publicar en la cola", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "Alerta publicada exitosamente"})
}
