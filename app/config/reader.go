package config

import (
	"flag"
	"fmt"
	"logbook/internal/utilities/reflux"
	"logbook/internal/utilities/strw"
	"logbook/internal/web/logger"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func getConfigPath() string {
	var configpath string
	flag.StringVar(&configpath, "config", "", "")
	flag.Parse()
	return configpath
}

func Read() Config {
	var log = logger.NewLogger("ConfigReader")

	configpath := getConfigPath()
	log.Println("Using config file:", configpath)

	f, err := os.Open(configpath)
	if err != nil {
		log.Fatalln(fmt.Errorf("could not open config file: %w", err))
	}
	config := Config{}
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		log.Fatalln(fmt.Errorf("could not decode config file: %w", err))
	}

	errs := reflux.FindZeroValues(config)
	if len(errs) > 0 {
		log.Fatalf(
			"config has missing values:\n%s\n",
			strw.IndentLines(strings.Join(errs, "\n"), 4),
		)
	}
	return config
}
