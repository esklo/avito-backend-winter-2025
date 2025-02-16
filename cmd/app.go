package main

import (
	"context"
	"log"

	"github.com/esklo/avito-backend-winter-2025/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.New(ctx)
	if err != nil {
		log.Panicf("failed to init app: %s", err)
	}

	defer func() { _ = a.Shutdown() }()

	err = a.Run(ctx)
	if err != nil {
		log.Panicf("failed to run app: %s", err)
	}
}
