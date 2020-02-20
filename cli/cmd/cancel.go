package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	cancelJobID string
)

var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a running job",
	Run: func(cmd *cobra.Command, args []string) {
		err := batchCli.CancelJob(cancelJobID)
		if err != nil {
			log.Fatalf("Unable to cancel job: %v", err)
		} else {
			fmt.Printf("Job %v is cancelled!\n", cancelJobID)

			// TODO: add flag to wait until the job is really cancel
		}
	},
}

func init() {
	cancelCmd.Flags().StringVarP(&cancelJobID, "id", "i", "", "Job ID")
	cancelCmd.MarkFlagRequired("id")
}
