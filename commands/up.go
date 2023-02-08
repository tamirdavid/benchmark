package commands

import (
	"benchmark/lib/benchmarkUtils"
	"fmt"
	"strconv"
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
			uploadConfig, err := setUploadConig(c)
			if err != nil {
				return err
			}
			return upCmd(c, uploadConfig)
		},
	}
}

func setUploadConig(c *components.Context) (*benchmarkUtils.BenchmarkConfig, error) {
	var uploadConfig = new(benchmarkUtils.BenchmarkConfig)
	uploadConfig.FilesSizesInMb = c.GetStringFlagValue("size")
	uploadConfig.Iterations = c.GetStringFlagValue("iterations")
	uploadConfig.RepositoryName = c.GetStringFlagValue("repo_name")
	uploadConfig.Operation = "upload"
	uploadConfig.Url = c.GetStringFlagValue("url")
	uploadConfig.UserName = c.GetStringFlagValue("username")
	uploadConfig.Password = c.GetStringFlagValue("password")
	err := benchmarkUtils.ValidateInput(uploadConfig)
	if err != nil {
		return nil, err
	}
	return uploadConfig, nil
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

func upCmd(c *components.Context, uploadConfig *benchmarkUtils.BenchmarkConfig) error {
	log.Info("Starting 'up' command to measure upload time to Artifactory...")

	servicesManager, serviceManagerError := benchmarkUtils.GetSvcManagerBasedOnAuthLogic(c, uploadConfig)
	if serviceManagerError != nil {
		return serviceManagerError
	}

	IterationsInt, _ := strconv.Atoi(uploadConfig.Iterations)
	FilesSizesInMbInt, _ := strconv.Atoi(uploadConfig.FilesSizesInMb)

	benchmarkUtils.CreateLocalRepository(uploadConfig.RepositoryName, servicesManager)
	filesNames, err := benchmarkUtils.GenerateFiles(IterationsInt, FilesSizesInMbInt)
	if err != nil {
		return err
	}
	uploadResults, measureError := benchmarkUtils.MeasureOperationTimes(uploadConfig, filesNames, servicesManager)
	if measureError != nil {
		return measureError
	}
	filePath := fmt.Sprintf("benchmark-upload-%s.output", time.Now().Format("2006-01-02T15:04:05"))
	writeResultsError := benchmarkUtils.WriteResult(filePath, uploadResults, uploadConfig.FilesSizesInMb)
	if writeResultsError != nil {
		return writeResultsError
	}
	log.Info("Finished 'up' command")
	deleteError := benchmarkUtils.DeleteRepository(uploadConfig.RepositoryName, servicesManager)
	if deleteError != nil {
		return deleteError
	}

	return nil
}
