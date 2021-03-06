package creek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type req struct {
	hc      *http.Client
	ctx     context.Context
	addr    string
	path    string
	par     params
	headers headers
}

type headers map[string]string

type params map[string][]string

func (p params) Get(key string) string {
	vs := p[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func (p params) Set(key, value string) {
	p[key] = []string{value}
}

func (p params) Encode() string {
	return url.Values(p).Encode()
}

func (r *req) url() *url.URL {
	u := url.URL{
		Scheme:   "https",
		Host:     r.addr,
		Path:     r.path,
		RawQuery: r.par.Encode(),
	}
	return &u
}

func (r *req) do(hr *http.Request) (*http.Response, func(), error) {
	if r.ctx != nil {
		hr = hr.WithContext(r.ctx)
	}
	for k, v := range r.headers {
		hr.Header.Set(k, v)
	}

	res, err := r.hc.Do(hr)
	if err != nil {
		return nil, func() {}, err
	}
	if err := responseError(res); err != nil {
		return nil, func() {}, err
	}

	return res, cleanup(res), nil
}

func (r *req) get() (*http.Response, func(), error) {
	hr, err := http.NewRequest("GET", r.url().String(), nil)
	if err != nil {
		return nil, func() {}, err
	}
	return r.do(hr)
}

func (r *req) post(ir io.Reader) (*http.Response, func(), error) {
	hr, err := http.NewRequest("POST", r.url().String(), ir)
	if err != nil {
		return nil, func() {}, err
	}
	return r.do(hr)
}

func (r *req) postJSON(data interface{}) (*http.Response, func(), error) {
	var body io.Reader
	if data != nil {
		var encoded bytes.Buffer
		err := json.NewEncoder(&encoded).Encode(data)
		if err != nil {
			return nil, func() {}, err
		}
		body = &encoded
		r.headers["Content-Type"] = "application/json"
		r.headers["Content-Length"] = strconv.Itoa(encoded.Len())
	}

	return r.post(body)
}

func (r *req) delete() (*http.Response, func(), error) {
	hr, err := http.NewRequest("DELETE", r.url().String(), nil)
	if err != nil {
		return nil, func() {}, err
	}
	return r.do(hr)
}

func cleanup(res *http.Response) func() {
	return func() {
		if res == nil || res.Body == nil {
			return
		}
		res.Body.Close()
	}
}

type ResponseError struct {
	Message    string
	StatusCode int
	URL        string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("request failed: %s", e.Message)
}

func newResponseError(err error, res *http.Response) error {
	rerr := &ResponseError{
		Message: err.Error(),
	}

	if res != nil {
		rerr.StatusCode = res.StatusCode
		rerr.URL = res.Request.URL.String()
	}

	return rerr
}

func responseError(res *http.Response) error {
	if res == nil {
		return &ResponseError{
			Message: "no response found",
		}
	}
	if res.StatusCode/100 == 2 {
		return nil
	}

	rerr := &ResponseError{
		StatusCode: res.StatusCode,
		URL:        res.Request.URL.String(),
	}

	if res.Body == nil {
		return rerr
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		rerr.Message = fmt.Sprintf("unable to read response body: %v", err)
		return rerr
	}

	var serr Error
	err = json.Unmarshal(body, &serr)
	if err != nil {
		rerr.Message = fmt.Sprintf("unable to unmarshal error response: %v", err)
		return rerr
	}

	rerr.Message = serr.Error
	return rerr
}
