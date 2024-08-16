// Code generated by MockGen. DO NOT EDIT.
// Source: JobTypes.go
//
// Generated by this command:
//
//	mockgen -source=JobTypes.go -destination=mocks/JobTypes.go -package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
	url "net/url"
	reflect "reflect"
	time "time"

	force "github.com/pflege-de/go-force/force"
	gomock "go.uber.org/mock/gomock"
)

// MockBulkClient is a mock of BulkClient interface.
type MockBulkClient struct {
	ctrl     *gomock.Controller
	recorder *MockBulkClientMockRecorder
}

// MockBulkClientMockRecorder is the mock recorder for MockBulkClient.
type MockBulkClientMockRecorder struct {
	mock *MockBulkClient
}

// NewMockBulkClient creates a new mock instance.
func NewMockBulkClient(ctrl *gomock.Controller) *MockBulkClient {
	mock := &MockBulkClient{ctrl: ctrl}
	mock.recorder = &MockBulkClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBulkClient) EXPECT() *MockBulkClientMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockBulkClient) Do(req *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", req)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockBulkClientMockRecorder) Do(req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockBulkClient)(nil).Do), req)
}

// MockForceApiInterface is a mock of ForceApiInterface interface.
type MockForceApiInterface struct {
	ctrl     *gomock.Controller
	recorder *MockForceApiInterfaceMockRecorder
}

// MockForceApiInterfaceMockRecorder is the mock recorder for MockForceApiInterface.
type MockForceApiInterfaceMockRecorder struct {
	mock *MockForceApiInterface
}

// NewMockForceApiInterface creates a new mock instance.
func NewMockForceApiInterface(ctrl *gomock.Controller) *MockForceApiInterface {
	mock := &MockForceApiInterface{ctrl: ctrl}
	mock.recorder = &MockForceApiInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockForceApiInterface) EXPECT() *MockForceApiInterfaceMockRecorder {
	return m.recorder
}

// CheckJobStatus mocks base method.
func (m *MockForceApiInterface) CheckJobStatus(op force.JobOperation, tickerSeconds time.Duration) (force.JobOperation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckJobStatus", op, tickerSeconds)
	ret0, _ := ret[0].(force.JobOperation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckJobStatus indicates an expected call of CheckJobStatus.
func (mr *MockForceApiInterfaceMockRecorder) CheckJobStatus(op, tickerSeconds any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckJobStatus", reflect.TypeOf((*MockForceApiInterface)(nil).CheckJobStatus), op, tickerSeconds)
}

// Delete mocks base method.
func (m *MockForceApiInterface) Delete(path string, params url.Values) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", path, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockForceApiInterfaceMockRecorder) Delete(path, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockForceApiInterface)(nil).Delete), path, params)
}

// DeleteSObject mocks base method.
func (m *MockForceApiInterface) DeleteSObject(id string, in force.SObject) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSObject", id, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSObject indicates an expected call of DeleteSObject.
func (mr *MockForceApiInterfaceMockRecorder) DeleteSObject(id, in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSObject", reflect.TypeOf((*MockForceApiInterface)(nil).DeleteSObject), id, in)
}

// DescribeSObject mocks base method.
func (m *MockForceApiInterface) DescribeSObject(in force.SObject) (*force.SObjectDescription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeSObject", in)
	ret0, _ := ret[0].(*force.SObjectDescription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeSObject indicates an expected call of DescribeSObject.
func (mr *MockForceApiInterfaceMockRecorder) DescribeSObject(in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeSObject", reflect.TypeOf((*MockForceApiInterface)(nil).DescribeSObject), in)
}

// DescribeSObjects mocks base method.
func (m *MockForceApiInterface) DescribeSObjects() (map[string]*force.SObjectMetaData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeSObjects")
	ret0, _ := ret[0].(map[string]*force.SObjectMetaData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeSObjects indicates an expected call of DescribeSObjects.
func (mr *MockForceApiInterfaceMockRecorder) DescribeSObjects() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeSObjects", reflect.TypeOf((*MockForceApiInterface)(nil).DescribeSObjects))
}

// Get mocks base method.
func (m *MockForceApiInterface) Get(path string, params url.Values, out any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", path, params, out)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockForceApiInterfaceMockRecorder) Get(path, params, out any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockForceApiInterface)(nil).Get), path, params, out)
}

// GetAccessToken mocks base method.
func (m *MockForceApiInterface) GetAccessToken() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccessToken")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetAccessToken indicates an expected call of GetAccessToken.
func (mr *MockForceApiInterfaceMockRecorder) GetAccessToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccessToken", reflect.TypeOf((*MockForceApiInterface)(nil).GetAccessToken))
}

// GetInstanceURL mocks base method.
func (m *MockForceApiInterface) GetInstanceURL() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInstanceURL")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetInstanceURL indicates an expected call of GetInstanceURL.
func (mr *MockForceApiInterfaceMockRecorder) GetInstanceURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstanceURL", reflect.TypeOf((*MockForceApiInterface)(nil).GetInstanceURL))
}

// GetLimits mocks base method.
func (m *MockForceApiInterface) GetLimits() (force.Limits, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLimits")
	ret0, _ := ret[0].(force.Limits)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLimits indicates an expected call of GetLimits.
func (mr *MockForceApiInterfaceMockRecorder) GetLimits() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLimits", reflect.TypeOf((*MockForceApiInterface)(nil).GetLimits))
}

// GetOauth mocks base method.
func (m *MockForceApiInterface) GetOauth() *force.ForceOauth {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOauth")
	ret0, _ := ret[0].(*force.ForceOauth)
	return ret0
}

// GetOauth indicates an expected call of GetOauth.
func (mr *MockForceApiInterfaceMockRecorder) GetOauth() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOauth", reflect.TypeOf((*MockForceApiInterface)(nil).GetOauth))
}

// GetSObject mocks base method.
func (m *MockForceApiInterface) GetSObject(id string, fields []string, out force.SObject) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSObject", id, fields, out)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetSObject indicates an expected call of GetSObject.
func (mr *MockForceApiInterfaceMockRecorder) GetSObject(id, fields, out any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSObject", reflect.TypeOf((*MockForceApiInterface)(nil).GetSObject), id, fields, out)
}

// InsertSObject mocks base method.
func (m *MockForceApiInterface) InsertSObject(in force.SObject) (*force.SObjectResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertSObject", in)
	ret0, _ := ret[0].(*force.SObjectResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertSObject indicates an expected call of InsertSObject.
func (mr *MockForceApiInterfaceMockRecorder) InsertSObject(in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertSObject", reflect.TypeOf((*MockForceApiInterface)(nil).InsertSObject), in)
}

// NewRequest mocks base method.
func (m *MockForceApiInterface) NewRequest(method, path string, params url.Values) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRequest", method, path, params)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRequest indicates an expected call of NewRequest.
func (mr *MockForceApiInterfaceMockRecorder) NewRequest(method, path, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRequest", reflect.TypeOf((*MockForceApiInterface)(nil).NewRequest), method, path, params)
}

// Patch mocks base method.
func (m *MockForceApiInterface) Patch(path string, params url.Values, payload, out any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Patch", path, params, payload, out)
	ret0, _ := ret[0].(error)
	return ret0
}

// Patch indicates an expected call of Patch.
func (mr *MockForceApiInterfaceMockRecorder) Patch(path, params, payload, out any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockForceApiInterface)(nil).Patch), path, params, payload, out)
}

// Post mocks base method.
func (m *MockForceApiInterface) Post(path string, params url.Values, payload, out any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Post", path, params, payload, out)
	ret0, _ := ret[0].(error)
	return ret0
}

// Post indicates an expected call of Post.
func (mr *MockForceApiInterfaceMockRecorder) Post(path, params, payload, out any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockForceApiInterface)(nil).Post), path, params, payload, out)
}

// Put mocks base method.
func (m *MockForceApiInterface) Put(path string, params url.Values, payload, out any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", path, params, payload, out)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockForceApiInterfaceMockRecorder) Put(path, params, payload, out any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockForceApiInterface)(nil).Put), path, params, payload, out)
}

// Query mocks base method.
func (m *MockForceApiInterface) Query(query string, out any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", query, out)
	ret0, _ := ret[0].(error)
	return ret0
}

// Query indicates an expected call of Query.
func (mr *MockForceApiInterfaceMockRecorder) Query(query, out any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockForceApiInterface)(nil).Query), query, out)
}

// QueryAll mocks base method.
func (m *MockForceApiInterface) QueryAll(query string, out any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAll", query, out)
	ret0, _ := ret[0].(error)
	return ret0
}

// QueryAll indicates an expected call of QueryAll.
func (mr *MockForceApiInterfaceMockRecorder) QueryAll(query, out any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAll", reflect.TypeOf((*MockForceApiInterface)(nil).QueryAll), query, out)
}

// QueryNext mocks base method.
func (m *MockForceApiInterface) QueryNext(uri string, out any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryNext", uri, out)
	ret0, _ := ret[0].(error)
	return ret0
}

// QueryNext indicates an expected call of QueryNext.
func (mr *MockForceApiInterfaceMockRecorder) QueryNext(uri, out any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryNext", reflect.TypeOf((*MockForceApiInterface)(nil).QueryNext), uri, out)
}

// RefreshToken mocks base method.
func (m *MockForceApiInterface) RefreshToken() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken")
	ret0, _ := ret[0].(error)
	return ret0
}

// RefreshToken indicates an expected call of RefreshToken.
func (mr *MockForceApiInterfaceMockRecorder) RefreshToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockForceApiInterface)(nil).RefreshToken))
}

// TraceOff mocks base method.
func (m *MockForceApiInterface) TraceOff() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "TraceOff")
}

// TraceOff indicates an expected call of TraceOff.
func (mr *MockForceApiInterfaceMockRecorder) TraceOff() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TraceOff", reflect.TypeOf((*MockForceApiInterface)(nil).TraceOff))
}

// TraceOn mocks base method.
func (m *MockForceApiInterface) TraceOn(prefix string, logger force.ApiLogger) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "TraceOn", prefix, logger)
}

// TraceOn indicates an expected call of TraceOn.
func (mr *MockForceApiInterfaceMockRecorder) TraceOn(prefix, logger any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TraceOn", reflect.TypeOf((*MockForceApiInterface)(nil).TraceOn), prefix, logger)
}

// UpdateSObject mocks base method.
func (m *MockForceApiInterface) UpdateSObject(id string, in force.SObject) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSObject", id, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSObject indicates an expected call of UpdateSObject.
func (mr *MockForceApiInterfaceMockRecorder) UpdateSObject(id, in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSObject", reflect.TypeOf((*MockForceApiInterface)(nil).UpdateSObject), id, in)
}

// MockForceApiResponse is a mock of ForceApiResponse interface.
type MockForceApiResponse struct {
	ctrl     *gomock.Controller
	recorder *MockForceApiResponseMockRecorder
}

// MockForceApiResponseMockRecorder is the mock recorder for MockForceApiResponse.
type MockForceApiResponseMockRecorder struct {
	mock *MockForceApiResponse
}

// NewMockForceApiResponse creates a new mock instance.
func NewMockForceApiResponse(ctrl *gomock.Controller) *MockForceApiResponse {
	mock := &MockForceApiResponse{ctrl: ctrl}
	mock.recorder = &MockForceApiResponseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockForceApiResponse) EXPECT() *MockForceApiResponseMockRecorder {
	return m.recorder
}

// GetDone mocks base method.
func (m *MockForceApiResponse) GetDone() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDone")
	ret0, _ := ret[0].(bool)
	return ret0
}

// GetDone indicates an expected call of GetDone.
func (mr *MockForceApiResponseMockRecorder) GetDone() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDone", reflect.TypeOf((*MockForceApiResponse)(nil).GetDone))
}

// GetNextRecordsUri mocks base method.
func (m *MockForceApiResponse) GetNextRecordsUri() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextRecordsUri")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetNextRecordsUri indicates an expected call of GetNextRecordsUri.
func (mr *MockForceApiResponseMockRecorder) GetNextRecordsUri() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextRecordsUri", reflect.TypeOf((*MockForceApiResponse)(nil).GetNextRecordsUri))
}

// GetRecords mocks base method.
func (m *MockForceApiResponse) GetRecords() any {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecords")
	ret0, _ := ret[0].(any)
	return ret0
}

// GetRecords indicates an expected call of GetRecords.
func (mr *MockForceApiResponseMockRecorder) GetRecords() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecords", reflect.TypeOf((*MockForceApiResponse)(nil).GetRecords))
}

// MockForceApiEvent is a mock of ForceApiEvent interface.
type MockForceApiEvent struct {
	ctrl     *gomock.Controller
	recorder *MockForceApiEventMockRecorder
}

// MockForceApiEventMockRecorder is the mock recorder for MockForceApiEvent.
type MockForceApiEventMockRecorder struct {
	mock *MockForceApiEvent
}

// NewMockForceApiEvent creates a new mock instance.
func NewMockForceApiEvent(ctrl *gomock.Controller) *MockForceApiEvent {
	mock := &MockForceApiEvent{ctrl: ctrl}
	mock.recorder = &MockForceApiEventMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockForceApiEvent) EXPECT() *MockForceApiEventMockRecorder {
	return m.recorder
}

// AllFields mocks base method.
func (m *MockForceApiEvent) AllFields() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllFields")
	ret0, _ := ret[0].(string)
	return ret0
}

// AllFields indicates an expected call of AllFields.
func (mr *MockForceApiEventMockRecorder) AllFields() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllFields", reflect.TypeOf((*MockForceApiEvent)(nil).AllFields))
}

// GetID mocks base method.
func (m *MockForceApiEvent) GetID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetID indicates an expected call of GetID.
func (mr *MockForceApiEventMockRecorder) GetID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetID", reflect.TypeOf((*MockForceApiEvent)(nil).GetID))
}
