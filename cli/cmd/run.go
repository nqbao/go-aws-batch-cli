package cmd

import (
	"fmt"
	"io"
	"log"
	"strings"

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
	runCommand       string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an AWS batch job",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		params := make(map[string]string)
		envs := make(map[string]string)

		for _, paramStr := range runJobParameters {
			bits := strings.SplitN(paramStr, "=", 2)

			if len(bits) == 2 {
				params[bits[0]] = bits[1]
			} else {
				params[bits[0]] = bits[0]
			}
		}

		for _, envStr := range runEnvironment {
			bits := strings.SplitN(envStr, "=", 2)

			if len(bits) == 2 {
				envs[bits[0]] = bits[1]
			} else {
				envs[bits[0]] = bits[0]
			}
		}

		request := &batch.SubmitRequest{
			Name:        runJobName,
			Definition:  runJobDefinition,
			Queue:       runJobQueue,
			Parameters:  params,
			Environment: envs,
			Retries:     runJobRetries,
		}

		if runCommand != "" {
			request.SetCommandString(runCommand)
		} else if len(args) > 0 {
			request.Command = args
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
	runCmd.Flags().StringVarP(&runCommand, "command", "c", "", "Override container command")
	runCmd.Flags().StringArrayVarP(&runJobParameters, "parameter", "p", []string{}, "")
	runCmd.Flags().StringArrayVarP(&runEnvironment, "env", "e", []string{}, "")
	runCmd.Flags().IntVarP(&runJobTimeout, "timeout", "", 0, "Timeout")
	runCmd.Flags().IntVarP(&runJobRetries, "num-retries", "r", 0, "Job retries")
	runCmd.Flags().BoolVarP(&runFollowFlag, "follow", "f", false, "Follow job log")

	runCmd.MarkFlagRequired("queue")
	runCmd.MarkFlagRequired("job")
}

func followJob(jobId string) {
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
