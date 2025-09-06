-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
	INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING *
)
SELECT
	inserted_feed_follow.*,
	feeds.name as feed_name,
	users.name as user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT
    ff.*,
    f.name AS feed_name,
    u.name AS user_name
FROM feed_follows ff
INNER JOIN users u ON ff.user_id = u.id
INNER JOIN feeds f ON ff.feed_id = f.id
WHERE u.name = $1;
