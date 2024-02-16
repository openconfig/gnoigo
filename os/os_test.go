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

package os_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	ospb "github.com/openconfig/gnoi/os"
	"github.com/openconfig/gnoigo/internal"
	gos "github.com/openconfig/gnoigo/os"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

type fakeOSClient struct {
	ospb.OSClient
	InstallFn  func(context.Context, ...grpc.CallOption) (ospb.OS_InstallClient, error)
	ActivateFn func(context.Context, *ospb.ActivateRequest, ...grpc.CallOption) (*ospb.ActivateResponse, error)
	VerifyFn   func(context.Context, *ospb.VerifyRequest, ...grpc.CallOption) (*ospb.VerifyResponse, error)
}

func (fg *fakeOSClient) OS() ospb.OSClient {
	return fg
}

func (fg *fakeOSClient) Install(ctx context.Context, opts ...grpc.CallOption) (ospb.OS_InstallClient, error) {
	return fg.InstallFn(ctx, opts...)
}

func (fg *fakeOSClient) Activate(ctx context.Context, in *ospb.ActivateRequest, opts ...grpc.CallOption) (*ospb.ActivateResponse, error) {
	return fg.ActivateFn(ctx, in, opts...)
}

func (fg *fakeOSClient) Verify(ctx context.Context, in *ospb.VerifyRequest, opts ...grpc.CallOption) (*ospb.VerifyResponse, error) {
	return fg.VerifyFn(ctx, in, opts...)
}

type fakeInstallClient struct {
	ospb.OS_InstallClient
	gotSent  []*ospb.InstallRequest
	stubRecv []*ospb.InstallResponse
}

func (ic *fakeInstallClient) Send(req *ospb.InstallRequest) error {
	ic.gotSent = append(ic.gotSent, req)
	return nil
}

func (ic *fakeInstallClient) Recv() (*ospb.InstallResponse, error) {
	if len(ic.stubRecv) == 0 {
		return nil, errors.New("no more stub responses")
	}
	resp := ic.stubRecv[0]
	ic.stubRecv[0] = nil
	ic.stubRecv = ic.stubRecv[1:]
	return resp, nil
}

func (*fakeInstallClient) CloseSend() error {
	return nil
}

func TestInstall(t *testing.T) {
	const version = "1.2.3"

	// Make a temp file to test specifying a file by file path.
	file, err := os.CreateTemp("", "package")
	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if err := os.WriteFile(file.Name(), []byte{0}, os.ModePerm); err != nil {
		t.Fatalf("error writing temp file: %v", err)
	}

	tests := []struct {
		desc          string
		op            *gos.InstallOperation
		resps         []*ospb.InstallResponse
		want          *ospb.InstallResponse
		installErr    string
		wantErr       string
		cancelContext bool
	}{
		{
			desc: "install with version",
			op:   gos.NewInstallOperation().Version(version),
			resps: []*ospb.InstallResponse{
				{Response: &ospb.InstallResponse_Validated{Validated: &ospb.Validated{Version: version}}},
			},
			want: &ospb.InstallResponse{Response: &ospb.InstallResponse_Validated{Validated: &ospb.Validated{Version: version}}},
		},
		{
			desc:       "install returns error",
			op:         gos.NewInstallOperation().Version(version),
			resps:      []*ospb.InstallResponse{},
			installErr: "install error",
			wantErr:    "install error",
		},
		{
			desc: "install with context cancel",
			op:   gos.NewInstallOperation().Version(version),
			resps: []*ospb.InstallResponse{
				{Response: &ospb.InstallResponse_Validated{Validated: &ospb.Validated{Version: version}}},
			},
			wantErr:       "context",
			cancelContext: true,
		},
		{
			desc: "install without ioreader returns error",
			op:   gos.NewInstallOperation().Version(version),
			resps: []*ospb.InstallResponse{
				{Response: &ospb.InstallResponse_TransferReady{TransferReady: &ospb.TransferReady{}}},
				{Response: &ospb.InstallResponse_Validated{Validated: &ospb.Validated{Version: version}}},
			},
			wantErr: "reader",
		},
		{
			desc: "install with ioreader",
			op:   gos.NewInstallOperation().Version(version).Reader(bytes.NewReader([]byte{0})),
			resps: []*ospb.InstallResponse{
				{Response: &ospb.InstallResponse_TransferReady{TransferReady: &ospb.TransferReady{}}},
				{Response: &ospb.InstallResponse_TransferProgress{TransferProgress: &ospb.TransferProgress{}}},
				{Response: &ospb.InstallResponse_Validated{Validated: &ospb.Validated{Version: version}}},
			},
			want: &ospb.InstallResponse{Response: &ospb.InstallResponse_Validated{Validated: &ospb.Validated{Version: version}}},
		},
		{
			desc: "install with mismatch version error",
			op:   gos.NewInstallOperation().Version(version).Reader(bytes.NewReader([]byte{0})),
			resps: []*ospb.InstallResponse{
				{Response: &ospb.InstallResponse_TransferReady{TransferReady: &ospb.TransferReady{}}},
				{Response: &ospb.InstallResponse_Validated{Validated: &ospb.Validated{Version: version + "new"}}},
			},
			wantErr: "version",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var fakeClient internal.Clients
			fakeClient.OSClient = &fakeOSClient{InstallFn: func(context.Context, ...grpc.CallOption) (ospb.OS_InstallClient, error) {
				if tt.installErr != "" {
					return nil, fmt.Errorf(tt.installErr)
				}
				return &fakeInstallClient{stubRecv: tt.resps}, nil
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
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("Execute() got unexpected response diff (-want +got): %s", diff)
			}
		})
	}
}

func TestActivate(t *testing.T) {
	tests := []struct {
		desc    string
		op      *gos.ActivateOperation
		want    *ospb.ActivateResponse
		wantErr string
	}{
		{
			desc: "Test Activate",
			op:   gos.NewActivateOperation(),
			want: &ospb.ActivateResponse{Response: &ospb.ActivateResponse_ActivateOk{}},
		},
		{
			desc:    "Activate returns error",
			op:      gos.NewActivateOperation(),
			wantErr: "Activate operation error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var fakeClient internal.Clients
			fakeClient.OSClient = &fakeOSClient{ActivateFn: func(context.Context, *ospb.ActivateRequest, ...grpc.CallOption) (*ospb.ActivateResponse, error) {
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

func TestVerify(t *testing.T) {
	tests := []struct {
		desc    string
		op      *gos.VerifyOperation
		want    *ospb.VerifyResponse
		wantErr string
	}{
		{
			desc: "Test Verify",
			op:   gos.NewVerifyOperation(),
			want: &ospb.VerifyResponse{Version: "1.2.3"},
		},
		{
			desc:    "Verify returns error",
			op:      gos.NewVerifyOperation(),
			wantErr: "Verify operation error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var fakeClient internal.Clients
			fakeClient.OSClient = &fakeOSClient{VerifyFn: func(context.Context, *ospb.VerifyRequest, ...grpc.CallOption) (*ospb.VerifyResponse, error) {
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
