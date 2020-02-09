package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/nqbao/go-aws-batch-cli/batch"
	"github.com/spf13/cobra"
)

var (
	runJobName       string
	runJobQueue      string
	runJobDefinition string
	runJobParameters []string
	runEnvironment   []string
	runJobTimeout    int
	runJobRetries    int
	runFollowFlag    bool
)

var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		if runJobName == "" {
			runJobName = runJobDefinition
		}

		params := make(map[string]*string)
		envs := make(map[string]string)

		for _, paramStr := range runJobParameters {
			bits := strings.SplitN(paramStr, "=", 2)
			params[bits[0]] = aws.String(bits[1])
		}

		for _, envStr := range runEnvironment {
			bits := strings.SplitN(envStr, "=", 2)
			envs[bits[0]] = bits[1]
		}

		request := &batch.SubmitRequest{
			Name:        runJobName,
			Definition:  runJobDefinition,
			Queue:       runJobQueue,
			Parameters:  params,
			Environment: envs,
			Retries:     runJobRetries,
		}

		jobID, err := batchCli.SubmitJob(request)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Job ID: %s\n", jobID)

		if runFollowFlag {
			followJob(jobID)
		}
	},
}

func init() {
	runCmd.Flags().StringVarP(&runJobName, "name", "", "", "Job name. Leave blank to autogenerate")
	runCmd.Flags().StringVarP(&runJobQueue, "queue", "q", "", "Queue")
	runCmd.Flags().StringVarP(&runJobDefinition, "job", "j", "", "Job Definition")
	runCmd.Flags().StringArrayVarP(&runJobParameters, "parameter", "p", []string{}, "")
	runCmd.Flags().StringArrayVarP(&runEnvironment, "env", "e", []string{}, "")
	runCmd.Flags().IntVarP(&runJobTimeout, "timeout", "", 0, "Timeout")
	runCmd.Flags().IntVarP(&runJobRetries, "num-retries", "r", 0, "Job retries")
	runCmd.Flags().BoolVarP(&runFollowFlag, "follow", "f", false, "Follow job log")

	runCmd.MarkFlagRequired("queue")
	runCmd.MarkFlagRequired("job")
}
