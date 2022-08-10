package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/DarkSoul94/money-processing-service/app"
	apphttp "github.com/DarkSoul94/money-processing-service/app/delivery/http"
	apprepo "github.com/DarkSoul94/money-processing-service/app/repo/mysql"
	appusecase "github.com/DarkSoul94/money-processing-service/app/usecase"
	"github.com/DarkSoul94/money-processing-service/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/DarkSoul94/dbconnectors/mysql"
)

// App ...
type App struct {
	appUC      app.Usecase
	appRepo    app.Repository
	httpServer *http.Server
}

// NewApp ...
func NewApp(conf config.Config) *App {
	db, err := mysql.InitMysqlDB(
		conf.DBLogin,
		conf.DBPass,
		conf.DbHost,
		conf.DbPort,
		conf.DbName,
		conf.DBArgs,
		"file://migrations",
	)
	if err != nil {
		panic(err)
	}

	repo := apprepo.NewMySQLRepo(db)
	uc := appusecase.NewUsecase(repo)
	return &App{
		appUC:   uc,
		appRepo: repo,
	}
}

// Run run application
func (a *App) Run(conf config.Config) {
	defer a.appRepo.Close()

	router := gin.New()
	if viper.GetBool("app.release") {
		gin.SetMode(gin.ReleaseMode)
	} else {
		router.Use(gin.Logger())
	}

	apiRouter := router.Group("/api")
	apphttp.RegisterHTTPEndpoints(apiRouter, a.appUC)

	a.httpServer = &http.Server{
		Addr:           ":" + conf.HTTPport,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	var l net.Listener
	var err error
	l, err = net.Listen("tcp", a.httpServer.Addr)
	if err != nil {
		panic(err)
	}

	if err := a.httpServer.Serve(l); err != nil {
		log.Fatalf("Failed to listen and serve: %+v", err)
	}
}

func (a *App) Stop() error {
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
