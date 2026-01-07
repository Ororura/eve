package main

import (
	httpDelivery "eve/internal/delivery/http"
	"eve/internal/repository/postgres"
	"eve/internal/usecase"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	db, _ := sqlx.Connect("postgres",
		"postgres://reviews_user:reviews_pass@localhost:5432/reviews_db?sslmode=disable")

	repo := postgres.NewUserRepo(db)

	createUC := usecase.NewCreateUserUseCase(repo)
	listUC := usecase.NewGetUserUseCase(repo)

	h := httpDelivery.NewHandler(createUC, listUC)

	e := echo.New()
	e.POST("/user", h.Create)
	e.GET("/user", h.List)

	log.Fatal(e.Start(":8080"))
}
