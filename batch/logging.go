package batch

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type LogFollower struct {
	Out       chan string
	Error     chan error
	Terminate chan bool
}

func FollowCloudWatchLog(sess *session.Session, group string, stream string) *LogFollower {
	follower := &LogFollower{
		Out:       make(chan string),
		Error:     make(chan error),
		Terminate: make(chan bool, 1), // so it won't block when we ask to terminate the follower
	}

	go func() {
		svc := cloudwatchlogs.New(sess)

		input := &cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  aws.String(group),
			LogStreamName: aws.String(stream),
			StartFromHead: aws.Bool(true),
		}

		// TODO: add retry due to network error
		err := svc.GetLogEventsPages(input, func(out *cloudwatchlogs.GetLogEventsOutput, last bool) bool {
			for _, e := range out.Events {
				select {
				case <-follower.Terminate:
					return false
				default:
					follower.Out <- *e.Message
				}
			}

			return true
		})

		close(follower.Out)

		if err != nil {
			follower.Error <- err
		} else {
			follower.Error <- io.EOF
		}
	}()

	return follower
}
