package core

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
)

type Conn_MySQL struct {
	Db  *sql.DB
	Err string
}

func (c *Conn_MySQL) ExecuteQuery(query string, id int) (any, any) {
	panic("unimplemented")
}

var (
	instance *Conn_MySQL
	once     sync.Once
)

func GetDBPool() *Conn_MySQL {
	once.Do(func() {
		dsn := "root:nocbro123@tcp(54.87.5.19:3306)/pizza"
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Error al conectar con MySQL: %v", err)
		}

		if err = db.Ping(); err != nil {
			log.Fatalf("Error al hacer ping a MySQL: %v", err)
		}

		instance = &Conn_MySQL{
			Db:  db,
			Err: "",
		}
	})

	return instance
}

func (c *Conn_MySQL) ExecutePreparedQuery(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := c.Db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(args...)
}

func (c *Conn_MySQL) FetchRows(query string, args ...interface{}) *sql.Rows {
	rows, err := c.Db.Query(query, args...)
	if err != nil {
		log.Fatalf("Error al ejecutar la consulta: %v", err)
	}
	return rows
}
