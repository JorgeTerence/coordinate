package main

import (
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Title   string
	Color   string
	Alt     string
	Upload  bool
	QrCode  bool
	Filters []string
}

func loadConfig() (config *Config) {
	home, err := os.UserHomeDir()
	if err != nil {
		record("WARN", err.Error())
	}

	f, err := os.ReadFile(path.Join(home, ".config/coordinate/coordinate.yaml"))
	if err != nil {
		record("WARN", err.Error())
	}

	if err := yaml.Unmarshal(f, &config); err != nil {
		record("WARN", err.Error())
	}

	return
}
