-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id) 
  VALUES ($1, $2, $3, $4, $5) RETURNING *
)
SELECT ff.*, u.name AS user_name, f.name AS feed_name, f.url AS feed_url
FROM inserted_feed_follow ff
INNER JOIN users u ON ff.user_id = u.id
INNER JOIN feeds f ON ff.feed_id = f.id;

-- name: GetFeedFollowsByUserID :many
SELECT ff.*, u.name AS user_name, f.name AS feed_name, f.url AS feed_url
FROM feed_follows ff
INNER JOIN users u ON ff.user_id = u.id
INNER JOIN feeds f ON ff.feed_id = f.id
WHERE ff.user_id = $1
ORDER BY ff.created_at DESC;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;