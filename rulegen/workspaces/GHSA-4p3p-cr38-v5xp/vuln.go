package main

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/siderolabs/omni/internal/backend/runtime/omni/audit"
	"github.com/siderolabs/omni/internal/pkg/auth"
	"github.com/siderolabs/omni/internal/pkg/ctxstore"
// Unary returns a new unary GRPC interceptor.
func (c *AuthConfig) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx = c.intercept(ctx, info.FullMethod)

		return handler(ctx, req)
	}
// Stream returns a new streaming GRPC interceptor.
func (c *AuthConfig) Stream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := c.intercept(ss.Context(), info.FullMethod)

		return handler(srv, &grpc_middleware.WrappedServerStream{
			ServerStream:   ss,
	}
}

func (c *AuthConfig) intercept(ctx context.Context, method string) context.Context {
	ctx = ctxstore.WithValue(ctx, auth.EnabledAuthContextKey{Enabled: c.enabled})

	if !c.enabled {
		md = metadata.New(nil)
	}

	msg := message.NewGRPC(md, method)

	auditData, ok := ctxstore.Value[*audit.Data](ctx)
	if ok {
		grpc_ctxtags.UnaryServerInterceptor(),
		logLevelOverrideUnaryInterceptor,
		grpc_zap.UnaryServerInterceptor(s.logger, grpc_zap.WithMessageProducer(messageProducer)),
		grpcutil.SetUserAgent(),
		grpcutil.SetRealPeerAddress(),
		grpcutil.SetAuditData(),
			),
			1024,
		),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_recovery.UnaryServerInterceptor(recoveryOpt),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(),
		logLevelOverrideStreamInterceptor,
		grpc_zap.StreamServerInterceptor(s.logger, grpc_zap.WithMessageProducer(messageProducer)),
		grpcutil.StreamSetUserAgent(),
		grpcutil.StreamSetRealPeerAddress(),
		grpcutil.StreamSetAuditData(),
				),
			},
		),
		grpc_prometheus.StreamServerInterceptor,
		grpc_recovery.StreamServerInterceptor(recoveryOpt),
	}

	authInterceptors, err := s.getAuthInterceptors()
func isSensitiveResource(res *v1alpha1.Resource) bool {
	protoR, err := protobuf.Unmarshal(res)
	if err != nil {
		return false
	}

	properResource, err := protobuf.UnmarshalResource(protoR)
	if err != nil {
		return false
	}

	resDef, ok := properResource.(meta.ResourceDefinitionProvider)
func isSensitiveSpec(resource *resapi.Resource) bool {
	res, err := grpcomni.CreateResource(resource)
	if err != nil {
		return false
	}

	resDef, ok := res.(meta.ResourceDefinitionProvider)

func getSource(ctx context.Context) common.Runtime {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		source := md.Get(message.RuntimeHeaderHey)
		if source != nil {
			if res, ok := common.Runtime_value[source[0]]; ok {
				return common.Runtime(res)

// CreateResource creates a resource from a resource proto representation.
func CreateResource(resource *resources.Resource) (cosiresource.Resource, error) { //nolint:ireturn
	if resource.Metadata.Version == "" {
		resource.Metadata.Version = "1"
	}
