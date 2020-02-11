package batch

import (
	"fmt"
	"io"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type JobFollower struct {
	Status  chan string
	Error   chan error
	Logging chan string
}

func (b *BatchCli) FollowJob(jobID string) *JobFollower {
	follower := &JobFollower{}
	follower.Error = make(chan error, 1)
	follower.Status = make(chan string, 1)
	follower.Logging = make(chan string, 100)

	// TODO: add some retry when querying job

	go func() {
		prevStatus := ""
		var logFollower *LogFollower

		for {
			job, err := b.GetJob(jobID)
			if err != nil {
				follower.Error <- err
				break
			} else {
				newStatus := *job.Status

				// when job is being retried, it will go back to RUNNABLE
				// RUNNABLE -> STARTING -> RUNNING
				if newStatus != prevStatus {
					follower.Status <- *job.Status

					if newStatus == "RUNNING" {
						runningContainer, err := b.GetRunningContainer(jobID)

						if err != nil {
							follower.Error <- fmt.Errorf("Unable to locate running container: %v", err)
							break
						}

						bits := strings.Split(*job.JobDefinition, "/")
						bits = strings.Split(bits[len(bits)-1], ":")
						streamName := fmt.Sprintf("%v/%v", bits[0], runningContainer)

						go func(streamName string) {
							if logFollower != nil {
								logFollower.Terminate <- true
							}

							logFollower = b.followRunningJobStream(follower, streamName)
						}(streamName)
					}
				}

				if newStatus == "SUCCEEDED" || newStatus == "FAILED" {
					break
				}

				prevStatus = newStatus

				log.Debug("Job status", newStatus)
			}

			// delay a bit, TODO add some jittering here
			<-time.After(1000 * time.Millisecond)
		}

		// give it 1s to finish
		if logFollower != nil {
			<-time.After(1 * time.Second)
			logFollower.Terminate <- true
		}

		// clean up
		close(follower.Logging)
		follower.Error <- io.EOF
		// close(follower.Status)
	}()

	return follower
}

func (b *BatchCli) followRunningJobStream(follower *JobFollower, stream string) *LogFollower {
	// fmt.Printf(">> Start log stream %v\n", stream)

	ff := FollowCloudWatchLog(b.Session, "/aws/batch/job", stream, true)

	go func() {
		for {
			select {
			case msg := <-ff.Out:
				if msg != "" {
					follower.Logging <- msg
				}

			case err := <-ff.Error:
				if err == io.EOF {
					fmt.Printf(">> Log stream %v end\n", stream)
				} else {
					fmt.Printf(">> Log stream error %v\n", err)
					follower.Error <- err
				}

				break
			}
		}
	}()

	return ff
}
