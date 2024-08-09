package errors

type FailedResultsError struct {
	ApiError
	SfId string `json:"sf__Id"`
}

func (e FailedResultsError) Validate() bool {
	return len(e.Fields) != 0 || len(e.Message) != 0 || len(e.ErrorCode) != 0 ||
		len(e.ErrorName) != 0 || len(e.ErrorDescription) != 0 || len(e.SfId) != 0
}
