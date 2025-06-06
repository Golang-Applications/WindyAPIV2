package api

import (
	"Windy-API/config"
	"Windy-API/persistence"
)

type Api interface {
	ProcessRequest()
}

func NewApi(cfg config.Config, dal persistence.Persistence) (Api, error) {
	a := &api{
		globalCfg: cfg,
		dal:       dal,
	}
	return a, nil
}

type api struct {
	globalCfg config.Config
	dal       persistence.Persistence
}

func (a *api) ProcessRequest() {
	a.processRequests()
}
