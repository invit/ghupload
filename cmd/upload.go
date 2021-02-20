package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-github/v33/github"
	"github.com/invit/ghupload/internal/lib/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringP("owner", "o", "", "Repository owner (user/organisation) (required)")
	uploadCmd.Flags().StringP("repo", "r", "", "Repository (required)")
	uploadCmd.Flags().StringP("branch", "b", "", "Commit to branch (default branch if empty)")
	uploadCmd.Flags().StringP("message", "m", "", "Commit message (required)")

	_ = uploadCmd.MarkFlagRequired("owner")
	_ = uploadCmd.MarkFlagRequired("repo")
	_ = uploadCmd.MarkFlagRequired("message")
}

var uploadCmd = &cobra.Command{
	Use:   "upload <local-path> <remote-path>",
	Short: "Uploads file to github repository",
	Args:  cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if os.Getenv("GITHUB_TOKEN") == "" {
			return errors.New("Missing GITHUB_TOKEN environment variable")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		upload, err := os.Open(args[0])
		defer func(){ _ = upload.Close() }()

		if err != nil {
			return err
		}

		ctx := context.Background()
		c, err := client.New(ctx, os.Getenv("GITHUB_TOKEN"))

		if err != nil {
			return err
		}

		owner, _ := cmd.Flags().GetString("owner")
		repo, _ := cmd.Flags().GetString("repo")
		msg, _ := cmd.Flags().GetString("msg")
		branch, _ := cmd.Flags().GetString("branch")
		path := args[1]

		opts := &github.RepositoryContentFileOptions{
			Message:   github.String(msg),
		}

		if(branch != ""){
			opts.Branch = github.String(branch)
		}

		// check if file exists in repo already
		f, d, _, err := c.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{})

		if err != nil {
			return err
		}

		if len(d) > 0 {
			return fmt.Errorf("Target expects path to file, %s is a directory", path)
		}

		if f != nil {
			opts.SHA = f.SHA
		}

		// attach file content
		content, err := ioutil.ReadAll(upload)

		if err != nil {
			return err
		}

		opts.Content = content

		// upload file
		_, _, err = c.Repositories.CreateFile(ctx, owner, repo, path, opts)

		return err
	},
}
