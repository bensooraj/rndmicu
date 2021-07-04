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

	r.S3.QueueJob(&s3engine.AudioFileJob{
		TaskType:    s3engine.UPLOAD_TASK,
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

func (r *mutationResolver) AudioshortUpdate(ctx context.Context, id string, audioshort model.UpdatedAudioShort) (*models.AudioShort, error) {
	q := models.New(r.DB)
	existingAS, err := q.FindAudioShortByID(ctx, uuid.MustParse(id))
	if err != nil {
		return nil, fmt.Errorf("The Audio Short ID %s doesn't exist.", id)
	}

	updateReq := models.UpdateAudioShortByIDParams{
		Title:       existingAS.Title,
		Description: existingAS.Description,
		Category:    existingAS.Category,
		DateUpdated: time.Now(),
		ID:          uuid.MustParse(id),
	}
	if audioshort.Title != nil {
		updateReq.Title = *audioshort.Title
	}
	if audioshort.Description != nil {
		updateReq.Description = *audioshort.Description
	}
	if audioshort.Category != nil {
		updateReq.Category = *audioshort.Category
	}

	newAS, err := q.UpdateAudioShortByID(ctx, updateReq)
	if err != nil {
		return nil, fmt.Errorf("There was an error updating the audio shortID.")
	}

	return &newAS, nil
}

func (r *mutationResolver) AudioshortDelete(ctx context.Context, id string) (*models.AudioShort, error) {
	q := models.New(r.DB)
	// Check if the audio short exists
	as, err := q.DeleteAudioShortByID(ctx, uuid.MustParse(id))
	if err != nil {
		return nil, fmt.Errorf("The Audio Short ID %s doesn't exist.", id)
	}
	r.S3.QueueJob(&s3engine.AudioFileJob{
		TaskType: s3engine.DELETE_TASK,
		Filename: filepath.Base(as.AudioFileUrl),
		KeyName:  fmt.Sprintf("%s/%d", "short", as.CreatorID),
	})
	return &as, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
