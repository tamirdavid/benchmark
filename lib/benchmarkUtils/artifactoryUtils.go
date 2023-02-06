package benchmarkUtils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	rtUtils "github.com/jfrog/jfrog-cli-core/v2/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/v2/common/commands"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/config"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

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
func CreateLocalRepository(repoName string, servicesManager artifactory.ArtifactoryServicesManager) {
	params := services.NewGenericLocalRepositoryParams()
	params.Key = repoName
	err := servicesManager.CreateLocalRepository().Generic(params)
	if err != nil {
		if strings.Contains(err.Error(), "Case insensitive repository key already exists") {
			ReCreateLocalRepository(repoName, servicesManager)
			return
		}
		log.Error(err)
	}
}

func ReCreateLocalRepository(repoName string, servicesManager artifactory.ArtifactoryServicesManager) {
	log.Info("Recreating [" + repoName + "] Because it is already exists")
	DeleteLocalRepository(repoName, servicesManager)
	CreateLocalRepository(repoName, servicesManager)

}
func DeleteLocalRepository(repoName string, servicesManager artifactory.ArtifactoryServicesManager) {
	err := servicesManager.DeleteRepository(repoName)
	if err != nil {
		log.Error("Not able to delete repository %v", repoName)
		log.Error(err)
		os.Exit(1)
	}
}

func getSvcManagerAfterValidation(serverDetails *config.ServerDetails) artifactory.ArtifactoryServicesManager {
	servicesManager, err := rtUtils.CreateServiceManagerWithThreads(serverDetails, false, 1, 1, 1)
	if err != nil {
		log.Error("Failed to create ServiceManager ", err)
		os.Exit(1)
	}
	version, err := servicesManager.GetVersion()
	if err != nil || version == "" {
		log.Error("")
	}

	return servicesManager
}

func GetSvcManagerBasedOnAuthLogic(c *components.Context, cliConfig *BenchmarkConfig) artifactory.ArtifactoryServicesManager {
	customServer := IsCustomCredsProvided(cliConfig)
	if customServer {
		serverDetails := config.ServerDetails{ArtifactoryUrl: cliConfig.Url, Password: cliConfig.Password, User: cliConfig.UserName}
		serverDetails.ArtifactoryUrl = clientutils.AddTrailingSlashIfNeeded(serverDetails.ArtifactoryUrl)
		serverDetails.ArtifactoryUrl = AddTrailingArtifactoryIfNeeded(serverDetails.ArtifactoryUrl)
		config.CreateInitialRefreshableTokensIfNeeded(&serverDetails)
		serviceManger := getSvcManagerAfterValidation(&serverDetails)
		return serviceManger
	} else {
		confDetails, err := getRtDetails(c)
		if err != nil {
			log.Error("Failed to get server details using default server-id")
		}
		serviceManger := getSvcManagerAfterValidation(confDetails)
		return serviceManger
	}
}

func AddTrailingArtifactoryIfNeeded(url string) string {
	if url != "" && !strings.HasSuffix(url, "artifactory/") {
		url += "artifactory/"
	}
	return url
}

func UploadFiles(fileName string, repositoryName string, servicesManager artifactory.ArtifactoryServicesManager) (timeTaken time.Duration) {
	up := services.NewUploadParams()
	up.CommonParams = &utils.CommonParams{Pattern: filepath.Join(fileName), Recursive: false, Target: repositoryName}
	up.Flat = true
	start := time.Now()
	totalSucceeded, totalFailed, err := servicesManager.UploadFiles(up)
	end := time.Since(start)
	if totalFailed > 0 && totalSucceeded == 0 || err != nil {
		log.Error("Failed to upload the files to artifactory")
		os.Exit(1)
	}
	return end
}

func DownloadFiles(fileName string, repositoryName string, servicesManager artifactory.ArtifactoryServicesManager) (timeTaken time.Duration) {
	dl := services.NewDownloadParams()
	dl.CommonParams = &utils.CommonParams{Pattern: filepath.Join(filepath.Base(fileName)), Recursive: false, Target: repositoryName}
	start := time.Now()
	totalSucceeded, totalFailed, err := servicesManager.DownloadFiles(dl)
	end := time.Since(start)
	if totalFailed > 0 && totalSucceeded == 0 || err != nil {
		log.Error("Failed to download files from Artifactory")
		os.Exit(1)
	}
	return end
}

func DeleteRepository(repo string, servicesManager artifactory.ArtifactoryServicesManager) {
	log.Info("Cleanup created resources [repositories/binaries]")
	err := servicesManager.DeleteRepository(repo)
	if err != nil {
		log.Error("Failed to delete repository ["+repo+"]", err)
		os.Exit(1)
	}
}
