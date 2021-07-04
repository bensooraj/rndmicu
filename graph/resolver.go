package graph

import (
	"github.com/bensooraj/rndmicu/s3engine"
	"github.com/jmoiron/sqlx"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB         *sqlx.DB
	S3         *s3engine.S3Engine
	CdnBaseURL string
}
