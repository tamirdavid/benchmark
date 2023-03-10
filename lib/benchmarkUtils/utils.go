package benchmarkUtils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
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
	Append         string
	SameFile       bool
}

func GenerateFiles(numberOfFiles int, sizeOfFilesInMB int, sameFile bool) ([]string, error) {
	log.Info("Starting to generate files locally")
	sliceOfFileNames := []string{}
	directoryName := CreateDirectory("/tmp/", "testfiles/")
	for i := 1; i < numberOfFiles+1; i++ {
		if sameFile && i != 1 {
			sliceOfFileNames = append(sliceOfFileNames, fmt.Sprintf("%v/File%v.txt", directoryName, 1))
			continue
		}
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

func MeasureOperationTimes(st *BenchmarkConfig, fileNames []string, servicesManager artifactory.ArtifactoryServicesManager,
	benchmarkResults *[]BenchmarkResult) error {
	// firstFile := ""
	for _, file := range fileNames {
		// if sameFile {
		// 	if i == 0 {
		// 		firstFile = file
		// 	}
		// 	file = firstFile
		// }
		if st.Operation == "upload" {
			uploadError := MeasureSingleOperation(file, st, servicesManager, *&benchmarkResults, UploadFiles)
			if uploadError != nil {
				return uploadError
			}
		}
		if st.Operation == "download" {
			downloadError := MeasureSingleOperation(file, st, servicesManager, *&benchmarkResults, DownloadFiles)
			if downloadError != nil {
				return downloadError
			}
		}
	}
	return nil
}

type runFunc func(fileName string, repositoryName string, servicesManager artifactory.ArtifactoryServicesManager) (time.Duration, error)

func MeasureSingleOperation(file string, st *BenchmarkConfig, serviceManager artifactory.ArtifactoryServicesManager,
	benchmarkResults *[]BenchmarkResult, operation runFunc) error {
	duration, downloadError := operation(file, st.RepositoryName, serviceManager)
	if downloadError != nil {
		return downloadError
	}
	sizeMbIntFormat, _ := strconv.Atoi(st.FilesSizesInMb)
	uploadedMB := int64(sizeMbIntFormat)
	speed := float64(uploadedMB) / duration.Seconds()
	*benchmarkResults = append(*benchmarkResults, *NewBenchmarkResult(file, st.FilesSizesInMb, fmt.Sprintf("%s", duration), fmt.Sprintf("%.2f", speed)))
	return nil
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

func IsCustomCredsProvided(cliConfig *BenchmarkConfig) (bool, error) {
	if cliConfig.Password != "" && cliConfig.UserName != "" && cliConfig.Url != "" {
		return true, nil
	}
	if cliConfig.Password != "" || cliConfig.UserName != "" || cliConfig.Url != "" {
		return false, errors.New("To use custom server with credentials, you must insert url + username + password ..")
	}
	return false, nil
}

func UrlStartsWithHttpMethod(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func ValidateInput(cliConfig *BenchmarkConfig) error {
	isCustomCredsProvided, customCredsErr := IsCustomCredsProvided(cliConfig)
	if customCredsErr != nil {
		return customCredsErr
	}
	if isCustomCredsProvided {
		err := validateUrlInput(cliConfig)
		if err != nil {
			return err
		}
	}
	StringsIntLikeErr := validateIntStringsLikeInput(cliConfig)
	if StringsIntLikeErr != nil {
		return StringsIntLikeErr
	}
	RepoNameNotValidError := ValidateRepoNameInput(cliConfig.RepositoryName)
	if RepoNameNotValidError != nil {
		return RepoNameNotValidError
	}
	if cliConfig.Append != "" {
		_, err := FileExists(cliConfig.Append)
		if err != nil {
			return err
		}
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
	CustomCredsProvided, customCredsErr := IsCustomCredsProvided(cliConfig)
	if customCredsErr != nil {
		return customCredsErr
	}
	if CustomCredsProvided {
		url := cliConfig.Url
		if !UrlStartsWithHttpMethod(url) {
			return errors.New("The url [" + url + "] not starting with http/https")
		}
		err := ValidateUrlUsingReadiness(cliConfig.Url)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func ValidateRepoNameInput(s string) error {
	match, _ := regexp.MatchString("^[a-zA-Z0-9-._]+$", s)
	if !match {
		return errors.New("Repository name can containletters, numbers, dashes, dots, and underscores only")
	}
	if len(s) > 63 {
		return errors.New("Repository name must be maximum length of 64 characters")
	}
	return nil
}

func CleanupCliResources(config *BenchmarkConfig, servicesManager artifactory.ArtifactoryServicesManager) error {
	log.Info("Starting to cleanup CLI created resources")
	deleteRepoError := DeleteRepository(config.RepositoryName, servicesManager)
	if deleteRepoError != nil {
		return deleteRepoError
	}
	deleteFilesError := DeleteLocalFilesAndTestDirectory(config.Iterations, config.SameFile)
	if deleteFilesError != nil {
		return deleteFilesError
	}
	log.Info("Finished cleanup CLI resources")
	return nil
}

func DeleteLocalFilesAndTestDirectory(Iterations string, sameFile bool) error {
	log.Info("Deleting files generated for test")
	IterationsInt, _ := strconv.Atoi(Iterations)
	for i := 1; i < IterationsInt+1; i++ {
		if sameFile && i != 1 {
			continue
		}
		fileName := fmt.Sprintf("/tmp/testfiles/File%v.txt", i)
		removingErr := os.Remove(fileName)
		if removingErr != nil {
			return removingErr
		}
	}
	deletingDirError := os.RemoveAll("/tmp/testfiles")
	if deletingDirError != nil {
		return deletingDirError
	}
	return nil
}

func ReadFileAndPrint(filename string) error {
	pwd, _ := os.Getwd()
	log.Info("Read [" + pwd + "/" + filename + "]")
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
