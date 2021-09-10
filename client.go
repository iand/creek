package creek

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ipfs/go-cid"
)

const (

	// Version is the current version of the client library.
	Version = "0.1.0"

	// DefaultUserAgent is the default user agent header used by the client library.
	DefaultUserAgent = "estuary-client/" + Version
	DefaultAddr      = "api.estuary.tech"
)

// Client is the base client used for interacting with services that do not
// require authentication.
type Client struct {
	// never modified once they have been set
	hc   *http.Client
	addr string
	ua   string
}

// New creates a new client that will use the supplied HTTP client and connect
// via the specified API host address.
func New(client *http.Client, addr string) *Client {
	c := &Client{
		hc:   client,
		addr: addr,
	}
	return c
}

// NewDefault creates a new client that will use the default HTTP client and connect
// to api.estuary.tech.
func NewDefault() *Client {
	return New(http.DefaultClient, DefaultAddr)
}

func (c *Client) newReq(path string) req {
	return req{
		hc:   c.hc,
		addr: c.addr,
		path: path,
		headers: headers{
			"User-Agent": c.userAgent(),
		},
		par: params{},
	}
}

func (c *Client) userAgent() string {
	if c.ua == "" {
		return DefaultUserAgent
	}

	return DefaultUserAgent + " " + c.ua
}

// WithToken creates a Client with the supplied authentication token,
// copying options set on the receiver.
func (c *Client) WithToken(token string) *AuthedClient {
	ac := NewAuthedClient(c.hc, c.addr, token)
	ac.ua = c.ua
	return ac
}

// PublicNodeInfo prepares a request for the health of the Estuary node.
func (c *Client) Health() *HealthReq {
	return &HealthReq{
		client: c,
		req:    c.newReq("/health"),
	}
}

type HealthReq struct {
	req
	client *Client
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *HealthReq) Context(ctx context.Context) *HealthReq {
	r.req.ctx = ctx
	return r
}

// Send sends the prepared request and returns the health status reported by the Estuary node.
func (r *HealthReq) Send() (*Health, error) {
	res, cleanup, err := r.req.get()
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var data Health
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, newResponseError(err, res)
	}

	return &data, nil
}

func (c *Client) PublicStats() *StatsReq {
	return &StatsReq{
		client: c,
		req:    c.newReq("/public/stats"),
	}
}

type StatsReq struct {
	req
	client *Client
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *StatsReq) Context(ctx context.Context) *StatsReq {
	r.req.ctx = ctx
	return r
}

// Send sends the prepared request and decodes the response.
func (r *StatsReq) Send() (*PublicStats, error) {
	res, cleanup, err := r.req.get()
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var data PublicStats
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, newResponseError(err, res)
	}

	return &data, nil
}

// PublicNodeInfo prepares a request for information about the Estuary node.
func (c *Client) PublicNodeInfo() *PublicNodeInfoReq {
	return &PublicNodeInfoReq{
		client: c,
		req:    c.newReq("/public/info"),
	}
}

type PublicNodeInfoReq struct {
	req
	client *Client
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *PublicNodeInfoReq) Context(ctx context.Context) *PublicNodeInfoReq {
	r.req.ctx = ctx
	return r
}

// Send sends the prepared request and returns public information about the Estuary node.
func (r *PublicNodeInfoReq) Send() (*PublicNodeInfo, error) {
	res, cleanup, err := r.req.get()
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var data PublicNodeInfo
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, newResponseError(err, res)
	}

	return &data, nil
}

// PublicNodeInfo prepares a request for information about content by its cid
func (c *Client) PublicContentByCid(ci cid.Cid) *PublicContentByCidReq {
	return &PublicContentByCidReq{
		client: c,
		req:    c.newReq("/public/by-cid/" + ci.String()),
	}
}

type PublicContentByCidReq struct {
	req
	client *Client
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *PublicContentByCidReq) Context(ctx context.Context) *PublicContentByCidReq {
	r.req.ctx = ctx
	return r
}

// Send sends the prepared request and returns public information about the content.
func (r *PublicContentByCidReq) Send() ([]ContentInfo, error) {
	res, cleanup, err := r.req.get()
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var data []ContentInfo
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, newResponseError(err, res)
	}

	return data, nil
}
