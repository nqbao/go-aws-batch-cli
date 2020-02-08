package cmd

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use: "aws-batch-cli",
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
