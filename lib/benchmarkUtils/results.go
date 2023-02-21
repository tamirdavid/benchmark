package benchmarkUtils

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Duration = time took for download/upload operations
// Speed = MB / sec

type BenchMarkResults struct {
	Results     []BenchmarkResult
	ColumnNames string
}

type BenchmarkResult struct {
	FileName string
	Size     string
	Duration string
	Speed    string
}

func NewBenchMarkResults(results []BenchmarkResult) *BenchMarkResults {
	return &BenchMarkResults{Results: results, ColumnNames: "file,size (MB),time taken (sec),speed (MB/sec)"}
}

func NewBenchmarkResult(file string, size string, duration string, speed string) *BenchmarkResult {
	return &BenchmarkResult{FileName: file, Size: size, Duration: duration, Speed: speed}
}

func WriteResults(filePath string, results []BenchmarkResult) error {
	finalResults := NewBenchMarkResults(results)
	var err error
	var file *os.File
	var writer *bufio.Writer

	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		writer = bufio.NewWriter(file)
		defer writer.Flush()
		fmt.Fprintln(writer, finalResults.ColumnNames)
	} else { // if file already exists append the results to it
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		writer = bufio.NewWriter(file)
		defer writer.Flush()
	}

	for _, result := range results {
		fmt.Fprintf(writer, "%s,%s,%s,%s\n", result.FileName, result.Size, result.Duration, result.Speed)
	}
	return nil
}

func GetFilePath(operation string, append string) string {
	if append != "" {
		return append
	} else {
		return fmt.Sprintf("benchmark-%s-%s.csv", operation, time.Now().Format("2006-01-02T15:04:05"))
	}
}

func FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, fmt.Errorf("File not exists not able to append the results")
	}
	return true, nil
}
