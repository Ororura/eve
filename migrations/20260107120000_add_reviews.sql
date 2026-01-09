-- +goose Up
BEGIN;

-- Table: reviews
-- Generic polymorphic reviews table that can be attached to any "reviewable" entity
CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    reviewable_type TEXT NOT NULL,             -- e.g. 'product', 'vendor', 'user', etc.
    reviewable_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    title TEXT,
    body TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Indexes for efficient lookup
CREATE INDEX idx_reviews_reviewable ON reviews (reviewable_type, reviewable_id);
CREATE INDEX idx_reviews_user_id ON reviews (user_id);
CREATE INDEX idx_reviews_created_at ON reviews (created_at);

-- Trigger to auto-update `updated_at` on rows modification
CREATE OR REPLACE FUNCTION reviews_updated_at_trigger() RETURNS trigger AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_reviews_updated_at
BEFORE UPDATE ON reviews
FOR EACH ROW EXECUTE FUNCTION reviews_updated_at_trigger();

-- Table: review_photos
-- Photos attached to reviews. Store path/URL and optional JSON metadata.
CREATE TABLE review_photos (
    id SERIAL PRIMARY KEY,
    review_id INTEGER NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,       -- path or URL to the photo
    metadata JSONB,                -- optional metadata (width/height/mime/etc)
    sort_order INTEGER DEFAULT 0,  -- ordering of photos for a single review
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_review_photos_review_id ON review_photos (review_id);
CREATE INDEX idx_review_photos_review_id_sort ON review_photos (review_id, sort_order);

-- Table: review_comments
-- Comments left by other users on a review
CREATE TABLE review_comments (
    id SERIAL PRIMARY KEY,
    review_id INTEGER NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_review_comments_review_id ON review_comments (review_id);
CREATE INDEX idx_review_comments_user_id ON review_comments (user_id);
CREATE INDEX idx_review_comments_created_at ON review_comments (created_at);

-- Trigger to auto-update `updated_at` on comment edits
CREATE OR REPLACE FUNCTION review_comments_updated_at_trigger() RETURNS trigger AS $$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_review_comments_updated_at
BEFORE UPDATE ON review_comments
FOR EACH ROW EXECUTE FUNCTION review_comments_updated_at_trigger();

COMMIT;

-- +goose Down
BEGIN;

-- Remove triggers and functions first
DROP TRIGGER IF EXISTS trg_review_comments_updated_at ON review_comments;
DROP FUNCTION IF EXISTS review_comments_updated_at_trigger();

DROP TRIGGER IF EXISTS trg_reviews_updated_at ON reviews;
DROP FUNCTION IF EXISTS reviews_updated_at_trigger();

-- Drop tables (drop in order that respects FKs)
DROP TABLE IF EXISTS review_comments;
DROP TABLE IF EXISTS review_photos;
DROP TABLE IF EXISTS reviews;

COMMIT;
