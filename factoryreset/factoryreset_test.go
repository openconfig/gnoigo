package factoryreset_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"google.golang.org/grpc"

	frpb "github.com/openconfig/gnoi/factory_reset"
	"github.com/openconfig/gnoigo/factoryreset"
	"github.com/openconfig/gnoigo/internal"
)

type fakeFactoryResetClient struct {
	frpb.FactoryResetClient
	StartFn func(context.Context, *frpb.StartRequest, ...grpc.CallOption) (*frpb.StartResponse, error)
}

func (fg *fakeFactoryResetClient) FactoryReset() frpb.FactoryResetClient {
	return fg
}

func (fg *fakeFactoryResetClient) Start(ctx context.Context, in *frpb.StartRequest, opts ...grpc.CallOption) (*frpb.StartResponse, error) {
	return fg.StartFn(ctx, in, opts...)
}

func TestFactoryResetStart(t *testing.T) {
	tests := []struct {
		desc    string
		op      *factoryreset.StartOperation
		want    *frpb.StartResponse
		wantErr string
	}{
		{
			desc: "Test factoryReset start success",
			op:   factoryreset.NewStartOperation().ZeroFill(true).FactoryOS(true),
			want: &frpb.StartResponse{Response: &frpb.StartResponse_ResetSuccess{ResetSuccess: &frpb.ResetSuccess{}}},
		},
		{
			desc:    "Test factoryReset start error",
			op:      factoryreset.NewStartOperation().ZeroFill(true).FactoryOS(true),
			wantErr: "Factory reset operation error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var fakeClient internal.Clients
			fakeClient.FactoryResetClient = &fakeFactoryResetClient{StartFn: func(context.Context, *frpb.StartRequest, ...grpc.CallOption) (*frpb.StartResponse, error) {
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
