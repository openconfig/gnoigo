// Copyright 2024 Google LLC
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

package file_test

import (
	"context"
	"crypto/sha256"
	"os"
	"path"
	"testing"

	fpb "github.com/openconfig/gnoi/file"
	tpb "github.com/openconfig/gnoi/types"

	"github.com/google/go-cmp/cmp"
	"github.com/openconfig/gnoigo/file"
	"github.com/openconfig/gnoigo/internal"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

type fakeFileClient struct {
	fpb.FileClient
	PutFn func(ctx context.Context, opts ...grpc.CallOption) (fpb.File_PutClient, error)
}

func (f *fakeFileClient) File() fpb.FileClient {
	return f
}

func (f *fakeFileClient) Put(ctx context.Context, opts ...grpc.CallOption) (fpb.File_PutClient, error) {
	return f.PutFn(ctx, opts...)
}

type fakePutClient struct {
	fpb.File_PutClient
	gotReq []*fpb.PutRequest
}

func (fc *fakePutClient) Send(req *fpb.PutRequest) error {
	fc.gotReq = append(fc.gotReq, req)
	return nil
}

func (fc *fakePutClient) Recv() (*fpb.PutResponse, error) {
	return &fpb.PutResponse{}, nil
}

func (fc *fakePutClient) CloseAndRecv() (*fpb.PutResponse, error) {
	return &fpb.PutResponse{}, nil
}

func (*fakePutClient) CloseSend() error {
	return nil
}

func generateFile(t *testing.T, data string) string {
	// Create a temporary file
	fileName := path.Join(t.TempDir(), "data")

	// Write some text to the file
	if err := os.WriteFile(fileName, []byte(data), 0644); err != nil {
		t.Fatalf("unable to write temp file contents: %v", err)
	}

	return fileName
}

func TestPut(t *testing.T) {
	const data = "some really important data"
	fileName := generateFile(t, data)
	hash := sha256.New()
	_, err := hash.Write([]byte(`some really important data`))
	if err != nil {
		t.Fatalf("Unable to hash string: %v", err)
	}
	tests := []struct {
		desc    string
		op      *file.PutOperation
		wantReq []*fpb.PutRequest
		wantErr bool
	}{
		{
			desc:    "put-with-no-file",
			op:      file.NewPutOperation(),
			wantErr: true,
		},
		{
			desc: "put-with-file",
			op:   file.NewPutOperation().SourceFile(fileName),
			wantReq: []*fpb.PutRequest{
				{
					Request: &fpb.PutRequest_Open{
						Open: &fpb.PutRequest_Details{},
					},
				},
				{
					Request: &fpb.PutRequest_Contents{
						Contents: []byte(data),
					},
				},
				{
					Request: &fpb.PutRequest_Hash{
						Hash: &tpb.HashType{
							Method: tpb.HashType_SHA256,
							Hash:   hash.Sum(nil),
						},
					},
				},
			},
		},
		{
			desc: "put-with-all-details",
			op:   file.NewPutOperation().SourceFile(fileName).RemoteFile("/tmp/here").Perms(644),
			wantReq: []*fpb.PutRequest{
				{
					Request: &fpb.PutRequest_Open{
						Open: &fpb.PutRequest_Details{
							RemoteFile:  "/tmp/here",
							Permissions: 644,
						},
					},
				},
				{
					Request: &fpb.PutRequest_Contents{
						Contents: []byte(data),
					},
				},
				{
					Request: &fpb.PutRequest_Hash{
						Hash: &tpb.HashType{
							Method: tpb.HashType_SHA256,
							Hash:   hash.Sum(nil),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			fpc := &fakePutClient{}
			var fakeClient internal.Clients
			fakeClient.FileClient = &fakeFileClient{
				PutFn: func(ctx context.Context, opts ...grpc.CallOption) (fpb.File_PutClient, error) {
					return fpc, nil
				},
			}

			_, err := tt.op.Execute(context.Background(), &fakeClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() got unexpected error %v", err)
			}

			if diff := cmp.Diff(fpc.gotReq, tt.wantReq, protocmp.Transform()); diff != "" {
				t.Errorf("Execute returned diff (-got, +want):\n%s", diff)
			}
		})
	}
}
