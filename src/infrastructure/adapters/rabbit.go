package adapters

import (
    "log"
    amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    queue   amqp.Queue
}

func NewRabbitMQ(url, queueName string) (*RabbitMQ, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }

    channel, err := conn.Channel()
    if err != nil {
        return nil, err
    }

    queue, err := channel.QueueDeclare(
        queueName,
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return nil, err
    }

    return &RabbitMQ{
        conn:    conn,
        channel: channel,
        queue:   queue,
    }, nil
}

func (r *RabbitMQ) Publish(message string) error {
    err := r.channel.Publish(
        "",
        r.queue.Name,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        []byte(message),
        },
    )
    if err != nil {
        log.Printf("Failed to publish message: %v", err)
        return err
    }
    return nil
}

func (r *RabbitMQ) PublishToQueue(queueName, message string) error {
    _, err := r.channel.QueueDeclare(
        queueName,
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return err
    }

    err = r.channel.Publish(
        "",
        queueName,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        []byte(message),
        },
    )
    if err != nil {
        log.Printf("Failed to publish message to queue %s: %v", queueName, err)
        return err
    }
    return nil
}

func (r *RabbitMQ) Consume(consumerTag string) (<-chan amqp.Delivery, error) {
    msgs, err := r.channel.Consume(
        r.queue.Name,
        consumerTag,
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return nil, err
    }
    return msgs, nil
}

func (r *RabbitMQ) Close() {
    r.channel.Close()
    r.conn.Close()
}