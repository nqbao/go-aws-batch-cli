package cmd

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nqbao/go-aws-batch-cli/batch"
	"github.com/spf13/cobra"
)

var (
	batchCli   *batch.BatchCli
	awsSession *session.Session

	rootCmd = &cobra.Command{
		Use: "aws-batch-cli",
		PersistentPreRun: func(*cobra.Command, []string) {
			awsSession = batch.InitAwsSession()

			batchCli = &batch.BatchCli{
				Session: awsSession,
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(cancelCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
