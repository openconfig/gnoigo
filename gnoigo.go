package gnoigo

import (
	"context"

	"google.golang.org/grpc"

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

// NewClients constructs all the gNOI clients.
func NewClients(conn *grpc.ClientConn) Clients {
	return &clients{
		bgp:               bpb.NewBGPClient(conn),
		certManagment:     cmpb.NewCertificateManagementClient(conn),
		diag:              dpb.NewDiagClient(conn),
		file:              fpb.NewFileClient(conn),
		factoryReset:      frpb.NewFactoryResetClient(conn),
		healthz:           hpb.NewHealthzClient(conn),
		l2:                lpb.NewLayer2Client(conn),
		linkQualification: plqpb.NewLinkQualificationClient(conn),
		mpls:              mpb.NewMPLSClient(conn),
		os:                ospb.NewOSClient(conn),
		otdr:              otpb.NewOTDRClient(conn),
		system:            spb.NewSystemClient(conn),
		wavelengthRouter:  wrpb.NewWavelengthRouterClient(conn),
	}
}

// Clients is a set of gNOI clients.
type Clients interface {
	BGP() bpb.BGPClient
	CertificateManagement() cmpb.CertificateManagementClient
	Diag() dpb.DiagClient
	FactoryReset() frpb.FactoryResetClient
	File() fpb.FileClient
	Healthz() hpb.HealthzClient
	Layer2() lpb.Layer2Client
	LinkQualification() plqpb.LinkQualificationClient
	MPLS() mpb.MPLSClient
	OS() ospb.OSClient
	OTDR() otpb.OTDRClient
	System() spb.SystemClient
	WavelengthRouter() wrpb.WavelengthRouterClient
}

type clients struct {
	bgp               bpb.BGPClient
	certManagment     cmpb.CertificateManagementClient
	diag              dpb.DiagClient
	factoryReset      frpb.FactoryResetClient
	file              fpb.FileClient
	healthz           hpb.HealthzClient
	l2                lpb.Layer2Client
	linkQualification plqpb.LinkQualificationClient
	mpls              mpb.MPLSClient
	os                ospb.OSClient
	otdr              otpb.OTDRClient
	system            spb.SystemClient
	wavelengthRouter  wrpb.WavelengthRouterClient
}

func (c *clients) BGP() bpb.BGPClient {
	return c.bgp
}

func (c *clients) CertificateManagement() cmpb.CertificateManagementClient {
	return c.certManagment
}

func (c *clients) Diag() dpb.DiagClient {
	return c.diag
}

func (c *clients) FactoryReset() frpb.FactoryResetClient {
	return c.factoryReset
}

func (c *clients) File() fpb.FileClient {
	return c.file
}

func (c *clients) Healthz() hpb.HealthzClient {
	return c.healthz
}

func (c *clients) Layer2() lpb.Layer2Client {
	return c.l2
}

func (c *clients) LinkQualification() plqpb.LinkQualificationClient {
	return c.linkQualification
}

func (c *clients) MPLS() mpb.MPLSClient {
	return c.mpls
}

func (c *clients) OS() ospb.OSClient {
	return c.os
}

func (c *clients) OTDR() otpb.OTDRClient {
	return c.otdr
}

func (c *clients) System() spb.SystemClient {
	return c.system
}

func (c *clients) WavelengthRouter() wrpb.WavelengthRouterClient {
	return c.wavelengthRouter
}

// Operation represents any operation in gNOI clients.
type Operation[T any] interface {
	Execute(context.Context, Clients) (T, error)
}

// Execute performs an operation and returns the response proto for the operation.
// For example, executing a PingOperation returns a PingResponse proto message.
func Execute[T any](ctx context.Context, c Clients, op Operation[T]) (T, error) {
	return op.Execute(ctx, c)
}
