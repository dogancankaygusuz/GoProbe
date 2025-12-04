package proto

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion7

type ProbeServiceClient interface {
	// Master, Worker'a Check, Worker sonucunu döner.
	CheckUrl(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
}

type probeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProbeServiceClient(cc grpc.ClientConnInterface) ProbeServiceClient {
	return &probeServiceClient{cc}
}

func (c *probeServiceClient) CheckUrl(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	out := new(CheckResponse)
	err := c.cc.Invoke(ctx, "/probe.ProbeService/CheckUrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type ProbeServiceServer interface {
	CheckUrl(context.Context, *CheckRequest) (*CheckResponse, error)
	mustEmbedUnimplementedProbeServiceServer()
}

type UnimplementedProbeServiceServer struct {
}

func (UnimplementedProbeServiceServer) CheckUrl(context.Context, *CheckRequest) (*CheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckUrl not implemented")
}
func (UnimplementedProbeServiceServer) mustEmbedUnimplementedProbeServiceServer() {}

type UnsafeProbeServiceServer interface {
	mustEmbedUnimplementedProbeServiceServer()
}

func RegisterProbeServiceServer(s grpc.ServiceRegistrar, srv ProbeServiceServer) {
	s.RegisterService(&ProbeService_ServiceDesc, srv)
}

func _ProbeService_CheckUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProbeServiceServer).CheckUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/probe.ProbeService/CheckUrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProbeServiceServer).CheckUrl(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var ProbeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "probe.ProbeService",
	HandlerType: (*ProbeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckUrl",
			Handler:    _ProbeService_CheckUrl_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/grpc/proto/probe.proto",
}
