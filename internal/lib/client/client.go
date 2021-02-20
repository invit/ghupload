package client

import (
	"context"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

func New(ctx context.Context, token string) (*github.Client, error){
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), nil
}
