package cfgs

import (
	"fmt"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/args"
	"logbook/internal/utilities/reflux"
	"logbook/internal/web/logger"
)

func Read() (*args.ServiceArgs, *deployment.Config, *api.Config, error) {
	flags, err := args.Service()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing args: %w", err)
	}
	l := logger.NewLogger("readConfigs")

	deplcfg, err := deployment.ReadConfig(flags.Deployment)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("reading deployment environment config: %w", err)
	}
	l.Println("deployment config:")
	reflux.Print(deplcfg)

	apicfg, err := api.ReadConfig(flags.Api)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("reading api config: %w", err)
	}
	l.Println("api config:")
	reflux.Print(apicfg)

	return &flags, deplcfg, apicfg, nil
}
