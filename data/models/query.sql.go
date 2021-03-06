// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createAudioShort = `-- name: CreateAudioShort :one
INSERT INTO audio_shorts (
  id, creator_id, title, description, category, audio_file_url, date_created, date_updated
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, creator_id, title, description, category, audio_file_url, date_created, date_updated
`

type CreateAudioShortParams struct {
	ID           uuid.UUID `json:"id"`
	CreatorID    int32     `json:"creator_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	AudioFileUrl string    `json:"audio_file_url"`
	DateCreated  time.Time `json:"date_created"`
	DateUpdated  time.Time `json:"date_updated"`
}

func (q *Queries) CreateAudioShort(ctx context.Context, arg CreateAudioShortParams) (AudioShort, error) {
	row := q.db.QueryRowContext(ctx, createAudioShort,
		arg.ID,
		arg.CreatorID,
		arg.Title,
		arg.Description,
		arg.Category,
		arg.AudioFileUrl,
		arg.DateCreated,
		arg.DateUpdated,
	)
	var i AudioShort
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.AudioFileUrl,
		&i.DateCreated,
		&i.DateUpdated,
	)
	return i, err
}

const createCreator = `-- name: CreateCreator :one
INSERT INTO creators (
  name, email
) VALUES (
  $1, $2
)
RETURNING id, name, email
`

type CreateCreatorParams struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (q *Queries) CreateCreator(ctx context.Context, arg CreateCreatorParams) (Creator, error) {
	row := q.db.QueryRowContext(ctx, createCreator, arg.Name, arg.Email)
	var i Creator
	err := row.Scan(&i.ID, &i.Name, &i.Email)
	return i, err
}

const deleteAudioShortByID = `-- name: DeleteAudioShortByID :one
DELETE FROM audio_shorts WHERE id=$1 RETURNING id, creator_id, title, description, category, audio_file_url, date_created, date_updated
`

func (q *Queries) DeleteAudioShortByID(ctx context.Context, id uuid.UUID) (AudioShort, error) {
	row := q.db.QueryRowContext(ctx, deleteAudioShortByID, id)
	var i AudioShort
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.AudioFileUrl,
		&i.DateCreated,
		&i.DateUpdated,
	)
	return i, err
}

const findAudioShortByID = `-- name: FindAudioShortByID :one
SELECT id, creator_id, title, description, category, audio_file_url, date_created, date_updated FROM audio_shorts WHERE id=$1 LIMIT 1
`

func (q *Queries) FindAudioShortByID(ctx context.Context, id uuid.UUID) (AudioShort, error) {
	row := q.db.QueryRowContext(ctx, findAudioShortByID, id)
	var i AudioShort
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.AudioFileUrl,
		&i.DateCreated,
		&i.DateUpdated,
	)
	return i, err
}

const findCreatorByID = `-- name: FindCreatorByID :one
SELECT id, name, email FROM creators WHERE id=$1 LIMIT 1
`

func (q *Queries) FindCreatorByID(ctx context.Context, id int32) (Creator, error) {
	row := q.db.QueryRowContext(ctx, findCreatorByID, id)
	var i Creator
	err := row.Scan(&i.ID, &i.Name, &i.Email)
	return i, err
}

const listAudioShorts = `-- name: ListAudioShorts :many
SELECT id, creator_id, title, description, category, audio_file_url, date_created, date_updated FROM audio_shorts
ORDER BY date_created
OFFSET $1 LIMIT $2
`

type ListAudioShortsParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

func (q *Queries) ListAudioShorts(ctx context.Context, arg ListAudioShortsParams) ([]AudioShort, error) {
	rows, err := q.db.QueryContext(ctx, listAudioShorts, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AudioShort
	for rows.Next() {
		var i AudioShort
		if err := rows.Scan(
			&i.ID,
			&i.CreatorID,
			&i.Title,
			&i.Description,
			&i.Category,
			&i.AudioFileUrl,
			&i.DateCreated,
			&i.DateUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listCreators = `-- name: ListCreators :many
SELECT id, name, email FROM creators
ORDER BY id
OFFSET $1 LIMIT $2
`

type ListCreatorsParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

func (q *Queries) ListCreators(ctx context.Context, arg ListCreatorsParams) ([]Creator, error) {
	rows, err := q.db.QueryContext(ctx, listCreators, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Creator
	for rows.Next() {
		var i Creator
		if err := rows.Scan(&i.ID, &i.Name, &i.Email); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAudioShortByID = `-- name: UpdateAudioShortByID :one
UPDATE audio_shorts
SET title = $1,
    description = $2,
    category = $3,
    date_updated = $4
WHERE id=$5 RETURNING id, creator_id, title, description, category, audio_file_url, date_created, date_updated
`

type UpdateAudioShortByIDParams struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	DateUpdated time.Time `json:"date_updated"`
	ID          uuid.UUID `json:"id"`
}

func (q *Queries) UpdateAudioShortByID(ctx context.Context, arg UpdateAudioShortByIDParams) (AudioShort, error) {
	row := q.db.QueryRowContext(ctx, updateAudioShortByID,
		arg.Title,
		arg.Description,
		arg.Category,
		arg.DateUpdated,
		arg.ID,
	)
	var i AudioShort
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.AudioFileUrl,
		&i.DateCreated,
		&i.DateUpdated,
	)
	return i, err
}
