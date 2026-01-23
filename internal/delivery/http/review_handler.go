package httpDelivery

import (
	"net/http"
	"strconv"

	"eve/domain"
	"eve/internal/usecase"

	"github.com/labstack/echo/v4"
)

// ReviewHandler holds use-cases for reviews and comments.
type ReviewHandler struct {
	createReview  *usecase.CreateReviewUseCase
	createComment *usecase.CreateCommentUseCase
	listReviews   *usecase.ListReviewsUseCase
	getReview     *usecase.GetReviewUseCase
}

// NewReviewHandler constructs a ReviewHandler.
func NewReviewHandler(
	cr *usecase.CreateReviewUseCase,
	cc *usecase.CreateCommentUseCase,
	lr *usecase.ListReviewsUseCase,
	gr *usecase.GetReviewUseCase,
) *ReviewHandler {
	return &ReviewHandler{
		createReview:  cr,
		createComment: cc,
		listReviews:   lr,
		getReview:     gr,
	}
}

// CreateReview handles POST /reviews
// Expects JSON body matching domain.CreateReviewRequest and an "X-User-ID" header.
func (h *ReviewHandler) CreateReview(c echo.Context) error {
	var req domain.CreateReviewRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body: " + err.Error()})
	}

	userID, err := extractUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	id, err := h.createReview.Execute(req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]int{"id": id})
}

// CreateComment handles POST /reviews/comments
// Expects JSON body matching domain.CreateCommentRequest and an "X-User-ID" header.
func (h *ReviewHandler) CreateComment(c echo.Context) error {
	var req domain.CreateCommentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body: " + err.Error()})
	}

	userID, err := extractUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	id, err := h.createComment.Execute(req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]int{"id": id})
}

// ListReviews handles GET /reviews?reviewable_type=...&reviewable_id=...
func (h *ReviewHandler) ListReviews(c echo.Context) error {
	rt := c.QueryParam("reviewable_type")
	ridStr := c.QueryParam("reviewable_id")
	if rt == "" || ridStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "reviewable_type and reviewable_id query params are required"})
	}
	rid, err := strconv.Atoi(ridStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid reviewable_id"})
	}

	reviews, err := h.listReviews.Execute(rt, rid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, reviews)
}

// GetReview handles GET /reviews/:id
// Returns the review and its comments.
func (h *ReviewHandler) GetReview(c echo.Context) error {
	idStr := c.Param("id")
	if idStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id path parameter is required"})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	review, comments, err := h.getReview.Execute(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"review":   review,
		"comments": comments,
	})
}

// extractUserID reads the X-User-ID header and returns it as an int.
// If header missing or invalid, returns an error.
func extractUserID(c echo.Context) (int, error) {
	h := c.Request().Header.Get("X-User-ID")
	if h == "" {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "missing X-User-ID header")
	}
	id, err := strconv.Atoi(h)
	if err != nil || id <= 0 {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "invalid X-User-ID header")
	}
	return id, nil
}
