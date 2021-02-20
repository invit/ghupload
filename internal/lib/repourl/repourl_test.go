package repourl_test

import (
	"testing"

	"github.com/invit/ghupload/internal/lib/repourl"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		want      *repourl.RepoURL
		wantError bool
	}{
		{
			"valid-https",
			"https://github.com/owner/repo/path/to/file",
			&repourl.RepoURL{
				Owner:      "owner",
				Repository: "repo",
				Path:       "path/to/file",
			},
			false,
		},
		{
			"valid-http",
			"http://github.com/owner/repo/path/to/file",
			&repourl.RepoURL{
				Owner:      "owner",
				Repository: "repo",
				Path:       "path/to/file",
			},
			false,
		},
		{
			"valid-ssh",
			"git@github.com:owner/repo.git/path/to/file",
			&repourl.RepoURL{
				Owner:      "owner",
				Repository: "repo",
				Path:       "path/to/file",
			},
			false,
		},
		{
			"valid-absolute-path",
			"/owner/repo/path/to/file",
			&repourl.RepoURL{
				Owner:      "owner",
				Repository: "repo",
				Path:       "path/to/file",
			},
			false,
		},
		{
			"valid-relative-path",
			"owner/repo/path/to/file",
			&repourl.RepoURL{
				Owner:      "owner",
				Repository: "repo",
				Path:       "path/to/file",
			},
			false,
		},
		{
			"missing-repo",
			"owner",
			nil,
			true,
		},
		{
			"missing-file",
			"owner/repo",
			nil,
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repoURL, err := repourl.Parse(test.url)

			if test.wantError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.EqualValues(t, test.want, repoURL)
		})
	}
}
