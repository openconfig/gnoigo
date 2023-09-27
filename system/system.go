// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// KillProcessOperation represents the parameters of a KillProcess operation.
type KillProcessOperation struct {
	req *spb.KillProcessRequest
}

// NewKillProcessOperation creates an empty KillProcessOperation.
func NewKillProcessOperation() *KillProcessOperation {
	return &KillProcessOperation{req: &spb.KillProcessRequest{}}
}

// PID specifies the process ID of the process to be killed.
func (k *KillProcessOperation) PID(pid uint32) *KillProcessOperation {
	k.req.Pid = pid
	return k
}

// Name specifies the name of the process to be killed.
func (k *KillProcessOperation) Name(n string) *KillProcessOperation {
	k.req.Name = n
	return k
}

// Signal specifies the termination signal sent to the process.
func (k *KillProcessOperation) Signal(s spb.KillProcessRequest_Signal) *KillProcessOperation {
	k.req.Signal = s
	return k
}

// Restart specifies whether the process should be restarted after termination.
func (k *KillProcessOperation) Restart(r bool) *KillProcessOperation {
	k.req.Restart = r
	return k
}

// Execute performs the KillProcess operation.
func (k *KillProcessOperation) Execute(ctx context.Context, c *internal.Clients) (*spb.KillProcessResponse, error) {
	return c.System().KillProcess(ctx, k.req)
}

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

// RebootOperation represents the parameters of a Reboot operation.
type RebootOperation struct {
	rebootMethod  spb.RebootMethod
	delay         time.Duration
	message       string
	subcomponents []*tpb.Path
	force         bool
	wait          time.Duration
	rebootStatus  bool
	cancelStatus  bool
}

// NewRebootOperation creates an empty RebootOperation.
func NewRebootOperation() *RebootOperation {
	return &RebootOperation{}
}

// RebootMethod specifies method to reboot.
func (r *RebootOperation) RebootMethod(rebootMethod spb.RebootMethod) *RebootOperation {
	r.rebootMethod = rebootMethod
	return r
}

// Delay specifies time in nanoseconds to wait before issuing reboot.
func (r *RebootOperation) Delay(delay time.Duration) *RebootOperation {
	r.delay = delay
	return r
}

// Message specifies informational reason for the reboot or cancel reboot.
func (r *RebootOperation) Message(message string) *RebootOperation {
	r.message = message
	return r
}

// Subcomponents specifies the sub-components to reboot.
func (r *RebootOperation) Subcomponents(subcomponents []*tpb.Path) *RebootOperation {
	r.subcomponents = subcomponents
	return r
}

// Force reboot if sanity checks fail.
func (r *RebootOperation) Force(force bool) *RebootOperation {
	r.force = force
	return r
}

// Wait specifies the duration to wait before checking status for reboot or cancel.
func (r *RebootOperation) Wait(wait time.Duration) *RebootOperation {
	r.wait = wait
	return r
}

// RebootWithStatus reboots the subcomponents and returns the status on Execute.
func (r *RebootOperation) RebootWithStatus() *RebootOperation {
	r.rebootStatus = true
	r.cancelStatus = false
	return r
}

// CancelWithStatus cancels the reboot of the subcomponents and returns the status on Execute.
func (r *RebootOperation) CancelWithStatus() *RebootOperation {
	r.rebootStatus = false
	r.cancelStatus = true
	return r
}

// Execute performs the Reboot or Cancel operation.
func (r *RebootOperation) Execute(ctx context.Context, c *internal.Clients) (*spb.RebootStatusResponse, error) {
	if r.rebootStatus {
		_, err := c.System().Reboot(ctx, &spb.RebootRequest{
			Method:        r.rebootMethod,
			Delay:         uint64(r.delay.Nanoseconds()),
			Message:       r.message,
			Subcomponents: r.subcomponents,
			Force:         r.force,
		})
		if err != nil {
			return nil, err
		}
		time.Sleep(r.wait)
	}
	if r.cancelStatus {
		_, err := c.System().CancelReboot(ctx, &spb.CancelRebootRequest{Subcomponents: r.subcomponents})
		if err != nil {
			return nil, err
		}
		time.Sleep(r.wait)
	}
	return c.System().RebootStatus(ctx, &spb.RebootStatusRequest{Subcomponents: r.subcomponents})
}

// SwitchControlProcessorOperation represents the parameters of a SwitchControlProcessor operation.
type SwitchControlProcessorOperation struct {
	req *spb.SwitchControlProcessorRequest
}

// NewSwitchControlProcessorOperation creates an empty SwitchControlProcessorOperation.
func NewSwitchControlProcessorOperation() *SwitchControlProcessorOperation {
	return &SwitchControlProcessorOperation{req: &spb.SwitchControlProcessorRequest{}}
}

// PathFromSubcomponentName sets the path of the target route processor to `/openconfig/components/component[name=<n>]`.
func (s *SwitchControlProcessorOperation) PathFromSubcomponentName(n string) *SwitchControlProcessorOperation {
	return s.Path(&tpb.Path{
		Origin: "openconfig",
		Elem: []*tpb.PathElem{
			{Name: "components"},
			{Name: "component", Key: map[string]string{"name": n}},
		},
	})
}

// Path sets the path of the target route processor.
func (s *SwitchControlProcessorOperation) Path(p *tpb.Path) *SwitchControlProcessorOperation {
	s.req.ControlProcessor = p
	return s
}

// Execute performs the SwitchControlProcessor operation.
func (s *SwitchControlProcessorOperation) Execute(ctx context.Context, c *internal.Clients) (*spb.SwitchControlProcessorResponse, error) {
	return c.System().SwitchControlProcessor(ctx, s.req)
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
