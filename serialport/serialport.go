package serialport

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/joshuarubin/marantz/pubsub"

	serial "github.com/tarm/goserial"
)

type SerialPort struct {
	Write  chan<- string
	Config serial.Config
	port   io.ReadWriteCloser
	opened bool
	ps     *pubsub.PubSub
}

func (s *SerialPort) reader() {
	rd := bufio.NewReader(s.port)

	for {
		str, err := rd.ReadString('\r')

		s.ps.Pub(strings.Trim(str, "@\r"))

		if err != nil {
			log.Println("SerialPort::reader error", err)
		}
	}
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

func (s *SerialPort) Open() (err error) {
	if s.opened {
		return nil
	}

	s.port, err = serial.OpenPort(&s.Config)
	if err != nil {
		return err
	}

	if s.ps == nil {
		s.ps = pubsub.New()
	}

	go s.reader()

	if s.Write == nil {
		s.Write = s.writer()
	}

	s.opened = true
	return nil
}

func (s *SerialPort) Sub() <-chan string {
	return s.ps.Sub()
}

func (s *SerialPort) UnSub(ch chan string) {
	s.ps.UnSub(ch)
}
