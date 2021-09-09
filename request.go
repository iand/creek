package creek

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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

func (r *req) get() (*http.Response, func(), error) {
	req, err := http.NewRequest("GET", r.url().String(), nil)
	if err != nil {
		return nil, func() {}, err
	}
	if r.ctx != nil {
		req = req.WithContext(r.ctx)
	}
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	res, err := r.hc.Do(req)
	if err != nil {
		return nil, func() {}, err
	}
	if err := responseError(res); err != nil {
		return nil, func() {}, err
	}

	return res, cleanup(res), nil
}

func (r *req) post(d io.Reader) (*http.Response, func(), error) {
	req, err := http.NewRequest("POST", r.url().String(), d)
	if err != nil {
		return nil, func() {}, err
	}
	if r.ctx != nil {
		req = req.WithContext(r.ctx)
	}
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	res, err := r.hc.Do(req)
	if err != nil {
		return nil, func() {}, err
	}
	if err := responseError(res); err != nil {
		return nil, func() {}, err
	}
	return res, cleanup(res), nil
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
