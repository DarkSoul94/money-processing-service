package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/DarkSoul94/money-processing-service/app"
	apphttp "github.com/DarkSoul94/money-processing-service/app/delivery/http"
	apprepo "github.com/DarkSoul94/money-processing-service/app/repo/postgresql"
	appusecase "github.com/DarkSoul94/money-processing-service/app/usecase"
	"github.com/DarkSoul94/money-processing-service/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // required
	"github.com/spf13/viper"
)

// App ...
type App struct {
	appUC      app.Usecase
	appRepo    app.Repository
	httpServer *http.Server
}

// NewApp ...
func NewApp(conf config.Config) *App {
	db := initDB(conf)

	repo := apprepo.NewPostgreRepo(db)
	uc := appusecase.NewUsecase(repo)
	return &App{
		appUC:   uc,
		appRepo: repo,
	}
}

// Run run application
func (a *App) Run(conf config.Config) {
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

	a.appRepo.Close()

	return a.httpServer.Shutdown(ctx)
}

func initDB(conf config.Config) *sql.DB {
	dbString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DbHost,
		conf.DbPort,
		conf.DBLogin,
		conf.DBPass,
		conf.DbName,
	)
	db, err := sql.Open(
		"postgres",
		dbString,
	)
	if err != nil {
		panic(err)
	}
	runMigrations(db, conf.DbName)
	return db
}

func runMigrations(db *sql.DB, dbName string) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		dbName,
		driver)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange && err != migrate.ErrNilVersion {
		fmt.Println(err)
	}
}
