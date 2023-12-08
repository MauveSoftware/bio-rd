package ris

import (
	"context"
	"fmt"
	"time"

	"github.com/bio-routing/bio-rd/cmd/ris/api"
	"github.com/bio-routing/bio-rd/net"
	"github.com/bio-routing/bio-rd/routingtable/locRIB"
	"github.com/bio-routing/bio-rd/routingtable/vrf"
	"github.com/bio-routing/bio-rd/util/servicewrapper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	routeapi "github.com/bio-routing/bio-rd/route/api"
)

type server struct {
	api.UnimplementedRoutingInformationServiceServer
	risVRF *vrf.VRF
}

func Start(port, httpPort, grpcKeepaliveMinTime uint, v *vrf.VRF) error {
	unaryInterceptors := []grpc.UnaryServerInterceptor{}
	streamInterceptors := []grpc.StreamServerInterceptor{}
	srv, err := servicewrapper.New(
		uint16(port),
		servicewrapper.HTTP(uint16(httpPort)),
		unaryInterceptors,
		streamInterceptors,
		keepalive.EnforcementPolicy{
			MinTime:             time.Duration(grpcKeepaliveMinTime) * time.Second,
			PermitWithoutStream: true,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := &server{
		risVRF: v,
	}
	api.RegisterRoutingInformationServiceServer(srv.GRPC(), s)

	return nil
}

func (srv *server) LPM(ctx context.Context, in *api.LPMRequest) (*api.LPMResponse, error) {
	pfx := net.NewPrefixFromProtoPrefix(in.Pfx)
	rib := srv.ribForPrefix(pfx)

	routes := rib.LPM(pfx)
	res := &api.LPMResponse{
		Routes: make([]*routeapi.Route, 0, len(routes)),
	}
	for _, route := range routes {
		res.Routes = append(res.Routes, route.ToProto())
	}

	return res, nil
}

func (srv *server) ribForPrefix(pfx *net.Prefix) *locRIB.LocRIB {
	if pfx.Addr().IsIPv4() {
		return srv.risVRF.IPv4UnicastRIB()
	}

	return srv.risVRF.IPv6UnicastRIB()
}

func (srv *server) Get(ctx context.Context, in *api.GetRequest) (*api.GetResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (srv *server) GetRouters(ctx context.Context, in *api.GetRoutersRequest) (*api.GetRoutersResponse, error) {
	return &api.GetRoutersResponse{
		Routers: []*api.Router{
			{
				Address: "127.0.0.1",
				VrfIds:  []uint64{0},
			},
		},
	}, nil
}

func (srv *server) GetLonger(ctx context.Context, in *api.GetLongerRequest) (*api.GetLongerResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (srv *server) ObserveRIB(in *api.ObserveRIBRequest, stream api.RoutingInformationService_ObserveRIBServer) error {

	return fmt.Errorf("not implemented")
}

func (srv *server) DumpRIB(in *api.DumpRIBRequest, stream api.RoutingInformationService_DumpRIBServer) error {
	return fmt.Errorf("not implemented")
}
