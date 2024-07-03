package args

import (
	"flag"
	"fmt"
	"strings"
)

type Args struct {
	Api            string // config file path
	Deployment     string // config file path
	Service        string // config file path
	TlsKey         string // file path
	TlsCertificate string // file path
}

func WithServiceConfig() (Args, error) {
	var args Args
	flag.StringVar(&args.Api, "a", "", "-a <api config file>")
	flag.StringVar(&args.Deployment, "d", "", "-d <deployment config file>")
	flag.StringVar(&args.Service, "s", "", "-s <service config file>")
	flag.StringVar(&args.TlsCertificate, "cert", "", "(optional) path to tls certificate")
	flag.StringVar(&args.TlsKey, "key", "", "(optional) path to tls key")
	flag.Parse()

	errs := []string{}
	if args.Api == "" {
		errs = append(errs, "-a <api config file>")
	}
	if args.Deployment == "" {
		errs = append(errs, "-d <deployment config file>")
	}
	if args.Service == "" {
		errs = append(errs, "-s <service config file>")
	}
	if len(errs) > 0 {
		return args, fmt.Errorf("flags are missing: %s", strings.Join(errs, ", "))
	}
	return args, nil
}

func Read() (Args, error) {
	var args Args
	flag.StringVar(&args.Api, "a", "", "-a <api config file>")
	flag.StringVar(&args.Deployment, "d", "", "-d <deployment config file>")
	flag.StringVar(&args.TlsCertificate, "cert", "", "(optional) path to tls certificate")
	flag.StringVar(&args.TlsKey, "key", "", "(optional) path to tls key")
	flag.Parse()

	errs := []string{}
	if args.Api == "" {
		errs = append(errs, "-a <api config file>")
	}
	if args.Deployment == "" {
		errs = append(errs, "-d <deployment config file>")
	}
	if len(errs) > 0 {
		return args, fmt.Errorf("flags are missing: %s", strings.Join(errs, ", "))
	}
	return args, nil
}
