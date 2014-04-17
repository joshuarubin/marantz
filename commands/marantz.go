package commands

import (
	"os/user"

	"github.com/joshuarubin/marantz/server"
	"github.com/joshuarubin/viper"
	"github.com/spf13/cobra"
)

const (
	defaultServerHost = "localhost"
	defaultServerPort = uint(6932)
)

var (
	srv server.Server

	marantzCmd = &cobra.Command{
		Use:   "marantz",
		Short: "Marantz is an app for controlling Marantz receivers",
		Long:  `Contains both a client and server for communicating over local and remote networks`,
	}
)

func homeDir() string {
	defer recover() // user.Current is not supported on linux/arm

	u, err := user.Current()
	if err != nil {
		return ""
	}

	return u.HomeDir
}

func init() {
	viper.AddConfigPath("/etc/marantz")
	if hd := homeDir(); len(hd) > 0 {
		viper.AddConfigPath(hd + "/.marantz")
	}

	viper.SetConfigName("config")

	viper.SetDefault("server", map[string]interface{}{
		"host": defaultServerHost,
		"port": defaultServerPort,
	})
}

func initServerConfig() {
	s := viper.GetStringMap("server")

	host := s["host"]
	if host == nil {
		srv.Config.Host = defaultServerHost
	} else {
		srv.Config.Host = host.(string)
	}

	port := s["port"]
	if port == nil {
		srv.Config.Port = defaultServerPort
	} else {
		srv.Config.Port = port.(uint)
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
