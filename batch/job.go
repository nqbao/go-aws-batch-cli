package batch

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
)

var (
	DefaultCancelReason = "Requested by user"
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

func (b *BatchCli) CancelJob(jobID string) error {
	batchSvc := batch.New(b.Session)

	job, err := b.GetJob(jobID)

	if err != nil {
		return fmt.Errorf("Unable to find job: %v", err)
	}

	if *job.Status == "SUCCEEDED" || *job.Status == "FAILED" {
		return fmt.Errorf("Invalid job status: %v", *job.Status)
	}

	if *job.Status == "STARTING" || *job.Status == "RUNNING" {
		_, err = batchSvc.TerminateJob(&batch.TerminateJobInput{
			JobId:  job.JobId,
			Reason: aws.String(DefaultCancelReason),
		})
	} else {
		_, err = batchSvc.CancelJob(&batch.CancelJobInput{
			JobId:  job.JobId,
			Reason: aws.String(DefaultCancelReason),
		})
	}

	return err
}
