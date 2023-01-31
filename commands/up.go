package commands

import (
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
			uploadConfig := setUploadConig(c)
			return upCmd(c, uploadConfig)
		},
	}
}

func setUploadConig(c *components.Context) *BenchmarkConfig {
	var uploadConfig = new(BenchmarkConfig)
	uploadConfig.FilesSizesInMb = c.GetStringFlagValue("size")
	uploadConfig.Iterations = c.GetStringFlagValue("iterations")
	uploadConfig.repositoryName = c.GetStringFlagValue("repo_name")
	uploadConfig.operation = "upload"
	return uploadConfig
}

func UploadCommandFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:         "size",
			Description:  "Size of file in MB to preform download tests for",
			DefaultValue: "50",
			Mandatory:    true,
		},
		components.StringFlag{
			Name:         "iterations",
			Description:  "Number of download iterations",
			DefaultValue: "30",
			Mandatory:    true,
		},
		components.StringFlag{
			Name:         "repo_name",
			Description:  "Repository name to check the upload against",
			DefaultValue: "benchmark-up-tests",
			Mandatory:    true,
		},
	}
}

func upCmd(c *components.Context, st *BenchmarkConfig) error {
	log.Info("Starting upload measurement command")
	IterationsInt, err := strconv.Atoi(st.Iterations)
	FilesSizesInMbInt, err2 := strconv.Atoi(st.FilesSizesInMb)
	if err != nil || err2 != nil {
		fmt.Println("Error converting Iterations and Files sizes from string to int:", err)
	}
	CreateLocalRepository(c, st.repositoryName)
	filesNames := GenerateFiles(IterationsInt, FilesSizesInMbInt)
	uploadResults := MeasureOperationTimes(c, st, filesNames)
	filePath := fmt.Sprintf("benchmark-upload-%s.output", time.Now().Format("2006-01-02T15:04:05"))
	WriteResult(filePath, uploadResults)
	log.Info("Finished upload mesurement command")

	return nil
}
