package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"code.google.com/p/goprotobuf/proto"
	"github.com/joshuarubin/marantz/msg"
)

func SendCmd(host string, cmd *msg.Cmd) {
	data, err := proto.Marshal(cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	req, err := http.NewRequest("PUT", "http://"+host+"/cmd", bytes.NewBuffer(data))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	// req.Header.Set("Content-Type", bodyType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	data, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	fmt.Printf("%s\n", data)
}
