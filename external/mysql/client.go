package mysql

import (
	"fmt"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client interface {
	Ping() error
}

type client struct {
	DbConnection *gorm.DB
}

func (c *client) Ping() error {
	var result int64
	tx := c.DbConnection.Raw("select 1").Scan(&result)
	if tx.Error != nil {
		return fmt.Errorf("mysql unable to serve basic query. %v", tx.Error)
	}
	return nil
}

func NewMysqlClient(config ClientConfig) *client {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s",
		config.Username,
		config.Password,
		config.Host,
		config.DBName,
	)
	pgUrl, err := pq.ParseURL(connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbConn, err := gorm.Open(postgres.Open(pgUrl), &gorm.Config{
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("unable to initiate elephantsql connection. %v", err)
	}
	return &client{
		DbConnection: dbConn,
	}
}
