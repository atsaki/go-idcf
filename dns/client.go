package dns

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const EndPoint = "https://dns.idcfcloud.com"

type Client struct {
	EndPoint  string
	APIKey    string
	SecretKey string

	httpClient *http.Client
	reqURL     *url.URL
}

func NewClient(apikey, secretkey string) (*Client, error) {
	c := &Client{
		EndPoint:  EndPoint,
		APIKey:    apikey,
		SecretKey: secretkey,
	}

	c.httpClient = &http.Client{}
	u, err := url.Parse(EndPoint)
	if err != nil {
		return nil, err
	}
	c.reqURL = u

	return c, nil
}

func (c *Client) Request(method, path string, param map[string]interface{}) ([]byte, error) {

	c.reqURL.Path = path

	expires := strconv.FormatInt(time.Now().Add(600*time.Second).Unix(), 10)

	h := hmac.New(sha256.New, []byte(c.SecretKey))
	h.Write([]byte(strings.Join([]string{method, path, c.APIKey, expires, ""}, "\n")))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	b, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	paramReader := strings.NewReader(string(b))

	req, err := http.NewRequest(method, c.reqURL.String(), paramReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-IDCF-APIKEY", c.APIKey)
	req.Header.Set("X-IDCF-Expires", expires)
	req.Header.Set("X-IDCF-Signature", signature)

	if method == "POST" || method == "PUT" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
