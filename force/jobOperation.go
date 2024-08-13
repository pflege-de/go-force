package force

import "io"

type JobOperation struct {
	Operation string
	Object    string
	Fields    []string

	NumberRecordsFailed    int
	NumberRecordsProcessed int
	ResponseMessages       []string
	JobIDs                 []string
	WriteLine              func(w io.Writer) bool
	ProgressReporter       func(msg string)
}

func (op *JobOperation) AddJobID(id string) {
	op.JobIDs = append(op.JobIDs, id)
}
