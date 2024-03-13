package args

import (
	"flag"
	"fmt"
	"strings"
)

type Args struct {
	Config      string
	Environment string
}

func Parse() (Args, error) {
	var args Args
	flag.StringVar(&args.Config, "c", "", "")
	flag.StringVar(&args.Environment, "e", "", "")
	flag.Parse()

	errs := []string{}
	if args.Config == "" {
		errs = append(errs, "-c config_file")
	}
	if args.Environment == "" {
		errs = append(errs, "-e env_file")
	}
	if len(errs) > 0 {
		return args, fmt.Errorf("flags are missing: %s", strings.Join(errs, ", "))
	}
	return args, nil
}
