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

// Package internal provides gnoigo structs that are not part of the public API.
package internal

import (
	bpb "github.com/openconfig/gnoi/bgp"
	cmpb "github.com/openconfig/gnoi/cert"
	cpb "github.com/openconfig/gnoi/containerz"
	dpb "github.com/openconfig/gnoi/diag"
	frpb "github.com/openconfig/gnoi/factory_reset"
	fpb "github.com/openconfig/gnoi/file"
	hpb "github.com/openconfig/gnoi/healthz"
	lpb "github.com/openconfig/gnoi/layer2"
	mpb "github.com/openconfig/gnoi/mpls"
	ospb "github.com/openconfig/gnoi/os"
	otpb "github.com/openconfig/gnoi/otdr"
	plqpb "github.com/openconfig/gnoi/packet_link_qualification"
	spb "github.com/openconfig/gnoi/system"
	wrpb "github.com/openconfig/gnoi/wavelength_router"
)

// Clients contains all gNOI clients.
type Clients struct {
	BGPClient              bpb.BGPClient
	CertMgmtClient         cmpb.CertificateManagementClient
	ContainerzClient       cpb.ContainerzClient
	DiagClient             dpb.DiagClient
	FactoryResetClient     frpb.FactoryResetClient
	FileClient             fpb.FileClient
	HealthzClient          hpb.HealthzClient
	Layer2Client           lpb.Layer2Client
	LinkQualClient         plqpb.LinkQualificationClient
	MPLSClient             mpb.MPLSClient
	OSClient               ospb.OSClient
	OTDRClient             otpb.OTDRClient
	SystemClient           spb.SystemClient
	WavelengthRouterClient wrpb.WavelengthRouterClient
}

// BGP returns the client for gNOI BGP service.
func (c *Clients) BGP() bpb.BGPClient {
	return c.BGPClient
}

// CertificateManagement returns the client for gNOI Certificate Management service.
func (c *Clients) CertificateManagement() cmpb.CertificateManagementClient {
	return c.CertMgmtClient
}

// Containerz returns the client for gNOI Containerz service.
func (c *Clients) Containerz() cpb.ContainerzClient {
	return c.ContainerzClient
}

// Diag returns the client for gNOI Diag service.
func (c *Clients) Diag() dpb.DiagClient {
	return c.DiagClient
}

// FactoryReset returns the client for gNOI FactoryReset service.
func (c *Clients) FactoryReset() frpb.FactoryResetClient {
	return c.FactoryResetClient
}

// File returns the client for gNOI File service.
func (c *Clients) File() fpb.FileClient {
	return c.FileClient
}

// Healthz returns the client for gNOI Healthz service.
func (c *Clients) Healthz() hpb.HealthzClient {
	return c.HealthzClient
}

// Layer2 returns the client for gNOI Layer2 service.
func (c *Clients) Layer2() lpb.Layer2Client {
	return c.Layer2Client
}

// LinkQualification returns the client for gNOI LinkQualification service.
func (c *Clients) LinkQualification() plqpb.LinkQualificationClient {
	return c.LinkQualClient
}

// MPLS returns the client for gNOI MPLS service.
func (c *Clients) MPLS() mpb.MPLSClient {
	return c.MPLSClient
}

// OS returns the client for gNOI OS service.
func (c *Clients) OS() ospb.OSClient {
	return c.OSClient
}

// OTDR returns the client for gNOI OTDR service.
func (c *Clients) OTDR() otpb.OTDRClient {
	return c.OTDRClient
}

// System returns the client for gNOI System service.
func (c *Clients) System() spb.SystemClient {
	return c.SystemClient
}

// WavelengthRouter returns the client for gNOI WavelengthRouter service.
func (c *Clients) WavelengthRouter() wrpb.WavelengthRouterClient {
	return c.WavelengthRouterClient
}
