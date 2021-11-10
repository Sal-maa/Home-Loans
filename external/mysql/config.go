package mysql

import "github.com/rysmaadit/go-template/config"

type ClientConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	DBName   string
}

func MysqlInit() *ClientConfig {
	mysqlConfig := &ClientConfig{
		Username: config.Init().DBUsername,
		Password: config.Init().DBPassword,
		Host:     config.Init().DBHost,
		Port:     config.Init().DBPort,
		DBName:   config.Init().DBName,
	}
	return mysqlConfig
}
