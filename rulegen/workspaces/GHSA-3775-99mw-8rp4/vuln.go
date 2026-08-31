package main

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
	assert.Equal(t, wfv1.WorkflowError, woc.wf.Status.Phase)
}

func TestReferenceModeBlocksPodSpecPatch(t *testing.T) {
	wf := wfv1.MustUnmarshalWorkflow("@testdata/workflow-template-ref.yaml")
	wfTmpl := wfv1.MustUnmarshalWorkflowTemplate("@testdata/workflow-template-submittable.yaml")

		woc := newWorkflowOperationCtx(ctx, wfWithPatch, controller)
		woc.operate(ctx)
		assert.Equal(t, wfv1.WorkflowError, woc.wf.Status.Phase)
		assert.Contains(t, woc.wf.Status.Message, "podSpecPatch is not permitted")
	})

	t.Run("Secure rejects podSpecPatch", func(t *testing.T) {
		woc := newWorkflowOperationCtx(ctx, wfWithPatch, controller)
		woc.operate(ctx)
		assert.Equal(t, wfv1.WorkflowError, woc.wf.Status.Phase)
		assert.Contains(t, woc.wf.Status.Message, "podSpecPatch is not permitted")
	})

	t.Run("No restrictions allows podSpecPatch", func(t *testing.T) {
		assert.Equal(t, wfv1.WorkflowRunning, woc.wf.Status.Phase)
	})

	t.Run("Without podSpecPatch still works in Strict mode", func(t *testing.T) {
		cancel, controller := newController(logging.TestContext(t.Context()), wf, wfTmpl)
		defer cancel()


import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	wfv1 "github.com/argoproj/argo-workflows/v4/pkg/apis/workflow/v1alpha1"
)

// MergeTo will merge one workflow (the "patch" workflow) into another (the "target" workflow.
// If the target workflow defines a field, this take precedence over the patch.
func MergeTo(patch, target *wfv1.Workflow) error {
package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
