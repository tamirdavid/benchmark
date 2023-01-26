package commands

import (
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-cli-core/v2/utils/coreutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"os"
	"strings"
)

func GetHelloCommand() components.Command {
	return components.Command{
		Name:        "hello",
		Description: "Says Hello.",
		Aliases:     []string{"hi"},
		Arguments:   getHelloArguments(),
		Flags:       getHelloFlags(),
		EnvVars:     getHelloEnvVar(),
		Action: func(c *components.Context) error {
			return helloCmd(c)
		},
	}
}

func getHelloArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "addressee",
			Description: "The name of the person you would like to greet.",
		},
	}
}

func getHelloFlags() []components.Flag {
	return []components.Flag{
		components.BoolFlag{
			Name:         "shout",
			Description:  "Makes output uppercase.",
			DefaultValue: false,
		},
	}
}

func getHelloEnvVar() []components.EnvVar {
	return []components.EnvVar{
		{
			Name:        "HELLO_FROG_GREET_PREFIX",
			Default:     "A new greet from your plugin template: ",
			Description: "Adds a prefix to every greet.",
		},
	}
}

type helloConfiguration struct {
	addressee string
	shout     bool
	prefix    string
}

func helloCmd(c *components.Context) error {
	if len(c.Arguments) == 0 {
		message := "Hello :) Now try adding an argument to the 'hi' command"
		// You log messages using the following log levels.
		log.Output(message)
		log.Debug(message)
		log.Info(message)
		log.Warn(message)
		log.Error(message)
		return nil
	}
	if len(c.Arguments) > 1 {
		return errors.New("too many arguments received. Now run the command again, with one argument only")
	}

	var conf = new(helloConfiguration)
	conf.addressee = c.Arguments[0]
	conf.shout = c.GetBoolFlagValue("shout")
	conf.prefix = os.Getenv("HELLO_FROG_GREET_PREFIX")
	if conf.prefix == "" {
		conf.prefix = "New greeting: "
	}

	log.Info(doGreet(conf))

	if !conf.shout {
		message := "Now try adding the --shout option to the command"
		log.Info(message)
		return nil
	}

	if os.Getenv(coreutils.LogLevel) == "" {
		message := fmt.Sprintf("Now try setting the %s environment variable to %s and run the command again", coreutils.LogLevel, "DEBUG")
		log.Info(message)
	}
	return nil
}

func doGreet(c *helloConfiguration) string {
	greet := c.prefix + "Hello " + c.addressee + "\n"

	if c.shout {
		greet = strings.ToUpper(greet)
	}

	return strings.TrimSpace(greet)
}
