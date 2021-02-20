package cmd

import (
	"fmt"
	"os"

	"github.com/invit/ghupload/internal/lib/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "ghupload",
	Short:         "Uploads a file to github repository",
	SilenceErrors: true,
	Version:       fmt.Sprintf("%s-%s\n", version.Version, version.Commit),
}

// Execute runs root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
