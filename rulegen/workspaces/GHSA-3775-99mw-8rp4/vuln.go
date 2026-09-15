package main


import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	wfv1 "github.com/argoproj/argo-workflows/v4/pkg/apis/workflow/v1alpha1"
)

// MergeTo will merge one workflow (the "patch" workflow) into another (the "target" workflow.
// If the target workflow defines a field, this take precedence over the patch.
func MergeTo(patch, target *wfv1.Workflow) error {
func (woc *wfOperationCtx) setExecWorkflow(ctx context.Context) error {
	if woc.wf.Spec.WorkflowTemplateRef != nil { // not-woc-misuse
		// When workflow restrictions require template referencing (Strict/Secure mode),
		// reject workflows that include a podSpecPatch as it could override security
		// settings defined in the WorkflowTemplate.
		if woc.controller.Config.WorkflowRestrictions.MustUseReference() && woc.wf.Spec.HasPodSpecPatch() { // not-woc-misuse: intentionally checking the user-submitted spec
			err := fmt.Errorf("podSpecPatch is not permitted when using workflowTemplateRef with templateReferencing restriction")
			woc.markWorkflowError(ctx, err)
			return err
		}
		err := woc.setStoredWfSpec(ctx)
		if err != nil {
	}
	// Update the Entrypoint, ShutdownStrategy and Suspend
	if woc.needsStoredWfSpecUpdate() {
		// Join workflow, workflow template, and workflow default metadata to workflow spec.
		mergedWf, err := wfutil.JoinWorkflowSpec(&woc.wf.Spec, workflowTemplateSpec, &wfDefault.Spec) // not-woc-misuse
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		mergedWf, err := wfutil.JoinWorkflowSpec(&woc.wf.Spec, wftHolder.GetWorkflowSpec(), &wfDefault.Spec) // not-woc-misuse
		if err != nil {
			return err
		}
