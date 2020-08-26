package http

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/ontio/oscore/models/tables"
	"io/ioutil"
	"net/http"
	"time"
)

var DefClient *Client

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	tr := &http.Transport{ //x509: certificate signed by unknown authority
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: tr, //x509: certificate signed by unknown authority
	}
	return &Client{
		client: client,
	}
}

func (this *Client) GetWithHeader(url string, headers []*tables.ApiHeadValues) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	for _, h := range headers {
		req.Header.Set(h.HeaderKey, h.HeaderValue)
	}

	resp, err := this.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("[Get] send http get request error:%s", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("[Get] read http body error:%s", err)
	}
	return data, resp.StatusCode, nil
}

func (this *Client) PostWithHeader(url string, headers []*tables.ApiHeadValues, bodyParam []byte) ([]byte, int, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyParam))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	for _, h := range headers {
		req.Header.Set(h.HeaderKey, h.HeaderValue)
	}

	resp, err := this.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("[Get] send http get request error:%s", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("[Post] read http body error:%s", err)
	}
	return data, resp.StatusCode, nil
}

func (this *Client) Get(url string) ([]byte, error) {
	resp, err := this.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("[Get] send http get request error:%s", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[Get] read http body error:%s", err)
	}
	return data, nil
}

func (this *Client) Post(url string, bodyParam []byte) ([]byte, error) {
	// to do. other common type
	resp, err := this.client.Post(url, "application/json", bytes.NewReader(bodyParam))
	if err != nil {
		return nil, fmt.Errorf("[Post] send http post request error:%s", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[Post] read http body error:%s", err)
	}
	return data, nil
}
