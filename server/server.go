package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/joshuarubin/marantz/msg"
	"github.com/joshuarubin/marantz/serialport"
)

const (
	GET       = ":?"
	CMD_VOL   = "VOL"
	CMD_PWR   = "PWR"
	CMD_SRC   = "SRC"
	VOL_SCALE = 72
	VOL_MAX   = 90
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
		log.Fatalln("serial port open error", err)
	}

	srv.Serial.Write <- "AST:F"

	http.HandleFunc("/cmd", srv.cmdHandler)
	log.Fatalln(http.ListenAndServe(srv.Config.String(), nil))
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
	case msg.Cmd_CMD_SRC:
		srv.onCmdSrc(cmd)
	}

	srv.sendResponse(w, cmd)
}

func (srv *Server) sendResponse(w http.ResponseWriter, cmd *msg.Cmd) {
	serialCh, _ := srv.Serial.Sub()
	defer srv.Serial.UnSub(serialCh)

	for {
		select {
		case iface := <-serialCh:
			if resp, ok := iface.(string); ok {
				parts := strings.Split(resp, ":")
				if len(parts) != 2 {
					continue
				}

				resp := strings.TrimSpace(parts[1])

				switch *cmd.Cmd {
				case msg.Cmd_CMD_PWR:
					if parts[0] == "PWR" {
						srv.onRespPwr(w, cmd, resp)
						return
					}
				case msg.Cmd_CMD_VOL:
					if parts[0] == "VOL" {
						srv.onRespVol(w, cmd, resp)
						return
					}
				case msg.Cmd_CMD_SRC:
					if parts[0] == "SRC" {
						srv.onRespSrc(w, cmd, resp)
						return
					}
				}
			}
		case <-time.After(time.Second):
			http.Error(w, "Timeout waiting for response from receiver", 500)
			log.Println("Timeout waiting for response from receiver (500)")
			return
		}
	}

}

func (srv *Server) onRespPwr(w http.ResponseWriter, cmd *msg.Cmd, sval string) {
	ival, err := strconv.Atoi(sval)
	if err != nil {
		http.Error(w, "Response Error", 500)
		log.Println("Response Error (500)")
		return
	}

	val := msg.Cmd_PwrValue(ival)

	// TODO(jrubin) send a protobuf response

	switch val {
	case msg.Cmd_PWR_ON:
		_, err = fmt.Fprintf(w, "pwr on")
	case msg.Cmd_PWR_OFF:
		_, err = fmt.Fprintf(w, "pwr off")
	default:
		http.Error(w, "pwr unknown: "+sval, 500)
		log.Println("pwr unknown", sval, "(500)")
		return
	}

	if err != nil {
		http.Error(w, "Write Error", 500)
		log.Println("Write Error (500)", err)
	}
}

func (srv *Server) onRespVol(w http.ResponseWriter, cmd *msg.Cmd, val string) {
	var res float32

	if val == "-ZZ" {
		res = -VOL_SCALE
	} else {
		intRes, err := strconv.Atoi(strings.TrimSpace(val))
		if err != nil {
			http.Error(w, "vol unknown: "+val, 500)
			log.Println("vol unknown", val, "(500)")
			return
		}
		res = float32(intRes)
	}

	res += VOL_SCALE
	if _, err := fmt.Fprintf(w, "vol %.0f", res/VOL_MAX*100); err != nil {
		http.Error(w, "Write Error", 500)
		log.Println("Write Error (500)", err)
	}
}

func srcToString(src uint8) string {
	switch src {
	case '1':
		return "tv"
	case '2':
		return "dvd"
	case '3':
		return "vcr"
	case '5':
		return "dss"
	case '9':
		return "aux1"
	case 'A':
		return "aux2"
	case 'C':
		return "cd"
	case 'E':
		return "tape"
	case 'F':
		return "tuner"
	case 'G':
		return "fm"
	case 'H':
		return "am"
	case 'J':
		return "xm"
	}

	return "unknown"
}

func (srv *Server) onRespSrc(w http.ResponseWriter, cmd *msg.Cmd, val string) {
	video := srcToString(val[0])
	audio := srcToString(val[1])

	if _, err := fmt.Fprintf(w, "video: %s, audio: %s", video, audio); err != nil {
		http.Error(w, "Write Error", 500)
		log.Println("Write Error (500)", err)
	}
}

func (srv *Server) onCmdPwr(cmd *msg.Cmd) {
	if cmd.Pwr == nil {
		srv.Serial.Write <- CMD_PWR + GET
		return
	}

	srv.Serial.Write <- fmt.Sprintf("%s:%d", CMD_PWR, *cmd.Pwr)
}

func (srv *Server) onCmdVol(cmd *msg.Cmd) {
	if cmd.IntValue != nil {
		val := float32(*cmd.IntValue)
		switch {
		case val <= 0:
			val = -VOL_SCALE
		case val >= 100:
			val = VOL_MAX - VOL_SCALE
		default:
			val = (val / 100 * VOL_MAX) - VOL_SCALE
		}

		srv.Serial.Write <- fmt.Sprintf("%s:0%+02.0f", CMD_VOL, val)
		return
	}

	if cmd.Vol == nil {
		srv.Serial.Write <- CMD_VOL + GET
		return
	}

	srv.Serial.Write <- fmt.Sprintf("%s:%d", CMD_VOL, *cmd.Vol)
}

func (srv *Server) onCmdSrc(cmd *msg.Cmd) {
	if cmd.Src == nil {
		srv.Serial.Write <- CMD_SRC + GET
		return
	}

	val := *cmd.Src

	if val < 10 {
		log.Printf("Selecting source: %s:%d\n", CMD_SRC, val)
		srv.Serial.Write <- fmt.Sprintf("%s:%d", CMD_SRC, val)
		return
	}

	log.Printf("Selecting source: %s:%c\n", CMD_SRC, val-10+'A')
	srv.Serial.Write <- fmt.Sprintf("%s:%c", CMD_SRC, val-10+'A')
}
