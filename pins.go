package creek

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
)

// PinServices provides access to pin related API services.
type PinServices struct {
	client *AuthedClient
}

func NewPinServices(a *AuthedClient) *PinServices { return &PinServices{client: a} }

// List prepares a request for a list of pins.
func (s *PinServices) List() *PinServicesListReq {
	return &PinServicesListReq{
		client: s.client,
		req:    s.client.newReq("/pinning/pins"),
	}
}

type PinServicesListReq struct {
	req
	client *AuthedClient
}

// Context sets the context to be used during this request.
func (r *PinServicesListReq) Context(ctx context.Context) *PinServicesListReq {
	r.req.ctx = ctx
	return r
}

// Send sends the prepared request and returns a list of pins.
func (r *PinServicesListReq) Send() (*PinList, error) {
	res, cleanup, err := r.req.get()
	defer cleanup()
	if err != nil {
		return nil, err
	}

	var data PinList
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, newResponseError(err, res)
	}

	return &data, nil
}

// Add prepares a request to add a pin.
func (s *PinServices) Add(ci cid.Cid) *PinServicesAddReq {
	return &PinServicesAddReq{
		client: s.client,
		req:    s.client.newReq("/pinning/pins"),
		data: IpfsPin{
			Cid:  ci.String(),
			Meta: make(map[string]interface{}),
		},
	}
}

type PinServicesAddReq struct {
	req
	client *AuthedClient
	data   IpfsPin
}

// Context sets the context to be used during this request.
func (r *PinServicesAddReq) Context(ctx context.Context) *PinServicesAddReq {
	r.req.ctx = ctx
	return r
}

// Name sets a name to be associated with the pin.
func (r *PinServicesAddReq) Name(v string) *PinServicesAddReq {
	r.data.Name = v
	return r
}

// Collection sets a collection to be associated with the pin.
func (r *PinServicesAddReq) Collection(v string) *PinServicesAddReq {
	r.data.Meta["collection"] = v
	return r
}

// Origins sets one or more origin addresses to be associated with the pin.
func (r *PinServicesAddReq) Origins(addrs ...peer.AddrInfo) *PinServicesAddReq {
	for _, addr := range addrs {
		r.data.Origins = append(r.data.Origins, addr.String())
	}
	return r
}

// Meta sets additional metadata to be associated with the pin.
func (r *PinServicesAddReq) Meta(meta map[string]interface{}) *PinServicesAddReq {
	for k, v := range meta {
		r.data.Meta[k] = v
	}
	return r
}

// Send sends the prepared request and returns the status of the added pin.
func (r *PinServicesAddReq) Send() (*IpfsPinStatus, error) {
	res, cleanup, err := r.req.postJSON(r.data)
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

// Get prepares a request to get the status of a pin
func (s *PinServices) Get(requestId string) *PinServicesGetReq {
	return &PinServicesGetReq{
		client: s.client,
		req:    s.client.newReq("/pinning/pins/" + url.PathEscape(requestId)),
	}
}

type PinServicesGetReq struct {
	req
	client *AuthedClient
}

// Context sets the context to be used during this request.
func (r *PinServicesGetReq) Context(ctx context.Context) *PinServicesGetReq {
	r.req.ctx = ctx
	return r
}

// Send sends the prepared request and returns the status of the pin.
func (r *PinServicesGetReq) Send() (*IpfsPinStatus, error) {
	res, cleanup, err := r.req.get()
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

// Replace prepares a request to replace a pin.
func (s *PinServices) Replace(requestId string, ci cid.Cid) *PinServicesReplaceReq {
	return &PinServicesReplaceReq{
		client: s.client,
		req:    s.client.newReq("/pinning/pins/" + url.PathEscape(requestId)),
		data: IpfsPin{
			Cid:  ci.String(),
			Meta: make(map[string]interface{}),
		},
	}
}

type PinServicesReplaceReq struct {
	req
	client *AuthedClient
	data   IpfsPin
}

// Context sets the context to be used during this request.
func (r *PinServicesReplaceReq) Context(ctx context.Context) *PinServicesReplaceReq {
	r.req.ctx = ctx
	return r
}

// Name sets a name to be associated with the pin.
func (r *PinServicesReplaceReq) Name(v string) *PinServicesReplaceReq {
	r.data.Name = v
	return r
}

// Collection sets a collection to be associated with the pin.
func (r *PinServicesReplaceReq) Collection(v string) *PinServicesReplaceReq {
	r.data.Meta["collection"] = v
	return r
}

// Origins sets one or more origin addresses to be associated with the pin.
func (r *PinServicesReplaceReq) Origins(addrs ...peer.AddrInfo) *PinServicesReplaceReq {
	for _, addr := range addrs {
		r.data.Origins = append(r.data.Origins, addr.String())
	}
	return r
}

// Meta sets additional metadata to be associated with the pin.
func (r *PinServicesReplaceReq) Meta(meta map[string]interface{}) *PinServicesReplaceReq {
	for k, v := range meta {
		r.data.Meta[k] = v
	}
	return r
}

// Send sends the prepared request and returns the status of the replaced pin.
func (r *PinServicesReplaceReq) Send() (*IpfsPinStatus, error) {
	res, cleanup, err := r.req.postJSON(r.data)
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

// Get prepares a request to delete a pin
func (s *PinServices) Delete(requestId string) *PinServicesDeleteReq {
	return &PinServicesDeleteReq{
		client: s.client,
		req:    s.client.newReq("/pinning/pins/" + url.PathEscape(requestId)),
	}
}

type PinServicesDeleteReq struct {
	req
	client *AuthedClient
}

// Context sets the context to be used during this request.
func (r *PinServicesDeleteReq) Context(ctx context.Context) *PinServicesDeleteReq {
	r.req.ctx = ctx
	return r
}

// Send sends the prepared request and returns the status of the pin.
func (r *PinServicesDeleteReq) Send() error {
	_, cleanup, err := r.req.delete()
	defer cleanup()

	return err
}
