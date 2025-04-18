package container

import (
	"context"
	"fmt"
	"github.com/rednek21/MTSTechHack/chat/config"
	"github.com/rednek21/MTSTechHack/chat/src/managers"
	"github.com/rednek21/MTSTechHack/chat/src/transport/grpc/clients"
	"github.com/rednek21/go-toolkit/logger"
	"go.uber.org/zap"
)

type Clients struct {
	LLMClient *clients.AnyClient
}

type Managers struct {
	ChatConnManager *managers.ChatConnManager
}

type Container struct {
	Clients  Clients
	Managers Managers
	Logger   *zap.Logger
}

func NewContainer(ctx context.Context, cfg *config.Config) (*Container, error) {
	// logger initialization
	loggerCfg := &logger.Config{
		Level:      logger.INFO,
		Format:     logger.JSON,
		Output:     logger.STDOUT,
		FilePath:   cfg.Chat.LogFile,
		MaxSizeMB:  cfg.Logger.MaxSizeMB,
		MaxAgeDays: cfg.Logger.MaxAgeDays,
		MaxBackups: cfg.Logger.MaxBackups,
	}
	l, err := logger.New(loggerCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// ChatConnManager initialize
	chatConnManager := managers.NewChatConnManager()

	// LLMClient initialization
	llmClient := clients.NewAnyClient()

	return &Container{
		Logger: l,

		Clients: Clients{
			LLMClient: llmClient,
		},
		Managers: Managers{
			ChatConnManager: chatConnManager,
		},
	}, nil
}
