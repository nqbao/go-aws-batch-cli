package batch

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type LogFollower struct {
	Out       chan string
	Error     chan error
	Terminate chan bool
}

func FollowCloudWatchLog(sess *session.Session, group string, stream string, repeat bool) *LogFollower {
	follower := &LogFollower{
		Out:       make(chan string),
		Error:     make(chan error),
		Terminate: make(chan bool, 1), // so it won't block when we ask to terminate the follower
	}

	go func() {
		svc := cloudwatchlogs.New(sess)

		running := true
		isErr := false
		var lastEvent *cloudwatchlogs.OutputLogEvent = nil

		for running {
			input := &cloudwatchlogs.GetLogEventsInput{
				LogGroupName:  aws.String(group),
				LogStreamName: aws.String(stream),
			}

			if lastEvent != nil {
				input.StartTime = lastEvent.IngestionTime
			} else {
				input.StartFromHead = aws.Bool(true)
			}

			// TODO: add retry due to network error
			err := svc.GetLogEventsPages(input, func(out *cloudwatchlogs.GetLogEventsOutput, last bool) bool {
				for _, e := range out.Events {
					select {
					case <-follower.Terminate:
						running = false // this will break the loop
						return false
					default:
						if lastEvent == nil {
							follower.Out <- *e.Message
							lastEvent = e
						} else if *lastEvent.Message != *e.Message || *lastEvent.IngestionTime != *e.IngestionTime {
							// it's a new message
							follower.Out <- *e.Message
							lastEvent = e
						}
					}
				}

				return true
			})

			if err != nil {
				isErr = true
				follower.Error <- err
				running = false
			} else if !repeat {
				running = false
			} else {
				// add some jitter
				<-time.After(200 * time.Microsecond)
			}
		}

		close(follower.Out)

		if !isErr {
			follower.Error <- io.EOF
		}
	}()

	return follower
}
