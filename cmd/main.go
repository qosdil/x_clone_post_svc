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
	config "x_clone_post_svc/config"

	"github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Set up MongoDB connection
	mongoURI := config.GetEnv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database(config.GetEnv("DB_NAME"))
	repo := app.NewMongoRepository(db)
	var (
		httpAddr = flag.String("http.addr", ":"+config.GetEnv("PORT"), "HTTP listen address")
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
