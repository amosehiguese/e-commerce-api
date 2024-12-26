package api

import (
	"github.com/amosehiguese/ecommerce-api/pkg/config"
	"github.com/amosehiguese/ecommerce-api/query"
)

type API struct {
	Q   query.Query
	Cfg *config.Config
}

func NewAPI(q query.Query, cfg *config.Config) API {
	return API{
		Q:   q,
		Cfg: cfg,
	}
}
