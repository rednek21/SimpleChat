package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rednek21/SimpleChat/src/managers"
	"github.com/rednek21/SimpleChat/src/transport/http/cors"
	"go.uber.org/zap"
	"net/http"
)

type Client interface {
	SendMessage(ctx context.Context, msg string) (string, error)
}

type Handler struct {
	logger        *zap.Logger
	client        Client
	connManager   *managers.ChatConnManager
	originChecker cors.OriginChecker
}

func NewChatHandler(
	chatManager Client,
	connManager *managers.ChatConnManager,
	logger *zap.Logger,
	checker cors.OriginChecker,
) *Handler {
	return &Handler{
		client:        chatManager,
		logger:        logger,
		connManager:   connManager,
		originChecker: checker,
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	publicEndpoints := r.Group("")
	{
		publicEndpoints.GET("/unicast", func(c *gin.Context) {
			h.Unicast(c)
		})
		publicEndpoints.GET("/multicast", func(c *gin.Context) {
			h.Multicast(c)
		})
	}
}

func (h *Handler) Unicast(c *gin.Context) {
	upgrader := h.getUpgrader()
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Upgrade failed", zap.Error(err))
		return
	}
	defer conn.Close()

	ip := conn.RemoteAddr().String()
	h.connManager.Add(ip, conn)
	defer h.connManager.Remove(ip)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			h.logger.Error("Read error", zap.Error(err))
			break
		}

		response, err := h.client.SendMessage(c.Request.Context(), string(msg))
		if err != nil {
			h.logger.Error("Process failed", zap.Error(err))
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
			h.logger.Error("Write error", zap.Error(err))
			break
		}
	}
}

func (h *Handler) Multicast(c *gin.Context) {
	upgrader := h.getUpgrader()
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Upgrade failed", zap.Error(err))
		return
	}
	defer conn.Close()

	ip := conn.RemoteAddr().String()
	h.connManager.Add(ip, conn)
	defer h.connManager.Remove(ip)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			h.logger.Error("Read error", zap.Error(err))
			break
		}

		for _, c := range h.connManager.GetConns() {
			if c.RemoteAddr().String() != ip {
				go func(conn *websocket.Conn) {
					if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
						h.logger.Error("Write error", zap.Error(err))
					}
				}(c)
			}
		}
	}
}

func (h *Handler) getUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			allowed := h.originChecker.IsAllowed(origin)
			if !allowed {
				h.logger.Warn("Origin rejected", zap.String("origin", origin))
			}
			return allowed
		},
	}
}
