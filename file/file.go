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

// Package file provides gNOI file operations.
package file

import (
	"context"
	"crypto/sha256"
	"io"
	"os"

	fpb "github.com/openconfig/gnoi/file"
	tpb "github.com/openconfig/gnoi/types"
	"github.com/openconfig/gnoigo/internal"
)

const (
	// chunkSize is the maximal size of a file chunk as defined by the spec.
	chunkSize = 64000
)

// PutOperation represents the parameters of a Put operation.
type PutOperation struct {
	sourceFile string
	req        *fpb.PutRequest
}

// NewPutOperation creates an empty PutOperation.
func NewPutOperation() *PutOperation {
	return &PutOperation{
		req: &fpb.PutRequest{
			Request: &fpb.PutRequest_Open{
				Open: &fpb.PutRequest_Details{},
			},
		},
	}
}

// Perms specifies the permissions to apply to the copied file.
func (p *PutOperation) Perms(perms uint32) *PutOperation {
	p.req.GetOpen().Permissions = perms
	return p
}

// RemoteFile specifies the name of the file on the target.
func (p *PutOperation) RemoteFile(file string) *PutOperation {
	p.req.GetOpen().RemoteFile = file
	return p
}

// SourceFile represents the source file to copy.
func (p *PutOperation) SourceFile(file string) *PutOperation {
	p.sourceFile = file
	return p
}

// Execute executes the Put operation.
func (p *PutOperation) Execute(ctx context.Context, c *internal.Clients) (*fpb.PutResponse, error) {
	pclient, err := c.File().Put(ctx)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(p.sourceFile)
	if err != nil {
		return nil, err
	}

	if err := pclient.Send(p.req); err != nil {
		return nil, err
	}

	hasher := sha256.New()
	buf := make([]byte, chunkSize)
	for i, done := 0, false; !done; i++ {
		n, err := f.ReadAt(buf, int64(i*chunkSize))
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			done = true
		}
		content := buf[:n]

		if _, err = hasher.Write(content); err != nil {
			return nil, err
		}

		req := &fpb.PutRequest{
			Request: &fpb.PutRequest_Contents{
				Contents: content,
			},
		}
		if err := pclient.Send(req); err != nil {
			return nil, err
		}
	}

	req := &fpb.PutRequest{
		Request: &fpb.PutRequest_Hash{
			Hash: &tpb.HashType{
				Hash:   hasher.Sum(nil),
				Method: tpb.HashType_SHA256,
			},
		},
	}
	if err := pclient.Send(req); err != nil {
		return nil, err
	}

	return pclient.CloseAndRecv()
}
