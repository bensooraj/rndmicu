// Code generated by sqlc. DO NOT EDIT.

package models

import (
	"time"

	"github.com/google/uuid"
)

type AudioShort struct {
	ID           uuid.UUID `json:"id"`
	CreatorID    int32     `json:"creator_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	AudioFileUrl string    `json:"audio_file_url"`
	DateCreated  time.Time `json:"date_created"`
	DateUpdated  time.Time `json:"date_updated"`
}

type Creator struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
