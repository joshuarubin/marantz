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
	defaultListenAddress = "0.0.0.0"
	defaultListenPort    = uint(6932)
	defaultSerialPort    = "/dev/ttyUSB0"
	defaultBaudRate      = 9600
)

var (
	serverCfg struct {
		listen struct {
			address string
			port    uint
		}

		serial struct {
			port     string
			baudRate int
		}
	}

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
	viper.SetDefault("server", map[string]interface{}{
		"listen": map[string]interface{}{
			"address": defaultListenAddress,
			"port":    defaultListenPort,
		},
		"serial": map[string]interface{}{
			"port": defaultSerialPort,
			"baud": defaultBaudRate,
		},
	})

	serverCmd.Flags().StringVarP(&serverCfg.listen.address, "address", "a", defaultListenAddress, "listen address")
	serverCmd.Flags().UintVarP(&serverCfg.listen.port, "port", "p", defaultListenPort, "listen port")
	serverCmd.Flags().StringVarP(&serverCfg.serial.port, "serial", "s", defaultSerialPort, "serial port")
	serverCmd.Flags().IntVarP(&serverCfg.serial.baudRate, "baud", "b", defaultBaudRate, "serial port baud rate")

	serverCmdP = serverCmd

	marantzCmd.AddCommand(serverCmd)
}

func initConfig() {
	server := viper.GetStringMap("server")

	l := server["listen"].(map[string]interface{})
	s := server["serial"].(map[string]interface{})

	if !serverCmdP.Flags().Lookup("address").Changed {
		serverCfg.listen.address = l["address"].(string)
	}

	if !serverCmdP.Flags().Lookup("port").Changed {
		serverCfg.listen.port = l["port"].(uint)
	}

	if !serverCmdP.Flags().Lookup("serial").Changed {
		serverCfg.serial.port = s["port"].(string)
	}

	if !serverCmdP.Flags().Lookup("baud").Changed {
		serverCfg.serial.baudRate = s["baud"].(int)
	}
}

func openSerial() (io.ReadWriteCloser, error) {
	serialCh.read = make(chan []byte, 128)
	serialCh.write = make(chan []byte, 128)
	serialCh.err = make(chan error, 128)

	c := &serial.Config{Name: serverCfg.serial.port, Baud: serverCfg.serial.baudRate}
	return serial.OpenPort(c)
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
	initConfig()

	s, err := openSerial()
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
