package commands

import (
	"benchmark/lib/benchmarkUtils"
	"strconv"

	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func DownloadCommand() components.Command {
	return components.Command{
		Name:        "dl",
		Description: "Download artifacts tests",
		Flags:       DownloadCommandFlags(),
		Action: func(c *components.Context) error {
			downloadConfig, err := setDownloadConfig(c)
			if err != nil {
				return err
			}
			return dlCmd(c, downloadConfig)
		},
	}
}

func setDownloadConfig(c *components.Context) (*benchmarkUtils.BenchmarkConfig, error) {
	var downloadConfig = new(benchmarkUtils.BenchmarkConfig)
	downloadConfig.FilesSizesInMb = c.GetStringFlagValue("size")
	downloadConfig.Iterations = c.GetStringFlagValue("iterations")
	downloadConfig.RepositoryName = c.GetStringFlagValue("repo_name")
	downloadConfig.Operation = "download"
	downloadConfig.Url = c.GetStringFlagValue("url")
	downloadConfig.UserName = c.GetStringFlagValue("username")
	downloadConfig.Password = c.GetStringFlagValue("password")
	downloadConfig.Append = c.GetStringFlagValue("append")
	err := benchmarkUtils.ValidateInput(downloadConfig)
	if err != nil {
		return nil, err
	}
	return downloadConfig, nil
}

func DownloadCommandFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:         "size",
			Description:  "The value provided for this flag will determine the size of the files that will be generated for testing the download process.",
			DefaultValue: "50",
			Mandatory:    true,
		},
		components.StringFlag{
			Name:         "iterations",
			Description:  "This flag specify how many files will be created for testing the download process.",
			DefaultValue: "30",
			Mandatory:    true,
		},
		components.StringFlag{
			Name:         "repo_name",
			Description:  "The value provided for this flag will determine which repository the tests will be executed on.",
			DefaultValue: "benchmark-dl-tests",
		},
		components.StringFlag{
			Name:         "url",
			DefaultValue: "",
			Description:  "[ONLY ONCE USING CUSTOM SERVER] url of Artifactory server",
		},
		components.StringFlag{
			Name:         "username",
			DefaultValue: "",
			Description:  "[ONLY ONCE USING CUSTOM SERVER] username for Artifactory server",
		},
		components.StringFlag{
			Name:         "password",
			DefaultValue: "",
			Description:  "[ONLY ONCE USING CUSTOM SERVER] password for Artifacory server",
		},
		components.StringFlag{
			Name:         "append",
			DefaultValue: "",
			Description:  "Append the results to existing file",
		},
	}
}

func dlCmd(c *components.Context, downloadConfig *benchmarkUtils.BenchmarkConfig) error {
	log.Info("Starting 'dl' command to measure download time from Artifactory...")
	var benchmarkResults []benchmarkUtils.BenchmarkResult
	servicesManager, serviceManagerError := benchmarkUtils.GetSvcManagerBasedOnAuthLogic(c, downloadConfig)
	if serviceManagerError != nil {
		return serviceManagerError
	}

	IterationsInt, _ := strconv.Atoi(downloadConfig.Iterations)
	FilesSizesInMbInt, _ := strconv.Atoi(downloadConfig.FilesSizesInMb)

	// Creating a repository and upload files that will be used to measure the download time.
	localRepoError := benchmarkUtils.CreateLocalRepository(downloadConfig.RepositoryName, servicesManager)
	if localRepoError != nil {
		return localRepoError
	}
	filesNames, err := benchmarkUtils.GenerateFiles(IterationsInt, FilesSizesInMbInt)
	if err != nil {
		return err
	}
	for _, file := range filesNames {
		_, err := benchmarkUtils.UploadFiles(file, downloadConfig.RepositoryName, servicesManager)
		if err != nil {
			return err
		}
	}
	measureError := benchmarkUtils.MeasureOperationTimes(downloadConfig, filesNames, servicesManager, &benchmarkResults)
	if measureError != nil {
		return measureError
	}
	path := benchmarkUtils.GetFilePath(downloadConfig.Operation, downloadConfig.Append)
	writeResultsError := benchmarkUtils.WriteResults(path, benchmarkResults)
	if writeResultsError != nil {
		return writeResultsError
	}
	log.Info("Finished 'dl' command.")
	cleanupErr := benchmarkUtils.CleanupCliResources(downloadConfig, servicesManager)
	if cleanupErr != nil {
		return cleanupErr
	}
	summriseError := benchmarkUtils.ReadFileAndPrint(path)
	if summriseError != nil {
		return summriseError
	}
	return nil
}
