package commands

import (
	"benchmark/lib/benchmarkUtils"
	"fmt"
	"strconv"
	"time"

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
	}
}

func dlCmd(c *components.Context, downloadConfig *benchmarkUtils.BenchmarkConfig) error {
	log.Info("Starting 'dl' command to measure download time from Artifactory...")

	servicesManager := benchmarkUtils.GetSvcManagerBasedOnAuthLogic(c, downloadConfig)

	IterationsInt, _ := strconv.Atoi(downloadConfig.Iterations)
	FilesSizesInMbInt, _ := strconv.Atoi(downloadConfig.FilesSizesInMb)

	// Creating a repository and upload files that will be used to measure the download time.
	benchmarkUtils.CreateLocalRepository(downloadConfig.RepositoryName, servicesManager)
	filesNames, err := benchmarkUtils.GenerateFiles(IterationsInt, FilesSizesInMbInt)
	if err != nil {
		return err
	}
	for _, file := range filesNames {
		benchmarkUtils.UploadFiles(file, downloadConfig.RepositoryName, servicesManager)
	}

	uploadResults := benchmarkUtils.MeasureOperationTimes(downloadConfig, filesNames, servicesManager)
	filePath := fmt.Sprintf("benchmark-download-%s.output", time.Now().Format("2006-01-02T15:04:05"))
	benchmarkUtils.WriteResult(filePath, uploadResults, downloadConfig.FilesSizesInMb)
	log.Info("Finished 'dl' command.")
	benchmarkUtils.DeleteRepository(downloadConfig.RepositoryName, servicesManager)

	return nil
}
