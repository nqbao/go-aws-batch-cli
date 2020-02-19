package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	cancelJobId string
)

var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a running job",
	Run: func(cmd *cobra.Command, args []string) {
		job, err := batchCli.GetJob(cancelJobId)

		if err != nil {
			log.Fatalf("Can not find job: %v", err)
		}

		if *job.Status == "SUCCEEDED" || *job.Status == "FAILED" {
			log.Fatalf("Invalid job status: %v", *job.Status)
		}

		err = batchCli.CancelJob(*job.JobId)
		if err != nil {
			log.Fatalf("Unable to cancel job: %v", err)
		} else {
			fmt.Printf("Job %v is cancelled!\n", cancelJobId)

			// TODO: wait until it is really cancelled
		}
	},
}

func init() {
	cancelCmd.Flags().StringVarP(&cancelJobId, "id", "i", "", "Job ID")
	cancelCmd.MarkFlagRequired("id")
}
