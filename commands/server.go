package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultSerialPort = "/dev/ttyUSB0"
	defaultBaudRate   = 9600
)

var (
	quitCh = make(chan bool)

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		Long:  `Start the server on the system with the serial connection to the reciever`,
		Run:   serverMain,
	}

	serverCmdP *cobra.Command
)

func init() {
	viper.SetDefault("serial", map[string]interface{}{
		"port": defaultSerialPort,
		"baud": defaultBaudRate,
	})

	serverCmd.Flags().StringVarP(&srv.Serial.Config.Name, "serial", "s", defaultSerialPort, "serial port")
	serverCmd.Flags().IntVarP(&srv.Serial.Config.Baud, "baud", "b", defaultBaudRate, "serial port baud rate")

	serverCmdP = serverCmd

	marantzCmd.AddCommand(serverCmd)
}

func initSerialConfig() {
	initServerConfig()

	s := viper.GetStringMap("serial")

	port := s["port"]
	if port == nil {
		srv.Serial.Config.Name = defaultSerialPort
	} else {
		srv.Serial.Config.Name = port.(string)
	}

	baud := s["baud"]
	if baud == nil {
		srv.Serial.Config.Baud = defaultBaudRate
	} else {
		srv.Serial.Config.Baud = baud.(int)
	}
}

func serverMain(cmd *cobra.Command, args []string) {
	initSerialConfig()
	srv.Start()
	srv.Serial.Write <- "AST:F"
	<-quitCh
}
