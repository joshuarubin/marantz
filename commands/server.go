package commands

import (
	"log"
	"os"

	"github.com/joshuarubin/marantz/serialport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultSerialPort = "/dev/ttyUSB0"
	defaultBaudRate   = 9600
)

var (
	serialPort serialport.SerialPort

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		Long:  `Start the server on the system with the serial connection to the reciever`,
		Run:   server,
	}

	serverCmdP *cobra.Command
)

func init() {
	viper.SetDefault("serial", map[string]interface{}{
		"port": defaultSerialPort,
		"baud": defaultBaudRate,
	})

	serverCmd.Flags().StringVarP(&serialPort.Config.Name, "serial", "s", defaultSerialPort, "serial port")
	serverCmd.Flags().IntVarP(&serialPort.Config.Baud, "baud", "b", defaultBaudRate, "serial port baud rate")

	serverCmdP = serverCmd

	marantzCmd.AddCommand(serverCmd)
}

func initSerialConfig() {
	initServerConfig()

	s := viper.GetStringMap("serial")

	if !serverCmdP.Flags().Lookup("serial").Changed {
		serialPort.Config.Name = s["port"].(string)
	}

	if !serverCmdP.Flags().Lookup("baud").Changed {
		serialPort.Config.Baud = s["baud"].(int)
	}
}

func server(cmd *cobra.Command, args []string) {
	initSerialConfig()

	if err := serialPort.Start(); err != nil {
		log.Println("serial port start error", err)
		os.Exit(-1)
	}

	serialPort.Write <- "AST:F"

	for {
		select {
		case data := <-serialPort.Read:
			log.Println("serial read", len(data), data)
		}
	}
}
