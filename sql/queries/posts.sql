-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUser :many
SELECT * FROM posts p
WHERE p.feed_id IN (
	SELECT ff.feed_id FROM feed_follows ff
	WHERE ff.user_id = $1
)
ORDER BY p.published_at DESC
LIMIT $2;
