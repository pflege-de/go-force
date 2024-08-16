package force

import (
	"net/http"
	"net/url"
	"time"
)

type OptionsFunc func(*Job)

type BulkClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ObjectMapper func(objects any) [][]string

type Job struct {
	info         *JobInfo
	operation    JobOperation
	forceApi     *ForceApiInterface
	objectMapper ObjectMapper
	client       BulkClient
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

type ForceApiInterface interface {
	Get(path string, params url.Values, out any) error
	Post(path string, params url.Values, payload, out any) error
	Put(path string, params url.Values, payload, out interface{}) error
	Patch(path string, params url.Values, payload, out any) error
	Delete(path string, params url.Values) error
	NewRequest(method, path string, params url.Values) (*http.Request, error)

	Query(query string, out any) error
	QueryAll(query string, out interface{}) (err error)
	QueryNext(uri string, out interface{}) (err error)

	InsertSObject(in SObject) (resp *SObjectResponse, err error)
	DeleteSObject(id string, in SObject) (err error)
	GetSObject(id string, fields []string, out SObject) (err error)
	UpdateSObject(id string, in SObject) (err error)
	DescribeSObject(in SObject) (resp *SObjectDescription, err error)
	DescribeSObjects() (map[string]*SObjectMetaData, error)

	GetInstanceURL() string
	GetAccessToken() string
	RefreshToken() error

	TraceOn(prefix string, logger ApiLogger)
	TraceOff()
	GetOauth() *ForceOauth

	GetLimits() (limits Limits, err error)

	CheckJobStatus(op JobOperation, tickerSeconds time.Duration) (*JobOperation, error)
}

// ForceApiResponse represents a response from salesforce to a fapi.Query() or fapi.Get() request.
type ForceApiResponse interface {
	GetDone() bool
	GetNextRecordsUri() string
	GetRecords() interface{}
}

// ForceApiEvent is an interface that any object from SF satisfies. We can use it to query objects other than "EventStream2__c". Kinda legacy
type ForceApiEvent interface {
	AllFields() string
	GetID() string
}
