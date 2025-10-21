package main

import (
	"flag"
	"log"

	"github.com/lapotkin/file-storage/internal/auth/jwt"
	"github.com/lapotkin/file-storage/internal/domain/yaml"
	"github.com/lapotkin/file-storage/internal/logging"
	"github.com/lapotkin/file-storage/pkg/user-management/repository/postgres"
	"github.com/lapotkin/file-storage/pkg/user-management/router"
	"github.com/lapotkin/file-storage/pkg/user-management/services"

	pg "github.com/lapotkin/file-storage/internal/adapter/db/postgres"
)

const defaultConfigPath = "./cmd/user-management/config.yaml"

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

	logger := logging.NewLogger(level, c.Logging.Formatter).WithField("app", "user-management")

	logger.Infoln("Trying to connect to database...")
	db, err := pg.NewDB(c.DB.URL(), c.DB.IsolationLevel)
	if err != nil {
		logger.Fatalf("Failed to connect to database. Err: %s", err.Error())
	}
	logger.Infoln("Successfully connected to database")

	tokenManager := jwt.NewTokenManager(
		c.JWT.SecretKey,
		c.JWT.AccessTokenTTL,
		c.JWT.RefreshTokenTTL,
	)
	svc := services.NewService(db, postgres.NewRepoManager(), tokenManager, logger)
	r := router.NewRouter(svc, tokenManager, logger)

	logger.Infof("Server is starting on %s", c.Addr)
	logger.Fatal(r.Run(c.Addr))
}
