package repositorys

import (
    "delivery-service/src/infrastructure/adapters"
    "log"
)

type RabbitRepository struct {
    rabbitMQ *adapters.RabbitMQ
}

func NewRabbitRepository(rabbitMQ *adapters.RabbitMQ) *RabbitRepository {
    return &RabbitRepository{rabbitMQ: rabbitMQ}
}

func (repo *RabbitRepository) PublishDeliveryAlert(alert string) error {
    err := repo.rabbitMQ.Publish(alert)
    if err != nil {
        log.Printf("Failed to publish delivery alert: %v", err)
        return err
    }
    return nil
}

func (repo *RabbitRepository) ConsumeDeliveryAlerts() (<-chan string, error) {
    msgs, err := repo.rabbitMQ.Consume("delivery_alerts")
    if err != nil {
        log.Printf("Failed to consume delivery alerts: %v", err)
        return nil, err
    }

    alerts := make(chan string)
    go func() {
        for msg := range msgs {
            alerts <- string(msg.Body)
        }
        close(alerts)
    }()
    return alerts, nil
}