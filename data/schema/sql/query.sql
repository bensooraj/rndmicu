-- name: CreateCreator :one
INSERT INTO creators (
  name, email
) VALUES (
  $1, $2
)
RETURNING *;

-- name: FindCreatorByID :one
SELECT * FROM creators WHERE id=$1 LIMIT 1;

-- name: ListCreators :many
SELECT * FROM creators
ORDER BY id
OFFSET $1 LIMIT $2;

-- name: FindAudioShortByID :one
SELECT * FROM audio_shorts WHERE id=$1 LIMIT 1;

-- name: ListAudioShorts :many
SELECT * FROM audio_shorts
ORDER BY date_created
OFFSET $1 LIMIT $2;




-- name: CreateAudioShort :one
INSERT INTO audio_shorts (
  id, creator_id, title, description, category, audio_file_url, date_created, date_updated
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: DeleteAudioShortByID :one
DELETE FROM audio_shorts WHERE id=$1 RETURNING *;