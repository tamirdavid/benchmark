package commands

import (
	"benchmark/lib/benchmarkUtils"
	"fmt"
	"time"

	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func UploadCommand() components.Command {
	return components.Command{
		Name:        "up",
		Description: "Upload artifacts tests",
		Flags:       UploadCommandFlags(),
		Action: func(c *components.Context) error {
			uploadConfig := setUploadConig(c)
			return upCmd(c, uploadConfig)
		},
	}
}

func setUploadConig(c *components.Context) *benchmarkUtils.BenchmarkConfig {
	var uploadConfig = new(benchmarkUtils.BenchmarkConfig)
	uploadConfig.FilesSizesInMb = c.GetStringFlagValue("size")
	uploadConfig.Iterations = c.GetStringFlagValue("iterations")
	uploadConfig.RepositoryName = c.GetStringFlagValue("repo_name")
	uploadConfig.Operation = "upload"
	return uploadConfig
}

func UploadCommandFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:         "size",
			Description:  "Determine the size of the files (in MB) that will be generated for testing the upload process.",
			DefaultValue: "50",
			Mandatory:    true,
		},
		components.StringFlag{
			Name:         "iterations",
			Description:  "This flag specify how many files will be created for testing the upload process.",
			DefaultValue: "30",
			Mandatory:    true,
		},
		components.StringFlag{
			Name:         "repo_name",
			Description:  "The value provided for this flag will determine which repository the tests will be executed on.",
			DefaultValue: "benchmark-up-tests",
			Mandatory:    true,
		},
	}
}

func upCmd(c *components.Context, uploadConfig *benchmarkUtils.BenchmarkConfig) error {
	log.Info("Starting 'up' command to measure upload time to Artifactory...")

	servicesManager := benchmarkUtils.CreateServiceManagerWithThreads(c)
	IterationsInt := benchmarkUtils.CheckIntLikeString(uploadConfig.Iterations)
	FilesSizesInMbInt := benchmarkUtils.CheckIntLikeString(uploadConfig.FilesSizesInMb)

	benchmarkUtils.CreateLocalRepository(uploadConfig.RepositoryName, servicesManager)
	filesNames := benchmarkUtils.GenerateFiles(IterationsInt, FilesSizesInMbInt)
	uploadResults := benchmarkUtils.MeasureOperationTimes(uploadConfig, filesNames, servicesManager)
	filePath := fmt.Sprintf("benchmark-upload-%s.output", time.Now().Format("2006-01-02T15:04:05"))
	benchmarkUtils.WriteResult(filePath, uploadResults, uploadConfig.FilesSizesInMb)
	log.Info("Finished 'up' command")

	return nil
}
