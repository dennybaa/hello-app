package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AppEnv type
type AppEnv int

// RouterFunc type
type RouterFunc func(router *gin.Engine)

const (
	DevApp AppEnv = 0 + iota
	ProdApp
)

// App struct to hold the main app objects
type App struct {
	client           *mongo.Client // mongo db client
	router           *gin.Engine   // http router
	appEnv           AppEnv        // prod/dev
	databaseURI      string
	databaseName     string
	apiPort          int
	dbConnectTimeout int
}

func (app *App) initConf() {
	app.databaseURI = "mongodb://localhost:27017"
	app.apiPort = 8080

	// set uri from env
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		app.databaseURI = uri
	}

	// set uri from env
	if dbName := os.Getenv("MONGODB_DATABASE"); dbName != "" {
		app.databaseName = dbName
	} else {
		app.databaseName = "hello"
	}

	// set db connection timeout in seconds
	if timeout := os.Getenv("MONGODB_CONN_TIMEOUT"); timeout != "" {
		var err error
		app.dbConnectTimeout, err = strconv.Atoi(timeout)

		if err != nil {
			log.Fatal("Unable to use MONGODB_CONN_TIMEOUT env variable, number is expected, got: ", timeout)
		}
	} else {
		app.dbConnectTimeout = 15
	}

	// set api port from env
	if port := os.Getenv("PORT"); port != "" {
		var err error
		app.apiPort, err = strconv.Atoi(port)

		if err != nil {
			log.Fatal("Unable to use PORT env variable, port number is expected, got: ", port)
		}
	}

	// set app_env
	if env := os.Getenv("APP_ENV"); env == "prod" || env == "production" {
		app.appEnv = ProdApp
	}
}

func (app *App) serve() {
	// Start http router
	if app.appEnv == ProdApp {
		gin.SetMode(gin.ReleaseMode)
	}

	app.router = gin.Default()
	app.defineRoutes()
	err := app.router.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func (app *App) dbConnect() {
	// make configurable
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.dbConnectTimeout)*time.Second)
	defer cancel()

	// mongo uri
	// ex: mongodb://hello:P%40ssw0rd@mongodb-headless.default.svc.cluster.local/hello?replicaSet=rs0
	clientOptions := options.Client().ApplyURI(app.databaseURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)

	select {
	case <-ctx.Done():
		log.Fatal("Connection to MongoDB timed out!")
	default:
		log.Println("Connected to MongoDB!")
	}

	app.client = client
}
