package system_test

import (
	"context"
	"io"
	"testing"

	"google.golang.org/grpc"

	spb "github.com/openconfig/gnoi/system"
	"github.com/openconfig/gnoigo"
	"github.com/openconfig/gnoigo/system"
)

type fakeClients struct {
	gnoigo.Clients
	SystemFn func() spb.SystemClient
}

func (f *fakeClients) Client() gnoigo.Clients {
	return f
}

func (f *fakeClients) System() spb.SystemClient {
	return f.SystemFn()
}

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
	resp *spb.PingResponse
	err  error
}

func (pc *fakePingClient) Recv() (*spb.PingResponse, error) {
	if pc.resp == nil && pc.err == nil {
		return nil, io.EOF
	}
	resp := pc.resp
	pc.resp = nil
	return resp, pc.err
}

func TestPing(t *testing.T) {
	tests := []struct {
		desc, dst, src    string
		count, packetSize int32
	}{
		{desc: "ping with source", dst: "1.2.3.4", src: "5.6.7.8"},
		{desc: "ping with source, count and packetsize", dst: "1.2.3.4", src: "5.6.7.8", count: 7, packetSize: 1000},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var fakeClient fakeClients
			var got string
			fakeClient.SystemFn = func() spb.SystemClient {
				return &fakeSystemClient{
					PingFn: func(_ context.Context, req *spb.PingRequest, _ ...grpc.CallOption) (spb.System_PingClient, error) {
						got = req.GetDestination()
						return &fakePingClient{resp: &spb.PingResponse{Source: tt.src, Sent: tt.count, Received: tt.count}}, nil
					},
				}
			}

			want := tt.dst

			pingOp := system.NewPingOperation().
				Destination(tt.dst).
				Source(tt.src).
				Count(tt.count)

			responses, err := gnoigo.Execute(context.Background(), fakeClient.Client(), pingOp)

			if got != want {
				t.Errorf("Operate(t) got %s, want %s", got, want)
			}

			if err != nil {
				t.Errorf("Error on ping %v, want nil", err)
			} else {
				if len(responses) != 1 {
					t.Errorf("Got %d responses, want 1", len(responses))
				}
				if responses[0].Source != tt.src {
					t.Errorf("Response.Source error got %s, want %s", responses[0].Source, tt.src)
				}
			}
		})
	}

}
