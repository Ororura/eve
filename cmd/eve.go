package main

import (
	httpDelivery "eve/internal/delivery/http"
	"eve/internal/infrastructure"
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

	hasher := infrastructure.NewBcryptHasher()
	repo := postgres.NewUserRepo(db)

	createUC := usecase.NewCreateUserUseCase(repo, hasher)
	listUC := usecase.NewGetUserUseCase(repo)

	h := httpDelivery.NewHandler(createUC, listUC)

	// --- Reviews wiring ---
	reviewRepo := postgres.NewReviewRepo(db)

	createReviewUC := usecase.NewCreateReviewUseCase(reviewRepo)
	createCommentUC := usecase.NewCreateCommentUseCase(reviewRepo)
	listReviewsUC := usecase.NewListReviewsUseCase(reviewRepo)
	getReviewUC := usecase.NewGetReviewUseCase(reviewRepo)

	reviewHandler := httpDelivery.NewReviewHandler(createReviewUC, createCommentUC, listReviewsUC, getReviewUC)
	// -----------------------

	e := echo.New()
	e.POST("/user", h.Create)
	e.GET("/user", h.List)

	// Review endpoints
	e.POST("/reviews", reviewHandler.CreateReview)
	e.POST("/reviews/comments", reviewHandler.CreateComment)
	e.POST("/reviews/:id/comments", reviewHandler.CreateComment)
	e.GET("/reviews", reviewHandler.ListReviews)
	e.GET("/reviews/:id", reviewHandler.GetReview)

	log.Fatal(e.Start(":8080"))
}
