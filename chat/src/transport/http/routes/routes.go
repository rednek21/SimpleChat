package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Route struct {
	engine    *gin.Engine
	logger    *zap.Logger
	apiPrefix string
}

func NewRoute(engine *gin.Engine, apiPrefix string, logger *zap.Logger) *Route {
	return &Route{
		engine:    engine,
		logger:    logger,
		apiPrefix: apiPrefix,
	}
}

func (r *Route) SetupRoutes(routeGroups []RouteGroup) {
	baseGroup := r.engine.Group(r.apiPrefix)
	baseGroup.Use()

	for _, rg := range routeGroups {
		group := baseGroup.Group(rg.Prefix)
		rg.Router.RegisterRoutes(group)
	}
}
