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
package debug

import (
	"context"

	dbgpb "github.com/openconfig/gnoi/debug"
	"github.com/openconfig/gnoigo/internal"
)

// DebugOperation represents the parameters of a Debug operation.
// Debug will execute the set of commands provided in the request.
// The command will be executed in the mode provided.
// All command modes must support an exit code on completion of the
// command. (e.g. Cli modes must exit after sending a command)
// Errors:
//
//	InvalidArgument: for unspecified mode
type DebugOperation struct {
	req *dbgpb.DebugRequest
}

// Mode the commands will be executed in.
//
//	enum Mode {
//	    MODE_UNSPECIFIED = 0;
//	    MODE_SHELL = 1;
//	    MODE_CLI = 2;
//	  }
func (a *DebugOperation) Mode(mode dbgpb.DebugRequest_Mode) *DebugOperation {
	a.req.Mode = mode
	return a
}

// Raw bytes for the command to be executed.
func (a *DebugOperation) Command(command []byte) *DebugOperation {
	a.req.Command = command
	return a
}

// Truncate the amount of data returned for the command.
func (a *DebugOperation) ByteLimit(byteLimit int64) *DebugOperation {
	a.req.ByteLimit = byteLimit
	return a
}

// Timeout in nanoseconds.
func (a *DebugOperation) Timeout(timeout int64) *DebugOperation {
	a.req.Timeout = timeout
	return a
}

// Role account to use for the command.
func (a *DebugOperation) RoleAccount(roleAccount string) *DebugOperation {
	a.req.RoleAccount = roleAccount
	return a
}

// NewDebugOperation creates an empty DebugOperation.
func NewDebugOperation() *DebugOperation {
	return &DebugOperation{req: &dbgpb.DebugRequest{}}
}

// Execute performs the Debug operation.
func (a *DebugOperation) Execute(ctx context.Context, c *internal.Clients) (dbgpb.Debug_DebugClient, error) {
	return c.DEBUG().Debug(ctx, a.req)
}
