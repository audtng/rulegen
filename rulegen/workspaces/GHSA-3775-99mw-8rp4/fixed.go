package main

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
	assert.Equal(t, wfv1.WorkflowError, woc.wf.Status.Phase)
}

func TestReferenceModeBlocksDisallowedFields(t *testing.T) {
	wf := wfv1.MustUnmarshalWorkflow("@testdata/workflow-template-ref.yaml")
	wfTmpl := wfv1.MustUnmarshalWorkflowTemplate("@testdata/workflow-template-submittable.yaml")

		woc := newWorkflowOperationCtx(ctx, wfWithPatch, controller)
		woc.operate(ctx)
		assert.Equal(t, wfv1.WorkflowError, woc.wf.Status.Phase)
		assert.Contains(t, woc.wf.Status.Message, "PodSpecPatch")
		assert.Contains(t, woc.wf.Status.Message, "not permitted")
	})

	t.Run("Secure rejects podSpecPatch", func(t *testing.T) {
		woc := newWorkflowOperationCtx(ctx, wfWithPatch, controller)
		woc.operate(ctx)
		assert.Equal(t, wfv1.WorkflowError, woc.wf.Status.Phase)
		assert.Contains(t, woc.wf.Status.Message, "PodSpecPatch")
		assert.Contains(t, woc.wf.Status.Message, "not permitted")
	})

	t.Run("Strict rejects ServiceAccountName", func(t *testing.T) {
		wfCopy := wf.DeepCopy()
		wfCopy.Spec.ServiceAccountName = "admin"
		cancel, controller := newController(logging.TestContext(t.Context()), wfCopy, wfTmpl)
		defer cancel()

		ctx := logging.TestContext(t.Context())
		controller.Config.WorkflowRestrictions = &config.WorkflowRestrictions{
			TemplateReferencing: config.TemplateReferencingStrict,
		}
		woc := newWorkflowOperationCtx(ctx, wfCopy, controller)
		woc.operate(ctx)
		assert.Equal(t, wfv1.WorkflowError, woc.wf.Status.Phase)
		assert.Contains(t, woc.wf.Status.Message, "ServiceAccountName")
	})

	t.Run("Strict rejects SecurityContext", func(t *testing.T) {
		wfCopy := wf.DeepCopy()
		wfCopy.Spec.SecurityContext = &apiv1.PodSecurityContext{}
		cancel, controller := newController(logging.TestContext(t.Context()), wfCopy, wfTmpl)
		defer cancel()

		ctx := logging.TestContext(t.Context())
		controller.Config.WorkflowRestrictions = &config.WorkflowRestrictions{
			TemplateReferencing: config.TemplateReferencingStrict,
		}
		woc := newWorkflowOperationCtx(ctx, wfCopy, controller)
		woc.operate(ctx)
		assert.Equal(t, wfv1.WorkflowError, woc.wf.Status.Phase)
		assert.Contains(t, woc.wf.Status.Message, "SecurityContext")
	})

	t.Run("Strict rejects Templates", func(t *testing.T) {
		wfCopy := wf.DeepCopy()
		wfCopy.Spec.Templates = []wfv1.Template{{Name: "injected"}}
		cancel, controller := newController(logging.TestContext(t.Context()), wfCopy, wfTmpl)
		defer cancel()

		ctx := logging.TestContext(t.Context())
		controller.Config.WorkflowRestrictions = &config.WorkflowRestrictions{
			TemplateReferencing: config.TemplateReferencingStrict,
		}
		woc := newWorkflowOperationCtx(ctx, wfCopy, controller)
		woc.operate(ctx)
		assert.Equal(t, wfv1.WorkflowError, woc.wf.Status.Phase)
		assert.Contains(t, woc.wf.Status.Message, "Templates")
	})

	t.Run("Strict rejects Volumes", func(t *testing.T) {
		wfCopy := wf.DeepCopy()
		wfCopy.Spec.Volumes = []apiv1.Volume{{Name: "secret-vol"}}
		cancel, controller := newController(logging.TestContext(t.Context()), wfCopy, wfTmpl)
		defer cancel()

		ctx := logging.TestContext(t.Context())
		controller.Config.WorkflowRestrictions = &config.WorkflowRestrictions{
			TemplateReferencing: config.TemplateReferencingStrict,
		}
		woc := newWorkflowOperationCtx(ctx, wfCopy, controller)
		woc.operate(ctx)
		assert.Equal(t, wfv1.WorkflowError, woc.wf.Status.Phase)
		assert.Contains(t, woc.wf.Status.Message, "Volumes")
	})

	t.Run("Strict rejects HostNetwork", func(t *testing.T) {
		wfCopy := wf.DeepCopy()
		hostNet := true
		wfCopy.Spec.HostNetwork = &hostNet
		cancel, controller := newController(logging.TestContext(t.Context()), wfCopy, wfTmpl)
		defer cancel()

		ctx := logging.TestContext(t.Context())
		controller.Config.WorkflowRestrictions = &config.WorkflowRestrictions{
			TemplateReferencing: config.TemplateReferencingStrict,
		}
		woc := newWorkflowOperationCtx(ctx, wfCopy, controller)
		woc.operate(ctx)
		assert.Equal(t, wfv1.WorkflowError, woc.wf.Status.Phase)
		assert.Contains(t, woc.wf.Status.Message, "HostNetwork")
	})

	t.Run("No restrictions allows podSpecPatch", func(t *testing.T) {
		assert.Equal(t, wfv1.WorkflowRunning, woc.wf.Status.Phase)
	})

	t.Run("Allowed fields pass in Strict mode", func(t *testing.T) {
		wfCopy := wf.DeepCopy()
		// Set allowed fields only - entrypoint must match one defined in the template
		wfCopy.Spec.Entrypoint = "whalesay-template"
		wfCopy.Spec.Shutdown = wfv1.ShutdownStrategyTerminate
		cancel, controller := newController(logging.TestContext(t.Context()), wfCopy, wfTmpl)
		defer cancel()

		ctx := logging.TestContext(t.Context())
		controller.Config.WorkflowRestrictions = &config.WorkflowRestrictions{
			TemplateReferencing: config.TemplateReferencingStrict,
		}
		woc := newWorkflowOperationCtx(ctx, wfCopy, controller)
		woc.operate(ctx)
		// Shutdown=Terminate means it won't reach Running, but it shouldn't be Error
		assert.NotEqual(t, wfv1.WorkflowError, woc.wf.Status.Phase)
	})

	t.Run("Without disallowed fields still works in Strict mode", func(t *testing.T) {
		cancel, controller := newController(logging.TestContext(t.Context()), wf, wfTmpl)
		defer cancel()


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
package util

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	wfv1 "github.com/argoproj/argo-workflows/v4/pkg/apis/workflow/v1alpha1"
)
}

var wfDefault = `
metadata:
  annotations:
    testAnnotation: default
  labels:
    testLabel: default
spec:
  entrypoint: whalesay
  activeDeadlineSeconds: 7200
  arguments:
    artifacts:
      -
        name: message
        path: /tmp/message
    parameters:
      -
        name: message
        value: "hello world"
  onExit: whalesay-exit
  serviceAccountName: default
  templates:
    -
      container:
        args:
          - "hello from the default exit handler"
        command:
          - cowsay
        image: docker/whalesay
      name: whalesay-exit
  ttlStrategy:
    secondsAfterCompletion: 60
  volumes:
    -
      name: test
      secret:
        secretName: test
`

  namespace: default
spec:
  workflowMetaData:
    annotations:
      testAnnotation: wft
    labels:
      testLabel: wft
var wf = `
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: hello-world-
spec:
  entrypoint: whalesay
  templates:
    -
      container:
        args:
          - "hello world"
        command:
          - cowsay
        image: "docker/whalesay:latest"
      name: whalesay
var resultSpec = `
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: hello-world-
spec:
  activeDeadlineSeconds: 7200
  workflowMetadata:
    annotations:
      testAnnotation: wft
    labels:
      testLabel: wft
  arguments:
    artifacts:
      -
        name: message
        path: /tmp/message
    parameters:
      -
        name: message
        value: "hello world"
  entrypoint: whalesay
  onExit: whalesay-exit
  serviceAccountName: default
  templates:
    -
      container:
        args:
          - "hello world"
        command:
          - cowsay
        image: "docker/whalesay:latest"
      name: whalesay
    -
      container:
        args:
          - "{{inputs.parameters.message}}"
        command:
          - cowsay
        image: docker/whalesay
      inputs:
        parameters:
          -
            name: message
      name: whalesay-template
    -
      container:
        args:
          - "hello from the default exit handler"
        command:
          - cowsay
        image: docker/whalesay
      name: whalesay-exit
  ttlStrategy:
    secondsAfterCompletion: 60
  volumes:
    -
      name: test
      secret:
        secretName: test

`
		assert.Equal(t, "BASE", targetWf.Spec.WorkflowMetadata.LabelsFrom[`baz`].Expression)
	})
}

// blockedUserOverrideFields is used by TestAllWorkflowSpecFieldsAccountedFor
// to verify that every WorkflowSpec field is consciously classified as either
// allowed or blocked.
var blockedUserOverrideFields = map[string]bool{
	"Templates":                    true,
	"TemplateDefaults":             true,
	"ServiceAccountName":           true,
	"AutomountServiceAccountToken": true,
	"Executor":                     true,
	"Volumes":                      true,
	"VolumeClaimTemplates":         true,
	"Parallelism":                  true,
	"NodeSelector":                 true,
	"Affinity":                     true,
	"Tolerations":                  true,
	"ImagePullSecrets":             true,
	"HostNetwork":                  true,
	"DNSPolicy":                    true,
	"DNSConfig":                    true,
	"OnExit":                       true,
	"SchedulerName":                true,
	"PodPriorityClassName":         true,
	"HostAliases":                  true,
	"SecurityContext":              true,
	"PodSpecPatch":                 true,
	"PodDisruptionBudget":          true,
	"ArtifactRepositoryRef":        true,
	"Synchronization":              true,
	"RetryStrategy":                true,
	"PodMetadata":                  true,
	"Hooks":                        true,
}

func TestValidateUserOverrides_AllowedFields(t *testing.T) {
	spec := &wfv1.WorkflowSpec{
		Entrypoint: "main",
		Arguments: wfv1.Arguments{
			Parameters: []wfv1.Parameter{{Name: "msg", Value: wfv1.AnyStringPtr("hello")}},
		},
		Shutdown:            wfv1.ShutdownStrategyTerminate,
		Priority:            ptr.To[int32](10),
		WorkflowTemplateRef: &wfv1.WorkflowTemplateRef{Name: "my-template"},
	}
	err := ValidateUserOverrides(spec)
	assert.NoError(t, err)
}

func TestValidateUserOverrides_BlockedFields(t *testing.T) {
	tests := []struct {
		name  string
		spec  wfv1.WorkflowSpec
		field string
	}{
		{
			name:  "ServiceAccountName",
			spec:  wfv1.WorkflowSpec{ServiceAccountName: "admin"},
			field: "ServiceAccountName",
		},
		{
			name:  "SecurityContext",
			spec:  wfv1.WorkflowSpec{SecurityContext: &apiv1.PodSecurityContext{}},
			field: "SecurityContext",
		},
		{
			name:  "Templates",
			spec:  wfv1.WorkflowSpec{Templates: []wfv1.Template{{Name: "evil"}}},
			field: "Templates",
		},
		{
			name:  "Volumes",
			spec:  wfv1.WorkflowSpec{Volumes: []apiv1.Volume{{Name: "secret-vol"}}},
			field: "Volumes",
		},
		{
			name:  "HostNetwork",
			spec:  wfv1.WorkflowSpec{HostNetwork: ptr.To(true)},
			field: "HostNetwork",
		},
		{
			name:  "PodSpecPatch",
			spec:  wfv1.WorkflowSpec{PodSpecPatch: `{"containers":[]}`},
			field: "PodSpecPatch",
		},
		{
			name:  "OnExit",
			spec:  wfv1.WorkflowSpec{OnExit: "backdoor"},
			field: "OnExit",
		},
		{
			name:  "Hooks",
			spec:  wfv1.WorkflowSpec{Hooks: wfv1.LifecycleHooks{"exit": wfv1.LifecycleHook{Template: "evil"}}},
			field: "Hooks",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserOverrides(&tt.spec)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.field)
			assert.Contains(t, err.Error(), "not permitted")
		})
	}
}

func TestValidateUserOverrides_MultipleViolations(t *testing.T) {
	spec := &wfv1.WorkflowSpec{
		ServiceAccountName: "admin",
		HostNetwork:        ptr.To(true),
		PodSpecPatch:       `{}`,
		Templates:          []wfv1.Template{{Name: "evil"}},
	}
	err := ValidateUserOverrides(spec)
	require.Error(t, err)
	msg := err.Error()
	assert.Contains(t, msg, "ServiceAccountName")
	assert.Contains(t, msg, "HostNetwork")
	assert.Contains(t, msg, "PodSpecPatch")
	assert.Contains(t, msg, "Templates")
}

func TestValidateUserOverrides_NilSpec(t *testing.T) {
	assert.NoError(t, ValidateUserOverrides(nil))
}

func TestSanitizeUserWorkflowSpec(t *testing.T) {
	spec := &wfv1.WorkflowSpec{
		Entrypoint:         "main",
		ServiceAccountName: "admin",
		HostNetwork:        ptr.To(true),
		Arguments: wfv1.Arguments{
			Parameters: []wfv1.Parameter{{Name: "msg", Value: wfv1.AnyStringPtr("hello")}},
		},
		Volumes:             []apiv1.Volume{{Name: "secret-vol"}},
		Templates:           []wfv1.Template{{Name: "evil"}},
		Shutdown:            wfv1.ShutdownStrategyTerminate,
		WorkflowTemplateRef: &wfv1.WorkflowTemplateRef{Name: "my-template"},
	}

	sanitized := SanitizeUserWorkflowSpec(spec)

	// Allowed fields are preserved
	assert.Equal(t, "main", sanitized.Entrypoint)
	assert.Equal(t, wfv1.ShutdownStrategyTerminate, sanitized.Shutdown)
	assert.Len(t, sanitized.Arguments.Parameters, 1)
	assert.Equal(t, "my-template", sanitized.WorkflowTemplateRef.Name)

	// Blocked fields are zeroed
	assert.Empty(t, sanitized.ServiceAccountName)
	assert.Nil(t, sanitized.HostNetwork)
	assert.Nil(t, sanitized.Volumes)
	assert.Nil(t, sanitized.Templates)
}

func TestSanitizeUserWorkflowSpec_Nil(t *testing.T) {
	assert.Nil(t, SanitizeUserWorkflowSpec(nil))
}

// TestAllWorkflowSpecFieldsAccountedFor is a compile-time safety net.
// It ensures that every field in WorkflowSpec appears in either the
// allowed or blocked list, so new fields force a conscious decision.
func TestAllWorkflowSpecFieldsAccountedFor(t *testing.T) {
	specType := reflect.TypeOf(wfv1.WorkflowSpec{})
	for i := 0; i < specType.NumField(); i++ {
		fieldName := specType.Field(i).Name
		inAllowed := allowedUserOverrideFields[fieldName]
		inBlocked := blockedUserOverrideFields[fieldName]
		if !inAllowed && !inBlocked {
			t.Errorf("WorkflowSpec field %q is not classified in either allowedUserOverrideFields or blockedUserOverrideFields — add it to one of them", fieldName)
		}
		if inAllowed && inBlocked {
			t.Errorf("WorkflowSpec field %q appears in both allowedUserOverrideFields and blockedUserOverrideFields — it should be in exactly one", fieldName)
		}
	}
}
