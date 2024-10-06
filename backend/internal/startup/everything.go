package startup

import (
	"fmt"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/utils/reflux"
)

func ApiGateway(l *logger.Logger) (*ApiGatewayArgs, *deployment.Config, *api.Config, error) {
	l = l.Sub("ApiGateway")
	args, err := parseApiGatewayArgs()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing args: %w", err)
	}

	deplcfg, err := deployment.ReadConfig(args.Deployment)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("reading deployment environment config: %w", err)
	}
	l.Println("deployment config:")
	reflux.Print(deplcfg)

	apicfg, err := api.ReadConfig(args.Api)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("reading api config: %w", err)
	}
	l.Println("api config:")
	reflux.Print(apicfg)

	return &args, deplcfg, apicfg, nil
}

func InternalGateway(l *logger.Logger) (*InternalGatewayArgs, *deployment.Config, *api.Config, error) {
	l = l.Sub("InternalGateway")

	args, err := parseInternalGatewayArgs()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing args: %w", err)
	}

	deplcfg, err := deployment.ReadConfig(args.Deployment)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("reading deployment environment config: %w", err)
	}
	l.Println("deployment config:")
	reflux.Print(deplcfg)

	apicfg, err := api.ReadConfig(args.Api)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("reading api config: %w", err)
	}
	l.Println("api config:")
	reflux.Print(apicfg)

	return &args, deplcfg, apicfg, nil
}

func Service() (*ServiceArgs, *deployment.Config, *api.Config, error) {
	args, err := parseServiceArgs()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing args: %w", err)
	}
	l := logger.New("readConfigs")

	deplcfg, err := deployment.ReadConfig(args.Deployment)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("reading deployment environment config: %w", err)
	}
	l.Println("deployment config:")
	reflux.Print(deplcfg)

	apicfg, err := api.ReadConfig(args.Api)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("reading api config: %w", err)
	}
	l.Println("api config:")
	reflux.Print(apicfg)

	return &args, deplcfg, apicfg, nil
}

type ServiceConfigReader[Config any] func(path string) (*Config, error)

// with custom service config
func ServiceWithCustomConfig[C any](serviceConfigReader ServiceConfigReader[C], l *logger.Logger) (*ServiceArgs, *C, *deployment.Config, *api.Config, error) {
	l = l.Sub("startup.ServiceWithCustomConfig")

	args, err := parseServiceArgs()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("parsing args: %w", err)
	}

	srvcfg, err := serviceConfigReader(args.Service)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("reading service config: %w", err)
	}
	l.Println("service config:")
	reflux.Print(srvcfg)

	deplcfg, err := deployment.ReadConfig(args.Deployment)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("reading deployment environment config: %w", err)
	}
	l.Println("deployment config:")
	reflux.Print(deplcfg)

	apicfg, err := api.ReadConfig(args.Api)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("reading api config: %w", err)
	}
	l.Println("api config:")
	reflux.Print(apicfg)

	return &args, srvcfg, deplcfg, apicfg, nil
}
