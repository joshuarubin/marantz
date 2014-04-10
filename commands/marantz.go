package commands

import (
	"log"
	"os/user"

	"github.com/joshuarubin/viper"
	"github.com/spf13/cobra"
)

const (
	defaultServerHost = "localhost"
	defaultServerPort = uint(6932)
)

var (
	verbose bool // TODO(jrubin) use this flag

	serverCfg struct {
		host string
		port uint
	}

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

	viper.SetDefault("server", map[string]interface{}{
		"host": defaultServerHost,
		"port": defaultServerPort,
	})

	marantzCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	marantzCmd.PersistentFlags().StringVarP(&serverCfg.host, "host", "h", defaultServerHost, "server host name (remote for client commands, listen for server)")
	marantzCmd.PersistentFlags().UintVarP(&serverCfg.port, "port", "p", defaultServerPort, "server port")
}

func initServerConfig() {
	s := viper.GetStringMap("server")

	if !marantzCmd.PersistentFlags().Lookup("host").Changed {
		serverCfg.host = s["host"].(string)
	}

	if !marantzCmd.PersistentFlags().Lookup("port").Changed {
		serverCfg.port = s["port"].(uint)
	}
}

func Execute() error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.UnsupportedConfigError); !ok {
			return err
		}
	}

	return marantzCmd.Execute()
}
