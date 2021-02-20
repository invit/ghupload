package repourl

import (
	"errors"
	"strings"

	giturls "github.com/whilp/git-urls"
)

// RepoURL contains repository address to a single file
type RepoURL struct {
	Owner      string
	Repository string
	Path       string
}

// Parse transforms a remote URL to a repository address
func Parse(repoURL string) (*RepoURL, error) {
	remote, err := giturls.Parse(repoURL)

	if err != nil {
		return nil, err
	}

	parts := strings.Split(strings.TrimPrefix(remote.Path, "/"), "/")

	if len(parts) < 3 {
		return nil, errors.New("Invalid remote URL (use OWNER/REPO/FILE)")
	}

	return &RepoURL{
		Owner:      parts[0],
		Repository: strings.TrimSuffix(parts[1], ".git"),
		Path:       strings.Join(parts[2:], "/"),
	}, nil
}
