package cmd

import (
	"fmt"
	"io"
	"log"

	"github.com/nqbao/go-aws-batch-cli/batch"

	"github.com/spf13/cobra"
)

var (
	logJobId string
)

var logCmd = &cobra.Command{
	Use: "log",
	Run: func(cmd *cobra.Command, args []string) {
		job, err := batchCli.GetJob(logJobId)

		if err != nil {
			log.Fatalf("Can not find job: %v", err)
		}

		if *job.Status != "SUCCEEDED" && *job.Status != "FAILED" {
			log.Fatalf("Invalid job status: %v", *job.Status)
		}

		attempt := job.Attempts[len(job.Attempts)-1]

		follower := batch.FollowCloudWatchLog(awsSession, "/aws/batch/job", *attempt.Container.LogStreamName, false)

		running := true
		for running {
			select {
			case msg := <-follower.Out:
				fmt.Println(msg)
			case err := <-follower.Error:
				if err != io.EOF {
					log.Fatalf("Error while retriving log stream: %v", err)
				}

				running = false
			}
		}
	},
}

func init() {
	logCmd.Flags().StringVarP(&logJobId, "id", "i", "", "Job ID")
	logCmd.MarkFlagRequired("id")
}
