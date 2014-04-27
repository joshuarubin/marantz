package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
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

func (srv *Server) Start() {
	if err := srv.Serial.Open(); err != nil {
		log.Println("serial port open error", err)
		os.Exit(-1)
	}

	srv.Serial.Write <- "AST:F"

	http.HandleFunc("/cmd", srv.cmdHandler)
	log.Fatal(http.ListenAndServe(srv.Config.String(), nil))
}

func (srv *Server) cmdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Invalid Method: "+r.Method, 400)
		log.Println("Invalid Method: " + r.Method + " (400)")
		return
	}

	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Read Error", 500)
		log.Println("Read Error (500)", err)
		return
	}

	cmd := &msg.Cmd{}
	if err := proto.Unmarshal(data, cmd); err != nil {
		http.Error(w, "Unmarshal Error", 400)
		log.Println("Unmarshal Error (400)", err)
		return
	}

	switch *cmd.Cmd {
	case msg.Cmd_CMD_RAW:
		srv.Serial.Write <- *cmd.StrValue
	case msg.Cmd_CMD_PWR:
		srv.onCmdPwr(cmd)
	case msg.Cmd_CMD_VOL:
		srv.onCmdVol(cmd)
	}

	serialCh, _ := srv.Serial.Sub()
	defer srv.Serial.UnSub(serialCh)

	// TODO(jrubin) wait for the 'correct' response for the given command
	msg := <-serialCh

	// TODO(jrubin) send a protobuf response
	_, err = fmt.Fprintf(w, "%s\n", msg)
	if err != nil {
		http.Error(w, "Write Error", 500)
		log.Println("Write Error (500)", err)
		return
	}
}

/*
func (srv *Server) onConn(conn net.Conn) {
	serialCh, _ := srv.Serial.Sub()
	defer srv.Serial.UnSub(serialCh)

	ch := srv.connReader(conn)

	for {
		select {
		case cmd, ok := <-ch:
			if !ok {
				return
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
*/

func (srv *Server) onCmdPwr(cmd *msg.Cmd) {
	const CMD = "PWR"

	if cmd.Pwr == nil {
		srv.Serial.Write <- CMD + ":?"
		return
	}

	srv.Serial.Write <- fmt.Sprintf("%s:%d", CMD, *cmd.Pwr)
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

func (srv *Server) onCmdVol(cmd *msg.Cmd) {
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

		srv.Serial.Write <- fmt.Sprintf("%s:0%+02.0f", CMD, val)
		return
	}

	if cmd.Vol == nil {
		srv.Serial.Write <- CMD + ":?"
		return
	}

	srv.Serial.Write <- fmt.Sprintf("%s:%d", CMD, *cmd.Vol)
}
