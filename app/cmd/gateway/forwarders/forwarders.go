package forwarders

import (
	"fmt"
	servicereg "logbook/cmd/registry/client"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/args"
	"logbook/internal/web/balancer"
	"logbook/internal/web/discoveryctl"
	"logbook/internal/web/discoveryfile"
	"logbook/internal/web/forwarder"
	"logbook/models"
	"path/filepath"
	"time"
)

type Forwarders struct {
	internaldiscovery *discoveryfile.FileReader // config-based service discovery
	discoveryctl      *discoveryctl.Client
	Accounts          *forwarder.LoadBalancedReverseProxy
	Objectives        *forwarder.LoadBalancedReverseProxy
}

func New(flags *args.GatewayArgs, deplcfg *deployment.Config, apicfg *api.Config) (*Forwarders, error) {
	internaldiscovery := discoveryfile.NewFileReader(flags.Discovery, time.Second, discoveryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	})
	// NOTE: service registry is accesed over internal gateway
	discoveryctl := discoveryctl.New(servicereg.NewClient(
		apicfg,
		balancer.New(internaldiscovery),
		filepath.Join(apicfg.Internal.Path, apicfg.Internal.Services.Discovery.Path),
	), []models.Service{
		models.Account,
		models.Objectives,
	})

	accountsfwd, err := forwarder.New(discoveryctl.InstanceSource(models.Account), models.Account, api.PathFromInternet(apicfg.Public.Services.Account))
	if err != nil {
		return nil, fmt.Errorf("creating forwarder for accounts service: %w", err)
	}

	objectivesfwd, err := forwarder.New(discoveryctl.InstanceSource(models.Objectives), models.Objectives, api.PathFromInternet(apicfg.Public.Services.Objectives))
	if err != nil {
		return nil, fmt.Errorf("creating forwarder for objectives service: %w", err)
	}

	return &Forwarders{
		discoveryctl:      discoveryctl,
		internaldiscovery: internaldiscovery,
		Accounts:          accountsfwd,
		Objectives:        objectivesfwd,
	}, nil
}

func (f *Forwarders) Stop() {
	f.internaldiscovery.Stop()
	f.discoveryctl.Stop()
}
