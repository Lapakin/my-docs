package main

import (
	"github.com/lapotkin/file-storage/internal/adapter/db"
	"github.com/lapotkin/file-storage/internal/auth/jwt"
	"github.com/lapotkin/file-storage/internal/logging"
)

type config struct {
	Addr    string          `yaml:"addr"`
	DB      *db.Config      `yaml:"db"`
	JWT     *jwt.Config     `yaml:"jwt"`
	Logging *logging.Config `yaml:"logging"`
}
