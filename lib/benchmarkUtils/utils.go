package benchmarkUtils

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type BenchmarkConfig struct {
	FilesSizesInMb string
	Iterations     string
	RepositoryName string
	Operation      string
	Url            string
	UserName       string
	Password       string
}

func GenerateFiles(numberOfFiles int, sizeOfFilesInMB int) ([]string, error) {
	log.Info("Starting to generate files locally")
	sliceOfFileNames := []string{}
	directoryName := CreateDirectory("/tmp/", "testfiles/")
	for i := 1; i < numberOfFiles+1; i++ {
		fileName := fmt.Sprintf("%v/File%v.txt", directoryName, i)
		log.Info("Genarating file [" + fileName + "] In size of [" + fmt.Sprint(sizeOfFilesInMB) + "MB]")
		file, err := os.Create(fileName)
		if err != nil {
			return nil, errors.New("Failed to create files -" + err.Error())
		}
		defer file.Close()
		data := make([]byte, sizeOfFilesInMB*1024*1024)
		rand.Read(data)
		_, err = file.Write(data)
		if err != nil {
			return nil, errors.New("Failed to insert content into files")
		}
		sliceOfFileNames = append(sliceOfFileNames, fileName)
	}
	log.Info("Sucessfully finished with generating files")
	return sliceOfFileNames, nil
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

func CheckIntLikeString(str string) error {
	value, err := strconv.Atoi(str)
	if err != nil {
		return errors.New("Error: " + str + " is not an integer-like string.")
	}
	if value <= 0 {
		return errors.New("Iterations and size must be positive")
	}
	return nil
}

func ValidateHostByLookup(url string) error {
	testedString := url
	re := regexp.MustCompile(`(/artifactory/)|(/artifactory)`)
	testedString = re.ReplaceAllString(testedString, "")
	if strings.HasSuffix(testedString, "/") {
		testedString = testedString[:len(testedString)-1]
	}
	_, err := net.LookupHost(testedString)
	if err != nil {
		return errors.New("URL [" + url + "] is not valid")
	}
	return nil
}

func IsCustomCredsProvided(cliConfig *BenchmarkConfig) bool {
	if cliConfig.Password != "" && cliConfig.UserName != "" && cliConfig.Url != "" {
		return true
	}
	if cliConfig.Password != "" || cliConfig.UserName != "" || cliConfig.Url != "" {
		log.Error("To use custom server with credentials, you must insert url + username + password ..")
		os.Exit(1)
	}
	return false
}

func UrlStartsWithHttpMethod(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func ValidateInput(cliConfig *BenchmarkConfig) error {
	if IsCustomCredsProvided(cliConfig) {
		err := validateUrlInput(cliConfig)
		if err != nil {
			return err
		}
	}
	StringsIntLikeErr := validateIntStringsLikeInput(cliConfig)
	if StringsIntLikeErr != nil {
		return StringsIntLikeErr
	}
	return nil
}

func validateIntStringsLikeInput(cliConfig *BenchmarkConfig) error {
	filesSizeErr := CheckIntLikeString(cliConfig.FilesSizesInMb)
	if filesSizeErr != nil {
		return filesSizeErr
	}
	iterationsErr := CheckIntLikeString(cliConfig.Iterations)
	if iterationsErr != nil {
		return iterationsErr
	}
	return nil
}

func validateUrlInput(cliConfig *BenchmarkConfig) error {
	if IsCustomCredsProvided(cliConfig) {
		url := cliConfig.Url
		if !UrlStartsWithHttpMethod(url) {
			return errors.New("The url [" + url + "] not starting with http/https")
		}
		err := ValidateHostByLookup(cliConfig.Url)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
