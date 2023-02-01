package benchmarkUtils

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type BenchmarkConfig struct {
	FilesSizesInMb string
	Iterations     string
	RepositoryName string
	Operation      string
}

func GenerateFiles(numberOfFiles int, sizeOfFilesInMB int) []string {
	log.Info("Starting to generate files locally")
	sliceOfFileNames := []string{}
	directoryName := CreateDirectory("/tmp/", "testfiles/")
	for i := 1; i < numberOfFiles+1; i++ {
		fileName := fmt.Sprintf("%v/File%v.txt", directoryName, i)
		log.Info("Genarating file [" + fileName + "] In size of [" + fmt.Sprint(sizeOfFilesInMB) + "MB]")
		file, err := os.Create(fileName)
		if err != nil {
			log.Error("Failed to create files -", err)
			os.Exit(1)
		}
		defer file.Close()
		data := make([]byte, sizeOfFilesInMB*1024*1024)
		rand.Read(data)
		_, err = file.Write(data)
		if err != nil {
			log.Error("Failed to insert content into files")
			os.Exit(1)
		}
		sliceOfFileNames = append(sliceOfFileNames, fileName)
	}
	log.Info("Sucessfully finished with generating files")
	return sliceOfFileNames
}

func MeasureOperationTimes(st *BenchmarkConfig, fileNames []string, servicesManager artifactory.ArtifactoryServicesManager) map[string]time.Duration {
	results := make(map[string]time.Duration)
	for _, file := range fileNames {
		if st.Operation == "upload" {
			results[file] = UploadFiles(file, st.RepositoryName, servicesManager)
		}
		if st.Operation == "download" {
			results[file] = DownloadFiles(file, st.RepositoryName, servicesManager)
		}
	}
	return results
}

func WriteResult(filePath string, results map[string]time.Duration, size string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintln(writer, "file,size,time_taken")
	for file, timeTaken := range results {
		fmt.Fprintf(writer, "%s,%s,%s\n", filepath.Base(file), size+"MB", timeTaken)
	}
}

func CreateDirectory(path string, dirName string) string {
	newPath := filepath.Join(path, dirName)
	err := os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		log.Error("Not able to create the directory [" + path + dirName + "]")
		os.Exit(1)
	}
	return newPath
}

func CheckIntLikeString(str string) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		log.Error("Error: " + str + " is not an integer-like string.")
		os.Exit(1)
	}
	return value
}
