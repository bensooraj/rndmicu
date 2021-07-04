package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/bensooraj/rndmicu/data/models"
	"github.com/bensooraj/rndmicu/graph/generated"
	"github.com/bensooraj/rndmicu/graph/model"
	"github.com/bensooraj/rndmicu/s3engine"
	"github.com/google/uuid"
)

func (r *mutationResolver) AudioshortCreate(ctx context.Context, audioshort model.NewAudioShort) (*models.AudioShort, error) {
	log.Println("audioshort: ", audioshort.AudioFile.ContentType)
	log.Println("audioshort: ", audioshort.AudioFile.Filename)
	log.Println("audioshort: ", ByteCountSI(audioshort.AudioFile.Size))

	audioID := uuid.New()
	fileName := audioID.String() + filepath.Ext(audioshort.AudioFile.Filename)
	keyName := fmt.Sprintf("%s/%d", "short", audioshort.CreatorID)

	r.S3.QueueUploadJob(&s3engine.AudioFileUploadJob{
		File:        audioshort.AudioFile.File,
		Filename:    fileName,
		Size:        audioshort.AudioFile.Size,
		ContentType: audioshort.AudioFile.ContentType,
		KeyName:     keyName,
	})

	q := models.New(r.DB)
	a, err := q.CreateAudioShort(ctx, models.CreateAudioShortParams{
		ID:           audioID,
		CreatorID:    int32(audioshort.CreatorID),
		Title:        audioshort.Title,
		Description:  audioshort.Description,
		Category:     audioshort.Category,
		AudioFileUrl: fmt.Sprintf("%s/%s/%s", r.CdnBaseURL, keyName, fileName),
		DateCreated:  time.Now(),
		DateUpdated:  time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
