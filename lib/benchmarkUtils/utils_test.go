package benchmarkUtils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestUrlStartWithHttpMethod(t *testing.T) {
	assert.Equal(t, UrlStartsWithHttpMethod("https://tamir_test.jfrog.io/artifactory/"), true)
	assert.Equal(t, UrlStartsWithHttpMethod("tamir_test.jfrog.io/artifactory/"), false)
}

func TestValidateRepoNameInput(t *testing.T) {
	assert.Equal(t, ValidateRepoNameInput("ThisIsTooLongRepositoryNameToBeInsertedForTheCliItIsSupposeFailedDuringValidation"),
		errors.New("Repository name must be maximum length of 64 characters"))
	assert.Equal(t, ValidateRepoNameInput("*₪%^₪%^$$$#$#$#$#$"),
		errors.New("Repository name can containletters, numbers, dashes, dots, and underscores only"))
}

func TestIsCustomCredsProvided(t *testing.T) {
	failureConf := BenchmarkConfig{UserName: "", Password: "", Url: "https://tamirtest.jfrog.io"}
	value, err := IsCustomCredsProvided(&failureConf)
	assert.Equal(t, err, errors.New("To use custom server with credentials, you must insert url + username + password .."))
	assert.Equal(t, value, false)
	sucessConf := BenchmarkConfig{UserName: "tamir", Password: "passwordpassword", Url: "https://tamirtest.jfrog.io"}
	value2, err2 := IsCustomCredsProvided(&sucessConf)
	assert.Equal(t, err2, nil)
	assert.Equal(t, value2, true)
}

func TestCreateDirectory(t *testing.T) {
	testPath := "./testDir"
	testDirName := "newDir"
	os.RemoveAll(testPath)
	newPath := CreateDirectory(testPath, testDirName)
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		t.Errorf("Expected directory to be created at %s, but it does not exist", newPath)
	}

	os.RemoveAll(testPath)
}

func TestDeleteLocalFiles(t *testing.T) {
	// Create test directory and files
	err := os.MkdirAll("/tmp/testfiles", os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	for i := 1; i <= 3; i++ {
		fileName := fmt.Sprintf("/tmp/testfiles/File%v.txt", i)
		err := ioutil.WriteFile(fileName, []byte("test"), os.ModePerm)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Call the function being tested
	err = DeleteLocalFilesAndTestDirectory("3")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Check if the test directory and files were deleted
	if _, err := os.Stat("/tmp/testfiles"); !os.IsNotExist(err) {
		t.Error("Expected test directory to be deleted, but it still exists")
	}
	for i := 1; i <= 3; i++ {
		fileName := fmt.Sprintf("/tmp/testfiles/File%v.txt", i)
		if _, err := os.Stat(fileName); !os.IsNotExist(err) {
			t.Errorf("Expected test file %s to be deleted, but it still exists", fileName)
		}
	}
}
func TestReadFileAndPrint(t *testing.T) {
	// Define test variables
	testFile := "test.txt"
	testContent := "This is a test file."

	// Create a test file with some content
	err := ioutil.WriteFile(testFile, []byte(testContent), os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(testFile)

	// Redirect stdout to a file
	outputFile, err := os.Create("output.txt")
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer os.Remove("output.txt")
	defer outputFile.Close()
	old := os.Stdout
	os.Stdout = outputFile
	defer func() { os.Stdout = old }()

	// Call the function being tested
	err = ReadFileAndPrint(testFile)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Check if the output is correct
	expectedOutput := testContent + "\n"
	outputBytes, err := ioutil.ReadFile("output.txt")
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	if string(outputBytes) != expectedOutput {
		t.Errorf("Expected output '%s', but got '%s'", expectedOutput, outputBytes)
	}
}
