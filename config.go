package main

import (
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Color   string
	Active  string
	Title   string
	Error   string
	Filters []string
	Upload  bool
	QrCode  bool
}

var configPaths = []string{
	".config/coordinate/coordinate.yaml", // prefered
	".config/coordinate/coordinate.yml",
	"coordinate.yaml",
	"coordinate.yml",
	"./coordinate.yaml",
	"./coordinate.yml",
}

func loadConfig() Config {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}

	for _, target := range configPaths {
		f, err := os.ReadFile(path.Join(home, target))
		if err != nil {
			continue
		}

		if err := yaml.Unmarshal(f, &config); err != nil {
			record("WARN", err.Error())
		}

		break
	}

	return config
}
