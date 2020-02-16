package batch

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
)

// Simple wrapper of AWS SubmitJobInput
type SubmitRequest struct {
	Name        string
	Queue       string
	Definition  string
	Command     []string
	Parameters  map[string]*string
	Environment map[string]string
	Timeout     int
	Retries     int
}

func (s *SubmitRequest) SetCommandString(cmd string) {
	if cmd == "" {
		s.Command = nil
	} else {
		s.Command = strings.Split(cmd, " ")
	}
}

type BatchCli struct {
	Session *session.Session
}

// SubmitJob submits a job to batch cli
func (b *BatchCli) SubmitJob(request *SubmitRequest) (string, error) {
	batchSvc := batch.New(b.Session)

	input := &batch.SubmitJobInput{
		JobQueue:           aws.String(request.Queue),
		JobDefinition:      aws.String(request.Definition),
		Parameters:         request.Parameters,
		ContainerOverrides: &batch.ContainerOverrides{},
	}

	if request.Timeout > 0 {
		input.Timeout = &batch.JobTimeout{
			AttemptDurationSeconds: aws.Int64(int64(request.Timeout)),
		}
	}

	if request.Retries > 0 {
		input.RetryStrategy = &batch.RetryStrategy{
			Attempts: aws.Int64(int64(request.Retries)),
		}
	}

	if len(request.Command) > 0 {
		input.ContainerOverrides.Command = []*string{}

		for _, c := range request.Command {
			input.ContainerOverrides.Command = append(input.ContainerOverrides.Command, aws.String(c))
		}
	}

	if request.Name == "" {
		input.JobName = aws.String("test")
	} else {
		input.JobName = aws.String(request.Name)
	}

	if request.Environment != nil && len(request.Environment) > 0 {
		input.ContainerOverrides.Environment = []*batch.KeyValuePair{}

		for k, v := range request.Environment {
			p := &batch.KeyValuePair{
				Name:  aws.String(k),
				Value: aws.String(v),
			}
			input.ContainerOverrides.Environment = append(input.ContainerOverrides.Environment, p)
		}
	}

	output, err := batchSvc.SubmitJob(input)

	if err != nil {
		return "", err
	}

	return *output.JobId, nil
}
