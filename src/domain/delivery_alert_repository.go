package domain

import "delivery-service/src/domain/entities"

type DeliveryAlertRepository interface {
    Create(alert *entities.DeliveryAlert) error
    GetAll() ([]entities.DeliveryAlert, error)
}

type RabbitRepository interface {
    PublishDeliveryAlert(alert string) error
}