package force

import (
	"github.com/pflege-de/go-force/sobjects"
	"net/http"
	"testing"
)

func TestCheckJobStatus(t *testing.T) {
	forceApi := createTest()
	accObj := insertSAccount(forceApi, t)

	ops := JobOperation{
		Operation: "update",
		Object:    accObj.ApiName(),
		Fields:    []string{"Id"},
		ProgressReporter: func(msg string, bytesTransferred int) {
			t.Logf("BulkOps Update: Account\nstate: %s\nsent bytes: %d", msg, bytesTransferred)
		},
	}
	job := CreateJob(
		JobWithHTTPClient(http.DefaultClient),
		JobWithForceApi(forceApi),
		JobWithApiVersion(DefaultAPIVersion),
		JobWithOperation(ops),
		JobWithMapper(func(objects any) [][]string {
			objs := objects.([]*sobjects.Account)
			records := make([][]string, len(objs))
			for i, o := range objs {
				records[i] = []string{
					o.Id,
				}
			}
			return records
		}),
	)

	err := job.Start()
	if err != nil {
		t.Fatalf("Could not start the job: %v", err)
	}

	t.Log("Job started")

	err = job.Run([]*sobjects.Account{accObj})
	if err != nil {
		t.Fatalf("Could not run the job: %v ", err)
	}

	_, err = forceApi.CheckJobStatus(ops, 3, 0)

	if err != nil {
		deleteSObject(forceApi, t, accObj.Id)
		t.Fatalf("Could not check job status: %v", err)
	}

	deleteSObject(forceApi, t, accObj.Id)

	t.Log("Job finished")
}

func insertSAccount(forceApi *ForceApi, t *testing.T) *sobjects.Account {
	// Need some random text for name field.
	someText := randomString(10)

	// Test Standard Object
	acc := &sobjects.Account{}
	acc.Name = someText

	resp, err := forceApi.InsertSObject(acc)
	if err != nil {
		t.Fatalf("Insert SObject Account failed: %v", err)
	}

	if len(resp.Id) == 0 {
		t.Fatalf("Insert SObject Account failed to return Id: %+v", resp)
	}

	acc.Id = resp.Id

	return acc
}
