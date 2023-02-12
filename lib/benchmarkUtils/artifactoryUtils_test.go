package benchmarkUtils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ValidateUrlInputFailuresProvider = []struct {
	Url string
}{
	{"https://jfrog.com"},
	{"https://google.com"},
	{"https://yahoo.com"},
}

func TestValidateUrlInput(t *testing.T) {
	for _, sample := range ValidateUrlInputFailuresProvider {
		err := ValidateUrlUsingReadiness(sample.Url)
		if err != nil {
			if strings.Contains(err.Error(), "failed") && strings.Contains(err.Error(), "Readiness") {
				t.Log(sample.Url + " readiness failure check succeed")
				continue
			}
		}
		t.Error("Readiness failure check has been failed.")
	}
}

func TestGetReadinessEndpointPerUrl(t *testing.T) {
	assert.Equal(t, GetReadinessEndpointPerUrl("https://tamir_test.jfrog.io"), "/artifactory/api/v1/system/readiness")
	assert.Equal(t, GetReadinessEndpointPerUrl("https://tamir_test.jfrog.io/"), "artifactory/api/v1/system/readiness")
	assert.Equal(t, GetReadinessEndpointPerUrl("https://tamir_test.jfrog.io/artifactory"), "/api/v1/system/readiness")
	assert.Equal(t, GetReadinessEndpointPerUrl("https://tamir_test.jfrog.io/artifactory/"), "api/v1/system/readiness")
}
