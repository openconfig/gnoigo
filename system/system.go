// Package system provides gNOI system operations.
package system

import (
	"context"
	"io"
	"time"

	spb "github.com/openconfig/gnoi/system"
	tpb "github.com/openconfig/gnoi/types"

	"github.com/openconfig/gnoigo/internal"
)

// PingOperation represents the parameters of a Ping operation.
type PingOperation struct {
	req *spb.PingRequest
}

// NewPingOperation creates an empty PingOperation.
func NewPingOperation() *PingOperation {
	return &PingOperation{req: &spb.PingRequest{}}
}

// Destination specifies the address to ping.
func (p *PingOperation) Destination(dst string) *PingOperation {
	p.req.Destination = dst
	return p
}

// Source specifies the address to ping from.
func (p *PingOperation) Source(src string) *PingOperation {
	p.req.Source = src
	return p
}

// Count specifies the number of packets.
func (p *PingOperation) Count(c int32) *PingOperation {
	p.req.Count = c
	return p
}

// Interval specifies the duration between ping requests.
func (p *PingOperation) Interval(i time.Duration) *PingOperation {
	p.req.Interval = i.Nanoseconds()
	return p
}

// Wait specifies the duration to wait for a response.
func (p *PingOperation) Wait(w time.Duration) *PingOperation {
	p.req.Wait = w.Nanoseconds()
	return p
}

// Size specifies the size of request packet (excluding ICMP header).
func (p *PingOperation) Size(s int32) *PingOperation {
	p.req.Size = s
	return p
}

// DoNotFragment sets the do not fragment bit (IPv4 destinations).
func (p *PingOperation) DoNotFragment(dnf bool) *PingOperation {
	p.req.DoNotFragment = dnf
	return p
}

// DoNotResolve specifies if address returned should be resolved.
func (p *PingOperation) DoNotResolve(dnr bool) *PingOperation {
	p.req.DoNotResolve = dnr
	return p
}

// L3Protocol specifies layer3 protocol for the ping.
func (p *PingOperation) L3Protocol(l3p tpb.L3Protocol) *PingOperation {
	p.req.L3Protocol = l3p
	return p
}

// Execute performs the Ping operation.
func (p *PingOperation) Execute(ctx context.Context, c *internal.Clients) ([]*spb.PingResponse, error) {
	ping, err := c.System().Ping(ctx, p.req)
	if err != nil {
		return nil, err
	}

	var pingResp []*spb.PingResponse

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

// TimeOperation represents the parameters of a Time operation.
type TimeOperation struct {
	req *spb.TimeRequest
}

// NewTimeOperation creates an empty TimeOperation.
func NewTimeOperation() *TimeOperation {
	return &TimeOperation{req: &spb.TimeRequest{}}
}

// Execute performs the Time operation.
func (t *TimeOperation) Execute(ctx context.Context, c *internal.Clients) (*spb.TimeResponse, error) {
	return c.System().Time(ctx, t.req)
}

// TracerouteOperation represents the parameters of a Traceroute operation.
type TracerouteOperation struct {
	req *spb.TracerouteRequest
}

// NewTracerouteOperation creates an empty TracerouteOperation.
func NewTracerouteOperation() *TracerouteOperation {
	return &TracerouteOperation{req: &spb.TracerouteRequest{}}
}

// Source specifies address to perform traceroute from.
func (t *TracerouteOperation) Source(src string) *TracerouteOperation {
	t.req.Source = src
	return t
}

// Destination specifies address to perform traceroute to.
func (t *TracerouteOperation) Destination(dst string) *TracerouteOperation {
	t.req.Destination = dst
	return t
}

// InitialTTL specifies traceroute ttl (default is 1).
func (t *TracerouteOperation) InitialTTL(ttl uint32) *TracerouteOperation {
	t.req.InitialTtl = ttl
	return t
}

// MaxTTL specifies maximum number of hops.
func (t *TracerouteOperation) MaxTTL(ttl int32) *TracerouteOperation {
	t.req.MaxTtl = ttl
	return t
}

// Wait specifies the duration to wait for a response.
func (t *TracerouteOperation) Wait(wait time.Duration) *TracerouteOperation {
	t.req.Wait = wait.Nanoseconds()
	return t
}

// DoNotFragment sets the do not fragment bit. (IPv4 destinations)
func (t *TracerouteOperation) DoNotFragment(dnf bool) *TracerouteOperation {
	t.req.DoNotFragment = dnf
	return t
}

// DoNotResolve specifies if address returned should be resolved.
func (t *TracerouteOperation) DoNotResolve(dnr bool) *TracerouteOperation {
	t.req.DoNotFragment = dnr
	return t
}

// L3Protocol specifies layer3 protocol for the traceroute.
func (t *TracerouteOperation) L3Protocol(l3 tpb.L3Protocol) *TracerouteOperation {
	t.req.L3Protocol = l3
	return t
}

// L4Protocol specifies layer3 protocol for the traceroute.
func (t *TracerouteOperation) L4Protocol(l4 spb.TracerouteRequest_L4Protocol) *TracerouteOperation {
	t.req.L4Protocol = l4
	return t
}

// DoNotLookupASN specifies if traceroute should try to lookup ASN.
func (t *TracerouteOperation) DoNotLookupASN(asn bool) *TracerouteOperation {
	t.req.DoNotLookupAsn = asn
	return t
}

// Execute performs the Traceroute operation.
func (t *TracerouteOperation) Execute(ctx context.Context, c *internal.Clients) ([]*spb.TracerouteResponse, error) {
	traceroute, err := c.System().Traceroute(ctx, t.req)
	if err != nil {
		return nil, err
	}

	var tracerouteResp []*spb.TracerouteResponse

	for {
		resp, err := traceroute.Recv()
		switch {
		case err == io.EOF:
			return tracerouteResp, nil
		case err != nil:
			return nil, err
		default:
			tracerouteResp = append(tracerouteResp, resp)
		}
	}
}
