package main

import (
	arg "github.com/jeffjen/ambd/ambctl/arg"

	ctx "golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"

	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	Endpoint []string
)

type Response struct {
	Data []byte
	Err  error
}

func CreateReq(pflag *arg.Info) (output chan *Response) {
	var (
		wg sync.WaitGroup

		root = ctx.Background()

		v = make(chan *Response, 1)
	)

	go func() {
		defer close(v)

		for _, endpoint := range Endpoint {
			wg.Add(1)
			go func(ep string) {
				defer wg.Done()

				var body *bytes.Reader
				if buf, err := json.Marshal(pflag); err != nil {
					v <- &Response{Err: err}
					return
				} else {
					body = bytes.NewReader(buf)
				}

				wk, abort := ctx.WithTimeout(root, 100*time.Millisecond)
				defer abort()

				resp, err := ctxhttp.Post(wk, nil, ep+"/proxy", "application/json", body)
				if err != nil {
					v <- &Response{Err: err}
					return
				}
				defer resp.Body.Close()

				inn := new(bytes.Buffer)
				io.Copy(inn, resp.Body)
				if ans := inn.String(); ans != "done" {
					v <- &Response{Err: errors.New(ans)}
					return
				}

				v <- &Response{}

			}(endpoint)
		}
		wg.Wait()
	}()

	return v
}

func CancelReq(src string) (output chan *Response) {
	var (
		wg sync.WaitGroup

		root = ctx.Background()

		v = make(chan *Response, 1)
	)

	go func() {
		defer close(v)

		for _, endpoint := range Endpoint {
			wg.Add(1)
			go func(ep string) {
				defer wg.Done()

				wk, abort := ctx.WithTimeout(root, 100*time.Millisecond)
				defer abort()

				var cli = new(http.Client)
				req, err := http.NewRequest("DELETE", ep+"/proxy/"+src, nil)
				if err != nil {
					v <- &Response{Err: err}
					return
				}

				resp, err := ctxhttp.Do(wk, cli, req)
				if err != nil {
					v <- &Response{Err: err}
					return
				}
				defer resp.Body.Close()

				var inn = new(bytes.Buffer)
				io.Copy(inn, resp.Body)
				if ans := inn.String(); ans != "done" {
					v <- &Response{Err: errors.New(ans)}
					return
				}

				v <- &Response{}

			}(endpoint)
		}
		wg.Wait()
	}()

	return v
}

func ConfigReq(proxycfg string) (output chan *Response) {
	var (
		wg sync.WaitGroup

		root = ctx.Background()

		v = make(chan *Response, 1)
	)

	go func() {
		defer close(v)

		for _, endpoint := range Endpoint {
			wg.Add(1)
			go func(ep string) {
				defer wg.Done()

				wk, abort := ctx.WithTimeout(root, 100*time.Millisecond)
				defer abort()

				var cli = new(http.Client)
				req, err := http.NewRequest("PUT", ep+"/proxy/app-config?key="+proxycfg, nil)
				if err != nil {
					v <- &Response{Err: err}
					return
				}

				resp, err := ctxhttp.Do(wk, cli, req)
				if err != nil {
					v <- &Response{Err: err}
					return
				}
				defer resp.Body.Close()

				var inn = new(bytes.Buffer)
				io.Copy(inn, resp.Body)
				if ans := inn.String(); ans != "done" {
					v <- &Response{Err: errors.New(ans)}
					return
				}

				v <- &Response{}

			}(endpoint)
		}
		wg.Wait()
	}()

	return v
}

func InfoReq() (output chan *Response) {
	var (
		wg sync.WaitGroup

		root = ctx.Background()

		v = make(chan *Response, 1)
	)

	go func() {
		defer close(v)

		for _, endpoint := range Endpoint {
			wg.Add(1)
			go func(ep string) {
				defer wg.Done()

				wk, abort := ctx.WithTimeout(root, 100*time.Millisecond)
				defer abort()

				resp, err := ctxhttp.Get(wk, nil, ep+"/info")
				if err != nil {
					v <- &Response{Err: err}
					return
				}
				defer resp.Body.Close()

				inn := new(bytes.Buffer)
				_, err = inn.ReadFrom(resp.Body)
				if err != nil {
					v <- &Response{Err: err}
					return
				}

				v <- &Response{Data: inn.Bytes()}

			}(endpoint)
		}
		wg.Wait()
	}()

	return v
}

func ListProxyReq() (output chan *Response) {
	var (
		wg sync.WaitGroup

		root = ctx.Background()

		v = make(chan *Response, 1)
	)

	go func() {
		defer close(v)

		for _, endpoint := range Endpoint {
			wg.Add(1)
			go func(ep string) {
				defer wg.Done()

				wk, abort := ctx.WithTimeout(root, 100*time.Millisecond)
				defer abort()

				resp, err := ctxhttp.Get(wk, nil, ep+"/proxy/list")
				if err != nil {
					v <- &Response{Err: err}
					return
				}
				defer resp.Body.Close()

				inn := new(bytes.Buffer)
				_, err = inn.ReadFrom(resp.Body)
				if err != nil {
					v <- &Response{Err: err}
					return
				}

				v <- &Response{Data: inn.Bytes()}

			}(endpoint)
		}
		wg.Wait()
	}()

	return v
}
