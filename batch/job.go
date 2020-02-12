package batch

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
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
