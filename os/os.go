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

// Package os provides gNOI os operations.
package os

import (
	"context"
	"fmt"
	"io"

	log "github.com/golang/glog"
	ospb "github.com/openconfig/gnoi/os"
	"github.com/openconfig/gnoigo/internal"
)

// InstallOperation represents the parameters of a Install operation.
type InstallOperation struct {
	req    *ospb.TransferRequest
	reader io.Reader
}

// Version identifies the OS version.
func (i *InstallOperation) Version(version string) *InstallOperation {
	i.req.Version = version
	return i
}

// Standby specifies if supervisor is on standby.
func (i *InstallOperation) Standby(standby bool) *InstallOperation {
	i.req.StandbySupervisor = standby
	return i
}

// Reader specifies the package reader for the OS file.
func (i *InstallOperation) Reader(reader io.Reader) *InstallOperation {
	i.reader = reader
	return i
}

// NewInstallOperation creates an empty InstallOperation.
func NewInstallOperation() *InstallOperation {
	return &InstallOperation{req: &ospb.TransferRequest{}}
}

// awaitPackageInstall receives messages from the client until either
// (a) the package is installed and validated, in which case it returns the InstallResponse message
// (b) the device does not have the package, in which case it returns a nil response
// (c) an error occurs, in which case it returns the error
// (d) context is cancelled, in which case it returns the context error
func awaitPackageInstall(ctx context.Context, ic ospb.OS_InstallClient) (*ospb.InstallResponse, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		cresp, err := ic.Recv()
		if err != nil {
			return nil, err
		}
		switch v := cresp.GetResponse().(type) {
		case *ospb.InstallResponse_Validated:
			return cresp, nil
		case *ospb.InstallResponse_TransferReady:
			return nil, nil
		case *ospb.InstallResponse_InstallError:
			errName := ospb.InstallError_Type_name[int32(v.InstallError.Type)]
			return nil, fmt.Errorf("installation error %q: %s", errName, v.InstallError.GetDetail())
		case *ospb.InstallResponse_TransferProgress:
			log.Infof("installation progress: %v bytes received from client", v.TransferProgress.GetBytesReceived())
		case *ospb.InstallResponse_SyncProgress:
			log.Infof("installation progress: %v%% synced from supervisor", v.SyncProgress.GetPercentageTransferred())
		default:
			return nil, fmt.Errorf("unexpected client install response: %v (%T)", v, v)
		}
	}
}

func transferContent(ctx context.Context, ic ospb.OS_InstallClient, reader io.Reader) error {
	// The gNOI SetPackage operation sets the maximum chunk size at 64K,
	// so assuming the install operation allows for up to the same size.
	buf := make([]byte, 64*1024)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		n, err := reader.Read(buf)
		if n > 0 {
			if err := ic.Send(
				&ospb.InstallRequest{
					Request: &ospb.InstallRequest_TransferContent{
						TransferContent: buf[0:n],
					},
				},
			); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return ic.Send(
		&ospb.InstallRequest{
			Request: &ospb.InstallRequest_TransferEnd{
				TransferEnd: &ospb.TransferEnd{},
			},
		},
	)
}

// Execute performs the Install operation.
func (i *InstallOperation) Execute(ctx context.Context, c *internal.Clients) (*ospb.InstallResponse, error) {
	ic, icErr := c.OS().Install(ctx)
	if icErr != nil {
		return nil, icErr
	}

	installReq := &ospb.InstallRequest{
		Request: &ospb.InstallRequest_TransferRequest{
			TransferRequest: i.req,
		},
	}

	if err := ic.Send(installReq); err != nil {
		return nil, err
	}

	installResp, err := awaitPackageInstall(ctx, ic)
	if err != nil {
		return nil, err
	}
	if installResp != nil {
		return installResp, nil
	}
	if i.reader == nil {
		return nil, fmt.Errorf("no reader specified for install operation")
	}
	awaitChan := make(chan error)
	go func() {
		installResp, err = awaitPackageInstall(ctx, ic)
		awaitChan <- err
	}()
	if err := transferContent(ctx, ic, i.reader); err != nil {
		return nil, err
	}
	if err := <-awaitChan; err != nil {
		return nil, err
	}
	if gotVersion := installResp.GetValidated().GetVersion(); gotVersion != i.req.Version {
		return nil, fmt.Errorf("installed version %q does not match requested version %q", gotVersion, i.req.Version)
	}
	return installResp, nil
}
