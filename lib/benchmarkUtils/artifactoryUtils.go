package benchmarkUtils

import (
	"errors"
	"net/http"
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

func CreateLocalRepository(repoName string, servicesManager artifactory.ArtifactoryServicesManager) error {
	params := services.NewGenericLocalRepositoryParams()
	params.Key = repoName
	err := servicesManager.CreateLocalRepository().Generic(params)
	if err != nil {
		if strings.Contains(err.Error(), "Case insensitive repository key already exists") {
			recreateError := ReCreateLocalRepository(repoName, servicesManager)
			if recreateError != nil {
				return recreateError
			}
		}
		return err
	}
	return nil
}

func ReCreateLocalRepository(repoName string, servicesManager artifactory.ArtifactoryServicesManager) error {
	log.Info("Recreating [" + repoName + "] Because it is already exists")
	deleteError := DeleteLocalRepository(repoName, servicesManager)
	if deleteError != nil {
		return deleteError
	}
	createRepoErr := CreateLocalRepository(repoName, servicesManager)
	if createRepoErr != nil {
		return createRepoErr
	}
	return nil
}
func DeleteLocalRepository(repoName string, servicesManager artifactory.ArtifactoryServicesManager) error {
	deleteRepositoryErr := servicesManager.DeleteRepository(repoName)
	if deleteRepositoryErr != nil {
		log.Error("Not able to delete repository %v", repoName)
		return deleteRepositoryErr
	}
	return nil
}

func getSvcManagerAfterValidation(serverDetails *config.ServerDetails) (artifactory.ArtifactoryServicesManager, error) {
	servicesManager, err := rtUtils.CreateServiceManagerWithThreads(serverDetails, false, 1, 1, 1)
	if err != nil {
		log.Error("Failed to create ServiceManager ", err)
		return nil, err
	}
	version, err := servicesManager.GetVersion()
	if err != nil || version == "" {
		log.Error("Failed to use serviceManager to get server version")
		return nil, err
	}

	return servicesManager, nil
}

func GetSvcManagerBasedOnAuthLogic(c *components.Context, cliConfig *BenchmarkConfig) (artifactory.ArtifactoryServicesManager, error) {
	customServer := IsCustomCredsProvided(cliConfig)
	if customServer {
		serverDetails := config.ServerDetails{ArtifactoryUrl: cliConfig.Url, Password: cliConfig.Password, User: cliConfig.UserName}
		serverDetails.ArtifactoryUrl = clientutils.AddTrailingSlashIfNeeded(serverDetails.ArtifactoryUrl)
		serverDetails.ArtifactoryUrl = AddTrailingArtifactoryIfNeeded(serverDetails.ArtifactoryUrl)
		tokenError := config.CreateInitialRefreshableTokensIfNeeded(&serverDetails)
		if tokenError != nil {
			return nil, tokenError
		}
		serviceManger, serviceMngrErr := getSvcManagerAfterValidation(&serverDetails)
		if serviceMngrErr != nil {
			return nil, serviceMngrErr
		}
		return serviceManger, nil
	} else {
		confDetails, err := getRtDetails(c)
		if err != nil {
			log.Error("Failed to get server details using default server-id")
			return nil, err
		}
		serviceManger, serviceMngrErr := getSvcManagerAfterValidation(confDetails)
		if serviceMngrErr != nil {
			return nil, serviceMngrErr
		}
		return serviceManger, nil
	}
}

func AddTrailingArtifactoryIfNeeded(url string) string {
	if url != "" && !strings.HasSuffix(url, "artifactory/") {
		url += "artifactory/"
	}
	return url
}

func UploadFiles(fileName string, repositoryName string, servicesManager artifactory.ArtifactoryServicesManager) (time.Duration, error) {
	up := services.NewUploadParams()
	up.CommonParams = &utils.CommonParams{Pattern: filepath.Join(fileName), Recursive: false, Target: repositoryName}
	up.Flat = true
	start := time.Now()
	totalSucceeded, totalFailed, err := servicesManager.UploadFiles(up)
	end := time.Since(start)
	if totalFailed > 0 && totalSucceeded == 0 || err != nil {
		return 0, err
	}
	return end, nil
}

func DownloadFiles(fileName string, repositoryName string, servicesManager artifactory.ArtifactoryServicesManager) (time.Duration, error) {
	dl := services.NewDownloadParams()
	dl.CommonParams = &utils.CommonParams{Pattern: filepath.Join(filepath.Base(fileName)), Recursive: false, Target: repositoryName}
	start := time.Now()
	totalSucceeded, totalFailed, err := servicesManager.DownloadFiles(dl)
	end := time.Since(start)
	if totalFailed > 0 && totalSucceeded == 0 || err != nil {
		return 0, errors.New("Failed to download files from Artifactory")
	}
	return end, nil
}

func DeleteRepository(repo string, servicesManager artifactory.ArtifactoryServicesManager) error {
	log.Info("Deleting the repository " + repo)
	err := servicesManager.DeleteRepository(repo)
	if err != nil {
		log.Error("Failed to delete repository ["+repo+"]", err)
		return err
	}
	return nil
}

func ValidateUrlUsingReadiness(url string) error {
	readinessEndpoint := getReadinessEndpointPerUrl(url)
	log.Info("Validate url is an Artifactory server by sending readiness request")
	resp, err := http.Get(url + readinessEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Readiness" + url + readinessEndpoint + " failed")
	}
	return nil
}

func getReadinessEndpointPerUrl(url string) string {
	if strings.HasSuffix(url, "/artifactory") {
		return "/api/v1/system/readiness"
	} else if strings.HasSuffix(url, "/artifactory/") {
		return "api/v1/system/readiness"
	} else if strings.HasSuffix(url, "/") {
		return "artifactory/api/v1/system/readiness"
	} else {
		return "/artifactory/api/v1/system/readiness"
	}
}
