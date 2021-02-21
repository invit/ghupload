package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/google/go-github/v33/github"
	"github.com/invit/ghupload/internal/lib/client"
	"github.com/invit/ghupload/internal/lib/repourl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringP("branch", "b", "", "Commit to branch (default branch if empty)")
	uploadCmd.Flags().StringP("message", "m", "", "Commit message (required)")

	_ = uploadCmd.MarkFlagRequired("message")
}

var uploadCmd = &cobra.Command{
	Use:   "upload <local-path> <remote-url>",
	Short: "Uploads file to github repository",
	Args:  cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if os.Getenv("GITHUB_TOKEN") == "" {
			return errors.New("Missing GITHUB_TOKEN environment variable")
		}

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
		c, err := client.New(ctx, os.Getenv("GITHUB_TOKEN"))

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
		f, d, _, err := c.Repositories.GetContents(
			ctx,
			repo.Owner,
			repo.Repository,
			repo.Path,
			&github.RepositoryContentGetOptions{},
		)

		if err != nil {
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
