package commands

import (
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
		Arguments:   DownloadCommandArguments(),
		Flags:       DownloadCommandFlags(),
		EnvVars:     DownloadCommandEnvVar(),
		Action: func(c *components.Context) error {
			downloadConfig := setDownloadConfig(c)
			return dlCmd(c, downloadConfig)
		},
	}
}

func setDownloadConfig(c *components.Context) *BenchmarkConfig {
	var downloadConfig = new(BenchmarkConfig)
	downloadConfig.FilesSizesInMb = c.GetStringFlagValue("size")
	downloadConfig.Iterations = c.GetStringFlagValue("iterations")
	downloadConfig.repositoryName = c.GetStringFlagValue("repo_name")
	downloadConfig.operation = "download"
	return downloadConfig
}

func DownloadCommandArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "addressee",
			Description: "The name of the person you would like to greet.",
		},
	}
}

func DownloadCommandFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:         "size",
			Description:  "Size of file in MB to preform download tests for",
			DefaultValue: "1",
			Mandatory:    true,
		},
		components.StringFlag{
			Name:         "iterations",
			Description:  "Number of download iterations",
			DefaultValue: "5",
			Mandatory:    true,
		},
		components.StringFlag{
			Name:         "repo_name",
			Description:  "Repository name to check the downloads against",
			DefaultValue: "benchmark-dl-tests",
		},
	}
}

func DownloadCommandEnvVar() []components.EnvVar {
	return []components.EnvVar{
		{
			Name:        "HELLO_FROG_GREET_PREFIX",
			Default:     "A new greet from your plugin template: ",
			Description: "Adds a prefix to every greet.",
		},
	}
}

func dlCmd(c *components.Context, st *BenchmarkConfig) error {
	log.Info("Starting download measurement command")
	IterationsInt, err := strconv.Atoi(st.Iterations)
	FilesSizesInMbInt, err2 := strconv.Atoi(st.FilesSizesInMb)
	if err != nil || err2 != nil {
		fmt.Println("Error converting Iterations and Files sizes from string to int:", err)
	}
	CreateLocalRepository(c, st.repositoryName)
	filesNames := GenerateFiles(IterationsInt, FilesSizesInMbInt)
	confDetails, _ := getRtDetails(c)
	// upload files before downloading it
	for _, file := range filesNames {
		UploadFiles(confDetails, file, st.repositoryName)
	}

	uploadResults := MeasureOperationTimes(c, st, filesNames)
	filePath := fmt.Sprintf("benchmark-download-%s.output", time.Now().Format("2006-01-02T15:04:05"))
	WriteResult(filePath, uploadResults)
	log.Info("Finishing download measurement command")

	return nil
}
