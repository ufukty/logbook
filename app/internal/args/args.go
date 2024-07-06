package args

import (
	"flag"
)

type ServiceArgs struct {
	Api            string // config file path
	Deployment     string // config file path
	Service        string // config file path
	EnvMode        string // local, stage, production
	TlsKey         string // file path
	TlsCertificate string // file path
}

func Service() (ServiceArgs, error) {
	var args ServiceArgs
	flag.StringVar(&args.EnvMode,
		"e", "", "either from [ local | stage | production ]")
	flag.StringVar(&args.Api,
		"api", "", "path to api config")
	flag.StringVar(&args.Deployment,
		"deployment", "", "path to deployment config")
	flag.StringVar(&args.Service,
		"service", "", "path to service config")
	flag.Parse()
	return args, nil
}

type GatewayArgs struct {
	Api            string
	Deployment     string
	Discovery      string
	EnvMode        string
	TlsKey         string
	TlsCertificate string
}

func Gateway() (GatewayArgs, error) {
	var args GatewayArgs
	flag.StringVar(&args.EnvMode,
		"e", "", "either from [ local | stage | production ]")
	flag.StringVar(&args.Api,
		"api", "", "api config file")
	flag.StringVar(&args.Deployment,
		"deployment", "", "path to deployment config")
	flag.StringVar(&args.Discovery,
		"discovery", "", "path to service config")
	flag.StringVar(&args.TlsCertificate,
		"cert", "", "(optional) path to tls certificate")
	flag.StringVar(&args.TlsKey,
		"key", "", "(optional) path to tls key")
	flag.Parse()
	return args, nil
}
