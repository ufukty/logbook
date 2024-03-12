package config

import (
	"fmt"
	"logbook/internal/utilities/reflux"
	"logbook/internal/utilities/strw"
	"logbook/internal/web/logger"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func Read(configpath string) Config {
	var l = logger.NewLogger("ConfigReader")

	f, err := os.Open(configpath)
	if err != nil {
		l.Fatalln(fmt.Errorf("could not open config file: %w", err))
	}
	config := Config{}
	err = yaml.NewDecoder(f).Decode(&config)
	if err != nil {
		l.Fatalln(fmt.Errorf("could not decode config file: %w", err))
	}

	errs := reflux.FindZeroValues(config)
	if len(errs) > 0 {
		l.Fatalf(
			"config has missing values:\n%s\n",
			strw.IndentLines(strings.Join(errs, "\n"), 4),
		)
	}

	for _, line := range strings.Split(reflux.String(config), "\n") {
		l.Println(line)
	}

	return config
}
