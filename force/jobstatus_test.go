package force_test

import (
	"github.com/pflege-de/go-force/force"
	"github.com/pflege-de/go-force/force/mocks"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=JobTypes.go -destination=mocks/JobTypes.go -package mocks

func TestForceApi_checkJobStatus(t *testing.T) {
	ctrl := gomock.NewController(t)

	type fields struct {
		fApi func() force.ForceApiSObjectInterface
	}
	type args struct {
		op       force.JobOperation
		duration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test Run with empty Ops",
			fields: fields{
				fApi: func() force.ForceApiSObjectInterface {
					fApi := mocks.NewMockForceApiSObjectInterface(ctrl)
					fApi.EXPECT().CheckJobStatus(gomock.Any(), gomock.Any()).Return(force.JobOperation{}, nil)
					return fApi
				},
			},
			args: args{
				op:       force.JobOperation{},
				duration: 3,
			},
			wantErr: false,
		},
		{
			name: "Run with all normal response states",
			fields: fields{
				fApi: func() force.ForceApiSObjectInterface {
					fApi := mocks.NewMockForceApiSObjectInterface(ctrl)
					var jobInfo *force.JobInfo
					i := 0
					fApi.EXPECT().Get("/services/data/"+force.DefaultAPIVersion+"/jobs/ingest/12341234", gomock.Any(), gomock.AssignableToTypeOf(jobInfo)).DoAndReturn(func(url string, params any, jobinfo *force.JobInfo) error {
						states := []string{
							"Open",
							"InProgress",
							"UploadComplete",
							"JobComplete",
						}
						jobinfo.State = states[i]
						i++
						return nil
					}).AnyTimes()
					fApi.EXPECT().CheckJobStatus(gomock.Any(), gomock.Any()).Return(force.JobOperation{}, nil)
					return fApi
				},
			},
			args: args{
				op: force.JobOperation{
					JobIDs: []string{"12341234"},
					ProgressReporter: func(msg string) {
						t.Log("state: ", msg)
					},
				},
				duration: 3,
			},
			wantErr: false,
		},
		{
			name: "Run with unknown response state",
			fields: fields{
				fApi: func() force.ForceApiSObjectInterface {
					fApi := mocks.NewMockForceApiSObjectInterface(ctrl)
					var jobInfo *force.JobInfo
					i := 0
					fApi.EXPECT().Get("/services/data/"+force.DefaultAPIVersion+"/jobs/ingest/12341234", gomock.Any(), gomock.AssignableToTypeOf(jobInfo)).DoAndReturn(func(url string, params any, jobinfo *force.JobInfo) error {
						states := []string{
							"Unknown",
							"Failed",
						}
						jobinfo.State = states[i]
						i++
						return nil
					}).AnyTimes()
					failedResultErr := force.FailedResultsError{}
					fApi.EXPECT().Get("/services/data/"+force.DefaultAPIVersion+"/jobs/ingest/12341234/failedResults", gomock.Any(), gomock.AssignableToTypeOf(failedResultErr)).Return(nil)
					fApi.EXPECT().CheckJobStatus(gomock.Any(), gomock.Any()).Return(force.JobOperation{}, nil)
					return fApi
				},
			},
			args: args{
				op: force.JobOperation{
					JobIDs: []string{"12341234"},
					ProgressReporter: func(msg string) {
						t.Log("state: ", msg)
					},
				},
				duration: 3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.fields.fApi().CheckJobStatus(tt.args.op, tt.args.duration); (err != nil) != tt.wantErr {
				t.Errorf("SalesCareBearService.checkJobStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
