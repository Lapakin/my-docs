package main

import (
	"github.com/lapotkin/file-storage/internal/logging"
)

type config struct {
	Addr          string          `yaml:"addr"`
	KrakendURL    string          `yaml:"krakend_url"`
	SessionSecret string          `yaml:"session_secret"`
	Logging       *logging.Config `yaml:"logging"`
}
