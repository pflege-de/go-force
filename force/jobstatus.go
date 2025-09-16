package force

import (
	"fmt"
	"regexp"
	"time"

	"golang.org/x/sync/errgroup"
)

// errRegexp prüft, ob in einem (Error-)String CSV enthalten ist ([ -~] matched alle Zeichen vom Space bis zur Tilde)
var errRegexp = regexp.MustCompile(`[ -~].*CSV[ -~].*`)

func (forceApi *ForceApi) CheckJobStatus(op JobOperation, interval time.Duration) (JobOperation, error) {
	var g errgroup.Group

	if interval <= 0 {
		interval = 2 * time.Second
	}

	for _, jobID := range op.JobIDs {
		g.Go(func() error {
			tt := time.Tick(interval)
			statusURI := fmt.Sprintf("/services/data/%s/jobs/ingest/%s", forceApi.apiVersion, jobID)
			var status *JobInfo
			for range tt {
				status = &JobInfo{}
				err := forceApi.Get(statusURI, nil, status)
				if err != nil {
					return err
				}

				op.NumberRecordsFailed += status.NumberRecordsFailed
				op.NumberRecordsProcessed += status.NumberRecordsProcessed
				op.ResponseMessages = append(op.ResponseMessages, status.JobMessage)
				statePrefix := fmt.Sprintf("Status %s", status.State)

				switch status.State {
				case "Failed":
					jobFailed := FailedResultsError{}
					failedResultURI := fmt.Sprintf("/services/data/%s/jobs/ingest/%s/failedResults", forceApi.apiVersion, jobID)
					err = forceApi.Get(failedResultURI, nil, jobFailed)
					if err != nil {
						return err
					}

					op.ProgressReporter(statePrefix)

					if jobFailed.ErrorName == "InvalidBatch" && errRegexp.MatchString(jobFailed.ErrorDescription) {
						return jobFailed
					}

					return nil // NOTE (CARECON-1352): equivalent to the behavior prior to v1.1.8
				case "Aborted", "JobComplete":
					op.ProgressReporter(statePrefix)
					return nil
				default:
					executeProgReporter(op.ProgressReporter, status.State, statePrefix)
				}
			}

			return nil
		})
	}

	err := g.Wait()

	return op, err
}

func executeProgReporter(pr func(msg string), state, statePrefix string) {
	pr(fmt.Sprintf(getTextByState(state), statePrefix))
}

func getTextByState(state string) string {
	switch state {
	case "Open":
		return "%s: still open\n"
	case "UploadComplete":
		return "%s: upload complete\n"
	case "InProgress":
		return "%s: working on it\n"
	default:
		return "%s: unknown state\n"
	}
}
