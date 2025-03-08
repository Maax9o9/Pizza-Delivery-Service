package application

import (
    "delivery-service/src/domain/entities"
    "delivery-service/src/domain"
    "delivery-service/src/application/repositorys"
)

type DeliveryAlertUseCase struct {
    repo       domain.DeliveryAlertRepository
    rabbitRepo *repositorys.RabbitRepository
}

func NewDeliveryAlertUseCase(repo domain.DeliveryAlertRepository, rabbitRepo *repositorys.RabbitRepository) *DeliveryAlertUseCase {
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