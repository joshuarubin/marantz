package client

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"

	"code.google.com/p/goprotobuf/proto"

	"github.com/joshuarubin/marantz/msg"
)

func SendCmd(host string, cmd *msg.Cmd) {
	conn, err := net.Dial("tcp", host)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	defer conn.Close()
	buf := new(bytes.Buffer)

	cmds := &msg.Cmds{
		Cmd: []*msg.Cmd{cmd},
	}

	data, err := proto.Marshal(cmds)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	err = binary.Write(buf, binary.LittleEndian, int16(len(data)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	err = binary.Write(buf, binary.LittleEndian, data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	_, err = buf.WriteTo(conn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	// presuming all Cmd names begin with "CMD_"
	cmdName := msg.Cmd_Cmd_name[int32(cmd.GetCmd())][4:]

	rd := bufio.NewScanner(conn)

	// TODO(jrubin) parse pb response
	for rd.Scan() {
		parts := strings.Split(rd.Text(), ":")
		if parts[0] == cmdName {
			fmt.Println(parts[1])
			return
		}
	}

	if err := rd.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	fmt.Fprintln(os.Stderr, "Unknown error")
	os.Exit(-1)
}
