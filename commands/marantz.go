package commands

import (
	"log"
	"os/user"

	"github.com/joshuarubin/viper"
	"github.com/spf13/cobra"
)

var (
	verbose bool // TODO(jrubin) use this flag

	marantzCmd = &cobra.Command{
		Use:   "marantz",
		Short: "Marantz is an app for controlling Marantz receivers",
		Long:  `Contains both a client and server for communicating over local and remote networks`,
	}
)

func homeDir() string {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return u.HomeDir
}

func init() {
	viper.AddConfigPath("/etc/marantz")
	viper.AddConfigPath(homeDir() + "/.marantz")

	viper.SetConfigName("config")

	marantzCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func Execute() error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.UnsupportedConfigError); !ok {
			return err
		}
	}

	return marantzCmd.Execute()
}
