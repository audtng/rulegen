package main

package grpc_server

import (
	"strconv"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() {
	caddy.RegisterPlugin("grpc_server", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	err := parseGRPCServer(c)
	if err != nil {
		return plugin.Error("grpc_server", err)
	}
	return nil
}

func parseGRPCServer(c *caddy.Controller) error {
	config := dnsserver.GetConfig(c)

	// Skip the "grpc_server" directive itself
	c.Next()

	// Get any arguments on the "grpc_server" line
	args := c.RemainingArgs()
	if len(args) > 0 {
		return c.ArgErr()
	}

	// Process all nested directives in the block
	for c.NextBlock() {
		switch c.Val() {
		case "max_streams":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return c.ArgErr()
			}
			val, err := strconv.Atoi(args[0])
			if err != nil {
				return c.Errf("invalid max_streams value '%s': %v", args[0], err)
			}
			if val < 0 {
				return c.Errf("max_streams must be a non-negative integer: %d", val)
			}
			if config.MaxGRPCStreams != nil {
				return c.Err("max_streams already defined for this server block")
			}
			config.MaxGRPCStreams = &val
		case "max_connections":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return c.ArgErr()
			}
			val, err := strconv.Atoi(args[0])
			if err != nil {
				return c.Errf("invalid max_connections value '%s': %v", args[0], err)
			}
			if val < 0 {
				return c.Errf("max_connections must be a non-negative integer: %d", val)
			}
			if config.MaxGRPCConnections != nil {
				return c.Err("max_connections already defined for this server block")
			}
			config.MaxGRPCConnections = &val
		default:
			return c.Errf("unknown property '%s'", c.Val())
		}
	}

	return nil
}
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/miekg/dns"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/netutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const (
	// maxDNSMessageBytes is the maximum size of a DNS message on the wire.
	maxDNSMessageBytes = dns.MaxMsgSize

	// maxProtobufPayloadBytes accounts for protobuf overhead.
	// Field tag=1 (1 byte) + length varint for 65535 (3 bytes) = 4 bytes total
	maxProtobufPayloadBytes = maxDNSMessageBytes + 4

	// DefaultGRPCMaxStreams is the default maximum number of concurrent streams per connection.
	DefaultGRPCMaxStreams = 256

	// DefaultGRPCMaxConnections is the default maximum number of concurrent connections.
	DefaultGRPCMaxConnections = 200
)

// ServergRPC represents an instance of a DNS-over-gRPC server.
type ServergRPC struct {
	*Server
	*pb.UnimplementedDnsServiceServer
	grpcServer     *grpc.Server
	listenAddr     net.Addr
	tlsConfig      *tls.Config
	maxStreams     int
	maxConnections int
}

// NewServergRPC returns a new CoreDNS GRPC server and compiles all plugin in to it.
		tlsConfig.NextProtos = []string{"h2"}
	}

	maxStreams := DefaultGRPCMaxStreams
	if len(group) > 0 && group[0] != nil && group[0].MaxGRPCStreams != nil {
		maxStreams = *group[0].MaxGRPCStreams
	}

	maxConnections := DefaultGRPCMaxConnections
	if len(group) > 0 && group[0] != nil && group[0].MaxGRPCConnections != nil {
		maxConnections = *group[0].MaxGRPCConnections
	}

	return &ServergRPC{
		Server:         s,
		tlsConfig:      tlsConfig,
		maxStreams:     maxStreams,
		maxConnections: maxConnections,
	}, nil
}

// Compile-time check to ensure ServergRPC implements the caddy.GracefulServer interface
	s.listenAddr = l.Addr()
	s.m.Unlock()

	serverOpts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(maxProtobufPayloadBytes),
		grpc.MaxSendMsgSize(maxProtobufPayloadBytes),
	}

	// Only set MaxConcurrentStreams if not unbounded (0)
	if s.maxStreams > 0 {
		serverOpts = append(serverOpts, grpc.MaxConcurrentStreams(uint32(s.maxStreams)))
	}

	if s.Tracer() != nil {
		onlyIfParent := func(parentSpanCtx opentracing.SpanContext, method string, req, resp any) bool {
			return parentSpanCtx != nil
		}
		serverOpts = append(serverOpts, grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(s.Tracer(), otgrpc.IncludingSpans(onlyIfParent))))
	}

	s.grpcServer = grpc.NewServer(serverOpts...)

	pb.RegisterDnsServiceServer(s.grpcServer, s)

	if s.tlsConfig != nil {
		l = tls.NewListener(l, s.tlsConfig)
	}

	// Wrap listener to limit concurrent connections
	if s.maxConnections > 0 {
		l = netutil.LimitListener(l, s.maxConnections)
	}

	return s.grpcServer.Serve(l)
}

// any normal server. We use a custom responseWriter to pick up the bytes we need to write
// back to the client as a protobuf.
func (s *ServergRPC) Query(ctx context.Context, in *pb.DnsPacket) (*pb.DnsPacket, error) {
	if len(in.GetMsg()) > dns.MaxMsgSize {
		return nil, fmt.Errorf("dns message exceeds size limit: %d", len(in.GetMsg()))
	}
	msg := new(dns.Msg)
	err := msg.Unpack(in.GetMsg())
	if err != nil {
package https

import (
	"strconv"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() {
	caddy.RegisterPlugin("https", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	err := parseDOH(c)
	if err != nil {
		return plugin.Error("https", err)
	}
	return nil
}

func parseDOH(c *caddy.Controller) error {
	config := dnsserver.GetConfig(c)

	// Skip the "https" directive itself
	c.Next()

	// Get any arguments on the "https" line
	args := c.RemainingArgs()
	if len(args) > 0 {
		return c.ArgErr()
	}

	// Process all nested directives in the block
	for c.NextBlock() {
		switch c.Val() {
		case "max_connections":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return c.ArgErr()
			}
			val, err := strconv.Atoi(args[0])
			if err != nil {
				return c.Errf("invalid max_connections value '%s': %v", args[0], err)
			}
			if val < 0 {
				return c.Errf("max_connections must be a non-negative integer: %d", val)
			}
			if config.MaxHTTPSConnections != nil {
				return c.Err("max_connections already defined for this server block")
			}
			config.MaxHTTPSConnections = &val
		default:
			return c.Errf("unknown property '%s'", c.Val())
		}
	}

	return nil
}
