// Package system gives apis for system operations.
package system

import (
	"context"
	"io"

	spb "google3/third_party/openconfig/gnoi/system/system_go_proto"
	tpb "google3/third_party/openconfig/gnoi/types/types_go_proto"
)

// Client is the client used to access operations from system service.
type Client struct {
	Client spb.SystemClient
}

// Operation will be any X operation from the system client
type Operation[T any] interface {
	execute(context.Context, *Client) (T, error)
}

// Execute function takes input from "Operation" (which is equivalent to the request proto) 
// and returns the response proto based on Operation.
// E.g. `PingOperation` returns `spb.PingResponse`.
func Execute[T any](ctx context.Context, sc *Client, so Operation[T]) (T, error) {
	return so.execute(ctx, sc)
}

// PingOperation represents fields of `PingRequest` proto.
type PingOperation struct {
	destination   string
	source        string
	count         int32
	interval      int64
	wait          int64
	size          int32
	doNotFragment bool
	doNotResolve  bool
	l3Protocol    tpb.L3Protocol
}

func (p *PingOperation) execute(ctx context.Context, sc *Client) ([]*spb.PingResponse, error) {
	ping, err := sc.Client.Ping(ctx, &spb.PingRequest{
		Destination:   p.destination,
		Source:        p.source,
		Count:         p.count,
		Interval:      p.interval,
		Wait:          p.wait,
		Size:          p.size,
		DoNotFragment: p.doNotFragment,
		DoNotResolve:  p.doNotResolve,
		L3Protocol:    p.l3Protocol,
	})
	if err != nil {
		return nil, err
	}

	pingResp := []*spb.PingResponse{}
	for {
		resp, err := ping.Recv()
		switch {
		case err == io.EOF:
			return pingResp, nil
		case err != nil:
			return nil, err
		default:
			pingResp = append(pingResp, resp)
		}
	}
}

// TracerouteOperation represents fields of TracerouteRequest proto.
type TracerouteOperation struct {
	source         string
	destination    string
	initialTTL     uint32
	maxTTL         int32
	wait           int64
	doNotFragment  bool
	doNotResolve   bool
	l3Protocol     tpb.L3Protocol
	l4Protocol     spb.TracerouteRequest_L4Protocol
	doNotLookupAsn bool
}

func (t *TracerouteOperation) execute(ctx context.Context, sc *Client) ([]*spb.TracerouteResponse, error) {
	traceroute, err := sc.Client.Traceroute(ctx, &spb.TracerouteRequest{
		Source:        t.source,
		Destination:   t.destination,
		InitialTtl:    t.initialTTL,
		MaxTtl:        t.maxTTL,
		Wait:          t.wait,
		DoNotFragment: t.doNotFragment,
		DoNotResolve:  t.doNotResolve,
		L3Protocol:    t.l3Protocol,
		L4Protocol:    t.l4Protocol,
	})
	if err != nil {
		return nil, err
	}

	traceResp := []*spb.TracerouteResponse{}
	for {
		resp, err := traceroute.Recv()
		switch {
		case err == io.EOF:
			return traceResp, nil
		case err != nil:
			return nil, err
		default:
			traceResp = append(traceResp, resp)
		}
	}
}

