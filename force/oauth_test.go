package force

import (
	"testing"
)

func TestOauth(t *testing.T) {
	forceApi := createTest()
	// Verify oauth object is valid
	if err := forceApi.GetOauth().Validate(); err != nil {
		t.Fatalf("Oauth object is invlaid: %#v", err)
	}
}
