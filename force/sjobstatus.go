package force

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

type Job struct {
	info         *JobInfo
	operation    JobOperation
	forceApi     ForceApi
	objectMapper ObjectMapper
	client       BulkClient
	bytes        []byte
	apiVersion   string
}

type BulkClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ObjectMapper func(objects any) [][]string

type JobInfo struct {
	Id                     string `json:"id"`
	State                  string `json:"state"`
	NumberRecordsFailed    int    `json:"numberRecordsFailed"`
	NumberRecordsProcessed int    `json:"numberRecordsProcessed"`
	JobMessage             string `json:"errorMessage"`
	ContentURL             string `json:"contentUrl"`
}

type JobOperation struct {
	Operation string
	Object    string
	Fields    []string

	NumberRecordsFailed    int
	NumberRecordsProcessed int
	ResponseMessages       []string
	JobIDs                 []string
	WriteLine              func(w io.Writer) bool
	ProgressReporter       func(msg string, bytesTransferred int)
}

type FailedResultsError struct {
	ApiError
	SfId string `json:"sf__Id"`
}

func (e FailedResultsError) Validate() bool {
	return len(e.Fields) != 0 || len(e.Message) != 0 || len(e.ErrorCode) != 0 ||
		len(e.ErrorName) != 0 || len(e.ErrorDescription) != 0 || len(e.SfId) != 0
}

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
					jobFailed := FailedResultsError{}
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
