package startup

import (
	"fmt"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/cliargs"
	"logbook/internal/utilities/reflux"
	"logbook/internal/web/logger"
)

func Everything() (*cliargs.ApiGatewayArgs, *deployment.Config, *api.Config, error) {
	args, err := cliargs.ApiGateway()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing args: %w", err)
	}
	l := logger.NewLogger("readConfigs")

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
