package force

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// CreateJob creates a new pointer to an instance of Job. Can be Modified with the given JobOptionsFuncs
func CreateJob(fapi ForceApiInterface, opts ...OptionsFunc) *Job {
	job := &Job{
		forceApi:     fapi,
		operation:    JobOperation{},
		objectMapper: func(objects any) [][]string { return nil },
		info:         &JobInfo{},
		apiVersion:   DefaultAPIVersion,
		client:       http.DefaultClient,
	}
	for _, opt := range opts {
		opt(job)
	}

	return job
}

// JobWithHTTPClient adds a HTTPClient to the Job, to communicate with salesforce
func JobWithHTTPClient(client BulkClient) OptionsFunc {
	return func(job *Job) {
		job.client = client
	}
}

// JobWithJobInfo adds Job Information to the Job
func JobWithJobInfo(info *JobInfo) OptionsFunc {
	return func(job *Job) {
		job.info = info
	}
}

// JobWithApiVersion set the ApiVersion of a Job
func JobWithApiVersion(apiVersion string) OptionsFunc {
	return func(job *Job) {
		job.apiVersion = apiVersion
	}
}

// JobWithMapper adds a given ObjectMapper to the Job
func JobWithMapper(mapper ObjectMapper) OptionsFunc {
	return func(job *Job) {
		job.objectMapper = mapper
	}
}

// JobWithOperation adds a given JobOperation to the Job
func JobWithOperation(operation JobOperation) OptionsFunc {
	return func(job *Job) {
		job.operation = operation
	}
}

func (job *Job) Start() error {
	params := map[string]string{
		"object":    job.operation.Object,
		"operation": job.operation.Operation,
	}

	if err := job.forceApi.Post("/services/data/"+job.apiVersion+"/jobs/ingest", nil, params, job.info); err != nil {
		return err
	}
	job.operation.ProgressReporter("job created")
	return nil
}

// Run marshals the given payload to csv with the given ObjectMapper for the Job and sends the csv to the given SalesforceJob
func (job *Job) Run(payload any) error {

	if payload == nil {
		return errors.New("could not send payload because it is empty")
	}

	body, err := job.marshalCSV(payload)
	if err != nil {
		return fmt.Errorf("cannot marshal csv. %w", err)
	}

	urlFormat := "%s%s"
	instanceUrl := job.forceApi.GetInstanceURL()
	contentUrl := job.info.ContentURL
	if !strings.HasPrefix(contentUrl, "/") && !strings.HasSuffix(instanceUrl, "/") {
		urlFormat = "%s/%s"
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf(urlFormat, instanceUrl, contentUrl), body)
	if err != nil {
		return fmt.Errorf("could not create new HTTP Request. %w", err)
	}

	req.Header.Set("Content-Type", "text/csv")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", job.forceApi.GetAccessToken()))

	res, err := job.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not put csv bulk data. %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		errb, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("unexpected StatusCode on PUT batch: %d (%s), %s", res.StatusCode, res.Status, string(errb))
	}

	statusURI := fmt.Sprintf("/services/data/"+job.apiVersion+"/jobs/ingest/%s", job.info.Id)
	params := map[string]string{
		"state": "UploadComplete",
	}

	if err := job.forceApi.Patch(statusURI, nil, params, job.info); err != nil {
		return err
	}

	job.operation.AddJobID(job.info.Id)

	return nil
}

func (job *Job) marshalCSV(payload any) (io.Reader, error) {
	// Map Objects to a csv Reader, for bulk api
	var bulkData bytes.Buffer
	w := csv.NewWriter(&bulkData)
	var records [][]string
	records = append(records, job.operation.Fields)
	records = append(records, job.objectMapper(payload)...)
	if err := w.WriteAll(records); err != nil {
		return nil, fmt.Errorf("could not create csv from records. %w", err)
	}
	return bytes.NewReader(bulkData.Bytes()), nil
}

func (job *Job) GetForceApi() ForceApiInterface {
	return job.forceApi
}

func (job *Job) GetHTTPClient() BulkClient {
	return job.client
}

func (job *Job) GetMapper() ObjectMapper {
	return job.objectMapper
}

func (job *Job) GetOperation() JobOperation {
	return job.operation
}
