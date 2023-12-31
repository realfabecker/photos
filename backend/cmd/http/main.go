package main

import (
	"log"

	"github.com/realfabecker/photos/internal/adapters"
	cordom "github.com/realfabecker/photos/internal/core/domain"
	corpts "github.com/realfabecker/photos/internal/core/ports"
)

func main() {
	if err := container.Container.Invoke(func(
		app corpts.HttpHandler,
		walletConfig *cordom.Config,
	) error {
		if err := app.Register(); err != nil {
			return err
		}
		return app.Listen(walletConfig.AppPort)
	}); err != nil {
		log.Fatalln(err)
	}
}
