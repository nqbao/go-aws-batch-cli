package batch

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func InitAwsSession() *session.Session {
	region := os.Getenv("AWS_DEFAULT_REGION")

	if region == "" {
		region = os.Getenv("AWS_REGION")
	}

	sess := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String(region),
		},
	))

	return sess
}
