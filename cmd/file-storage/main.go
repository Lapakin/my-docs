package main

import (
	"flag"
	"log"

	"github.com/lapotkin/file-storage/internal/auth/jwt"
	"github.com/lapotkin/file-storage/internal/domain/yaml"
	"github.com/lapotkin/file-storage/internal/logging"
	"github.com/lapotkin/file-storage/pkg/file-storage/repository/minio"
	"github.com/lapotkin/file-storage/pkg/file-storage/repository/postgres"
	"github.com/lapotkin/file-storage/pkg/file-storage/router"
	"github.com/lapotkin/file-storage/pkg/file-storage/services"

	_ "github.com/jackc/pgx/v5/stdlib"

	pg "github.com/lapotkin/file-storage/internal/adapter/db/postgres"
	m "github.com/lapotkin/file-storage/internal/adapter/object/minio"
)

const defaultConfigPath = "./cmd/file-storage/config.yaml"

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

	logger := logging.NewLogger(level, c.Logging.Formatter).WithField("app", "file-storage")

	logger.Infoln("Trying to connect to database...")
	db, err := pg.NewDB(c.DB.URL(), c.DB.IsolationLevel)
	if err != nil {
		logger.Fatalf("Failed to connect to database. Err: %s", err.Error())
	}
	logger.Infoln("Successfully connected to database")

	logger.Infoln("Trying to connect to MinIO storage...")
	s3, err := m.NewMinIOClient(c.Minio)
	if err != nil {
		logger.Fatalf("Failed to create MinIO client: %s", err.Error())
	}
	logger.Infoln("Successfully connected to MinIO storage")

	tokenManager := jwt.NewTokenManager(
		c.JWT.SecretKey,
		c.JWT.AccessTokenTTL,
		c.JWT.RefreshTokenTTL,
	)

	svc := services.NewService(db, s3, postgres.NewRepoManager(), minio.NewRepoManager(), logger)
	r := router.NewRouter(svc, tokenManager, logger)

	logger.Infof("Server is starting on %s", c.Addr)
	logger.Fatal(r.Run(c.Addr))
}
