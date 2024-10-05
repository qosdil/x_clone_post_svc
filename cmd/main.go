package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	app "x_clone_post_svc"
	configs "x_clone_post_svc/configs"
	db "x_clone_post_svc/databases"

	"github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
)

func getDbRepo(repo string) app.Repository {
	switch repo {
	case "mongo":
		// Set up MongoDB connection
		mongoURI := configs.GetEnv("MONGODB_URI")
		clientOptions := options.Client().ApplyURI(mongoURI)
		var err error
		mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			panic(err)
		}
		mongoDB := mongoClient.Database(configs.GetEnv("DB_NAME"))
		return db.NewMongoRepository(mongoDB)
	default:
		return db.NewDummyRepository()
	}
}

func main() {
	// Load environment variables
	configs.LoadEnv()

	dbRepo := configs.GetEnv("DB_REPO")
	if dbRepo == "mongo" {
		defer mongoClient.Disconnect(context.TODO())
	}
	repo := getDbRepo(dbRepo)

	var (
		httpAddr = flag.String("http.addr", ":"+configs.GetEnv("PORT"), "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s app.Service
	{
		s = app.NewService(repo)
		s = app.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = app.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
