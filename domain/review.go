package domain

import "encoding/json"

// Review represents a user review attached to a reviewable entity.
type Review struct {
	ID             int    `db:"id" json:"id"`
	ReviewableType string `db:"reviewable_type" json:"reviewable_type"` // e.g. "product", "vendor", "user"
	ReviewableID   int    `db:"reviewable_id" json:"reviewable_id"`
	UserID         int    `db:"user_id" json:"user_id"` // author
	Rating         int    `db:"rating" json:"rating"`   // 1..5
	Title          string `db:"title" json:"title"`
	Body           string `db:"body" json:"body"`
	CreatedAt      string `db:"created_at" json:"created_at"`
	UpdatedAt      string `db:"updated_at" json:"updated_at"`
}

// ReviewPhoto represents a photo attached to a review.
type ReviewPhoto struct {
	ID        int             `db:"id" json:"id"`
	ReviewID  int             `db:"review_id" json:"review_id"`
	FilePath  string          `db:"file_path" json:"file_path"`   // path, URL or storage key
	Metadata  json.RawMessage `db:"metadata" json:"metadata"`     // optional JSON metadata (width/height/mime/etc)
	SortOrder int             `db:"sort_order" json:"sort_order"` // ordering within a review
	CreatedAt string          `db:"created_at" json:"created_at"`
}

// ReviewComment represents a comment left by a user on a review.
type ReviewComment struct {
	ID        int    `db:"id" json:"id"`
	ReviewID  int    `db:"review_id" json:"review_id"`
	UserID    int    `db:"user_id" json:"user_id"`
	Body      string `db:"body" json:"body"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

// DTOs used for HTTP binding / use-cases:

// CreateReviewRequest is the payload for creating a new review.
// Photos can be provided as a list of file paths / URLs (server-side should persist and map to ReviewPhoto).
type CreateReviewRequest struct {
	ReviewableType string            `json:"reviewable_type" binding:"required"`
	ReviewableID   int               `json:"reviewable_id" binding:"required"`
	Rating         int               `json:"rating" binding:"required"`
	Title          string            `json:"title,omitempty"`
	Body           string            `json:"body,omitempty"`
	PhotoPaths     []string          `json:"photo_paths,omitempty"` // list of uploaded file paths / storage keys
	PhotoMetadata  []json.RawMessage `json:"photo_metadata,omitempty"`
}

// CreateCommentRequest is the payload for creating a new comment on a review.
type CreateCommentRequest struct {
	ReviewID int    `json:"review_id" binding:"required"`
	Body     string `json:"body" binding:"required"`
}
