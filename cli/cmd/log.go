package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var (
	logJobId string
)

var logCmd = &cobra.Command{
	Use: "log",
	Run: func(cmd *cobra.Command, args []string) {
		followJob(logJobId)
	},
}

func init() {
	logCmd.Flags().StringVarP(&logJobId, "job", "j", "", "Job ID")
	logCmd.MarkFlagRequired("job")
}

func followJob(jobId string) {
	// follower := batch.FollowCloudWatchLog(awsSession, "/aws/batch/job", "hussmann-scoring/default/8ec60ff9-6e6a-4cc1-86f9-0cd6f70e7c92")

	// for msg := range follower.Out {
	// 	fmt.Println(msg)
	// 	break
	// }

	// select {
	// case err := <-follower.Err:
	// 	fmt.Printf("%v", err)
	// default:
	// 	// no error
	// }

	follower := batchCli.FollowJob(jobId)

	running := true
	for running {
		select {
		case msg := <-follower.Logging:
			fmt.Println(msg)

		case status := <-follower.Status:
			fmt.Printf("Status: %v\n", status)

		case err := <-follower.Error:
			if err != io.EOF {
				fmt.Printf("Error: %v", err)
			}

			running = false
		}
	}
}
