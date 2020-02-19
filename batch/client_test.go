package batch_test

import (
	"testing"

	"github.com/nqbao/go-aws-batch-cli/batch"
)

func TestPrepareBatchSubmitInput(t *testing.T) {
	b := &batch.BatchCli{}

	r := &batch.SubmitRequest{
		Definition: "my-job",
		Queue:      "test",
		Environment: map[string]string{
			"Env1": "Value1",
		},
		Parameters: map[string]string{
			"Par1": "Par2",
		},
	}

	i := b.PrepareBatchSubmitInput(r)

	if i.JobName == nil {
		t.Errorf("Job name should be filled automatically")
	}

	if *i.JobQueue != "test" {
		t.Errorf("Job queue is not set properly")
	}

	if len(i.ContainerOverrides.Environment) != 1 {
		t.Errorf("Job container environment is not set properly")
	} else {
		if *i.ContainerOverrides.Environment[0].Name != "Env1" && *i.ContainerOverrides.Environment[0].Value != "Value1" {
			t.Errorf("Job container environment value is not set properly")
		}
	}

	if len(i.Parameters) != 1 {
		t.Errorf("Job parameters is not set properly")
	} else {
		if i.Parameters["Par1"] == nil || *i.Parameters["Par1"] != "Par2" {
			t.Errorf("Job parameters value is not set properly")
		}
	}
}
