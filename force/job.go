package force

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// CreateJob creates a new pointer to an instance of Job. Can be Modified with the given JobOptionsFuncs
func CreateJob(opts ...OptionsFunc) *Job {
	job := &Job{
		info: &JobInfo{},
	}
	for _, opt := range opts {
		opt(job)
	}

	return job
}

// WithOperation adds a given JobOperation to the Job
func (_ *Job) WithOperation(op JobOperation) OptionsFunc {
	return func(job *Job) {
		job.operation = op
	}
}

// WithMapper adds a given ObjectMapper to the Job
func (_ *Job) WithMapper(mapper ObjectMapper) OptionsFunc {
	return func(job *Job) {
		job.objectMapper = mapper
	}
}

// WithHTTPClient adds a HTTPClient to the Job, to communicate with salesforce
func (_ *Job) WithHTTPClient(client BulkClient) OptionsFunc {
	return func(job *Job) {
		job.client = client
	}
}

// WithJobInfo adds Job Information to the Job
func (_ *Job) WithJobInfo(info *JobInfo) OptionsFunc {
	return func(job *Job) {
		job.info = info
	}
}

func (_ *Job) WithApiVersion(apiVersion string) OptionsFunc {
	return func(job *Job) {
		job.apiVersion = apiVersion
	}
}

func (job *Job) GetForceApi() ForceApi      { return job.forceApi }
func (job *Job) GetOperation() JobOperation { return job.operation }
func (job *Job) GetMapper() ObjectMapper    { return job.objectMapper }
func (job *Job) GetHTTPClient() BulkClient  { return job.client }

func (job *Job) Start() error {
	params := map[string]string{
		"object":    job.operation.Object,
		"operation": job.operation.Operation,
	}

	if err := job.forceApi.Post("/services/data/"+job.apiVersion+"/jobs/ingest", nil, params, job.info); err != nil {
		return err
	}
	job.operation.ProgressReporter("job created", -1)
	return nil
}

// Run marshals the given payload to csv with the given ObjectMapper
// for the Job and sends the csv to the given SalesforceJob
func (job *Job) Run(payload any) error {
	if payload == nil {
		return errors.New("could not send payload because it is empty")
	}

	body, err := job.marshalCSV(payload)
	if err != nil {
		return fmt.Errorf("cannot marshal csv. %w", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s", job.forceApi.GetInstanceURL(), job.info.ContentURL), body)
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

func (job *Job) GetByteCount() int {
	cr := bytes.NewReader(job.bytes)
	return int(cr.Size())
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
	job.bytes = bulkData.Bytes()
	return bytes.NewReader(bulkData.Bytes()), nil
}
