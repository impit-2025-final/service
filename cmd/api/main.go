package main

import (
	"log"
	"net/http"
	"service/internal/config"
	"service/internal/repository"
	"service/internal/router"
	"service/internal/usecase"
)

func main() {
	conf := config.LoadDBConfig()
	db, err := repository.InitDB(conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	repo := repository.NewGormRepository(db)

	useCase := usecase.NewContainerUseCase(repo)
	handler := router.NewHandler(useCase)

	mux := router.SetupRoutes(handler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("error: %v", err)
	}
}
