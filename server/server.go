package server

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"

	"code.google.com/p/goprotobuf/proto"

	"github.com/joshuarubin/marantz/msg"
	"github.com/joshuarubin/marantz/serialport"
)

type HostConfig struct {
	Host string
	Port uint
}

func (c *HostConfig) String() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type Server struct {
	listener net.Listener
	Config   HostConfig
	Serial   serialport.SerialPort
}

func (s *Server) Start() {
	var err error

	if err = s.Serial.Open(); err != nil {
		log.Println("serial port open error", err)
		os.Exit(-1)
	}

	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port))
	if err != nil {
		log.Println("tcp listen error", err)
		os.Exit(-1)
	}

	go s.listen()
}

func (s *Server) connReader(conn net.Conn) <-chan *msg.Cmd {
	ch := make(chan *msg.Cmd)

	go func() {
		var l int16
		err := binary.Read(conn, binary.LittleEndian, &l)
		if err != nil {
			close(ch)
			log.Println("binary read error", err)
			return
		}

		data := make([]byte, l)
		err = binary.Read(conn, binary.LittleEndian, &data)
		if err != nil {
			close(ch)
			log.Println("binary read error", err)
			return
		}

		cmds := &msg.Cmds{}
		if err := proto.Unmarshal(data, cmds); err != nil {
			close(ch)
			log.Println("unmarshal error", err)
			return
		}

		for _, cmd := range cmds.Cmd {
			ch <- cmd
		}
	}()

	return ch
}

func (s *Server) onConn(conn net.Conn) {
	serialCh, _ := s.Serial.Sub()
	defer s.Serial.UnSub(serialCh)

	ch := s.connReader(conn)

	for {
		select {
		case cmd, ok := <-ch:
			if !ok {
				return
			}

			switch *cmd.Cmd {
			case msg.Cmd_CMD_CLOSE:
				if err := conn.Close(); err != nil {
					log.Println("conn close error", err)
				}
				return
			case msg.Cmd_CMD_RAW:
				s.Serial.Write <- *cmd.StrValue
			case msg.Cmd_CMD_PWR:
				s.onCmdPwr(cmd)
			case msg.Cmd_CMD_VOL:
				s.onCmdVol(cmd)
			}
		case msg := <-serialCh:
			// TODO(jrubin) send a protobuf response
			_, err := fmt.Fprintf(conn, "%s\n", msg)
			if err != nil {
				log.Println("conn write error", err)
				return
			}
		}
	}
}

func (s *Server) listen() {
	defer log.Println("Server.listen returned") // TODO(jrubin) ensure there is a way to stop this

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("tcp connection error", err)
		} else {
			go s.onConn(conn)
		}
	}
}

func (s *Server) onCmdPwr(cmd *msg.Cmd) {
	const CMD = "PWR"

	if cmd.Pwr == nil {
		s.Serial.Write <- CMD + ":?"
		return
	}

	s.Serial.Write <- fmt.Sprintf("%s:%d", CMD, *cmd.Pwr)
}

/*
 *    if res == "-ZZ" {
 *        floatRes = -scale
 *    } else {
 *        intRes, err := strconv.Atoi(strings.TrimSpace(res))
 *        floatRes = float32(intRes)
 *        if err != nil {
 *            fmt.Fprintln(os.Stderr, "invalid response from server", res, err)
 *            os.Exit(-1)
 *        }
 *    }
 *
 *    floatRes += scale
 *    fmt.Printf("%.0f\n", floatRes/max*100)
 */

func (s *Server) onCmdVol(cmd *msg.Cmd) {
	const (
		CMD   = "VOL"
		SCALE = 72
		MAX   = 90
	)

	if cmd.IntValue != nil {
		val := float32(*cmd.IntValue)
		switch {
		case val <= 0:
			val = -SCALE
		case val >= 100:
			val = MAX - SCALE
		default:
			val = (val / 100 * MAX) - SCALE
		}

		s.Serial.Write <- fmt.Sprintf("%s:0%+02.0f", CMD, val)
		return
	}

	if cmd.Vol == nil {
		s.Serial.Write <- CMD + ":?"
		return
	}

	s.Serial.Write <- fmt.Sprintf("%s:%d", CMD, *cmd.Vol)
}
