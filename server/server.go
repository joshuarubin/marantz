package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joshuarubin/marantz/serialport"
)

type Server struct {
	listener net.Listener
	Config   struct {
		Host string
		Port uint
	}
	Serial serialport.SerialPort
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

func (s *Server) connReader(conn net.Conn) <-chan string {
	ch := make(chan string)

	go func() {
		scanner := bufio.NewScanner(conn)

		for scanner.Scan() {
			ch <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			close(ch)
			// *net.OpError is probably just a closed connection
			if _, ok := err.(*net.OpError); !ok {
				log.Println("conn read error", err)
			}
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
		case msg, ok := <-ch:
			if !ok {
				return
			}

			if len(msg) == 0 {
				if err := conn.Close(); err != nil {
					log.Println("conn close error", err)
				}
				return
			}

			s.Serial.Write <- msg
		case msg := <-serialCh:
			_, err := conn.Write([]byte(fmt.Sprintf("%s\n", msg)))
			if err != nil {
				log.Println("conn write error", err)
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
