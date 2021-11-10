package service

import (
	"fmt"

	"github.com/rysmaadit/go-template/external/mysql"
	log "github.com/sirupsen/logrus"
)

type checkService struct {
	mysqlClient mysql.Client
}

type CheckService interface {
	CheckMysql() ([]byte, error)
}

func NewCheckService(mysqlClient mysql.Client) *checkService {
	return &checkService{
		mysqlClient: mysqlClient,
	}
}

func (c *checkService) CheckMysql() ([]byte, error) {
	err := c.mysqlClient.Ping()
	if err != nil {
		log.Warning(fmt.Errorf("mysql ping failed: %v", err))
		return nil, err
	}
	return []byte("Mysql OK"), err
}
