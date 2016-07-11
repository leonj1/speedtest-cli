package speedtest

import (
	"net/http"
	"fmt"
	"runtime"
	"strings"
	"io"
	"net"
	"log"
)

type Client struct {
	http.Client
	opts *Opts
	config *Config
}

func NewClient(opts *Opts) *Client {
	dialer := &net.Dialer{
		Timeout: opts.Timeout,
		KeepAlive: opts.Timeout,
	}

	if len(opts.Source) != 0 {
		dialer.LocalAddr = &net.IPAddr{IP: net.ParseIP(opts.Source)}
		if dialer.LocalAddr == nil {
			log.Fatalf("Invalid source IP: %s\n", opts.Source)
		}
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: dialer.Dial,
		TLSHandshakeTimeout: opts.Timeout,
		ExpectContinueTimeout: opts.Timeout,
	}

	client := &Client{
		Client: http.Client{
			Transport: transport,
			Timeout: opts.Timeout,
		},
		opts: opts,
	}

	return client;
}

func (client *Client) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	if strings.HasPrefix(url, ":") {
		if client.opts.Secure {
			url = "https" + url
		} else {
			url = "http" + url
		}
	}
	req, err := http.NewRequest(method, url, body);
	if err == nil {
		req.Header.Set(
			"User-Agent",
			"Mozilla/5.0 " +
				fmt.Sprintf("(%s; U; %s; en-us)", runtime.GOOS, runtime.GOARCH) +
				fmt.Sprintf("Go/%s", runtime.Version()) +
				fmt.Sprintf("(KHTML, like Gecko) speedtest-cli/%s", Version))
	}
	return req, err;
}

func (client *Client) Get(url string) (*http.Response, error) {
	req, err := client.NewRequest("GET", url, nil);
	if err != nil {
		return nil, err
	}
	return client.Client.Do(req)
}