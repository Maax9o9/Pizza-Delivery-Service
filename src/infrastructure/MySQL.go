package infrastructure

import (
    "delivery-service/src/core"
    "delivery-service/src/domain"
    "delivery-service/src/domain/entities"
    "log"
)

type MySQL struct {
    conn *core.Conn_MySQL
}

func NewMySQL() domain.DeliveryAlertRepository {
    conn := core.GetDBPool()
    if conn.Err != "" {
        log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
    }

    return &MySQL{conn: conn}
}

func (mysql *MySQL) Create(alert *entities.DeliveryAlert) error {
    query := "INSERT INTO delivery_alerts (alert) VALUES (?)"
    _, err := mysql.conn.ExecutePreparedQuery(query, alert.Alert)
    if err != nil {
        log.Printf("Error al insertar delivery alert: %v", err)
        return err
    }
    return nil
}

func (mysql *MySQL) GetAll() ([]entities.DeliveryAlert, error) {
    query := "SELECT id, alert FROM delivery_alerts ORDER BY id DESC"
    rows := mysql.conn.FetchRows(query)
    defer rows.Close()

    var alerts []entities.DeliveryAlert
    for rows.Next() {
        var alert entities.DeliveryAlert
        if err := rows.Scan(&alert.ID, &alert.Alert); err != nil {
            return nil, err
        }
        alerts = append(alerts, alert)
    }
    return alerts, nil
}