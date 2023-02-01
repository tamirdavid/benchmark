package commands

import (
	"benchmark/lib/benchmarkUtils"
	"fmt"
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
			downloadConfig := setDownloadConfig(c)
			return dlCmd(c, downloadConfig)
		},
	}
}

func setDownloadConfig(c *components.Context) *benchmarkUtils.BenchmarkConfig {
	var downloadConfig = new(benchmarkUtils.BenchmarkConfig)
	downloadConfig.FilesSizesInMb = c.GetStringFlagValue("size")
	downloadConfig.Iterations = c.GetStringFlagValue("iterations")
	downloadConfig.RepositoryName = c.GetStringFlagValue("repo_name")
	downloadConfig.Operation = "download"
	return downloadConfig
}

func DownloadCommandFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:         "size",
			Description:  "The value provided for this flag will determine the size of the files that will be generated for testing the download process.",
			DefaultValue: "1",
			Mandatory:    true,
		},
		components.StringFlag{
			Name:         "iterations",
			Description:  "This flag specify how many files will be created for testing the download process.",
			DefaultValue: "5",
			Mandatory:    true,
		},
		components.StringFlag{
			Name:         "repo_name",
			Description:  "The value provided for this flag will determine which repository the tests will be executed on.",
			DefaultValue: "benchmark-dl-tests",
		},
	}
}

func dlCmd(c *components.Context, downloadConfig *benchmarkUtils.BenchmarkConfig) error {
	log.Info("Starting 'dl' command to measure download time from Artifactory...")

	servicesManager := benchmarkUtils.CreateServiceManagerWithThreads(c)

	IterationsInt := benchmarkUtils.CheckIntLikeString(downloadConfig.Iterations)
	FilesSizesInMbInt := benchmarkUtils.CheckIntLikeString(downloadConfig.FilesSizesInMb)

	// Create a repository and upload files that will be used to measure the download time.
	benchmarkUtils.CreateLocalRepository(downloadConfig.RepositoryName, servicesManager)
	filesNames := benchmarkUtils.GenerateFiles(IterationsInt, FilesSizesInMbInt)
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
