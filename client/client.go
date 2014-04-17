package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/joshuarubin/marantz/server"
)

type Client struct {
	Server *server.Server
}

func New(s *server.Server) *Client {
	return &Client{
		Server: s,
	}
}

func (c *Client) Pwr(val ...bool) {
	const cmd = "PWR"

	var res string

	if len(val) > 0 {
		strVal := "1" // off
		if val[0] {
			strVal = "2" // on
		}

		res = c.cmd(cmd, strVal)
	} else {
		res = c.cmd(cmd)
	}

	switch res {
	case "1":
		fmt.Println("off")
	case "2":
		fmt.Println("on")
	default:
		fmt.Fprintln(os.Stderr, "Unknown response: %s", res)
	}
}

func (c *Client) Vol(val ...int) {
	// TODO(jrubin) work with .0 and .5 values
	const cmd = "VOL"
	const scale = 72
	const max = 90

	var res string

	if len(val) > 0 {
		intVal := float32(val[0])
		if intVal <= 0 {
			intVal = -scale
		} else if intVal >= 100 {
			intVal = max - scale
		} else {
			intVal = (intVal / 100 * max) - scale
		}
		msg := fmt.Sprintf("0%+02.0f", intVal)
		res = c.cmd(cmd, msg)
	} else {
		res = c.cmd(cmd)
	}

	var floatRes float32
	if res == "-ZZ" {
		floatRes = -scale
	} else {
		intRes, err := strconv.Atoi(strings.TrimSpace(res))
		floatRes = float32(intRes)
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid response from server", res, err)
			os.Exit(-1)
		}
	}

	floatRes += scale
	fmt.Printf("%.0f\n", floatRes/max*100)
}

func (c *Client) cmd(cmd string, val ...string) string {
	conn, err := net.Dial("tcp", c.Server.Config.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	if len(val) > 0 {
		_, err = fmt.Fprintf(conn, "%s:%s\n", cmd, val[0])
	} else {
		_, err = fmt.Fprintf(conn, "%s:?\n", cmd)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	defer func() {
		_, err := fmt.Fprintf(conn, "\n") // indicates end of connection

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
	}()

	rd := bufio.NewScanner(conn)

	for rd.Scan() {
		parts := strings.Split(rd.Text(), ":")
		if parts[0] == cmd {
			return parts[1]
		}
	}

	if err := rd.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	fmt.Fprintln(os.Stderr, "Unknown error")
	os.Exit(-1)
	return ""
}
