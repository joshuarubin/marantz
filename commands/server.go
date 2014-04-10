package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	serial "github.com/tarm/goserial"
)

const (
	defaultSerialPort = "/dev/ttyUSB0"
	defaultBaudRate   = 9600
)

var (
	serialCfg serial.Config

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		Long:  `Start the server on the system with the serial connection to the reciever`,
		Run:   server,
	}

	serialCh struct {
		read  chan []byte
		write chan []byte
		err   chan error
	}

	serverCmdP *cobra.Command
)

func init() {
	serialCh.read = make(chan []byte, 128)
	serialCh.write = make(chan []byte, 128)
	serialCh.err = make(chan error, 128)

	viper.SetDefault("serial", map[string]interface{}{
		"port": defaultSerialPort,
		"baud": defaultBaudRate,
	})

	serverCmd.Flags().StringVarP(&serialCfg.Name, "serial", "s", defaultSerialPort, "serial port")
	serverCmd.Flags().IntVarP(&serialCfg.Baud, "baud", "b", defaultBaudRate, "serial port baud rate")

	serverCmdP = serverCmd

	marantzCmd.AddCommand(serverCmd)
}

func initSerialConfig() {
	initServerConfig()

	s := viper.GetStringMap("serial")

	if !serverCmdP.Flags().Lookup("serial").Changed {
		serialCfg.Name = s["port"].(string)
	}

	if !serverCmdP.Flags().Lookup("baud").Changed {
		serialCfg.Baud = s["baud"].(int)
	}
}

func serialWatcher(s io.ReadWriteCloser) {
	// serial reader
	go func() {
		buf := make([]byte, 128)
		n, err := s.Read(buf)
		if err != nil {
			serialCh.err <- err
		} else {
			serialCh.read <- buf[:n]
		}
	}()

	for val := range serialCh.write {
		_, err := s.Write(val)
		if err != nil {
			serialCh.err <- err
		}
	}
}

func server(cmd *cobra.Command, args []string) {
	initSerialConfig()

	s, err := serial.OpenPort(&serialCfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	go serialWatcher(s)

	for {
		select {
		case data := <-serialCh.read:
			fmt.Println("serial read", data)
		case err := <-serialCh.err:
			fmt.Println("serial err", err)
		}
	}
}
