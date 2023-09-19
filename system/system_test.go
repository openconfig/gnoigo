package system_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/google/go-cmp/cmp"
	spb "github.com/openconfig/gnoi/system"
	"github.com/openconfig/gnoigo/internal"
	"github.com/openconfig/gnoigo/system"
)

type fakeSystemClient struct {
	spb.SystemClient
	KillProcessFn func(context.Context, *spb.KillProcessRequest, ...grpc.CallOption) (*spb.KillProcessResponse, error)
	PingFn       func(context.Context, *spb.PingRequest, ...grpc.CallOption) (spb.System_PingClient, error)
	TimeFn       func(context.Context, *spb.TimeRequest, ...grpc.CallOption) (*spb.TimeResponse, error)
	TracerouteFn func(context.Context, *spb.TracerouteRequest, ...grpc.CallOption) (spb.System_TracerouteClient, error)
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

func (fg *fakeSystemClient) Time(ctx context.Context, in *spb.TimeRequest, opts ...grpc.CallOption) (*spb.TimeResponse, error) {
	return fg.TimeFn(ctx, in, opts...)
}

func (fg *fakeSystemClient) Traceroute(ctx context.Context, in *spb.TracerouteRequest, opts ...grpc.CallOption) (spb.System_TracerouteClient, error) {
	return fg.TracerouteFn(ctx, in, opts...)
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
