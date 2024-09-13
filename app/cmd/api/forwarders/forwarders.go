package forwarders

import (
	registry "logbook/cmd/registry/client"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/startup"
	"logbook/internal/web/balancer"
	"logbook/internal/web/forwarder"
	"logbook/internal/web/registryfile"
	"logbook/internal/web/sidecar"
	"logbook/models"
)

type Forwarders struct {
	internaldiscovery *registryfile.FileReader // config-based service discovery
	discoveryctl      *sidecar.Sidecar
	Accounts          *forwarder.LoadBalancedReverseProxy
	Objectives        *forwarder.LoadBalancedReverseProxy
}

func New(args *startup.ApiGatewayArgs, deplcfg *deployment.Config, apicfg *api.Config) (*Forwarders, error) {
	internalsd := registryfile.NewFileReader(args.InternalGateway, deplcfg.ServiceDiscovery.UpdatePeriod, registryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	})
	// NOTE: service registry needs to be accessed through internal gateway
	sc := sidecar.New(
		registry.NewClient(balancer.New(internalsd), apicfg, true),
		deplcfg.ServiceDiscovery.UpdatePeriod,
		[]models.Service{
			models.Account,
			models.Objectives,
		},
	)
	defer sc.Stop()

	return &Forwarders{
		discoveryctl:      sc,
		internaldiscovery: internalsd,

		Accounts:   forwarder.New(sc.InstanceSource(models.Account), models.Account, api.PathFromInternet(apicfg.Public.Services.Account)),
		Objectives: forwarder.New(sc.InstanceSource(models.Objectives), models.Objectives, api.PathFromInternet(apicfg.Public.Services.Objectives)),
	}, nil
}

func (f *Forwarders) Stop() {
	f.internaldiscovery.Stop()
	f.discoveryctl.Stop()
}
