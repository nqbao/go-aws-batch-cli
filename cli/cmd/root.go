package cmd

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fatih/color"
	"github.com/nqbao/go-aws-batch-cli/batch"
	"github.com/spf13/cobra"
)

var (
	batchCli   *batch.BatchCli
	awsSession *session.Session
	noColor    bool

	rootCmd = &cobra.Command{
		Use: "aws-batch-cli",
		PersistentPreRun: func(*cobra.Command, []string) {
			awsSession = batch.InitAwsSession()

			batchCli = &batch.BatchCli{
				Session: awsSession,
			}

			if noColor {
				color.NoColor = true
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(cancelCmd)

	rootCmd.Flags().BoolVar(&noColor, "--no-color", false, "Set to true to disable color")
}

func Execute() error {
	return rootCmd.Execute()
}
