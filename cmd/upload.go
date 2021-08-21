package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v33/github"
	"github.com/invit/ghupload/internal/lib/client"
	"github.com/invit/ghupload/internal/lib/repourl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringP("token", "t", "", fmt.Sprintf("File to read token from (default %s)", getTokenFilePath()))
	uploadCmd.Flags().StringP("branch", "b", "", "Commit to branch (default branch if empty)")
	uploadCmd.Flags().StringP("message", "m", "", "Commit message (required)")

	_ = uploadCmd.MarkFlagRequired("message")
}

func getTokenFilePath() string {
	c, err := os.UserConfigDir()

	if err != nil {
		return ""
	}

	return filepath.Join(c, "ghupload", "github-token")
}

func readToken(tokenFile string) (string, error) {
	if tokenFile != "" {
		t, err := os.ReadFile(getTokenFilePath())

		if err != nil {
			return "", err
		}

		return strings.TrimSpace(string(t)), nil
	}

	if t, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		return t, nil
	}

	t, err := os.ReadFile(getTokenFilePath())

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(t)), nil
}

var token = ""

var uploadCmd = &cobra.Command{
	Use:                   "upload -m <commit-msg> [-b <branch>] [-t <token-file>] <local-path> <remote-url>",
	DisableFlagsInUseLine: true,
	Short:                 "Uploads file to github repository",
	Long: "Uploads (commits) a local file to a github repository\n" +
		"\n" +
		"<local-path> is either a path to a local file or - for STDIN.\n" +
		"<remote-url> can be one of the following formats and has to include the repository owner, " +
		"the repository and the path to the file inside the repository:\n" +
		"* https://github.com/owner/repository/path/in/repo\n" +
		"* git@github.com:owner/repository.git/path/in/repo\n" +
		"* owner/repository/path/in/repo\n" +
		"\n" +
		"Command prints the commit SHA on success." +
		"\n\n" +
		"Authentication:\n" +
		"  Create personal access token on Github with read/write access to the repository and either provide\n" +
		"  * GITHUB_TOKEN environment variable with the token, or\n" +
		fmt.Sprintf("  * Store the token in the %s file, or\n", getTokenFilePath()) +
		"  * Provide a --token parameter pointing to a file containing the token.",
	Example: "* Upload local file\n" +
		"  $ ghupload upload -m \"commit msg\" README.md owner/repository/README.md\n" +
		"  b6cbb5b2ea041956c4ac8da17007f95d2312a461\n" +
		"* Upload data from STDIN\n" +
		"  $ ghupload upload -m \"commit msg\" - owner/repository/README.md\n" +
		"  this is the new \n" +
		"  content \n" +
		"  of the file\n" +
		"  ^D\n" +
		"  3be39e60c3ae44faa40f4efc31241f3564c396f1",
	Args: cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		tokenFile, _ := cmd.Flags().GetString("token")
		t, err := readToken(tokenFile)

		if err != nil {
			return fmt.Errorf("Missing/invalid authentication (%s)", err)
		}

		token = t

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var upload *os.File

		if args[0] == "-" {
			upload = os.Stdin
		} else {
			f, err := os.Open(args[0])
			defer func() { _ = f.Close() }()

			if err != nil {
				return err
			}

			upload = f
		}

		ctx := context.Background()
		c, err := client.New(ctx, token)

		if err != nil {
			return err
		}

		// parse remote url
		repo, err := repourl.Parse(args[1])

		if err != nil {
			return err
		}

		msg, _ := cmd.Flags().GetString("message")
		branch, _ := cmd.Flags().GetString("branch")

		opts := &github.RepositoryContentFileOptions{
			Message: github.String(msg),
		}

		if branch != "" {
			opts.Branch = github.String(branch)
		}

		// check if file exists in repo already
		f, d, resp, err := c.Repositories.GetContents(
			ctx,
			repo.Owner,
			repo.Repository,
			repo.Path,
			&github.RepositoryContentGetOptions{},
		)

		if resp == nil && err != nil {
			return err
		}

		if len(d) > 0 {
			return fmt.Errorf("Target expects path to file, %s is a directory", repo.Path)
		}

		if f != nil {
			opts.SHA = f.SHA
		}

		// attach file content
		content, err := io.ReadAll(upload)

		if err != nil {
			return err
		}

		opts.Content = content

		// upload file
		cr, _, err := c.Repositories.CreateFile(ctx, repo.Owner, repo.Repository, repo.Path, opts)

		if err != nil {
			return err
		}

		fmt.Println(*cr.SHA)

		return nil
	},
}
