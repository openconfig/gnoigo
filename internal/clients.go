package internal

import (
	bpb "github.com/openconfig/gnoi/bgp"
	cmpb "github.com/openconfig/gnoi/cert"
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
	Bgp            bpb.BGPClient
	CertManagement cmpb.CertificateManagementClient
	Diagnostic     dpb.DiagClient
	Factory        frpb.FactoryResetClient
	FileClient     fpb.FileClient
	Hz             hpb.HealthzClient
	L2             lpb.Layer2Client
	LinkQual       plqpb.LinkQualificationClient
	Mpls           mpb.MPLSClient
	Os             ospb.OSClient
	Otdr           otpb.OTDRClient
	Sys            spb.SystemClient
	WavelengthR    wrpb.WavelengthRouterClient
}

// BGP returns the client for gNOI BGP service.
func (c *Clients) BGP() bpb.BGPClient {
	return c.Bgp
}

// CertificateManagement returns the client for gNOI Certificate Management service.
func (c *Clients) CertificateManagement() cmpb.CertificateManagementClient {
	return c.CertManagement
}

// Diag returns the client for gNOI Diag service.
func (c *Clients) Diag() dpb.DiagClient {
	return c.Diagnostic
}

// FactoryReset returns the client for gNOI FactoryReset service.
func (c *Clients) FactoryReset() frpb.FactoryResetClient {
	return c.Factory
}

// File returns the client for gNOI File service.
func (c *Clients) File() fpb.FileClient {
	return c.FileClient
}

// Healthz returns the client for gNOI Healthz service.
func (c *Clients) Healthz() hpb.HealthzClient {
	return c.Hz
}

// Layer2 returns the client for gNOI Layer2 service.
func (c *Clients) Layer2() lpb.Layer2Client {
	return c.L2
}

// LinkQualification returns the client for gNOI LinkQualification service.
func (c *Clients) LinkQualification() plqpb.LinkQualificationClient {
	return c.LinkQual
}

// MPLS returns the client for gNOI MPLS service.
func (c *Clients) MPLS() mpb.MPLSClient {
	return c.Mpls
}

// OS returns the client for gNOI OS service.
func (c *Clients) OS() ospb.OSClient {
	return c.Os
}

// OTDR returns the client for gNOI OTDR service.
func (c *Clients) OTDR() otpb.OTDRClient {
	return c.Otdr
}

// System returns the client for gNOI System service.
func (c *Clients) System() spb.SystemClient {
	return c.Sys
}

// WavelengthRouter returns the client for gNOI WavelengthRouter service.
func (c *Clients) WavelengthRouter() wrpb.WavelengthRouterClient {
	return c.WavelengthR
}
