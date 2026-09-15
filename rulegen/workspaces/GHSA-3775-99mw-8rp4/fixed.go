package main


import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	wfv1 "github.com/argoproj/argo-workflows/v4/pkg/apis/workflow/v1alpha1"
)

// allowedUserOverrideFields lists WorkflowSpec fields that users may set when
// submitting a workflow via workflowTemplateRef under Strict/Secure mode.
// Any field NOT in this set is blocked by default, ensuring new fields added
// to WorkflowSpec are denied until explicitly reviewed.
var allowedUserOverrideFields = map[string]bool{
	"Arguments":             true,
	"Entrypoint":            true,
	"Shutdown":              true,
	"Suspend":               true,
	"ActiveDeadlineSeconds": true,
	"Priority":              true,
	"TTLStrategy":           true,
	"PodGC":                 true,
	"VolumeClaimGC":         true,
	"ArchiveLogs":           true,
	"WorkflowMetadata":      true,
	"WorkflowTemplateRef":   true,
	"Metrics":               true,
	"ArtifactGC":            true,
}

// ValidateUserOverrides checks that a user-submitted WorkflowSpec only sets
// fields from the allow-list. Returns an error listing all violations.
func ValidateUserOverrides(userSpec *wfv1.WorkflowSpec) error {
	if userSpec == nil {
		return nil
	}
	v := reflect.ValueOf(userSpec).Elem()
	t := v.Type()
	zero := reflect.New(t).Elem()

	var violations []string
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		if allowedUserOverrideFields[fieldName] {
			continue
		}
		if !reflect.DeepEqual(v.Field(i).Interface(), zero.Field(i).Interface()) {
			violations = append(violations, fieldName)
		}
	}
	if len(violations) > 0 {
		sort.Strings(violations)
		return fmt.Errorf("fields %v are not permitted when using workflowTemplateRef with templateReferencing restriction", violations)
	}
	return nil
}

// SanitizeUserWorkflowSpec returns a copy of userSpec with only allow-listed
// fields preserved. This provides defense-in-depth after validation.
func SanitizeUserWorkflowSpec(userSpec *wfv1.WorkflowSpec) *wfv1.WorkflowSpec {
	if userSpec == nil {
		return nil
	}
	sanitized := &wfv1.WorkflowSpec{}
	src := reflect.ValueOf(userSpec).Elem()
	dst := reflect.ValueOf(sanitized).Elem()
	t := src.Type()

	for i := 0; i < t.NumField(); i++ {
		if allowedUserOverrideFields[t.Field(i).Name] {
			dst.Field(i).Set(src.Field(i))
		}
	}
	return sanitized
}

// MergeTo will merge one workflow (the "patch" workflow) into another (the "target" workflow.
// If the target workflow defines a field, this take precedence over the patch.
func MergeTo(patch, target *wfv1.Workflow) error {
func (woc *wfOperationCtx) setExecWorkflow(ctx context.Context) error {
	if woc.wf.Spec.WorkflowTemplateRef != nil { // not-woc-misuse
		// When workflow restrictions require template referencing (Strict/Secure mode),
		// reject workflows that set any non-allowed fields, as they could override
		// security settings defined in the WorkflowTemplate.
		if woc.controller.Config.WorkflowRestrictions.MustUseReference() { // not-woc-misuse: intentionally checking the user-submitted spec
			if err := wfutil.ValidateUserOverrides(&woc.wf.Spec); err != nil { // not-woc-misuse
				woc.markWorkflowError(ctx, err)
				return err
			}
		}
		err := woc.setStoredWfSpec(ctx)
		if err != nil {
	}
	// Update the Entrypoint, ShutdownStrategy and Suspend
	if woc.needsStoredWfSpecUpdate() {
		// In reference mode, sanitize the user spec before merging so that
		// only allow-listed fields participate in the strategic merge patch.
		userSpec := &woc.wf.Spec // not-woc-misuse
		if woc.controller.Config.WorkflowRestrictions.MustUseReference() {
			userSpec = wfutil.SanitizeUserWorkflowSpec(&woc.wf.Spec) // not-woc-misuse
		}
		// Join workflow, workflow template, and workflow default metadata to workflow spec.
		mergedWf, err := wfutil.JoinWorkflowSpec(userSpec, workflowTemplateSpec, &wfDefault.Spec)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		userSpec := &woc.wf.Spec // not-woc-misuse
		if woc.controller.Config.WorkflowRestrictions.MustUseReference() {
			userSpec = wfutil.SanitizeUserWorkflowSpec(&woc.wf.Spec) // not-woc-misuse
		}
		mergedWf, err := wfutil.JoinWorkflowSpec(userSpec, wftHolder.GetWorkflowSpec(), &wfDefault.Spec)
		if err != nil {
			return err
		}
