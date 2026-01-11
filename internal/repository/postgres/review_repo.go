package postgres

import (
	"eve/domain"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// ReviewRepository defines storage operations for reviews, photos and comments.
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

	// DeleteReview deletes a review by id.
	DeleteReview(id int) error
}

// ReviewRepo is a Postgres implementation of ReviewRepository.
type ReviewRepo struct {
	db *sqlx.DB
}

func NewReviewRepo(db *sqlx.DB) *ReviewRepo {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) Create(review domain.Review) (int, error) {
	var id int
	query := `
		INSERT INTO reviews (reviewable_type, reviewable_id, user_id, rating, title, body)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.Get(&id, query,
		review.ReviewableType,
		review.ReviewableID,
		review.UserID,
		review.Rating,
		review.Title,
		review.Body,
	)
	if err != nil {
		return 0, fmt.Errorf("insert review: %w", err)
	}
	return id, nil
}

func (r *ReviewRepo) AddPhotos(reviewID int, photos []domain.ReviewPhoto) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		// If still in transaction and not committed, rollback.
		_ = tx.Rollback()
	}()

	stmt := `
		INSERT INTO review_photos (review_id, file_path, metadata, sort_order)
		VALUES ($1, $2, $3, $4)
	`
	for _, p := range photos {
		var metadata interface{}
		if len(p.Metadata) == 0 {
			metadata = nil
		} else {
			metadata = p.Metadata
		}
		if _, err := tx.Exec(stmt, reviewID, p.FilePath, metadata, p.SortOrder); err != nil {
			return fmt.Errorf("insert photo: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit photos tx: %w", err)
	}
	return nil
}

func (r *ReviewRepo) AddComment(comment domain.ReviewComment) (int, error) {
	var id int
	query := `
		INSERT INTO review_comments (review_id, user_id, body)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.db.Get(&id, query, comment.ReviewID, comment.UserID, comment.Body)
	if err != nil {
		return 0, fmt.Errorf("insert comment: %w", err)
	}
	return id, nil
}

func (r *ReviewRepo) GetByID(id int) (domain.Review, error) {
	var review domain.Review
	query := `
		SELECT id, reviewable_type, reviewable_id, user_id, rating, title, body, created_at, updated_at
		FROM reviews
		WHERE id = $1
	`
	if err := r.db.Get(&review, query, id); err != nil {
		return domain.Review{}, fmt.Errorf("get review by id: %w", err)
	}
	return review, nil
}

func (r *ReviewRepo) ListByReviewable(reviewableType string, reviewableID int) ([]domain.Review, error) {
	var reviews []domain.Review
	query := `
		SELECT id, reviewable_type, reviewable_id, user_id, rating, title, body, created_at, updated_at
		FROM reviews
		WHERE reviewable_type = $1 AND reviewable_id = $2
		ORDER BY created_at DESC
	`
	if err := r.db.Select(&reviews, query, reviewableType, reviewableID); err != nil {
		return nil, fmt.Errorf("list reviews by reviewable: %w", err)
	}
	return reviews, nil
}

func (r *ReviewRepo) ListComments(reviewID int) ([]domain.ReviewComment, error) {
	var comments []domain.ReviewComment
	query := `
		SELECT id, review_id, user_id, body, created_at, updated_at
		FROM review_comments
		WHERE review_id = $1
		ORDER BY created_at ASC
	`
	if err := r.db.Select(&comments, query, reviewID); err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	return comments, nil
}

func (r *ReviewRepo) DeleteReview(id int) error {
	_, err := r.db.Exec("DELETE FROM reviews WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete review: %w", err)
	}
	return nil
}
