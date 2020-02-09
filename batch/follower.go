package batch

import (
	"io"
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
		attempts := 0
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
				}

				if len(job.Attempts) > attempts {
					attempts = len(job.Attempts)

					go func(stream string) {
						// fmt.Printf(">> Start log stream %v\n", stream)

						ff := FollowCloudWatchLog(b.Session, "/aws/batch/job", stream)

						for {
							select {
							case msg := <-ff.Out:
								if msg != "" {
									follower.Logging <- msg
								}

							case err := <-ff.Error:
								if err == io.EOF {
									// fmt.Printf(">> Log stream %v end\n", stream)
								} else {
									// fmt.Printf(">> Log stream error %v\n", err)
									// TODO: propagate to job follower?
								}

								break
							}
						}
					}(*job.Attempts[attempts-1].Container.LogStreamName)

					// we have new attempt, make sure we have drain out
					// all the previous logs before starting new one
					if logFollower != nil {

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

		// clean up
		close(follower.Logging)
		follower.Error <- io.EOF
		// close(follower.Status)
	}()

	return follower
}
