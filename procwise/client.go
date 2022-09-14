package procwise

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Data map[string]any

func (d Data) Query() url.Values {
	q := make(url.Values)
	for k, v := range d {
		q.Add(k, fmt.Sprintf("%v", v))
	}
	return q
}

type Response struct {
	*http.Response
	Data Data
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientOptions struct {
	Client HttpClient
	Token  string
	Secret string
}

type Client struct {
	Client   HttpClient
	Token    string
	secret   string
	user     string
	password string
}

func NewClient(opts ClientOptions) *Client {
	if opts.Client == nil {
		opts.Client = &http.Client{}
	}
	return &Client{Client: opts.Client, Token: opts.Token, secret: opts.Secret}
}

func (c *Client) SignIn(domain, email, password string) (*Response, error) {
	host := hostFor(domain)
	p := Data{"user": Data{"email": email, "password": password}}
	res, err := c.Request("POST", "https://"+host+"/users/sign_in.json", p)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return res, fmt.Errorf(res.Status)
	}
	var ok bool
	fmt.Println(res.Data)
	c.Token, ok = res.Data["api_token"].(string)
	if !ok {
		return res, fmt.Errorf("missing api_token")
	}
	c.secret, ok = res.Data["secret_key"].(string)
	if !ok {
		return res, fmt.Errorf("missing secret_key")
	}
	c.user = email
	c.password = password
	return res, nil
}

func (c *Client) FindForReview(domain, id string) (*Response, error) {
	if c.Token == "" {
		res, err := c.SignIn(domain, c.user, c.password)
		if err != nil {
			return res, err
		}
	}
	stoken := strings.ReplaceAll(id, "-", "")
	p := Data{"id": stoken}

	host := hostFor(domain)

	return c.ApiRequest("GET", "https://"+host+"/api/v3/student_sessions/"+stoken+"/reviewable", p)
}

func (c *Client) ApiRequest(method string, url string, params Data) (*Response, error) {
	return c.doRequest(func() (*http.Request, error) { return c.newApiRequest(method, url, params) })
}

func (c *Client) Request(method string, url string, params Data) (*Response, error) {
	return c.doRequest(func() (*http.Request, error) { return c.newRequest(method, url, params) })
}

func (c *Client) doRequest(reqFunc func() (*http.Request, error)) (*Response, error) {
	req, err := reqFunc()
	if err != nil {
		return nil, err
	}
	res, err := c.Client.Do(req)
	fmt.Printf("\n====\ncurl -H 'Authorization: Token= %s' -H 'Content-Type: application/json' -H 'Accept: application/vnd.procwise.v3' %s?%s\n====", c.Token, req.URL.String(), req.URL.RawQuery)
	fmt.Printf("req: %+v\n", req.Header)
	fmt.Printf("res: %+v\n", res.Header)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var data Data
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	return &Response{res, data}, nil
}

func (c *Client) newRequest(method string, url string, params Data) (*http.Request, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest(method, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	return r, nil
}

func (c *Client) newApiRequest(method string, uri string, params Data) (r *http.Request, err error) {
	params["nonce"] = time.Now().UnixMilli()
	params["timestamp"] = params["nonce"]
	sig := c.signature(params)
	params["signature"] = sig

	switch method {
	case http.MethodPatch, http.MethodPost, http.MethodPut:
		var b []byte
		b, err = json.Marshal(params)
		if err != nil {
			return
		}
		r, err = http.NewRequest(method, uri, bytes.NewReader(b))
	default:
		r, err = http.NewRequest(method, fmt.Sprintf("%s?%s", uri, params.Query().Encode()), nil)
	}

	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Token token="+c.Token)
	r.Header.Add("Accept", "application/vnd.procwise.v3")
	return r, nil
}

func (c *Client) signature(params Data) string {
	ps := make([]string, 0)
	for k, v := range params {
		ps = append(ps, fmt.Sprintf("%s=%v", k, v))
	}
	data := strings.Join(ps, "?")
	h := hmac.New(sha256.New, []byte(c.secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func hostFor(domain string) string {
	host := domain + ".proctorexam.com"
	if domain == "localhost" {
		host += ":3001"
	}
	return host
}
