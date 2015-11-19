package command

import (
	pxy "github.com/jeffjen/docker-ambassador/proxy"

	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

const (
	Endpoint = "http://localhost:29091/info"

	EndpointProxy = "http://localhost:29091/proxy"
)

func CreateReq(pflag pxy.Info) error {
	var buf = new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(pflag); err != nil {
		return err
	}
	resp, err := http.Post(EndpointProxy, "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var ret = new(bytes.Buffer)
	io.Copy(ret, resp.Body)

	if ans := ret.String(); ans != "done" {
		return errors.New(ans)
	} else {
		return nil
	}
}

func CancelReq(src string) error {
	var cli = new(http.Client)
	req, err := http.NewRequest("DELETE", EndpointProxy+"/"+src, nil)
	if err != nil {
		return err
	}
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var ret = new(bytes.Buffer)
	io.Copy(ret, resp.Body)

	if ans := ret.String(); ans != "done" {
		return errors.New(ans)
	} else {
		return nil
	}
}

func InfoReq() error {
	resp, err := http.Get(Endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var out, inn bytes.Buffer

	_, err = inn.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	json.Indent(&out, inn.Bytes(), "", "    ")
	out.WriteTo(os.Stdout)
	return nil
}
