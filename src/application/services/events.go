package services

import (
    "log"
    "fmt"
    "time"
    "delivery-service/src/infrastructure/adapters"
    "delivery-service/src/domain/entities"
    "delivery-service/src/domain"
)

type DeliveryAlertService struct {
    rabbitMQ *adapters.RabbitMQ
    useCase  *DeliveryAlertUseCase
}

func NewDeliveryAlertService(rabbitMQ *adapters.RabbitMQ, useCase *DeliveryAlertUseCase) *DeliveryAlertService {
    return &DeliveryAlertService{
        rabbitMQ: rabbitMQ,
        useCase:  useCase,
    }
}

func (s *DeliveryAlertService) ProcessDeliveryAlerts() {
    consumerTag := fmt.Sprintf("consumer-%d", time.Now().UnixNano())
    msgs, err := s.rabbitMQ.Consume(consumerTag)
    if err != nil {
        log.Fatalf("Failed to consume messages: %v", err)
    }

    for msg := range msgs {
        log.Printf("Orden Entregada: %s", msg.Body)
        alert := entities.DeliveryAlert{
            Alert: string(msg.Body),
        }
        if err := s.useCase.CreateDeliveryAlert(alert); err != nil {
            log.Printf("No se pudo entregar la Orden: %v", err)
        }
    }
}

func (s *DeliveryAlertService) GetAllDeliveryAlerts() ([]entities.DeliveryAlert, error) {
    return s.useCase.GetAllDeliveryAlerts()
}

type DeliveryAlertUseCase struct {
    repo       domain.DeliveryAlertRepository
    rabbitRepo domain.RabbitRepository
}

func NewDeliveryAlertUseCase(repo domain.DeliveryAlertRepository, rabbitRepo domain.RabbitRepository) *DeliveryAlertUseCase {
    return &DeliveryAlertUseCase{
        repo:       repo,
        rabbitRepo: rabbitRepo,
    }
}

func (uc *DeliveryAlertUseCase) CreateDeliveryAlert(alert entities.DeliveryAlert) error {
    err := uc.repo.Create(&alert)
    if err != nil {
        return err
    }
    return uc.rabbitRepo.PublishDeliveryAlert(alert.Alert)
}

func (uc *DeliveryAlertUseCase) GetAllDeliveryAlerts() ([]entities.DeliveryAlert, error) {
    return uc.repo.GetAll()
}