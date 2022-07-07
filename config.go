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
	config = &Config{"", "", "", true, true, []string{}}

	home, err := os.UserHomeDir()
	if err != nil {
		record("WARN", err.Error())
		return
	}

	f, err := os.ReadFile(path.Join(home, ".config/coordinate/coordinate.yaml"))
	if err != nil {
		record("WARN", err.Error())
		return
	}

	if err := yaml.Unmarshal(f, &config); err != nil {
		record("WARN", err.Error())
	}

	return
}
