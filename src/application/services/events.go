package services

import (
    "encoding/json"
    "log"
    "fmt"
    "time"
    "delivery-service/src/infrastructure/adapters"
    "delivery-service/src/domain/entities"
)

type DeliveryAlertUseCase interface {
    CreateDeliveryAlert(alert entities.DeliveryAlert) error
    GetAllDeliveryAlerts() ([]entities.DeliveryAlert, error)
}

type DeliveryAlertService struct {
    rabbitMQ    *adapters.RabbitMQ
    useCase     DeliveryAlertUseCase
    latestAlert *entities.DeliveryAlert
}

func NewDeliveryAlertService(rabbitMQ *adapters.RabbitMQ, useCase DeliveryAlertUseCase) *DeliveryAlertService {
    return &DeliveryAlertService{
        rabbitMQ: rabbitMQ,
        useCase:  useCase,
    }
}

func (s *DeliveryAlertService) ProcessDeliveryAlerts() {
    consumerTag := fmt.Sprintf("consumer-%d", time.Now().UnixNano())
    msgs, err := s.rabbitMQ.Consume(consumerTag)
    if err != nil {
        log.Fatalf("Error al consumir mensajes: %v", err)
    }

    for msg := range msgs {
        log.Printf("Orden Entregada: %s", msg.Body)
        
        var alert entities.DeliveryAlert
        if err := json.Unmarshal(msg.Body, &alert); err != nil {
            log.Printf("Error al deserializar mensaje: %v", err)
            continue
        }

        if err := s.useCase.CreateDeliveryAlert(alert); err != nil {
            log.Printf("No se pudo entregar la Orden: %v", err)
        } else {
            s.latestAlert = &alert
            log.Printf("Última alerta actualizada: %+v", s.latestAlert)

            message, err := json.Marshal(alert)
            if err != nil {
                log.Printf("Error al serializar alerta: %v", err)
                continue
            }

            err = s.rabbitMQ.PublishToQueue("Alert", string(message))
            if err != nil {
                log.Printf("Error al publicar alerta en la cola 'Alert': %v", err)
            }
        }
    }
}

func (s *DeliveryAlertService) GetLatestDeliveryAlert() (*entities.DeliveryAlert, error) {
    if s.latestAlert == nil {
        return nil, fmt.Errorf("No hay alertas recientes")
    }
    return s.latestAlert, nil
}

func (s *DeliveryAlertService) GetAllDeliveryAlerts() ([]entities.DeliveryAlert, error) {
    alerts, err := s.useCase.GetAllDeliveryAlerts()
    if err != nil {
        log.Printf("Error al obtener todas las alertas: %v", err)
        return nil, err
    }
    return alerts, nil
}
