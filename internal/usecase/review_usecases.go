package usecase

import (
	"fmt"

	"eve/domain"
)

// ReviewRepository defines the methods the use-cases expect from a persistence layer.
// Implementations live in internal/repository (for example a Postgres implementation).
type ReviewRepository interface {
	// Create inserts a new review and returns its generated ID.
	Create(review domain.Review) (int, error)

	// AddPhotos attaches photos to an existing review.
	AddPhotos(reviewID int, photos []domain.ReviewPhoto) error

	// AddComment inserts a comment for a review and returns the comment ID.
	AddComment(comment domain.ReviewComment) (int, error)

	// GetByID loads a single review by ID.
	GetByID(id int) (domain.Review, error)

	// ListByReviewable returns reviews for a specific reviewable entity.
	ListByReviewable(reviewableType string, reviewableID int) ([]domain.Review, error)

	// ListComments returns comments for a review.
	ListComments(reviewID int) ([]domain.ReviewComment, error)
}

// CreateReviewUseCase handles the creation of reviews and optional photos.
type CreateReviewUseCase struct {
	repo ReviewRepository
}

// NewCreateReviewUseCase constructs a new CreateReviewUseCase.
func NewCreateReviewUseCase(r ReviewRepository) *CreateReviewUseCase {
	return &CreateReviewUseCase{repo: r}
}

// Execute creates a review authored by authorID from the given request.
// Returns the created review ID.
func (uc *CreateReviewUseCase) Execute(req domain.CreateReviewRequest, authorID int) (int, error) {
	// Basic validation
	if req.ReviewableType == "" {
		return 0, fmt.Errorf("reviewable_type is required")
	}
	if req.ReviewableID == 0 {
		return 0, fmt.Errorf("reviewable_id is required")
	}
	if req.Rating < 1 || req.Rating > 5 {
		return 0, fmt.Errorf("rating must be between 1 and 5")
	}

	rev := domain.Review{
		ReviewableType: req.ReviewableType,
		ReviewableID:   req.ReviewableID,
		UserID:         authorID,
		Rating:         req.Rating,
		Title:          req.Title,
		Body:           req.Body,
	}

	id, err := uc.repo.Create(rev)
	if err != nil {
		return 0, fmt.Errorf("create review: %w", err)
	}

	// Attach photos if provided
	if len(req.PhotoPaths) > 0 {
		photos := make([]domain.ReviewPhoto, 0, len(req.PhotoPaths))
		for i, p := range req.PhotoPaths {
			var meta []byte
			if len(req.PhotoMetadata) > i {
				meta = req.PhotoMetadata[i]
			}
			photos = append(photos, domain.ReviewPhoto{
				ReviewID:  id,
				FilePath:  p,
				Metadata:  meta,
				SortOrder: i,
			})
		}
		if err := uc.repo.AddPhotos(id, photos); err != nil {
			return id, fmt.Errorf("add photos: %w", err)
		}
	}

	return id, nil
}

// CreateCommentUseCase handles adding comments to reviews.
type CreateCommentUseCase struct {
	repo ReviewRepository
}

// NewCreateCommentUseCase constructs a new CreateCommentUseCase.
func NewCreateCommentUseCase(r ReviewRepository) *CreateCommentUseCase {
	return &CreateCommentUseCase{repo: r}
}

// Execute creates a comment authored by authorID on the specified review.
func (uc *CreateCommentUseCase) Execute(req domain.CreateCommentRequest, authorID int) (int, error) {
	if req.ReviewID == 0 {
		return 0, fmt.Errorf("review_id is required")
	}
	if req.Body == "" {
		return 0, fmt.Errorf("body is required")
	}

	c := domain.ReviewComment{
		ReviewID: req.ReviewID,
		UserID:   authorID,
		Body:     req.Body,
	}

	id, err := uc.repo.AddComment(c)
	if err != nil {
		return 0, fmt.Errorf("add comment: %w", err)
	}
	return id, nil
}

// ListReviewsUseCase returns reviews for a given reviewable entity.
type ListReviewsUseCase struct {
	repo ReviewRepository
}

// NewListReviewsUseCase constructs a new ListReviewsUseCase.
func NewListReviewsUseCase(r ReviewRepository) *ListReviewsUseCase {
	return &ListReviewsUseCase{repo: r}
}

// Execute returns reviews for the provided reviewable identifier.
func (uc *ListReviewsUseCase) Execute(reviewableType string, reviewableID int) ([]domain.Review, error) {
	if reviewableType == "" || reviewableID == 0 {
		return nil, fmt.Errorf("reviewable_type and reviewable_id are required")
	}
	return uc.repo.ListByReviewable(reviewableType, reviewableID)
}

// GetReviewUseCase loads a single review together with its comments.
type GetReviewUseCase struct {
	repo ReviewRepository
}

// NewGetReviewUseCase constructs a new GetReviewUseCase.
func NewGetReviewUseCase(r ReviewRepository) *GetReviewUseCase {
	return &GetReviewUseCase{repo: r}
}

// Execute returns the review and its comments.
func (uc *GetReviewUseCase) Execute(reviewID int) (domain.Review, []domain.ReviewComment, error) {
	if reviewID == 0 {
		return domain.Review{}, nil, fmt.Errorf("review id is required")
	}

	review, err := uc.repo.GetByID(reviewID)
	if err != nil {
		return domain.Review{}, nil, fmt.Errorf("get review: %w", err)
	}

	comments, err := uc.repo.ListComments(reviewID)
	if err != nil {
		return review, nil, fmt.Errorf("list comments: %w", err)
	}

	return review, comments, nil
}
