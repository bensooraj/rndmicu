package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/bensooraj/rndmicu/data/models"
	"github.com/bensooraj/rndmicu/graph/generated"
)

func (r *audioShortResolver) ID(ctx context.Context, obj *models.AudioShort) (string, error) {
	return obj.ID.String(), nil
}

func (r *audioShortResolver) Creator(ctx context.Context, obj *models.AudioShort) (*models.Creator, error) {
	q := models.New(r.DB)
	c, err := q.FindCreatorByID(ctx, obj.CreatorID)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// AudioShort returns generated.AudioShortResolver implementation.
func (r *Resolver) AudioShort() generated.AudioShortResolver { return &audioShortResolver{r} }

type audioShortResolver struct{ *Resolver }
