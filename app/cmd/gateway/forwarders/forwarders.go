package forwarders

import (
	servicereg "logbook/cmd/servicereg/client"
	"logbook/config/api"
	"logbook/config/deployment"
	"logbook/internal/args"
	"logbook/internal/web/balancer"
	"logbook/internal/web/discoveryfile"
	"logbook/models"
	"net/http"
	"path/filepath"
	"time"
)

// Forwarder forwards requests to an instance of the target service
// The service's address listed through discovery service
// Discovery service is accessed through internal gateway.
// Internal gateway might have multiple instances, addresses known by config-based service discovery
type Forwarders struct {
	// to stop tickers later
	internalstore   *discoveryfile.FileReader // config-based service discovery
	discoverystore  *servicereg.Discovery     // self-registration based service discovery client
	objectivesstore *servicereg.Discovery
}

// instances of services listed
func New(flags *args.GatewayArgs, deplcfg *deployment.Config, apicfg *api.Config) *Forwarders {
	internalstore := discoveryfile.NewFileReader(flags.Discovery, time.Second, discoveryfile.ServiceParams{
		Port: deplcfg.Ports.Internal,
		Tls:  true,
	})
	discoveryctl := servicereg.NewClient(
		apicfg,
		balancer.New(internalstore),
		filepath.Join(apicfg.Internal.Path, apicfg.Internal.Services.Discovery.Path),
	)

	return &Forwarders{
		internalstore:   internalstore,
		discoverystore:  servicereg.NewDiscoveryStore(discoveryctl, models.Discovery),
		objectivesstore: servicereg.NewDiscoveryStore(discoveryctl, models.Objectives),
	}
}

func (f *Forwarders) Stop() {
	f.internalstore.Stop()
	f.discoverystore.Stop()
	f.objectivesstore.Stop()
}

func (f *Forwarders) Objectives(w http.ResponseWriter, r *http.Request) {

}

func (f *Forwarders) Account(w http.ResponseWriter, r *http.Request) {

}
