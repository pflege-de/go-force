package force

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

// testInterval is used in place of the production default (2s) so the
// time.Tick-based poll loop in CheckJobStatus runs quickly during tests.
const testInterval = time.Millisecond

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func newTestForceApi(rt http.RoundTripper) *ForceApi {
	return &ForceApi{
		apiVersion:        DefaultAPIVersion,
		instance:          "https://example.my.salesforce.com",
		accessTokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test-token"}),
		httpClient:        &http.Client{Transport: rt},
	}
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func jobInfoJSON(state string) string {
	return fmt.Sprintf(`{"state":%q}`, state)
}

func statusURIFor(jobID string) string {
	return fmt.Sprintf("/services/data/%s/jobs/ingest/%s", DefaultAPIVersion, jobID)
}

func failedResultsURIFor(jobID string) string {
	return statusURIFor(jobID) + "/failedResults"
}

// netError is a stand-in for the kind of transient network failure that
// should trigger CheckJobStatus's retry behavior.
func netError() error {
	return &net.OpError{Op: "read", Net: "tcp", Err: errors.New("connection reset by peer")}
}

func collectingProgressReporter(t *testing.T) (func(string), func() []string) {
	var mu sync.Mutex
	var messages []string
	return func(msg string) {
			mu.Lock()
			defer mu.Unlock()
			messages = append(messages, msg)
			t.Log("state:", msg)
		}, func() []string {
			mu.Lock()
			defer mu.Unlock()
			return append([]string(nil), messages...)
		}
}

func TestForceApi_checkJobStatus(t *testing.T) {
	t.Run("empty JobIDs makes no HTTP calls", func(t *testing.T) {
		fApi := newTestForceApi(roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			t.Fatalf("unexpected HTTP call: %s", r.URL.Path)
			return nil, nil
		}))

		op, err := fApi.CheckJobStatus(JobOperation{}, testInterval)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(op.JobIDs) != 0 {
			t.Fatalf("expected no job IDs, got %v", op.JobIDs)
		}
	})

	t.Run("all normal response states", func(t *testing.T) {
		const jobID = "12341234"
		states := []string{"Open", "InProgress", "UploadComplete", "JobComplete"}
		var calls int32

		fApi := newTestForceApi(roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path != statusURIFor(jobID) {
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
			i := atomic.AddInt32(&calls, 1) - 1
			if int(i) >= len(states) {
				t.Fatalf("unexpected extra request after JobComplete (call #%d)", i+1)
			}
			return jsonResponse(jobInfoJSON(states[i])), nil
		}))

		reporter, messages := collectingProgressReporter(t)
		op, err := fApi.CheckJobStatus(JobOperation{
			JobIDs:           []string{jobID},
			ProgressReporter: reporter,
		}, testInterval)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := len(messages()); got != len(states) {
			t.Fatalf("expected %d progress messages, got %d: %v", len(states), got, messages())
		}
		_ = op
	})

	t.Run("Failed state with InvalidBatch and CSV description returns error", func(t *testing.T) {
		const jobID = "12341234"

		fApi := newTestForceApi(roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			switch r.URL.Path {
			case statusURIFor(jobID):
				return jsonResponse(jobInfoJSON("Failed")), nil
			case failedResultsURIFor(jobID):
				return jsonResponse(`{"error":"InvalidBatch","error_description":"Field name provided in CSV does not match any field"}`), nil
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
				return nil, nil
			}
		}))

		reporter, _ := collectingProgressReporter(t)
		_, err := fApi.CheckJobStatus(JobOperation{
			JobIDs:           []string{jobID},
			ProgressReporter: reporter,
		}, testInterval)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if _, ok := errors.AsType[FailedResultsError](err); !ok {
			t.Fatalf("expected a FailedResultsError, got %T: %v", err, err)
		}
	})

	t.Run("Failed state with non-matching error returns nil", func(t *testing.T) {
		const jobID = "12341234"

		fApi := newTestForceApi(roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			switch r.URL.Path {
			case statusURIFor(jobID):
				return jsonResponse(jobInfoJSON("Failed")), nil
			case failedResultsURIFor(jobID):
				return jsonResponse(`{"error":"SomeOtherError","error_description":"unrelated failure"}`), nil
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
				return nil, nil
			}
		}))

		reporter, _ := collectingProgressReporter(t)
		_, err := fApi.CheckJobStatus(JobOperation{
			JobIDs:           []string{jobID},
			ProgressReporter: reporter,
		}, testInterval)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("retries on net.Error and eventually succeeds", func(t *testing.T) {
		const jobID = "12341234"
		var calls int32

		fApi := newTestForceApi(roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path != statusURIFor(jobID) {
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
			n := atomic.AddInt32(&calls, 1)
			if n <= 2 {
				return nil, netError()
			}
			return jsonResponse(jobInfoJSON("JobComplete")), nil
		}))

		reporter, _ := collectingProgressReporter(t)
		_, err := fApi.CheckJobStatus(JobOperation{
			JobIDs:           []string{jobID},
			ProgressReporter: reporter,
		}, testInterval)
		if err != nil {
			t.Fatalf("expected no error after retrying, got: %v", err)
		}
		if got := atomic.LoadInt32(&calls); got != 3 {
			t.Fatalf("expected exactly 3 calls (2 failures + 1 success), got %d", got)
		}
	})

	t.Run("gives up after retryLimit consecutive net.Errors", func(t *testing.T) {
		const jobID = "12341234"
		var calls int32

		fApi := newTestForceApi(roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path != statusURIFor(jobID) {
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
			atomic.AddInt32(&calls, 1)
			return nil, netError()
		}))

		reporter, _ := collectingProgressReporter(t)
		_, err := fApi.CheckJobStatus(JobOperation{
			JobIDs:           []string{jobID},
			ProgressReporter: reporter,
		}, testInterval)
		if err == nil {
			t.Fatal("expected an error after exhausting retries, got nil")
		}
		if _, ok := errors.AsType[net.Error](err); !ok {
			t.Fatalf("expected the returned error to surface as a net.Error, got %T: %v", err, err)
		}
		const wantCalls = 11 // 10 retries permitted + the 11th failure that gives up
		if got := atomic.LoadInt32(&calls); got != wantCalls {
			t.Fatalf("expected exactly %d calls, got %d", wantCalls, got)
		}
	})

	t.Run("attempts resets after a successful poll", func(t *testing.T) {
		const jobID = "12341234"
		var calls int32

		fApi := newTestForceApi(roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Path != statusURIFor(jobID) {
				t.Fatalf("unexpected request path: %s", r.URL.Path)
			}
			n := atomic.AddInt32(&calls, 1)
			switch {
			case n <= 2:
				return nil, netError() // 2 failures, well under retryLimit
			case n == 3:
				return jsonResponse(jobInfoJSON("InProgress")), nil // success resets attempts
			case n <= 13:
				return nil, netError() // 10 more failures -- would exceed retryLimit if not reset
			default:
				return jsonResponse(jobInfoJSON("JobComplete")), nil
			}
		}))

		reporter, _ := collectingProgressReporter(t)
		_, err := fApi.CheckJobStatus(JobOperation{
			JobIDs:           []string{jobID},
			ProgressReporter: reporter,
		}, testInterval)
		if err != nil {
			t.Fatalf("expected no error (attempts should reset after the successful poll), got: %v", err)
		}
	})

	t.Run("concurrent job IDs retry independently", func(t *testing.T) {
		const (
			failingJobID    = "fail-1234"
			succeedingJobID = "succeed-5678"
		)
		var failingCalls, succeedingCalls int32

		fApi := newTestForceApi(roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			switch r.URL.Path {
			case statusURIFor(failingJobID):
				atomic.AddInt32(&failingCalls, 1)
				return nil, netError()
			case statusURIFor(succeedingJobID):
				n := atomic.AddInt32(&succeedingCalls, 1)
				if n == 1 {
					return nil, netError()
				}
				return jsonResponse(jobInfoJSON("JobComplete")), nil
			default:
				t.Fatalf("unexpected request path: %s", r.URL.Path)
				return nil, nil
			}
		}))

		reporter, _ := collectingProgressReporter(t)
		_, err := fApi.CheckJobStatus(JobOperation{
			JobIDs:           []string{failingJobID, succeedingJobID},
			ProgressReporter: reporter,
		}, testInterval)

		if err == nil {
			t.Fatal("expected the overall error from the failing job ID, got nil")
		}
		if _, ok := errors.AsType[net.Error](err); !ok {
			t.Fatalf("expected the returned error to surface as a net.Error, got %T: %v", err, err)
		}
		const wantFailingCalls = 11
		if got := atomic.LoadInt32(&failingCalls); got != wantFailingCalls {
			t.Fatalf("failing job: expected exactly %d calls, got %d (attempts budget may be shared across job IDs)", wantFailingCalls, got)
		}
		const wantSucceedingCalls = 2
		if got := atomic.LoadInt32(&succeedingCalls); got != wantSucceedingCalls {
			t.Fatalf("succeeding job: expected exactly %d calls, got %d (retry budget may have been starved by the other job ID)", wantSucceedingCalls, got)
		}
	})
}
