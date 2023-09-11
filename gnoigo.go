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
	"github.com/openconfig/gnoigo/internal"
)

// NewClients constructs all the gNOI clients.
func NewClients(conn *grpc.ClientConn) Clients {
	return &internal.Clients{
		BGPClient:              bpb.NewBGPClient(conn),
		CertMgmtClient:         cmpb.NewCertificateManagementClient(conn),
		DiagClient:             dpb.NewDiagClient(conn),
		FactoryResetClient:     frpb.NewFactoryResetClient(conn),
		FileClient:             fpb.NewFileClient(conn),
		HealthzClient:          hpb.NewHealthzClient(conn),
		Layer2Client:           lpb.NewLayer2Client(conn),
		LinkQualClient:         plqpb.NewLinkQualificationClient(conn),
		MPLSClient:             mpb.NewMPLSClient(conn),
		OSClient:               ospb.NewOSClient(conn),
		OTDRClient:             otpb.NewOTDRClient(conn),
		SystemClient:           spb.NewSystemClient(conn),
		WavelengthRouterClient: wrpb.NewWavelengthRouterClient(conn),
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

// Operation represents any gNOI operation.
type Operation[T any] interface {
	Execute(context.Context, *internal.Clients) (T, error)
}

// Execute performs an operation and returns one or more response protos.
// For example, a PingOperation returns a slice of PingResponse messages.
func Execute[T any](ctx context.Context, c Clients, op Operation[T]) (T, error) {
	return op.Execute(ctx, c)
}
