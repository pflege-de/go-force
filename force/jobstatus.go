package force

import (
	"fmt"
	"github.com/pflege-de/go-force/force/errors"
	"regexp"
	"time"
)

var isCsvError = regexp.MustCompile(`[ -~].*CSV[ -~].*`).MatchString

func (forceApi *ForceApi) CheckJobStatus(op JobOperation, tickerSeconds time.Duration, bytesCount int) (JobOperation, error) {
	tt := time.NewTicker(tickerSeconds * time.Second)
	defer tt.Stop()

	for _, jobID := range op.JobIDs {
		statusURI := fmt.Sprintf("/services/data/%s/jobs/ingest/%s", forceApi.apiVersion, jobID)
		var status *JobInfo

	STATUS:
		for {
			select {
			case <-tt.C:
				status = &JobInfo{}
				err := forceApi.Get(statusURI, nil, status)
				if err != nil {
					return op, err
				}

				statePrefix := fmt.Sprintf("Status %s", status.State)

				switch status.State {
				case "Failed":
					jobFailed := errors.FailedResultsError{}
					failedResultURI := fmt.Sprintf("/services/data/%s/jobs/ingest/%s/failedResults", forceApi.apiVersion, jobID)
					err = forceApi.Get(failedResultURI, nil, jobFailed)
					if err != nil {
						return op, err
					}

					op.ProgressReporter(statePrefix, 0)

					if jobFailed.ErrorName == "InvalidBatch" && isCsvError(jobFailed.ErrorDescription) {
						return op, jobFailed
					}

					break STATUS
				case "Aborted", "JobComplete":
					op.ProgressReporter(statePrefix, 0)
					break STATUS
				default:
					executeProgReporter(op.ProgressReporter, status.State, statePrefix, bytesCount)
				}
			}
		}

		op.NumberRecordsFailed += status.NumberRecordsFailed
		op.NumberRecordsProcessed += status.NumberRecordsProcessed
		op.ResponseMessages = append(op.ResponseMessages, status.JobMessage)
	}
	return op, nil
}

func executeProgReporter(pr func(msg string, bytesTransferred int), state, statePrefix string, bytesCount int) {
	pr(fmt.Sprintf(getTextByState(state), statePrefix), bytesCount)
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
