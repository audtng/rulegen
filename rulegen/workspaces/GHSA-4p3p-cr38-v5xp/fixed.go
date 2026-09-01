package main

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	resapi "github.com/siderolabs/omni/client/api/omni/resources"
	authres "github.com/siderolabs/omni/client/pkg/omni/resources/auth"
	"github.com/siderolabs/omni/internal/backend/runtime/omni/audit"
	"github.com/siderolabs/omni/internal/pkg/auth"
	"github.com/siderolabs/omni/internal/pkg/ctxstore"
// Unary returns a new unary GRPC interceptor.
func (c *AuthConfig) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		isGetAuthConfigRequest := false

		if req != nil && info != nil && info.FullMethod == resapi.ResourceService_Get_FullMethodName {
			if getReq, getReqOk := req.(*resapi.GetRequest); getReqOk && getReq.Type == authres.AuthConfigType {
				isGetAuthConfigRequest = true
			}
		}

		ctx = c.intercept(ctx, isGetAuthConfigRequest, info.FullMethod)

		return handler(ctx, req)
	}
// Stream returns a new streaming GRPC interceptor.
func (c *AuthConfig) Stream() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := c.intercept(ss.Context(), false, info.FullMethod)

		return handler(srv, &grpc_middleware.WrappedServerStream{
			ServerStream:   ss,
	}
}

func (c *AuthConfig) intercept(ctx context.Context, isGetAuthConfigRequest bool, method string) context.Context {
	ctx = ctxstore.WithValue(ctx, auth.EnabledAuthContextKey{Enabled: c.enabled})

	if !c.enabled {
		md = metadata.New(nil)
	}

	msg := message.NewGRPC(md, method, message.WithSignatureRequiredCheck(func() (bool, error) {
		return !isGetAuthConfigRequest, nil
	}))

	auditData, ok := ctxstore.Value[*audit.Data](ctx)
	if ok {
		grpc_ctxtags.UnaryServerInterceptor(),
		logLevelOverrideUnaryInterceptor,
		grpc_zap.UnaryServerInterceptor(s.logger, grpc_zap.WithMessageProducer(messageProducer)),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_recovery.UnaryServerInterceptor(recoveryOpt),
		grpcutil.SetUserAgent(),
		grpcutil.SetRealPeerAddress(),
		grpcutil.SetAuditData(),
			),
			1024,
		),
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(),
		logLevelOverrideStreamInterceptor,
		grpc_zap.StreamServerInterceptor(s.logger, grpc_zap.WithMessageProducer(messageProducer)),
		grpc_prometheus.StreamServerInterceptor,
		grpc_recovery.StreamServerInterceptor(recoveryOpt),
		grpcutil.StreamSetUserAgent(),
		grpcutil.StreamSetRealPeerAddress(),
		grpcutil.StreamSetAuditData(),
				),
			},
		),
	}

	authInterceptors, err := s.getAuthInterceptors()
func isSensitiveResource(res *v1alpha1.Resource) bool {
	protoR, err := protobuf.Unmarshal(res)
	if err != nil {
		return true
	}

	properResource, err := protobuf.UnmarshalResource(protoR)
	if err != nil {
		return true
	}

	resDef, ok := properResource.(meta.ResourceDefinitionProvider)
func isSensitiveSpec(resource *resapi.Resource) bool {
	res, err := grpcomni.CreateResource(resource)
	if err != nil {
		return true
	}

	resDef, ok := res.(meta.ResourceDefinitionProvider)

func getSource(ctx context.Context) common.Runtime {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		source := md.Get(message.RuntimeHeaderKey)
		if source != nil {
			if res, ok := common.Runtime_value[source[0]]; ok {
				return common.Runtime(res)

// CreateResource creates a resource from a resource proto representation.
func CreateResource(resource *resources.Resource) (cosiresource.Resource, error) { //nolint:ireturn
	if resource == nil {
		return nil, errors.New("resource is nil")
	}

	if resource.Metadata == nil {
		return nil, errors.New("resource metadata is nil")
	}

	if resource.Metadata.Version == "" {
		resource.Metadata.Version = "1"
	}
