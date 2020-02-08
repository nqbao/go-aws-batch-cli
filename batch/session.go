package batch

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func InitAwsSession() *session.Session {
	region := "ap-northeast-1" // XXX
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String(region),
		},
	))

	return sess
}
