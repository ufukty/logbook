package startup

import (
	"flag"
)

type ServiceArgs struct {
	PrivateNetworkIp string
	Api              string
	Deployment       string
	Service          string
	InternalGateway  string
	EnvMode          string
	TlsKey           string
	TlsCertificate   string
}

func parseServiceArgs() (ServiceArgs, error) {
	var args ServiceArgs
	flag.StringVar(&args.PrivateNetworkIp,
		"ip", "", "Host's IP in the private network, which will be used to register the service into service registry")
	flag.StringVar(&args.EnvMode,
		"e", "", "either from [ local | stage | production ]")
	flag.StringVar(&args.Api,
		"api", "", "path to api config")
	flag.StringVar(&args.Deployment,
		"deployment", "", "path to deployment config")
	flag.StringVar(&args.Service,
		"service", "", "path to service config")
	flag.StringVar(&args.InternalGateway,
		"internal", "", "path to a JSON file which contains list of internal gateway instances")
	flag.StringVar(&args.TlsCertificate,
		"cert", "", "(optional) path to tls certificate")
	flag.StringVar(&args.TlsKey,
		"key", "", "(optional) path to tls key")
	flag.Parse()
	return args, nil
}

type ApiGatewayArgs struct {
	Api             string
	Deployment      string
	InternalGateway string
	EnvMode         string
	TlsKey          string
	TlsCertificate  string
}

func parseApiGatewayArgs() (ApiGatewayArgs, error) {
	var args ApiGatewayArgs
	flag.StringVar(&args.EnvMode,
		"e", "", "either from [ local | stage | production ]")
	flag.StringVar(&args.Api,
		"api", "", "api config file")
	flag.StringVar(&args.Deployment,
		"deployment", "", "path to deployment config")
	flag.StringVar(&args.InternalGateway,
		"internal", "", "path to a JSON file which contains list of internal gateway instances")
	flag.StringVar(&args.TlsCertificate,
		"cert", "", "(optional) path to tls certificate")
	flag.StringVar(&args.TlsKey,
		"key", "", "(optional) path to tls key")
	flag.Parse()
	return args, nil
}

type InternalGatewayArgs struct {
	Api             string
	Deployment      string
	RegistryService string
	EnvMode         string
	TlsKey          string
	TlsCertificate  string
}

func parseInternalGatewayArgs() (InternalGatewayArgs, error) {
	var args InternalGatewayArgs
	flag.StringVar(&args.EnvMode,
		"e", "", "either from [ local | stage | production ]")
	flag.StringVar(&args.Api,
		"api", "", "api config file")
	flag.StringVar(&args.Deployment,
		"deployment", "", "path to deployment config")
	flag.StringVar(&args.RegistryService,
		"registry", "", "path to a JSON file which contains list of registry service instances")
	flag.StringVar(&args.TlsCertificate,
		"cert", "", "(optional) path to tls certificate")
	flag.StringVar(&args.TlsKey,
		"key", "", "(optional) path to tls key")
	flag.Parse()
	return args, nil
}
