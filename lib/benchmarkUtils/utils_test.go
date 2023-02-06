package benchmarkUtils

import (
	"strings"
	"testing"
)

var generateFilesProvider = []struct {
	numberOfFiles   int
	sizeOfFilesInMB int
}{
	{1, 2},
	{3, 4},
}

func TestGenerateFiles(t *testing.T) {
	for _, sample := range generateFilesProvider {
		t.Logf("Starting to generate %v files in size of %v", sample.numberOfFiles, sample.sizeOfFilesInMB)
		results, err := GenerateFiles(sample.numberOfFiles, sample.sizeOfFilesInMB)
		if err != nil {
			t.Error(err)
		}
		if len(results) > 0 {
			t.Logf("Generating %v file of in size of %v succeed", sample.numberOfFiles, sample.sizeOfFilesInMB)
		}
	}
}

var checkIntLikeStringProvider = []struct {
	numString string
}{
	{"8"},
	{"9"},
	{"-1"},
}

func TestCheckIntLikeString(t *testing.T) {
	for _, sample := range checkIntLikeStringProvider {
		err := CheckIntLikeString(sample.numString)
		if err != nil {
			if strings.Contains(err.Error(), "must be positive") {
				t.Logf("Test for positive number PASSED")
			} else {
				t.Error(err)
			}
		}
	}
}

var validateHostByLookupProvider = []struct {
	url string
}{
	{"www.moshe.com"},
	{"https://www.moshe.com"},
	{"https://productionautomation.jfrog.io"},
	{"https://mohammadt.jfrog.io"},
	{"https://mohammadt.jfrog.io/"},
	{"https://mohammadt.jfrog.io/artifactory/"},
	{"https://mohammadt.jfrog.io/artifactory"},
}

func TestValidateHostByLookup(t *testing.T) {
	for _, sample := range validateHostByLookupProvider {
		err := ValidateHostByLookup(sample.url)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}
