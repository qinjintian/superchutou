package http

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
)

func Request(method string, url string, body io.Reader, headers map[string]string, args ...interface{}) (*http.Request, *http.Response, error) {
	var (
		err     error
		req     *http.Request
		resp    *http.Response
		timeout int64 = 30
	)

	if len(args) > 0 {
		if v, ok := args[0].(int64); ok {
			timeout = v
		}
	}

	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return req, nil, err
	}

	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn, err := net.Dial(network, addr)
				if conn != nil {
					req.RemoteAddr = conn.RemoteAddr().String()
				}
				return conn, err
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	client.Jar, err = cookiejar.New(nil)
	if err != nil {
		return req, nil, err
	}

	resp, err = client.Do(req)
	if err != nil {
		return req, nil, err
	}

	return req, resp, nil
}
