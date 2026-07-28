package force

import (
	"errors"
	"fmt"
	"net"
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

	const retryLimit = 10

	for _, jobID := range op.JobIDs {
		g.Go(func() error {
			var attempts int // TODO: refactor the entire retry logic once this method accepts a [context.Context].

			tt := time.Tick(interval)
			statusURI := fmt.Sprintf("/services/data/%s/jobs/ingest/%s", forceApi.apiVersion, jobID)
			var status *JobInfo
			for range tt {
				status = &JobInfo{}

				if err := forceApi.Get(statusURI, nil, status); err != nil {
					if _, ok := errors.AsType[net.Error](err); ok && attempts < retryLimit {
						forceApi.trace("Network Error:", err.Error(), "%s")
						attempts += 1 // NOTE: prevent infinite retry loop.
						continue
					}

					return err
				}

				attempts = 0 // NOTE: `retryLimit` applies per (failed) request.

				op.NumberRecordsFailed += status.NumberRecordsFailed
				op.NumberRecordsProcessed += status.NumberRecordsProcessed
				op.ResponseMessages = append(op.ResponseMessages, status.JobMessage)
				statePrefix := fmt.Sprintf("Status %s", status.State)

				switch status.State {
				case "Failed":
					jobFailed := FailedResultsError{}
					failedResultURI := fmt.Sprintf("/services/data/%s/jobs/ingest/%s/failedResults", forceApi.apiVersion, jobID)
					err := forceApi.Get(failedResultURI, nil, &jobFailed)
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
