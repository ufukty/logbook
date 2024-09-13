package startup

import (
	"fmt"
	"logbook/cmd/account/service"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/cliargs"
	"logbook/internal/utilities/reflux"
	"logbook/internal/web/logger"
)

func Everything() (*cliargs.ServiceArgs, *service.Config, *deployment.Config, *api.Config, error) {
	args, err := cliargs.Service()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("parsing args: %w", err)
	}
	l := logger.NewLogger("readConfigs")

	srvcfg, err := service.ReadConfig(args.Service)
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
