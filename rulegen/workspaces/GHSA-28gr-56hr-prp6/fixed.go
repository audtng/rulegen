package main

	errors = append(errors, validateName(tempo.Name)...)
	addValidationResults(v.validateStorage(ctx, tempo))
	errors = append(errors, v.validateJaegerUI(tempo)...)
	errors = append(errors, v.validateMultitenancy(ctx, tempo)...)
	errors = append(errors, v.validateObservability(tempo)...)
	errors = append(errors, v.validateServiceAccount(ctx, tempo)...)
	errors = append(errors, v.validateConflictWithTempoStack(ctx, tempo)...)
	return nil
}

func (v *monolithicValidator) validateMultitenancy(ctx context.Context, tempo tempov1alpha1.TempoMonolithic) field.ErrorList {
	if tempo.Spec.Query != nil && tempo.Spec.Query.RBAC.Enabled && (tempo.Spec.Multitenancy == nil || !tempo.Spec.Multitenancy.Enabled) {
		return field.ErrorList{
			field.Invalid(field.NewPath("spec", "rbac", "enabled"), tempo.Spec.Query.RBAC.Enabled,

	multitenancyBase := field.NewPath("spec", "multitenancy")

	if tempo.Spec.Multitenancy != nil && tempo.Spec.Multitenancy.Mode == v1alpha1.ModeOpenShift {
		err := validateGatewayOpenShiftModeRBAC(ctx, v.client)
		if err != nil {
			return field.ErrorList{field.Invalid(
				multitenancyBase.Child("mode"),
				tempo.Spec.Multitenancy.Mode,
				fmt.Sprintf("Cannot enable OpenShift tenancy mode: %v", err),
			)}
		}
	}

	err := ValidateTenantConfigs(&tempo.Spec.Multitenancy.TenantsSpec, tempo.Spec.Multitenancy.IsGatewayEnabled())
	if err != nil {
		return field.ErrorList{field.Invalid(multitenancyBase.Child("enabled"), tempo.Spec.Multitenancy.Enabled, err.Error())}
	"context"
	"testing"

	configv1alpha1 "github.com/grafana/tempo-operator/api/config/v1alpha1"
	"github.com/grafana/tempo-operator/api/tempo/v1alpha1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	authorizationv1 "k8s.io/api/authorization/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestMonolithicValidate(t *testing.T) {
	ctx := admission.NewContextWithRequest(context.Background(), admission.Request{})

	tests := []struct {
		name       string
		ctrlConfig configv1alpha1.ProjectConfig

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &k8sFake{
				subjectAccessReview: &authorizationv1.SubjectAccessReview{
					Status: authorizationv1.SubjectAccessReviewStatus{
						Allowed: true,
					},
				},
			}
			v := &monolithicValidator{
				client:     client,
				ctrlConfig: test.ctrlConfig,
			}

			warnings, errors := v.validateTempoMonolithic(ctx, test.tempo)
			require.Equal(t, test.warnings, warnings)
			require.Equal(t, test.errors, errors)
		})
	return nil
}

func (v *validator) validateGateway(ctx context.Context, tempo v1alpha1.TempoStack) field.ErrorList {
	path := field.NewPath("spec").Child("template").Child("gateway").Child("enabled")
	if tempo.Spec.Template.Gateway.Enabled {
		if tempo.Spec.Template.QueryFrontend.JaegerQuery.Ingress.Type != v1alpha1.IngressTypeNone {
				"Cannot enable gateway and distributor TLS at the same time",
			)}
		}

		if tempo.Spec.Tenants != nil && tempo.Spec.Tenants.Mode == v1alpha1.ModeOpenShift {
			err := validateGatewayOpenShiftModeRBAC(ctx, v.client)
			if err != nil {
				return field.ErrorList{field.Invalid(
					field.NewPath("spec").Child("tenants").Child("mode"),
					tempo.Spec.Tenants.Mode,
					fmt.Sprintf("Cannot enable OpenShift tenancy mode: %v", err),
				)}
			}
		}
	}
	return nil
}

	allErrors = append(allErrors, v.validateReplicationFactor(*tempo)...)
	allErrors = append(allErrors, v.validateQueryFrontend(*tempo)...)
	allErrors = append(allErrors, v.validateGateway(ctx, *tempo)...)
	allErrors = append(allErrors, v.validateTenantConfigs(*tempo)...)
	allErrors = append(allErrors, v.validateObservability(*tempo)...)
	allErrors = append(allErrors, v.validateDeprecatedFields(*tempo)...)
	"time"

	"github.com/stretchr/testify/assert"
	authorizationv1 "k8s.io/api/authorization/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator := &validator{ctrlConfig: configv1alpha1.ProjectConfig{}}
			errs := validator.validateGateway(context.Background(), test.input)
			assert.Equal(t, test.expected, errs)
		})
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validator := &validator{ctrlConfig: configv1alpha1.ProjectConfig{}}
			errs := validator.validateGateway(context.Background(), test.input)
			assert.Equal(t, test.expected, errs)
		})
	}
}

type k8sFake struct {
	secret              *corev1.Secret
	configmap           *corev1.ConfigMap
	tempoStack          *v1alpha1.TempoStack
	tempoMonolithic     *v1alpha1.TempoMonolithic
	subjectAccessReview *authorizationv1.SubjectAccessReview
	client.Client
}

func (k *k8sFake) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	switch typed := obj.(type) {
	case *authorizationv1.SubjectAccessReview:
		if k.subjectAccessReview != nil {
			k.subjectAccessReview.DeepCopyInto(typed)
			return nil
		}
	}
	return fmt.Errorf("mock: fails always")
}

func (k *k8sFake) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	switch typed := obj.(type) {
	case *corev1.Secret:
package webhooks

import (
	"context"
	"fmt"

	"github.com/grafana/tempo-operator/internal/manifests/gateway"

	authenticationv1 "k8s.io/api/authentication/v1"
	authorizationv1 "k8s.io/api/authorization/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const maxLabelLength = 63
	}
	return allErrs
}

func subjectAccessReviewsForClusterRole(user authenticationv1.UserInfo, clusterRole rbacv1.ClusterRole) []authorizationv1.SubjectAccessReview {
	reviews := []authorizationv1.SubjectAccessReview{}
	for _, rule := range clusterRole.Rules {
		for _, apiGroup := range rule.APIGroups {
			for _, resource := range rule.Resources {
				for _, verb := range rule.Verbs {
					reviews = append(reviews, authorizationv1.SubjectAccessReview{
						Spec: authorizationv1.SubjectAccessReviewSpec{
							UID:    user.UID,
							User:   user.Username,
							Groups: user.Groups,
							ResourceAttributes: &authorizationv1.ResourceAttributes{
								Group:    apiGroup,
								Resource: resource,
								Verb:     verb,
							},
						},
					})
				}
			}
		}
	}

	return reviews
}

// validateGatewayOpenShiftModeRBAC checks if the user requesting the change on the CR
// has already the permissions which the operator would grant to the ServiceAccount of the Tempo instance
// when enabling the OpenShift tenancy mode.
//
// In other words, the operator should not grant e.g. TokenReview permissions to the ServiceAccount of the Tempo instance
// if the user creating or modifying the TempoStack or TempoMonolithic doesn't have these permissions.
func validateGatewayOpenShiftModeRBAC(ctx context.Context, client client.Client) error {
	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return err
	}

	user := req.UserInfo
	clusterRole := gateway.NewAccessReviewClusterRole("", map[string]string{})
	reviews := subjectAccessReviewsForClusterRole(user, *clusterRole)

	for _, sar := range reviews {
		err := client.Create(ctx, &sar)
		if err != nil {
			return fmt.Errorf("failed to create subject access review: %w", err)
		}

		if !sar.Status.Allowed {
			return fmt.Errorf("user %s does not have permission to %s %s.%s", user.Username, sar.Spec.ResourceAttributes.Verb, sar.Spec.ResourceAttributes.Resource, sar.Spec.ResourceAttributes.Group)
		}
	}

	return nil
}
