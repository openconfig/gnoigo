// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package factoryreset provides gNOI factoryreset operations.
package factoryreset

import (
	"context"

	frpb "github.com/openconfig/gnoi/factory_reset"

	"github.com/openconfig/gnoigo/internal"
)

// StartOperation represents the parameters of a FactoryReset Start operation.
type StartOperation struct {
	req *frpb.StartRequest
}

// NewStartOperation creates an empty StartOperation.
func NewStartOperation() *StartOperation {
	return &StartOperation{req: &frpb.StartRequest{}}
}

// FactoryOS instructs the target to rollback the OS to the
// same version as it shipped from factory.
func (s *StartOperation) FactoryOS(fos bool) *StartOperation {
	s.req.FactoryOs = fos
	return s
}

// ZeroFill instructs the target to zero fill persistent storage state data.
func (s *StartOperation) ZeroFill(zf bool) *StartOperation {
	s.req.ZeroFill = zf
	return s
}

// RetainCerts instructs the target to retain certificates.
func (s *StartOperation) RetainCerts(rc bool) *StartOperation {
	s.req.RetainCerts = rc
	return s
}

// Execute performs the Start operation.
func (s *StartOperation) Execute(ctx context.Context, c *internal.Clients) (*frpb.StartResponse, error) {
	return c.FactoryReset().Start(ctx, s.req)
}
