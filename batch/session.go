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

	profile := os.Getenv("AWS_PROFILE")
	if profile == "" {
		profile = "default"
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	}))

	return sess
}
