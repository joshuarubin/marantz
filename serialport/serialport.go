package serialport

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	serial "github.com/tarm/goserial"
)

type SerialPort struct {
	Read   <-chan string
	Write  chan<- string
	Config serial.Config
	port   io.ReadWriteCloser
	opened bool
}

func (s *SerialPort) open() (err error) {
	if s.opened {
		return nil
	}

	s.port, err = serial.OpenPort(&s.Config)
	if err != nil {
		return err
	}

	s.opened = true
	return nil
}

func (s *SerialPort) reader() <-chan string {
	ch := make(chan string)

	go func() {
		rd := bufio.NewReader(s.port)

		for {
			str, err := rd.ReadString('\r')
			str = strings.Trim(str, "@\r")
			ch <- str

			if err != nil {
				log.Println("SerialPort::reader error", err)
			}
		}
	}()

	return ch
}

func (s *SerialPort) writer() chan<- string {
	ch := make(chan string)

	go func() {
		for val := range ch {
			data := []byte(fmt.Sprintf("@%s\r", val))
			_, err := s.port.Write(data)
			if err != nil {
				log.Println("SerialPort::writer error", err)
			}
		}
	}()

	return ch
}

func (s *SerialPort) Start() error {
	if err := s.open(); err != nil {
		return err
	}

	if s.Read == nil {
		s.Read = s.reader()
	}

	if s.Write == nil {
		s.Write = s.writer()
	}

	return nil
}
