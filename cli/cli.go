package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rs/cors"
	"github.com/rysmaadit/go-template/app"
	"github.com/rysmaadit/go-template/contract"
	"github.com/rysmaadit/go-template/external/mysql"
	"github.com/rysmaadit/go-template/router"
	"github.com/rysmaadit/go-template/service"
	log "github.com/sirupsen/logrus"
)

type Cli struct {
	Args []string
}

func NewCli(args []string) *Cli {
	return &Cli{
		Args: args,
	}
}

func (c *Cli) Run(application *app.Application) {
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                                       // All origins
		AllowedMethods: []string{"GET", "OPTIONS", "POST", "PUT", "DELETE"}, // Allowing only get, just an example
	})

	log.SetLevel(log.InfoLevel)
	log.StandardLogger()
	log.SetOutput(os.Stdout)

	if strings.ToLower(application.Config.LogLevel) == log.DebugLevel.String() {
		log.SetLevel(log.DebugLevel)
	}

	log.SetReportCaller(true)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", application.Config.AppPort),
		Handler: cors.Handler(router.NewRouter(service.InstantiateDependencies(application))),
	}

	log.Println(fmt.Sprintf("starting application { %v } on port :%v", application.Config.AppName, application.Config.AppPort))

	db := mysql.NewMysqlClient(*mysql.MysqlInit())

	db.DbConnection.AutoMigrate(&contract.User{}, &contract.Identity{}, &contract.Submission{})

	go listenAndServe(srv)
	waitForShutdown(srv)
}

func listenAndServe(apiServer *http.Server) {
	err := apiServer.ListenAndServe()

	if err != nil {
		log.WithField("error", err.Error()).Fatal("unable to serve")
	}
}

func waitForShutdown(apiServer *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM)

	_ = <-sig

	log.Warn("shutting down")

	if err := apiServer.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}

	log.Warn("shutdown complete")
}
