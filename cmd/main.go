package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/ansakharov/lets_test/cmd/config"
	"github.com/ansakharov/lets_test/handler"
	"github.com/ansakharov/lets_test/logger"
	"github.com/ansakharov/lets_test/metrics"
	_ "github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func main() {
	// Get logger interface.
	log := logger.New()

	if err := mainNoExit(log); err != nil {
		log.Fatalf("fatal err: %s", err.Error())
	}
}

//
func mainNoExit(log logrus.FieldLogger) error {
	metrics.Init()

	// get application config
	confFlag := flag.String("conf", "", "config yaml file")
	flag.Parse()

	confString := *confFlag
	if confString == "" {
		return fmt.Errorf(" 'conf' flag required")
	}
	config, err := config.Parse(confString)
	if err != nil {
		return err
	}

	log.Println(config)
	log.Println("Starting the service...")

	ctx := context.Background()
	router, err := handler.Router(ctx, log, config)
	if err != nil {
		return fmt.Errorf("can't init router: %s", err.Error())
	}

	log.Print("The service is ready to listen and serve.")
	return http.ListenAndServe(
		config.AppPort, //port,
		router,
	)
}
