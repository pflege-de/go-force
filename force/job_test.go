package force

import (
	"github.com/pflege-de/go-force/sobjects"
	"testing"
)

func TestCheckJobStatus(t *testing.T) {
	fapi := createTest()
	accObj := insertSAccount(fapi, t)

	ops := JobOperation{
		Operation: "update",
		Object:    accObj.ApiName(),
		Fields:    []string{"Id"},
		ProgressReporter: func(msg string) {
			t.Logf("BulkOps Update: Account\nstate: %s", msg)
		},
	}
	job := CreateJob(
		fapi,
		JobWithOperation(ops),
		JobWithMapper(objMapper),
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

	_, err = fapi.CheckJobStatus(ops, 3)

	if err != nil {
		deleteSObject(fapi, t, accObj.Id)
		t.Fatalf("Could not check job status: %v", err)
	}

	deleteSObject(fapi, t, accObj.Id)

	t.Log("Job finished")
}

func objMapper(objects any) [][]string {
	objs := objects.([]*sobjects.Account)
	records := make([][]string, len(objs))
	for i, o := range objs {
		records[i] = []string{
			o.Id,
		}
	}
	return records
}

func insertSAccount(forceApi ForceApiInterface, t *testing.T) *sobjects.Account {
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
