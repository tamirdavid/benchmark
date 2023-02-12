package benchmarkUtils

import (
	"bufio"
	"os"
	"testing"
)

func TestNewBenchmarkResult(t *testing.T) {
	// Define test inputs
	file := "testfile.dat"
	size := "1"
	duration := "1.23"
	speed := "800"

	// Call the constructor being tested
	result := NewBenchmarkResult(file, size, duration, speed)

	// Check if the fields are set correctly
	if result.FileName != file {
		t.Errorf("Expected file name '%s', but got '%s'", file, result.FileName)
	}
	if result.Size != size {
		t.Errorf("Expected size '%s', but got '%s'", size, result.Size)
	}
	if result.Duration != duration {
		t.Errorf("Expected duration '%s', but got '%s'", duration, result.Duration)
	}
	if result.Speed != speed {
		t.Errorf("Expected speed '%s', but got '%s'", speed, result.Speed)
	}
}

func TestWriteResults(t *testing.T) {
	// Define test inputs
	filePath := "results.csv"
	results := []BenchmarkResult{
		{"file1.dat", "1", "1.23", "800"},
		{"file2.dat", "2", "2.34", "900"},
	}

	// Call the function being tested
	err := WriteResults(filePath, results)

	// Check if there is no error
	if err != nil {
		t.Errorf("Expected no error, but got '%v'", err)
	}

	// Read the written file and check its contents
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Cannot open the written file '%s'", filePath)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n') // read the first line containing column names
	if err != nil {
		t.Fatalf("Cannot read the first line of the written file '%s'", filePath)
	}
	expectedHeader := "file,size (MB),time taken (sec),speed (MB/sec)\n"
	if line != expectedHeader {
		t.Errorf("Expected header '%s', but got '%s'", expectedHeader, line)
	}
	line, err = reader.ReadString('\n') // read the first data line
	if err != nil {
		t.Fatalf("Cannot read the data from the written file '%s'", filePath)
	}
	expectedData := "file1.dat,1,1.23,800\n"
	if line != expectedData {
		t.Errorf("Expected data line '%s', but got '%s'", expectedData, line)
	}
	line, err = reader.ReadString('\n') // read the second data line
	if err != nil {
		t.Fatalf("Cannot read the data from the written file '%s'", filePath)
	}
	expectedData = "file2.dat,2,2.34,900\n"
	if line != expectedData {
		t.Errorf("Expected data line '%s', but got '%s'", expectedData, line)
	}
	err = os.Remove("results.csv")
	if err != nil {
		t.Errorf("not able to delete the csv file")
	}
}
