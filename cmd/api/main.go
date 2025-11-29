package main

import (
	"context"
	"log"

	"team-manager-leagues/internal/config"
	"team-manager-leagues/internal/repository"
	"team-manager-leagues/internal/service"
	transporthttp "team-manager-leagues/internal/transport/http"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := repository.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	store := repository.NewStore(pool)
	svc := service.NewLeaguesService(store)

	r := transporthttp.NewRouter(cfg, svc)

	log.Printf("Leagues service starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
