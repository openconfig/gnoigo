// Package system provides gNOI system operations.
package system

import (
	"context"
	"io"

	spb "github.com/openconfig/gnoi/system"
	tpb "github.com/openconfig/gnoi/types"

	"github.com/openconfig/gnoigo/internal"
)

// PingOperation represents input fields required to perform a Ping operation.
type PingOperation struct {
	req *spb.PingRequest
}

// NewPingOperation creates a PingOperation with empty PingRequest.
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

// Interval specifies the nanoseconds between ping requests.
func (p *PingOperation) Interval(i int64) *PingOperation {
	p.req.Interval = i
	return p
}

// Wait specifies nanoseconds to wait for a response.
func (p *PingOperation) Wait(w int64) *PingOperation {
	p.req.Wait = w
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

func (p *PingOperation) Execute(ctx context.Context, c internal.Clients) ([]*spb.PingResponse, error) {
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
