package main

	"sort"

	"google.golang.org/grpc/codes"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterwftmplpkg "github.com/argoproj/argo-workflows/v4/pkg/apiclient/clusterworkflowtemplate"
}

func (cwts *Server) GetClusterWorkflowTemplate(ctx context.Context, req *clusterwftmplpkg.ClusterWorkflowTemplateGetRequest) (*v1alpha1.ClusterWorkflowTemplate, error) {
	wfTmpl, err := cwts.getTemplateAndValidate(ctx, req.Name)
	if err != nil {
		return nil, serverutils.ToStatusError(err, codes.Internal)
	"strings"

	"google.golang.org/grpc/codes"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	workflowtemplatepkg "github.com/argoproj/argo-workflows/v4/pkg/apiclient/workflowtemplate"
}

func (wts *Server) GetWorkflowTemplate(ctx context.Context, req *workflowtemplatepkg.WorkflowTemplateGetRequest) (*v1alpha1.WorkflowTemplate, error) {
	wfTmpl, err := wts.getTemplateAndValidate(ctx, req.Namespace, req.Name)
	if err != nil {
		return nil, sutils.ToStatusError(err, codes.Internal)
