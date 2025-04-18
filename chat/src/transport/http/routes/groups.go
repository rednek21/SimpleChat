package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rednek21/SimpleChat/src/managers"
	"github.com/rednek21/SimpleChat/src/transport/grpc/clients"
	"github.com/rednek21/SimpleChat/src/transport/http/cors"
	"github.com/rednek21/SimpleChat/src/transport/http/handler"
	"go.uber.org/zap"
)

type RouterRegistrar interface {
	RegisterRoutes(group *gin.RouterGroup)
}

type RouteGroup struct {
	Prefix string
	Router RouterRegistrar
}

func GetV1RouteGroups(
	client *clients.AnyClient,
	manager *managers.ChatConnManager,
	checker cors.OriginChecker,
	logger *zap.Logger,
) []RouteGroup {
	return []RouteGroup{
		{
			Prefix: "/chat",
			Router: handler.NewChatHandler(client, manager, logger, checker),
		},
	}
}
