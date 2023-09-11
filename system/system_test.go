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
	PingFn func(context.Context, *spb.PingRequest, ...grpc.CallOption) (spb.System_PingClient, error)
}

func (fg *fakeSystemClient) System() spb.SystemClient {
	return fg
}

func (fg *fakeSystemClient) Ping(ctx context.Context, in *spb.PingRequest, opts ...grpc.CallOption) (spb.System_PingClient, error) {
	return fg.PingFn(ctx, in, opts...)
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
			fakeClient.SystemClient = &fakeSystemClient{PingFn: func(_ context.Context, req *spb.PingRequest, _ ...grpc.CallOption) (spb.System_PingClient, error) {
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
