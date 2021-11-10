package service

import (
	"github.com/rysmaadit/go-template/app"
	"github.com/rysmaadit/go-template/external/jwt_client"
	"github.com/rysmaadit/go-template/external/mysql"
)

type Dependencies struct {
	AuthService     AuthServiceInterface
	CheckService    CheckService
	UserService     UserServiceInterface
	CustomerService CustomerServiceInterface
	EmployeeService EmployeeServiceInterface
}

func InstantiateDependencies(application *app.Application) Dependencies {
	jwtWrapper := jwt_client.New()
	authService := NewAuthService(application.Config, jwtWrapper)
	mysqlClient := mysql.NewMysqlClient(mysql.ClientConfig{
		Username: application.Config.DBUsername,
		Password: application.Config.DBPassword,
		Host:     application.Config.DBHost,
		Port:     application.Config.DBPort,
		DBName:   application.Config.DBName,
	})
	checkService := NewCheckService(mysqlClient)
	userService := NewUserService(application.Config, jwtWrapper)
	customerService := NewCustomerService(application.Config, jwtWrapper)
	employeeService := NewEmployeeService(application.Config, jwtWrapper)

	return Dependencies{
		AuthService:     authService,
		CheckService:    checkService,
		UserService:     userService,
		CustomerService: customerService,
		EmployeeService: employeeService,
	}
}
