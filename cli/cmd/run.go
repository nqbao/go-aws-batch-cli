package cmd

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/nqbao/go-aws-batch-cli/batch"
	"github.com/spf13/cobra"
)

var (
	runJobName         string
	runJobQueue        string
	runJobDefinition   string
	runJobParameters   []string
	runEnvironment     []string
	runJobTimeout      int
	runJobRetries      int
	runFollowFlag      bool
	runCommand         string
	runContainerMemory int
	runContainerVcpus  int
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
			Name:            runJobName,
			Definition:      runJobDefinition,
			Queue:           runJobQueue,
			Parameters:      params,
			Environment:     envs,
			Retries:         runJobRetries,
			ContainerMemory: runContainerMemory,
			ContainerVcpus:  runContainerVcpus,
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

		fmt.Printf("Job ID: %s\n", color.New(color.Bold).Sprint(jobID))

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
	runCmd.Flags().IntVarP(&runContainerMemory, "memory", "", 0, "Override container memory (in MiB)")
	runCmd.Flags().IntVarP(&runContainerVcpus, "vcpus", "", 0, "Override container vcpus")
	runCmd.Flags().BoolVarP(&runFollowFlag, "follow", "f", false, "Follow job log")

	runCmd.MarkFlagRequired("queue")
	runCmd.MarkFlagRequired("job")
}

func followJob(jobId string) {
	follower := batchCli.FollowJob(jobId)

	running := true
	success := false

	// TODO: add a message when a new attempt is started
	for running {
		select {
		case msg := <-follower.Logging:
			fmt.Println(msg)

		case status := <-follower.Status:
			fmt.Println(color.CyanString("[Status]"), status)

			if status == "SUCCEEDED" {
				success = true
			}

		case err := <-follower.Error:
			if err != io.EOF {
				fmt.Println(color.RedString("[Error]"), err)
			}

			running = false
		}
	}

	if success {
		fmt.Println(color.BlueString("Job has completed successfully!"))
	} else {
		fmt.Println(color.RedString("Job has failed."))
	}
}
