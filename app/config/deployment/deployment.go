package deployment

import (
	"fmt"
	"logbook/internal/utilities/reflux"
	"logbook/internal/utilities/strw"
	"logbook/internal/web/logger"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type RouterParameters struct {
	RequestTimeout time.Duration `yaml:"request-timeout"`
	WriteTimeout   time.Duration `yaml:"write-timeout"`
	ReadTimeout    time.Duration `yaml:"read-timeout"`
	IdleTimeout    time.Duration `yaml:"idle-timeout"`
	GracePeriod    time.Duration `yaml:"grace-period"`
}

type Common struct {
	RouterParameters             `yaml:",inline"`
	ServiceDiscoveryConfig       string        `yaml:"service-discovery-config"`
	ServiceDiscoveryUpdatePeriod time.Duration `yaml:"service-discovery-update-period"`
	RouterPublic                 string        `yaml:"router-public"`
	RouterPrivate                string        `yaml:"router-private"`
}

type APIGateway struct {
	Common `yaml:",inline"`
}

type Captcha struct {
	Common `yaml:",inline"`
}

type Tasks struct {
	Common `yaml:",inline"`
}

type Config struct {
	Common     Common     `yaml:"common"`
	APIGateway APIGateway `yaml:"api_gateway"`
	Captcha    Captcha    `yaml:"captcha"`
	Tasks      Tasks      `yaml:"tasks"`
}

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
