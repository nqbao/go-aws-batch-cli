package batch

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func (b *BatchCli) GetJob(jobID string) (*batch.JobDetail, error) {
	batchSvc := batch.New(b.Session)
	out, err := batchSvc.DescribeJobs(&batch.DescribeJobsInput{
		Jobs: []*string{aws.String(jobID)},
	})

	if err != nil {
		return nil, err
	} else if len(out.Jobs) == 0 {
		return nil, errors.New("Job is not found")
	}

	job := out.Jobs[0]
	return job, nil
}

func (b *BatchCli) GetRunningContainer(jobID string) (string, error) {
	ecsSvc := ecs.New(b.Session)
	clusters := []string{}

	err := ecsSvc.ListClustersPages(&ecs.ListClustersInput{}, func(out *ecs.ListClustersOutput, last bool) bool {
		for _, arn := range out.ClusterArns {
			clusters = append(clusters, *arn)
		}
		return true
	})

	if err != nil {
		return "", err
	}

	// try to search for container in all of the running containers
	containerName := ""
	for _, cluster := range clusters {
		var nestedError error
		err = ecsSvc.ListTasksPages(&ecs.ListTasksInput{
			Cluster:    aws.String(cluster),
			MaxResults: aws.Int64(50),
		}, func(out *ecs.ListTasksOutput, last bool) bool {
			if len(out.TaskArns) > 0 {
				tasks, err := ecsSvc.DescribeTasks(&ecs.DescribeTasksInput{
					Cluster: aws.String(cluster),
					Tasks:   out.TaskArns,
				})

				if err != nil {
					nestedError = err
					return false
				}

				for _, task := range tasks.Tasks {
					containerJobID := ""

					// only 1 container is running for now
					for _, ev := range task.Overrides.ContainerOverrides[0].Environment {
						if *ev.Name == "AWS_BATCH_JOB_ID" {
							containerName = *task.Containers[0].Name
							containerJobID = *ev.Value
							break
						}
					}

					if containerJobID == jobID {
						bits := strings.Split(*task.TaskArn, "/")
						containerName = fmt.Sprintf("%v/%v", containerName, bits[len(bits)-1])
						return false
					}
				}
			}

			return true
		})

		if err != nil {
			return "", err
		} else if nestedError != nil {
			return "", nestedError
		}

		if containerName != "" {
			break
		}
	}

	return containerName, nil
}
