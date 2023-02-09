package benchmarkUtils

import (
	"bufio"
	"fmt"
	"os"
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
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintln(writer, finalResults.ColumnNames)
	for _, result := range results {
		fmt.Fprintf(writer, "%s,%s,%s,%s\n", result.FileName, result.Size, result.Duration, result.Speed)
	}
	return nil
}
