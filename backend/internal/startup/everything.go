package startup

import (
	"fmt"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/logger"
	"logbook/internal/utils/reflux"
	"os"
	"path/filepath"
)

func ApiGateway(loggername string) (*logger.Logger, *ApiGatewayArgs, *deployment.Config, *api.Config, error) {
	args, err := parseApiGatewayArgs()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("parsing args: %w", err)
	}

	deplcfg, err := deployment.ReadConfig(args.Deployment)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("reading deployment environment config: %w", err)
	}
	l := logger.New(deplcfg, loggername)
	l.Println("deployment config:")
	reflux.Print(deplcfg)

	apicfg, err := api.ReadConfig(args.Api)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("reading api config: %w", err)
	}
	l.Println("api config:")
	reflux.Print(apicfg)

	return l, &args, deplcfg, apicfg, nil
}

func InternalGateway() (*logger.Logger, *InternalGatewayArgs, *deployment.Config, *api.Config, error) {

	args, err := parseInternalGatewayArgs()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("parsing args: %w", err)
	}

	deplcfg, err := deployment.ReadConfig(args.Deployment)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("reading deployment environment config: %w", err)
	}
	l := logger.New(deplcfg, "internal-gateway")
	l.Println("deployment config:")
	reflux.Print(deplcfg)

	apicfg, err := api.ReadConfig(args.Api)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("reading api config: %w", err)
	}
	l.Println("api config:")
	reflux.Print(apicfg)

	return l, &args, deplcfg, apicfg, nil
}

func Service(loggername string) (*logger.Logger, *ServiceArgs, *deployment.Config, *api.Config, error) {
	args, err := parseServiceArgs()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("parsing args: %w", err)
	}

	deplcfg, err := deployment.ReadConfig(args.Deployment)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("reading deployment environment config: %w", err)
	}
	l := logger.New(deplcfg, loggername)
	l.Println("deployment config:")
	reflux.Print(deplcfg)

	apicfg, err := api.ReadConfig(args.Api)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("reading api config: %w", err)
	}
	l.Println("api config:")
	reflux.Print(apicfg)

	return l, &args, deplcfg, apicfg, nil
}

type ServiceConfigReader[Config any] func(path string) (*Config, error)

// with custom service config
func ServiceWithCustomConfig[C any](loggername string, serviceConfigReader ServiceConfigReader[C]) (*logger.Logger, *ServiceArgs, *C, *deployment.Config, *api.Config, error) {
	args, err := parseServiceArgs()
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("parsing args: %w", err)
	}

	deplcfg, err := deployment.ReadConfig(args.Deployment)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("reading deployment environment config: %w", err)
	}
	l := logger.New(deplcfg, loggername)

	l.Println("deployment config:")
	reflux.Print(deplcfg)

	srvcfg, err := serviceConfigReader(args.Service)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("reading service config: %w", err)
	}
	l.Println("service config:")
	reflux.Print(srvcfg)

	apicfg, err := api.ReadConfig(args.Api)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("reading api config: %w", err)
	}
	l.Println("api config:")
	reflux.Print(apicfg)

	return l, &args, srvcfg, deplcfg, apicfg, nil
}

func TestDependencies() (*logger.Logger, *deployment.Config, *api.Config, error) {
	workspace := os.Getenv("WORKSPACE")

	deplcfg, err := deployment.ReadConfig(filepath.Join(workspace, "platform/local/deployment.yml"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("deployment.ReadConfig: %w", err)
	}

	l := logger.New(deplcfg, "test")

	apicfg, err := api.ReadConfig(filepath.Join(workspace, "backend/api.yml"))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("api.ReadConfig: %w", err)
	}

	return l, deplcfg, apicfg, nil
}

func TestDependenciesWithServiceConfig[C any](servicename string, serviceConfigReader ServiceConfigReader[C]) (*logger.Logger, *C, *deployment.Config, *api.Config, error) {
	workspace := os.Getenv("WORKSPACE")

	deplcfg, err := deployment.ReadConfig(filepath.Join(workspace, "platform/local/deployment.yml"))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("deployment.ReadConfig: %w", err)
	}

	l := logger.New(deplcfg, "test")

	apicfg, err := api.ReadConfig(filepath.Join(workspace, "backend/api.yml"))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("api.ReadConfig: %w", err)
	}

	srvcnf, err := serviceConfigReader(filepath.Join(workspace, "backend/cmd", servicename, "local.yml"))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("reading service config: %w", err)
	}

	return l, srvcnf, deplcfg, apicfg, nil
}
