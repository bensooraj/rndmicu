package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/bensooraj/rndmicu/data/models"
	"github.com/bensooraj/rndmicu/graph/generated"
	"github.com/google/uuid"
)

func (r *queryResolver) Creator(ctx context.Context, id int) (*models.Creator, error) {
	q := models.New(r.DB)
	c, err := q.FindCreatorByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *queryResolver) Creators(ctx context.Context, limit *int, offset *int) ([]*models.Creator, error) {
	q := models.New(r.DB)
	cs, err := q.ListCreators(ctx, models.ListCreatorsParams{
		Limit:  int32(*limit),
		Offset: int32(*offset),
	})
	if err != nil {
		return nil, err
	}
	return starryNightC(&cs), nil
}

func (r *queryResolver) AudioShort(ctx context.Context, id string) (*models.AudioShort, error) {
	q := models.New(r.DB)
	a, err := q.FindAudioShortByID(ctx, uuid.MustParse(id))
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *queryResolver) AudioShorts(ctx context.Context, limit *int, offset *int) ([]*models.AudioShort, error) {
	q := models.New(r.DB)
	as, err := q.ListAudioShorts(ctx, models.ListAudioShortsParams{
		Limit:  int32(*limit),
		Offset: int32(*offset),
	})
	if err != nil {
		return nil, err
	}
	return starryNightA(&as), nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
