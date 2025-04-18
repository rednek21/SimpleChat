package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rednek21/SimpleChat/config"
	"github.com/rednek21/SimpleChat/src/container"
	cors2 "github.com/rednek21/SimpleChat/src/transport/http/cors"
	"github.com/rednek21/SimpleChat/src/transport/http/routes"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Cors struct {
	Origins          []string
	Headers          []string
	AllowCredentials bool
}

type Server struct {
	host          string
	port          int
	cors          *Cors
	logger        *zap.Logger
	di            *container.Container
	originChecker cors2.OriginChecker
}

func NewServer(cfg *config.Chat, container *container.Container) *Server {
	originChecker := cors2.NewChecker(cfg.Http.Cors)

	c := &Cors{
		Origins:          cfg.Http.Cors.Origins,
		Headers:          cfg.Http.Cors.Headers,
		AllowCredentials: cfg.Http.Cors.AllowCredentials,
	}

	return &Server{
		host:          cfg.Http.Host,
		port:          cfg.Http.Port,
		cors:          c,
		logger:        container.Logger,
		di:            container,
		originChecker: originChecker,
	}
}

func (s *Server) Run(ctx context.Context) error {
	engine := s.setupEngine()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: engine,
	}

	errChan := make(chan error)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()
	select {
	case err := <-errChan:
		s.logger.Fatal("Failed to start server", zap.Error(err))
		return err
	case <-ctx.Done():
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	s.logger.Info("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Fatal("Server forced to shutdown", zap.Error(err))
		return err
	}
	return nil
}

func (s *Server) setupEngine() *gin.Engine {
	engine := gin.Default()

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = s.cors.Origins
	corsCfg.AllowCredentials = s.cors.AllowCredentials
	corsCfg.AllowHeaders = s.cors.Headers

	engine.Use(cors.New(corsCfg))

	routeGroups := routes.GetV1RouteGroups(
		s.di.Clients.LLMClient,
		s.di.Managers.ChatConnManager,
		s.originChecker,
		s.logger,
	)

	v1Router := routes.NewRoute(engine, "", s.logger)
	v1Router.SetupRoutes(routeGroups)

	return engine
}
