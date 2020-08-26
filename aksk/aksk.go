package aksk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"github.com/ontio/oscore/oscoreconfig"

	"github.com/gin-gonic/gin"
)

var (
	client = http.Client{
		Timeout: time.Second * 15,
	}
)

const (
	maxContentLength = 1024 * 1024
	contentType      = "application/json"
)

type readerCloser struct {
	io.Reader
	io.Closer
}

// SignRequest generate aksk string for request
func SignRequest(req *http.Request) (string, error) {
	h := hmac.New(sha1.New, []byte(oscoreconfig.DefOscoreConfig.SecretKeyV))

	var sb strings.Builder

	// 3 lines
	// 1st
	sb.WriteString(req.Method)
	sb.WriteString(" ")
	sb.WriteString(req.URL.Path)
	if req.URL.RawQuery != "" {
		sb.WriteString("?")
		sb.WriteString(req.URL.RawQuery)
	}
	// 2nd
	sb.WriteString("\nHost: ")
	sb.WriteString(req.Host)
	contentType := req.Header.Get("Content-Type")
	if contentType != "" {
		// 3rd
		sb.WriteString("\nContent-Type: ")
		sb.WriteString(contentType)
	}
	sb.WriteString("\n\n")

	if includeBody(req, contentType) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return "", err
		}

		sb.Write(body)

		req.Body = readerCloser{Reader: bytes.NewBuffer(body), Closer: req.Body}
	}

	str := sb.String()
	io.WriteString(h, str)
	encodeString := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return encodeString, nil
}

func includeBody(req *http.Request, ctType string) bool {
	typeOk := ctType != "" && ctType != "application/octet-stream"
	lengthOk := req.ContentLength > 0 && req.ContentLength < maxContentLength
	return typeOk && lengthOk && req.Body != nil
}

// AkSk middleware
func AkSk(c *gin.Context) {
	authString := c.GetHeader("Authorization")
	if authString == "" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Authorization not found"})
		c.Abort()
		return
	}

	// authString is like: mt ak:signature
	// extract ak
	spaceIdx := strings.IndexByte(authString, ' ')
	if spaceIdx == -1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Authorization not valid"})
		c.Abort()
		return
	}
	colonIdx := strings.IndexByte(authString[spaceIdx:], ':')
	if colonIdx == -1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "Authorization not valid"})
		c.Abort()
		return
	}

	skIdx := spaceIdx + 1 + colonIdx
	ak := authString[spaceIdx+1 : skIdx-1]

	if oscoreconfig.DefOscoreConfig.AccessKeyV != ak {
		c.JSON(http.StatusForbidden, gin.H{"message": "App and Ak not match"})
		c.Abort()
		return
	}

	// verify
	sign, err := SignRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "SignRequest err"})
		c.Abort()
		return
	}

	if sign != authString[skIdx:] {
		c.JSON(http.StatusForbidden, gin.H{"message": "aksk sign not match"})
		c.Abort()
		return
	}

	c.Next()
}

func GetAuthorization(method string, url *url.URL, contentType string, body []byte) string {
	// 第一行Method Path
	firstLine := method + " " + url.Path + "\n"
	secondLine := "Host: " + url.Host + "\n"
	thirdLine := "Content-Type: " + contentType + "\n\n"

	str := firstLine + secondLine + thirdLine + string(body)

	//hmac ,use sha1
	key := []byte(oscoreconfig.DefOscoreConfig.SecretKeyR)
	h := hmac.New(sha1.New, key)

	io.WriteString(h, str)

	encodeString := base64.URLEncoding.EncodeToString(h.Sum(nil))

	code := "mt " + oscoreconfig.DefOscoreConfig.AccessKeyR + ":" + encodeString
	return code
}

func Post(ru *url.URL, data []byte) ([]byte, error) {
	br := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, ru.String(), br)

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", GetAuthorization(http.MethodPost, ru, contentType, data))
	req.Header.Set("App", oscoreconfig.DefOscoreConfig.AppID)

	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("status code expect 200, got: %d, %s", resp.StatusCode, string(b))
	}

	return ioutil.ReadAll(resp.Body)
}
