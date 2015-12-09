package main

import (
	arg "github.com/jeffjen/ambd/ambctl/arg"

	cli "github.com/codegangsta/cli"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	Endpoint string
)

func endpoint(ctx *cli.Context) error {
	var host string = ctx.String("host")
	Endpoint = fmt.Sprintf("http://%s", host)
	return nil
}

func CreateReq(pflag arg.Info) error {
	var buf = new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(pflag); err != nil {
		return err
	}
	resp, err := http.Post(Endpoint+"/proxy", "application/json", buf)
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
	req, err := http.NewRequest("DELETE", Endpoint+"/proxy/"+src, nil)
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

func ConfigReq(proxycfg string) error {
	var cli = new(http.Client)
	req, err := http.NewRequest("PUT", Endpoint+"/proxy/app-config?key="+proxycfg, nil)
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
	resp, err := http.Get(Endpoint + "/info")
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

func ListProxyReq() error {
	resp, err := http.Get(Endpoint + "/proxy/list")
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
