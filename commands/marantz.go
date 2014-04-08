package commands

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	//jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var marantzCmd = &cobra.Command{
	Use:   "marantz",
	Short: "Marantz is an app for controlling Marantz receivers",
	Long:  `Contains both a client and server for communicating over local and remote networks`,
}

var Config struct {
	Verbose bool
	//File    string
}

func homeDir() string {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return u.HomeDir
}

func init() {
	//jww.SetLogThreshold(jww.LevelTrace)
	//jww.SetStdoutThreshold(jww.LevelInfo)

	viper.AddConfigPath("/etc/marantz")
	viper.AddConfigPath(fmt.Sprintf("%s/.marantz", homeDir()))

	viper.SetConfigName("config")

	viper.SetDefault("server", map[string]interface{}{
		"listen": map[string]interface{}{
			"address": "0.0.0.0",
			"port":    6932,
		},
	})

	viper.ReadInConfig()

	marantzCmd.PersistentFlags().BoolVarP(&Config.Verbose, "verbose", "v", false, "verbose output")
	//marantzCmd.PersistentFlags().StringVar(&Config.File, "config", "", "config file (default is [/etc/marantz|~/.marantz]/config.[yaml|json|toml])")
}

func Execute() {
	//viper.Debug()

	marantzCmd.AddCommand(versionCmd)
	//marantzCmd.AddCommand(serverCmd)

	if err := marantzCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
