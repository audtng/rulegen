package main

	errors = append(errors, validateName(tempo.Name)...)
	addValidationResults(v.validateStorage(ctx, tempo))
	errors = append(errors, v.validateJaegerUI(tempo)...)
	errors = append(errors, v.validateMultitenancy(tempo)...)
	errors = append(errors, v.validateObservability(tempo)...)
	errors = append(errors, v.validateServiceAccount(ctx, tempo)...)
	errors = append(errors, v.validateConflictWithTempoStack(ctx, tempo)...)
	return nil
}

func (v *monolithicValidator) validateMultitenancy(tempo tempov1alpha1.TempoMonolithic) field.ErrorList {
	if tempo.Spec.Query != nil && tempo.Spec.Query.RBAC.Enabled && (tempo.Spec.Multitenancy == nil || !tempo.Spec.Multitenancy.Enabled) {
		return field.ErrorList{
			field.Invalid(field.NewPath("spec", "rbac", "enabled"), tempo.Spec.Query.RBAC.Enabled,

	multitenancyBase := field.NewPath("spec", "multitenancy")

	err := ValidateTenantConfigs(&tempo.Spec.Multitenancy.TenantsSpec, tempo.Spec.Multitenancy.IsGatewayEnabled())
	if err != nil {
		return field.ErrorList{field.Invalid(multitenancyBase.Child("enabled"), tempo.Spec.Multitenancy.Enabled, err.Error())}
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	configv1alpha1 "github.com/grafana/tempo-operator/api/config/v1alpha1"
	"github.com/grafana/tempo-operator/api/tempo/v1alpha1"
)

func TestMonolithicValidate(t *testing.T) {
	tests := []struct {
		name       string
		ctrlConfig configv1alpha1.ProjectConfig

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &k8sFake{}
			v := &monolithicValidator{
				client:     client,
				ctrlConfig: test.ctrlConfig,
			}

			warnings, errors := v.validateTempoMonolithic(context.Background(), test.tempo)
			require.Equal(t, test.warnings, warnings)
			require.Equal(t, test.errors, errors)
		})
	return nil
}

func (v *validator) validateGateway(tempo v1alpha1.TempoStack) field.ErrorList {
	path := field.NewPath("spec").Child("template").Child("gateway").Child("enabled")
	if tempo.Spec.Template.Gateway.Enabled {
		if tempo.Spec.Template.QueryFrontend.JaegerQuery.Ingress.Type != v1alpha1.IngressTypeNone {
				"Cannot enable gateway and distributor TLS at the same time",
			)}
		}
	}
	return nil
}

	allErrors = append(allErrors, v.validateReplicationFactor(*tempo)...)
	allErrors = append(allErrors, v.validateQueryFrontend(*tempo)...)
	allErrors = append(allErrors, v.validateGateway(*tempo)...)
	allErrors = append(allErrors, v.validateTenantConfigs(*tempo)...)
	allErrors = append(allErrors, v.validateObservability(*tempo)...)
	allErrors = append(allErrors, v.validateDeprecatedFields(*tempo)...)
	"time"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator := &validator{ctrlConfig: configv1alpha1.ProjectConfig{}}
			errs := validator.validateGateway(test.input)
			assert.Equal(t, test.expected, errs)
		})
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator := &validator{ctrlConfig: configv1alpha1.ProjectConfig{}}
			errs := validator.validateGateway(test.input)
			assert.Equal(t, test.expected, errs)
		})
	}
}

type k8sFake struct {
	secret          *corev1.Secret
	configmap       *corev1.ConfigMap
	tempoStack      *v1alpha1.TempoStack
	tempoMonolithic *v1alpha1.TempoMonolithic
	client.Client
}

func (k *k8sFake) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	switch typed := obj.(type) {
	case *corev1.Secret:
package webhooks

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const maxLabelLength = 63
	}
	return allErrs
}
