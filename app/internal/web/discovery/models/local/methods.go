package local

import (
	"fmt"
	"logbook/models"
)

func (l Local) ServicePool(service models.Service) ([]string, error) {
	switch service {
	case models.Account:
		return l.Accounts, nil
	case models.Discovery:
		return l.Discovery, nil
	case models.Internal:
		return l.Internal, nil
	case models.Objectives:
		return l.Objectives, nil
	default:
		return nil, fmt.Errorf("unrecognized service name %q", service)
	}
}
