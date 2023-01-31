package commands

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	rtUtils "github.com/jfrog/jfrog-cli-core/v2/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/v2/common/commands"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type BenchmarkConfig struct {
	FilesSizesInMb string
	Iterations     string
	repositoryName string
	operation      string
}

// Returns the Artifactory Details of the provided server-id, or the default one.
func getRtDetails(c *components.Context) (*config.ServerDetails, error) {
	details, err := commands.GetConfig(c.GetStringFlagValue("server-id"), false)
	if err != nil {
		return nil, err
	}
	if details.ArtifactoryUrl == "" {
		return nil, errors.New("no server-id was found, or the server-id has no Artifactory url.")
	}
	details.ArtifactoryUrl = clientutils.AddTrailingSlashIfNeeded(details.ArtifactoryUrl)
	err = config.CreateInitialRefreshableTokensIfNeeded(details)
	if err != nil {
		return nil, err
	}
	return details, nil
}

// Create local repository
func CreateLocalRepository(c *components.Context, repoName string) {
	confDetails, _ := getRtDetails(c)
	servicesManager, err := rtUtils.CreateServiceManager(confDetails, -1, 0, false)
	if err != nil {
		log.Error(err)
	}
	params := services.NewGenericLocalRepositoryParams()
	params.Key = repoName
	err = servicesManager.CreateLocalRepository().Generic(params)
	if err != nil {
		if strings.Contains(err.Error(), "Case insensitive repository key already exists") {
			ReCreateLocalRepository(c, repoName)
		}
	}
}

func ReCreateLocalRepository(c *components.Context, repoName string) {
	log.Info("Recreating [" + repoName + "] Because it is already exists")
	DeleteLocalRepository(c, repoName)
	CreateLocalRepository(c, repoName)

}
func DeleteLocalRepository(c *components.Context, repoName string) {
	confDetails, _ := getRtDetails(c)
	servicesManager, err := rtUtils.CreateServiceManager(confDetails, -1, 0, false)

	if err != nil {
		log.Error(err)
	}
	err = servicesManager.DeleteRepository(repoName)
	if err != nil {
		log.Error("Not able to delete repository %v", repoName)
		log.Error(err)
	}
}

func GenerateFiles(numberOfFiles int, sizeOfFilesInMB int) []string {
	log.Info("Starting to generate files locally")
	sliceOfFileNames := []string{}
	for i := 1; i < numberOfFiles+1; i++ {
		fileName := fmt.Sprintf("File%v.txt", i)
		log.Info("Genarate file [" + fileName + "] In size of [" + fmt.Sprint(sizeOfFilesInMB) + "MB]")
		file, err := os.Create(fileName)
		if err != nil {
			log.Error("Failed to generate files")
			os.Exit(1)
		}
		defer file.Close()
		data := make([]byte, sizeOfFilesInMB*1024*1024)
		rand.Read(data)
		_, err = file.Write(data)
		if err != nil {
			log.Error("Failed to generate files")
			os.Exit(1)
		}
		sliceOfFileNames = append(sliceOfFileNames, fileName)
	}
	log.Info("Sucessfully finished with generating files")
	return sliceOfFileNames
}

func MeasureOperationTimes(c *components.Context, st *BenchmarkConfig, fileNames []string) map[string]time.Duration {
	confDetails, _ := getRtDetails(c)
	results := make(map[string]time.Duration)
	for _, file := range fileNames {
		if st.operation == "upload" {
			results[file] = UploadFiles(confDetails, file, st.repositoryName)
		}
		if st.operation == "download" {
			results[file] = DownloadFiles(confDetails, file, st.repositoryName)
		}
	}
	return results
}

func WriteResult(filePath string, results map[string]time.Duration) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintln(writer, "file,time_taken")
	for file, timeTaken := range results {
		fmt.Fprintf(writer, "%s,%s\n", file, timeTaken)
	}
}

func UploadFiles(confDetails *config.ServerDetails, fileName string, repositoryName string) (timeTaken time.Duration) {
	servicesManager, _ := rtUtils.CreateServiceManagerWithThreads(confDetails, false, 1, 1, 1)
	up := services.NewUploadParams()
	up.CommonParams = &utils.CommonParams{Pattern: filepath.Join(fileName), Recursive: true, Target: repositoryName}
	start := time.Now()
	totalSucceeded, totalFailed, err := servicesManager.UploadFiles(up)
	end := time.Since(start)
	if totalFailed > 0 && totalSucceeded == 0 || err != nil {
		log.Error("Failed to upload the files to artifactory")
		os.Exit(1)
	}
	return end
}

func DownloadFiles(confDetails *config.ServerDetails, fileName string, repositoryName string) (timeTaken time.Duration) {
	servicesManager, _ := rtUtils.CreateServiceManagerWithThreads(confDetails, false, 1, 1, 1)
	dl := services.NewDownloadParams()
	dl.CommonParams = &utils.CommonParams{Pattern: filepath.Join(fileName), Recursive: true, Target: repositoryName}
	start := time.Now()
	totalSucceeded, totalFailed, err := servicesManager.DownloadFiles(dl)
	fmt.Printf(string(totalSucceeded))
	end := time.Since(start)
	if totalFailed > 0 && totalSucceeded == 0 || err != nil {
		log.Error("Failed to download the files from Artifactory")
		os.Exit(1)
	}
	return end
}
