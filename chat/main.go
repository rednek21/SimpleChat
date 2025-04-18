package main

import (
	"context"
	"github.com/rednek21/MTSTechHack/chat/config"
	"github.com/rednek21/MTSTechHack/chat/src/container"
	"github.com/rednek21/MTSTechHack/chat/src/transport/http"
	c "github.com/rednek21/go-toolkit/config"
	"log"
)

func main() {
	// Global context
	ctx := context.Background()

	// Reading config
	var cfg *config.Config
	loader := c.NewLoaderConfig("")
	if err := loader.Load(&cfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", cfg)

	// Initialize DI-container
	diContainer, err := container.NewContainer(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer diContainer.Logger.Sync()

	// Creating & running server
	s := http.NewServer(&cfg.Chat, diContainer)
	if err := s.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
