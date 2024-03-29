package client

import (
	"context"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

// New returns new GitHub client
func New(ctx context.Context, token string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), nil
}
