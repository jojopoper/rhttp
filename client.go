package rhttp

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// RClient http client struct
type RClient struct {
	timeout int
}

// GetOrigClient 获取原生的Http Client
func (ths *RClient) GetOrigClient() *http.Client {
	return http.DefaultClient
}

// GetClient 获取设定超时的Http client
// timeout 最小值为30 单位秒
func (ths *RClient) GetClient(timeout int) *http.Client {
	if timeout < 30 {
		timeout = 30
	}
	ths.timeout = timeout
	dialer := &net.Dialer{
		Timeout: time.Second * time.Duration(ths.timeout),
		// Deadline: time.Now().Add(time.Duration(ths.timeout) * time.Second),
		KeepAlive: 30 * time.Second,
	}
	trans := &http.Transport{
		Dial: dialer.Dial,
		ResponseHeaderTimeout: time.Duration(ths.timeout) * time.Second,

		TLSHandshakeTimeout: time.Duration(ths.timeout) * time.Second,
	}

	ret := &http.Client{
		Transport: trans,
		Timeout:   time.Duration(ths.timeout) * time.Second,
	}
	return ret
}

// GetProxyClient 获取代理Http client
// timeout 最小值为30 单位秒
func (ths *RClient) GetProxyClient(timeout int, proxyIP, proxyPort string, auth ...*proxy.Auth) (*http.Client, error) {
	if timeout < 30 {
		timeout = 30
	}
	ths.timeout = timeout
	proxyurl := proxyIP + ":" + proxyPort
	var author *proxy.Auth
	if auth != nil {
		author = auth[0]
	}
	// dialer, err := proxy.SOCKS5("tcp", proxyurl, author,
	// 	&net.Dialer{
	// 		Timeout:   time.Duration(ths.timeout) * time.Second,
	// 		KeepAlive: 30 * time.Second,
	// 	},
	// )
	dialer, err := proxy.SOCKS5("tcp", proxyurl, author, proxy.Direct)

	if err != nil {
		return nil, fmt.Errorf(" ** Socket5ProxyClient() Error\r\n\tproxy.SOCKS5: %s", err)
	}

	transport := &http.Transport{
		// Proxy: nil,
		Dial: dialer.Dial,
		ResponseHeaderTimeout: time.Duration(ths.timeout) * time.Second,
		TLSHandshakeTimeout:   time.Duration(ths.timeout) * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(ths.timeout) * time.Second,
	}, nil
}

// GetClientConn 获取设定超时的Http client connection
func (ths *RClient) GetClientConn(address string, timeout int, unsecure bool) (*httputil.ClientConn, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("Address is not valid URL!\r\n Error:%+v", err)
	}

	var conn net.Conn
	ths.timeout = timeout

	if strings.HasPrefix(address, "https://") {
		config := &tls.Config{}
		config.InsecureSkipVerify = unsecure
		conn, err = tls.Dial("tcp", u.Host, config)
	} else {
		conn, err = net.DialTimeout("tcp", u.Host, time.Duration(ths.timeout)*time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("Create connection has error!\r\n%+v", err)
	}
	return httputil.NewClientConn(conn, nil), nil
}
