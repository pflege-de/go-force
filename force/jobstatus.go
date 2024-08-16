package force

import (
	"fmt"
	"regexp"
	"time"
)

// errRegexp prüft, ob in einem (Error-)String CSV enthalten ist ([ -~] matched alle Zeichen vom Space bis zur Tilde)
var errRegexp = regexp.MustCompile(`[ -~].*CSV[ -~].*`)

func (forceApi *ForceApi) CheckJobStatus(op JobOperation, tickerSeconds time.Duration) (*JobOperation, error) {
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
					return &op, err
				}

				statePrefix := fmt.Sprintf("Status %s", status.State)

				switch status.State {
				case "Failed":
					jobFailed := FailedResultsError{}
					failedResultURI := fmt.Sprintf("/services/data/%s/jobs/ingest/%s/failedResults", forceApi.apiVersion, jobID)
					err = forceApi.Get(failedResultURI, nil, jobFailed)
					if err != nil {
						return &op, err
					}

					op.ProgressReporter(statePrefix)

					if jobFailed.ErrorName == "InvalidBatch" && errRegexp.MatchString(jobFailed.ErrorDescription) {
						return &op, jobFailed
					}

					break STATUS
				case "Aborted", "JobComplete":
					op.ProgressReporter(statePrefix)
					break STATUS
				default:
					executeProgReporter(op.ProgressReporter, status.State, statePrefix)
				}
			}
		}

		op.NumberRecordsFailed += status.NumberRecordsFailed
		op.NumberRecordsProcessed += status.NumberRecordsProcessed
		op.ResponseMessages = append(op.ResponseMessages, status.JobMessage)
	}
	return &op, nil
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
