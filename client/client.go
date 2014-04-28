package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"code.google.com/p/gogoprotobuf/proto"
	"github.com/joshuarubin/marantz/msg"
)

func SendCmd(host string, cmd *msg.Cmd) {
	data, err := proto.Marshal(cmd)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("PUT", "http://"+host+"/cmd", bytes.NewBuffer(data))
	if err != nil {
		log.Fatalln(err)
	}

	// req.Header.Set("Content-Type", bodyType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	data, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s\n", data)
}
