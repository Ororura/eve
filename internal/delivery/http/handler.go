package httpDelivery

import (
	"eve/domain"
	"eve/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	create *usecase.CreateUserUseCase
	getAll *usecase.GetUserUseCase
}

func NewHandler(cu *usecase.CreateUserUseCase, gu *usecase.GetUserUseCase) *Handler {
	return &Handler{
		create: cu,
		getAll: gu,
	}
}

func (h *Handler) Create(c echo.Context) error {
	var r domain.User
	if err := c.Bind(&r); err != nil {
		return err
	}
	if err := h.create.Execute(r); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) List(c echo.Context) error {
	data, _ := h.getAll.Execute()
	return c.JSON(http.StatusOK, data)
}
