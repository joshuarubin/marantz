package commands

import (
	"bufio"
	"fmt"
	"log"
	"net"
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
	serverListener net.Listener
	serialPort     serialport.SerialPort

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

func handleConnection(conn net.Conn) {
	msgCh := make(chan string)
	end := make(chan bool)
	serPortCh := serialPort.Sub()

	go func() {
		scanner := bufio.NewScanner(conn)

		for scanner.Scan() {
			msg := scanner.Text()
			if len(msg) > 0 {
				msgCh <- scanner.Text()
			} else {
				end <- true
				return
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println("conn read error", err)
			end <- true
		}
	}()

	for {
		select {
		case msg := <-msgCh:
			serialPort.Write <- msg
		case msg := <-serPortCh:
			_, err := conn.Write([]byte(fmt.Sprintf("%s\n", msg)))
			if err != nil {
				log.Println("conn write error", err)
			}
		case <-end:
			if err := conn.Close(); err != nil {
				log.Println("conn close error", err)
			}
			return
		}
	}
}

func listen() {
	for {
		conn, err := serverListener.Accept()
		if err != nil {
			log.Println("tcp connection error", err)
		} else {
			go handleConnection(conn)
		}
	}
}

func serverMain(cmd *cobra.Command, args []string) {
	initSerialConfig()

	var err error

	if err = serialPort.Open(); err != nil {
		log.Println("serial port open error", err)
		os.Exit(-1)
	}

	serverListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", serverCfg.host, serverCfg.port))
	if err != nil {
		log.Println("tcp listen error", err)
		os.Exit(-1)
	}

	go listen()

	serPortCh := serialPort.Sub()

	serialPort.Write <- "AST:F"

	for {
		select {
		case data := <-serPortCh:
			log.Println("serial read", len(data), data)
		}
	}
}
