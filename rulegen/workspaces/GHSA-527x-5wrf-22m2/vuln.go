package main

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/miekg/dns"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// ServergRPC represents an instance of a DNS-over-gRPC server.
type ServergRPC struct {
	*Server
	*pb.UnimplementedDnsServiceServer
	grpcServer *grpc.Server
	listenAddr net.Addr
	tlsConfig  *tls.Config
}

// NewServergRPC returns a new CoreDNS GRPC server and compiles all plugin in to it.
		tlsConfig.NextProtos = []string{"h2"}
	}

	return &ServergRPC{Server: s, tlsConfig: tlsConfig}, nil
}

// Compile-time check to ensure ServergRPC implements the caddy.GracefulServer interface
	s.listenAddr = l.Addr()
	s.m.Unlock()

	if s.Tracer() != nil {
		onlyIfParent := func(parentSpanCtx opentracing.SpanContext, method string, req, resp any) bool {
			return parentSpanCtx != nil
		}
		intercept := otgrpc.OpenTracingServerInterceptor(s.Tracer(), otgrpc.IncludingSpans(onlyIfParent))
		s.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(intercept))
	} else {
		s.grpcServer = grpc.NewServer()
	}

	pb.RegisterDnsServiceServer(s.grpcServer, s)

	if s.tlsConfig != nil {
		l = tls.NewListener(l, s.tlsConfig)
	}
	return s.grpcServer.Serve(l)
}

// any normal server. We use a custom responseWriter to pick up the bytes we need to write
// back to the client as a protobuf.
func (s *ServergRPC) Query(ctx context.Context, in *pb.DnsPacket) (*pb.DnsPacket, error) {
	msg := new(dns.Msg)
	err := msg.Unpack(in.GetMsg())
	if err != nil {
