package creek

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
)

// AuthedClient is the base client used for interacting with services that
// require authentication.
type AuthedClient struct {
	// never modified once they have been set
	hc    *http.Client
	addr  string
	ua    string
	token string
}

func (c *AuthedClient) userAgent() string {
	if c.ua == "" {
		return DefaultUserAgent
	}

	return DefaultUserAgent + " " + c.ua
}

func (c *AuthedClient) newReq(path string) req {
	return req{
		hc:   c.hc,
		addr: c.addr,
		path: path,
		headers: headers{
			"User-Agent":    c.userAgent(),
			"Authorization": "Bearer " + c.token,
		},
		par: params{},
	}
}

// New creates a new client that will use the supplied HTTP client and connect
// via the specified API host address.
func NewAuthedClient(client *http.Client, addr string, token string) *AuthedClient {
	c := &AuthedClient{
		hc:    client,
		addr:  addr,
		token: token,
	}
	return c
}

func (c *AuthedClient) ContentAdd(name string, r io.Reader) *ContentAddReq {
	return &ContentAddReq{
		client: c,
		req:    c.newReq("/content/add"),
		name:   name,
		r:      r,
	}
}

type ContentAddReq struct {
	req
	client *AuthedClient
	name   string
	r      io.Reader
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *ContentAddReq) Context(ctx context.Context) *ContentAddReq {
	r.req.ctx = ctx
	return r
}

func (r *ContentAddReq) Send() (*AddedContent, error) {
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	go func() {
		var outerr error
		defer func() {
			if outerr != nil {
				pw.CloseWithError(outerr)
			} else {
				pw.Close()
			}
		}()

		part, err := mw.CreateFormFile("data", r.name)
		if err != nil {
			outerr = err
			return
		}

		_, err = io.Copy(part, r.r)
		if err != nil {
			outerr = err
			return
		}
		mw.Close()
	}()

	r.req.headers["Content-Type"] = mw.FormDataContentType()

	res, cleanup, err := r.req.post(pr)
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var data AddedContent
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, newResponseError(err, res)
	}

	return &data, nil
}

func (c *AuthedClient) ContentAddFromIpfs(root cid.Cid) *ContentAddFromIpfsReq {
	req := c.newReq("/content/add-ipfs")
	req.par.Set("root", root.String())

	return &ContentAddFromIpfsReq{
		client: c,
		req:    req,
	}
}

type ContentAddFromIpfsReq struct {
	req
	client *AuthedClient
}

// Context sets the context to be used during this request. If no context is supplied then
// the request will use context.Background.
func (r *ContentAddFromIpfsReq) Context(ctx context.Context) *ContentAddFromIpfsReq {
	r.req.ctx = ctx
	return r
}

func (r *ContentAddFromIpfsReq) Name(v string) *ContentAddFromIpfsReq {
	r.req.par.Set("name", v)
	return r
}

func (r *ContentAddFromIpfsReq) Collection(v string) *ContentAddFromIpfsReq {
	r.req.par.Set("collection", v)
	return r
}

func (r *ContentAddFromIpfsReq) Peers(peers ...peer.AddrInfo) *ContentAddFromIpfsReq {
	s := make([]string, len(peers))
	for i := range peers {
		s[i] = peers[i].String()
	}

	r.req.par.Set("peers", strings.Join(s, ","))
	return r
}

func (r *ContentAddFromIpfsReq) Send() (*IpfsPinStatus, error) {
	res, cleanup, err := r.req.post(nil)
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var data IpfsPinStatus
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, newResponseError(err, res)
	}

	return &data, nil
}
