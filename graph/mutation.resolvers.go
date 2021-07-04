package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"
	"time"

	"github.com/bensooraj/rndmicu/data/models"
	"github.com/bensooraj/rndmicu/graph/generated"
	"github.com/bensooraj/rndmicu/graph/model"
	"github.com/google/uuid"
)

func (r *mutationResolver) AudioshortCreate(ctx context.Context, audioshort model.NewAudioShort) (*models.AudioShort, error) {
	log.Println("audioshort: ", audioshort.AudioFile.ContentType)
	log.Println("audioshort: ", audioshort.AudioFile.Filename)
	log.Println("audioshort: ", ByteCountSI(audioshort.AudioFile.Size))
	q := models.New(r.DB)
	a, err := q.CreateAudioShort(ctx, models.CreateAudioShortParams{
		ID:           uuid.New(),
		CreatorID:    int32(audioshort.CreatorID),
		Title:        audioshort.Title,
		Description:  audioshort.Description,
		Category:     audioshort.Category,
		AudioFileUrl: audioshort.AudioFileURL,
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
