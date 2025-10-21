package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/lapotkin/file-storage/internal/adapter/krakend"
	"github.com/lapotkin/file-storage/internal/domain/yaml"
	"github.com/lapotkin/file-storage/internal/logging"
	"github.com/lapotkin/file-storage/pkg/ui/router"
	"github.com/lapotkin/file-storage/pkg/ui/services"
)

const defaultConfigPath = "./cmd/ui/config.yaml"

func main() {
	var configPath string
	flag.StringVar(
		&configPath,
		"config-path",
		defaultConfigPath,
		"provides a path to configuration file with .yaml extension",
	)

	flag.Parse()

	var c config
	if err := yaml.UnmarshalYAML(configPath, &c); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	level, err := logging.ParseLevel(c.Logging.Level)
	if err != nil {
		log.Fatalf("failed to parse log level: %v", err)
	}

	logger := logging.NewLogger(level, c.Logging.Formatter).WithField("app", "ui")

	client := krakend.NewClient(c.KrakendURL)

	svc := services.NewServices(client)
	r := router.NewRouter(svc, logger)

	logger.Infof("Server is starting on %s", c.Addr)
	logger.Fatal(http.ListenAndServe(c.Addr, r))
}
