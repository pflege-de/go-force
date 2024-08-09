package force

import (
	"net/http"
)

type OptionsFunc func(*Job)

type BulkClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ObjectMapper func(objects any) [][]string

type Job struct {
	info         *JobInfo
	operation    JobOperation
	forceApi     ForceApi
	objectMapper ObjectMapper
	client       BulkClient
	bytes        []byte
	apiVersion   string
}

type JobInfo struct {
	Id                     string `json:"id"`
	State                  string `json:"state"`
	NumberRecordsFailed    int    `json:"numberRecordsFailed"`
	NumberRecordsProcessed int    `json:"numberRecordsProcessed"`
	JobMessage             string `json:"errorMessage"`
	ContentURL             string `json:"contentUrl"`
}
