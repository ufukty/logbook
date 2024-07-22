package forwarders

import (
	servicereg "logbook/cmd/registry/client"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/args"
	"logbook/internal/web/balancer"
	"logbook/internal/web/discoveryctl"
	"logbook/internal/web/discoveryfile"
	"logbook/internal/web/forwarder"
	"logbook/models"
)

type Forwarders struct {
	internaldiscovery *discoveryfile.FileReader // config-based service discovery
	discoveryctl      *discoveryctl.Client
	Accounts          *forwarder.LoadBalancedReverseProxy
	Objectives        *forwarder.LoadBalancedReverseProxy
}

func New(flags *args.GatewayArgs, deplcfg *deployment.Config, apicfg *api.Config) (*Forwarders, error) {
	internaldiscovery := discoveryfile.NewFileReader(flags.Discovery, deplcfg.ServiceDiscovery.UpdatePeriod, discoveryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	})
	// NOTE: service registry is needs to be accessed through internal gateway
	discovery := discoveryctl.New(
		servicereg.NewClient(balancer.New(internaldiscovery), apicfg, true),
		deplcfg.ServiceDiscovery.UpdatePeriod,
		[]models.Service{
			models.Account,
			models.Objectives,
		},
	)

	return &Forwarders{
		discoveryctl:      discovery,
		internaldiscovery: internaldiscovery,

		Accounts:   forwarder.New(discovery.InstanceSource(models.Account), models.Account, api.PathFromInternet(apicfg.Public.Services.Account)),
		Objectives: forwarder.New(discovery.InstanceSource(models.Objectives), models.Objectives, api.PathFromInternet(apicfg.Public.Services.Objectives)),
	}, nil
}

func (f *Forwarders) Stop() {
	f.internaldiscovery.Stop()
	f.discoveryctl.Stop()
}
