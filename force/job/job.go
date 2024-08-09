package job

import (
	"io"
	"net/http"
)

type OptionsFunc func(*Job)

// CreateJob creates a new pointer to an instane of Job. Can be Modified with the given JobOptionsFuncs
func CreateJob(opts ...OptionsFunc) *Job {
	job := &Job{
		info: &Info{},
	}
	for _, opt := range opts {
		opt(job)
	}

	return job
}

// WithOperation adds a given Operation to the Job
func WithOperation(op Operation) OptionsFunc {
	return func(job *Job) {
		job.operation = op
	}
}

// WithMapper adds a given ObjectMapper to the Job
func WithMapper(mapper ObjectMapper) OptionsFunc {
	return func(job *Job) {
		job.objectMapper = mapper
	}
}

// WithHTTPClient adds a HTTPClient to the Job, to communicate with salesforce
func WithHTTPClient(client BulkClient) OptionsFunc {
	return func(job *Job) {
		job.client = client
	}
}

// WithJobInfo adds Job Information to the Job
func WithJobInfo(info *Info) OptionsFunc {
	return func(job *Job) {
		job.info = info
	}
}

func WithApiVersion(apiVersion string) OptionsFunc {
	return func(job *Job) {
		job.apiVersion = apiVersion
	}
}

type Job struct {
	info         *Info
	operation    Operation
	objectMapper ObjectMapper
	client       BulkClient
	apiVersion   string
}

type BulkClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ObjectMapper func(objects any) [][]string

type Info struct {
	Id                     string `json:"id"`
	State                  string `json:"state"`
	NumberRecordsFailed    int    `json:"numberRecordsFailed"`
	NumberRecordsProcessed int    `json:"numberRecordsProcessed"`
	JobMessage             string `json:"errorMessage"`
	ContentURL             string `json:"contentUrl"`
}

type Operation struct {
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
