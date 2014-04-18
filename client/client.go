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

type VolMode int

const (
	VolModeGet VolMode = iota
	VolModeUp
	VolModeDown
	VolModeUpFast
	VolModeDownFast
	VolModeSet
)

type VolMsg struct {
	Mode VolMode
	Val  int
}

func (c *Client) Vol(msg VolMsg) {
	const (
		cmd   = "VOL"
		scale = 72
		max   = 90
	)

	var res string

	switch msg.Mode {
	case VolModeGet:
		res = c.cmd(cmd)
	case VolModeSet:
		intVal := float32(msg.Val)
		if intVal <= 0 {
			intVal = -scale
		} else if intVal >= 100 {
			intVal = max - scale
		} else {
			intVal = (intVal / 100 * max) - scale
		}
		res = c.cmd(cmd, fmt.Sprintf("0%+02.0f", intVal))
	default:
		res = c.cmd(cmd, fmt.Sprintf("%d", msg.Mode))
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

	msg := ""
	if len(val) > 0 {
		msg = fmt.Sprintf("%s:%s", cmd, val[0])
	} else {
		msg = fmt.Sprintf("%s:?", cmd)
	}

	_, err = fmt.Fprintf(conn, "%s\n", msg)

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
