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

package system_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/google/go-cmp/cmp"
	spb "github.com/openconfig/gnoi/system"
	tpb "github.com/openconfig/gnoi/types"
	"github.com/openconfig/gnoigo/internal"
	"github.com/openconfig/gnoigo/system"
)

type fakeSystemClient struct {
	spb.SystemClient
	KillProcessFn            func(context.Context, *spb.KillProcessRequest, ...grpc.CallOption) (*spb.KillProcessResponse, error)
	PingFn                   func(context.Context, *spb.PingRequest, ...grpc.CallOption) (spb.System_PingClient, error)
	RebootFn                 func(context.Context, *spb.RebootRequest, ...grpc.CallOption) (*spb.RebootResponse, error)
	RebootStatusFn           func(context.Context, *spb.RebootStatusRequest, ...grpc.CallOption) (*spb.RebootStatusResponse, error)
	SwitchControlProcessorFn func(context.Context, *spb.SwitchControlProcessorRequest, ...grpc.CallOption) (*spb.SwitchControlProcessorResponse, error)
	TimeFn                   func(context.Context, *spb.TimeRequest, ...grpc.CallOption) (*spb.TimeResponse, error)
	TracerouteFn             func(context.Context, *spb.TracerouteRequest, ...grpc.CallOption) (spb.System_TracerouteClient, error)
}

func (fg *fakeSystemClient) System() spb.SystemClient {
	return fg
}

func (fg *fakeSystemClient) KillProcess(ctx context.Context, in *spb.KillProcessRequest, opts ...grpc.CallOption) (*spb.KillProcessResponse, error) {
	return fg.KillProcessFn(ctx, in, opts...)
}

func (fg *fakeSystemClient) Ping(ctx context.Context, in *spb.PingRequest, opts ...grpc.CallOption) (spb.System_PingClient, error) {
	return fg.PingFn(ctx, in, opts...)
}

func (fg *fakeSystemClient) SwitchControlProcessor(ctx context.Context, in *spb.SwitchControlProcessorRequest, opts ...grpc.CallOption) (*spb.SwitchControlProcessorResponse, error) {
	return fg.SwitchControlProcessorFn(ctx, in, opts...)
}

func (fg *fakeSystemClient) Reboot(ctx context.Context, in *spb.RebootRequest, opts ...grpc.CallOption) (*spb.RebootResponse, error) {
	return fg.RebootFn(ctx, in, opts...)
}

func (fg *fakeSystemClient) RebootStatus(ctx context.Context, in *spb.RebootStatusRequest, opts ...grpc.CallOption) (*spb.RebootStatusResponse, error) {
	return fg.RebootStatusFn(ctx, in, opts...)
}

func (fg *fakeSystemClient) Time(ctx context.Context, in *spb.TimeRequest, opts ...grpc.CallOption) (*spb.TimeResponse, error) {
	return fg.TimeFn(ctx, in, opts...)
}

func (fg *fakeSystemClient) Traceroute(ctx context.Context, in *spb.TracerouteRequest, opts ...grpc.CallOption) (spb.System_TracerouteClient, error) {
	return fg.TracerouteFn(ctx, in, opts...)
}

func TestKillProcess(t *testing.T) {
	tests := []struct {
		desc    string
		op      *system.KillProcessOperation
		want    *spb.KillProcessResponse
		wantErr string
	}{
		{
			desc: "Test KillProcess",
			op:   system.NewKillProcessOperation().PID(1234),
			want: &spb.KillProcessResponse{},
		},
		{
			desc:    "KillProcess returns error",
			op:      system.NewKillProcessOperation(),
			wantErr: "KillProcess operation error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var fakeClient internal.Clients
			fakeClient.SystemClient = &fakeSystemClient{KillProcessFn: func(context.Context, *spb.KillProcessRequest, ...grpc.CallOption) (*spb.KillProcessResponse, error) {
				if tt.wantErr != "" {
					return nil, fmt.Errorf(tt.wantErr)
				}
				return tt.want, nil
			}}

			got, gotErr := tt.op.Execute(context.Background(), &fakeClient)
			if (gotErr == nil) != (tt.wantErr == "") || (gotErr != nil && !strings.Contains(gotErr.Error(), tt.wantErr)) {
				t.Errorf("Execute() got unexpected error %v want %s", gotErr, tt.wantErr)
			}
			if tt.want != got {
				t.Errorf("Execute() got unexpected response want %v got %v", tt.want, got)
			}
		})
	}
}

type fakePingClient struct {
	spb.System_PingClient
	resp []*spb.PingResponse
	err  error
}

func (pc *fakePingClient) Recv() (*spb.PingResponse, error) {
	if len(pc.resp) == 0 && pc.err == nil {
		return nil, io.EOF
	}
	resp := pc.resp[0]
	pc.resp = pc.resp[1:]
	return resp, pc.err
}

func TestPing(t *testing.T) {
	tests := []struct {
		desc    string
		op      *system.PingOperation
		want    []*spb.PingResponse
		wantErr string
	}{
		{
			desc: "ping with source",
			op:   system.NewPingOperation().Destination("1.2.3.4").Source("5.6.7.8"),
			want: []*spb.PingResponse{{Source: "5.6.7.8"}},
		},
		{
			desc: "ping with source and count",
			op:   system.NewPingOperation().Destination("1.2.3.4").Source("5.6.7.8").Count(7),
			want: []*spb.PingResponse{{Source: "5.6.7.8", Sent: 7, Received: 7}},
		},
		{
			desc: "ping with multiple response",
			op:   system.NewPingOperation().Destination("1.2.3.4").Source("5.6.7.8").Count(7),
			want: []*spb.PingResponse{{Source: "5.6.7.8", Sent: 1, Received: 1}, {Source: "5.6.7.8", Sent: 2, Received: 2}},
		},
		{
			desc:    "ping returns error",
			op:      system.NewPingOperation().Destination("1.2.3.4").Source("5.6.7.8").Count(7),
			wantErr: "ping operation error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var fakeClient internal.Clients
			fakeClient.SystemClient = &fakeSystemClient{PingFn: func(context.Context, *spb.PingRequest, ...grpc.CallOption) (spb.System_PingClient, error) {
				if tt.wantErr != "" {
					return nil, fmt.Errorf(tt.wantErr)
				}
				return &fakePingClient{resp: tt.want}, nil
			}}

			got, gotErr := tt.op.Execute(context.Background(), &fakeClient)
			if (gotErr == nil) != (tt.wantErr == "") || (gotErr != nil && !strings.Contains(gotErr.Error(), tt.wantErr)) {
				t.Errorf("Execute() got unexpected error %v want %s", gotErr, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Execute() got unexpected response diff (-want +got): %s", diff)
			}
		})
	}
}

func TestReboot(t *testing.T) {
	tests := []struct {
		desc               string
		op                 *system.RebootOperation
		want               *spb.RebootResponse
		rebootErr, wantErr string
		statusErrs         []error
		statusResps        []*spb.RebootStatusResponse
		cancelContext      bool
	}{
		{
			desc: "Test reboot",
			op: system.NewRebootOperation().RebootMethod(spb.RebootMethod_COLD).Subcomponents([]*tpb.Path{{
				Elem: []*tpb.PathElem{
					{Name: "components"},
					{Name: "component", Key: map[string]string{"name": "RP0"}},
				},
			}}),
			want: &spb.RebootResponse{},
		},
		{
			desc:        "Test reboot wait for active status",
			op:          system.NewRebootOperation().RebootMethod(spb.RebootMethod_COLD).WaitForActive(true),
			statusResps: []*spb.RebootStatusResponse{{Active: true}},
			want:        &spb.RebootResponse{},
		},
		{
			desc:        "Test reboot wait for active status and ignore unavailable error",
			op:          system.NewRebootOperation().RebootMethod(spb.RebootMethod_COLD).WaitForActive(true).IgnoreUnavailableErr(true),
			statusErrs:  []error{status.Errorf(codes.Unavailable, "unavailable")},
			statusResps: []*spb.RebootStatusResponse{{Active: true}},
			want:        &spb.RebootResponse{},
		},
		{
			desc:        "Test reboot with non active status response",
			op:          system.NewRebootOperation().RebootMethod(spb.RebootMethod_COLD).WaitForActive(true),
			statusResps: []*spb.RebootStatusResponse{{Active: false, Wait: 2}, {Active: true}},
			want:        &spb.RebootResponse{},
		},
		{
			desc:       "Test reboot wait for active status returns unknown error",
			op:         system.NewRebootOperation().RebootMethod(spb.RebootMethod_COLD).WaitForActive(true).IgnoreUnavailableErr(true),
			statusErrs: []error{status.Errorf(codes.Unknown, "unknown")},
			wantErr:    "unknown",
		},
		{
			desc:      "Test reboot returns error on reboot",
			op:        system.NewRebootOperation().RebootMethod(spb.RebootMethod_COLD),
			rebootErr: "Reboot operation error",
			wantErr:   "Reboot operation error",
		},
		{
			desc:       "Test reboot returns error on reboot status",
			op:         system.NewRebootOperation().RebootMethod(spb.RebootMethod_COLD).WaitForActive(true),
			statusErrs: []error{status.Errorf(codes.Unavailable, "unavailable")},
			wantErr:    "unavailable",
		},
		{
			desc:          "Test reboot with context cancel",
			op:            system.NewRebootOperation().RebootMethod(spb.RebootMethod_COLD).WaitForActive(true),
			statusResps:   []*spb.RebootStatusResponse{{Wait: 20}},
			wantErr:       "context",
			cancelContext: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var fakeClient internal.Clients
			fakeClient.SystemClient = &fakeSystemClient{
				RebootFn: func(context.Context, *spb.RebootRequest, ...grpc.CallOption) (*spb.RebootResponse, error) {
					if tt.rebootErr != "" {
						return nil, fmt.Errorf(tt.rebootErr)
					}
					return tt.want, nil
				},
				RebootStatusFn: func(context.Context, *spb.RebootStatusRequest, ...grpc.CallOption) (*spb.RebootStatusResponse, error) {
					if len(tt.statusErrs) > 0 {
						statusErr := tt.statusErrs[0]
						tt.statusErrs = tt.statusErrs[1:]
						return nil, statusErr
					}
					if len(tt.statusResps) > 0 {
						statusResp := tt.statusResps[0]
						tt.statusResps = tt.statusResps[1:]
						return statusResp, nil
					}
					return &spb.RebootStatusResponse{}, nil
				}}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tt.cancelContext {
				cancel()
			}

			got, gotErr := tt.op.Execute(ctx, &fakeClient)
			if (gotErr == nil) != (tt.wantErr == "") || (gotErr != nil && !strings.Contains(gotErr.Error(), tt.wantErr)) {
				t.Errorf("Execute() got unexpected error %v want %s", gotErr, tt.wantErr)
			}
			if tt.want != got {
				t.Errorf("Execute() got unexpected response want %v got %v", tt.want, got)
			}
		})
	}
}

func TestSwitchControlProcessor(t *testing.T) {
	tests := []struct {
		desc    string
		op      *system.SwitchControlProcessorOperation
		want    *spb.SwitchControlProcessorResponse
		wantErr string
	}{
		{
			desc: "Test SwitchControlProcessor with Path",
			op: system.NewSwitchControlProcessorOperation().Path(&tpb.Path{
				Origin: "openconfig",
				Elem: []*tpb.PathElem{
					{Name: "components"},
					{Name: "component", Key: map[string]string{"name": "RP0"}},
				},
			}),
			want: &spb.SwitchControlProcessorResponse{Version: "new"},
		},
		{
			desc: "Test SwitchControlProcessor with PathFromSubcomponentName",
			op:   system.NewSwitchControlProcessorOperation().PathFromSubcomponentName("RP0"),
			want: &spb.SwitchControlProcessorResponse{Version: "new"},
		},
		{
			desc: "Test SwitchControlProcessor with PathFromSubcomponentName and Path returns error",
			op: system.NewSwitchControlProcessorOperation().PathFromSubcomponentName("RP0").Path(&tpb.Path{
				Elem: []*tpb.PathElem{
					{Name: "components"},
					{Name: "component", Key: map[string]string{"name": "RP0"}},
				},
			}),
			want: &spb.SwitchControlProcessorResponse{Version: "new"},
		},
		{
			desc:    "SwitchControlProcessor returns error",
			op:      system.NewSwitchControlProcessorOperation(),
			wantErr: "SwitchControlProcessor operation error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var fakeClient internal.Clients
			fakeClient.SystemClient = &fakeSystemClient{SwitchControlProcessorFn: func(context.Context, *spb.SwitchControlProcessorRequest, ...grpc.CallOption) (*spb.SwitchControlProcessorResponse, error) {
				if tt.wantErr != "" {
					return nil, fmt.Errorf(tt.wantErr)
				}
				return tt.want, nil
			}}

			got, gotErr := tt.op.Execute(context.Background(), &fakeClient)
			if (gotErr == nil) != (tt.wantErr == "") || (gotErr != nil && !strings.Contains(gotErr.Error(), tt.wantErr)) {
				t.Errorf("Execute() got unexpected error %v want %s", gotErr, tt.wantErr)
			}
			if tt.want != got {
				t.Errorf("Execute() got unexpected response want %v got %v", tt.want, got)
			}
		})
	}
}

func TestTime(t *testing.T) {
	tests := []struct {
		desc    string
		op      *system.TimeOperation
		want    *spb.TimeResponse
		wantErr string
	}{
		{
			desc: "Test time",
			op:   system.NewTimeOperation(),
			want: &spb.TimeResponse{Time: 1234},
		},
		{
			desc:    "Time returns error",
			op:      system.NewTimeOperation(),
			wantErr: "Time operation error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var fakeClient internal.Clients
			fakeClient.SystemClient = &fakeSystemClient{TimeFn: func(context.Context, *spb.TimeRequest, ...grpc.CallOption) (*spb.TimeResponse, error) {
				if tt.wantErr != "" {
					return nil, fmt.Errorf(tt.wantErr)
				}
				return tt.want, nil
			}}

			got, gotErr := tt.op.Execute(context.Background(), &fakeClient)
			if (gotErr == nil) != (tt.wantErr == "") || (gotErr != nil && !strings.Contains(gotErr.Error(), tt.wantErr)) {
				t.Errorf("Execute() got unexpected error %v want %s", gotErr, tt.wantErr)
			}
			if tt.want != got {
				t.Errorf("Execute() got unexpected response want %v got %v", tt.want, got)
			}
		})
	}
}

type fakeTracerouteClient struct {
	spb.System_TracerouteClient
	resp []*spb.TracerouteResponse
	err  error
}

func (tc *fakeTracerouteClient) Recv() (*spb.TracerouteResponse, error) {
	if len(tc.resp) == 0 && tc.err == nil {
		return nil, io.EOF
	}
	resp := tc.resp[0]
	tc.resp = tc.resp[1:]
	return resp, tc.err
}

func TestTraceroute(t *testing.T) {
	tests := []struct {
		desc    string
		op      *system.TracerouteOperation
		want    []*spb.TracerouteResponse
		wantErr string
	}{
		{
			desc: "Traceroute with source, destination and L4Protocol",
			op:   system.NewTracerouteOperation().Source("5.6.7.8").Destination("1.2.3.4").L4Protocol(spb.TracerouteRequest_UDP),
			want: []*spb.TracerouteResponse{{DestinationAddress: "1.2.3.4"}},
		},
		{
			desc: "Traceroute with multiple response",
			op:   system.NewTracerouteOperation().Destination("1.2.3.4").Source("5.6.7.8").MaxTTL(2),
			want: []*spb.TracerouteResponse{{DestinationAddress: "1.2.3.4", Hop: 1}, {DestinationAddress: "1.2.3.4", Hop: 2}},
		},
		{
			desc:    "Traceroute returns error",
			op:      system.NewTracerouteOperation().Destination("1.2.3.4").Source("5.6.7.8"),
			wantErr: "Traceroute operation error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var fakeClient internal.Clients
			fakeClient.SystemClient = &fakeSystemClient{TracerouteFn: func(context.Context, *spb.TracerouteRequest, ...grpc.CallOption) (spb.System_TracerouteClient, error) {
				if tt.wantErr != "" {
					return nil, fmt.Errorf(tt.wantErr)
				}
				return &fakeTracerouteClient{resp: tt.want}, nil
			}}

			got, gotErr := tt.op.Execute(context.Background(), &fakeClient)
			if (gotErr == nil) != (tt.wantErr == "") || (gotErr != nil && !strings.Contains(gotErr.Error(), tt.wantErr)) {
				t.Errorf("Execute() got unexpected error %v want %s", gotErr, tt.wantErr)
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Execute() got unexpected response diff (-want +got): %s", diff)
			}
		})
	}
}
