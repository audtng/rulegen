package main

	extensionscmdcontroller "github.com/gardener/gardener/extensions/pkg/controller/cmd"
	"github.com/gardener/gardener/extensions/pkg/util"
	extensionscmdwebhook "github.com/gardener/gardener/extensions/pkg/webhook/cmd"
	gardencoreinstall "github.com/gardener/gardener/pkg/apis/core/install"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	securityinstall "github.com/gardener/gardener/pkg/apis/security/install"
	gardenerhealthz "github.com/gardener/gardener/pkg/healthz"
	admissioncmd "github.com/gardener/gardener/pkg/provider-local/admission/cmd"
	localinstall "github.com/gardener/gardener/pkg/provider-local/apis/local/install"
				return fmt.Errorf("could not instantiate manager: %w", err)
			}

			gardencoreinstall.Install(mgr.GetScheme())
			securityinstall.Install(mgr.GetScheme())

			if err := localinstall.AddToScheme(mgr.GetScheme()); err != nil {
				return fmt.Errorf("could not update manager scheme: %w", err)
	"github.com/gardener/gardener/pkg/admissioncontroller/webhook/admission/internaldomainsecret"
	"github.com/gardener/gardener/pkg/admissioncontroller/webhook/admission/kubeconfigsecret"
	"github.com/gardener/gardener/pkg/admissioncontroller/webhook/admission/namespacedeletion"
	"github.com/gardener/gardener/pkg/admissioncontroller/webhook/admission/providersecretlabels"
	"github.com/gardener/gardener/pkg/admissioncontroller/webhook/admission/resourcesize"
	"github.com/gardener/gardener/pkg/admissioncontroller/webhook/admission/seedrestriction"
	"github.com/gardener/gardener/pkg/admissioncontroller/webhook/admission/shootkubeconfigsecretref"
		return fmt.Errorf("failed adding %s webhook handler: %w", namespacedeletion.HandlerName, err)
	}

	if err := (&providersecretlabels.Handler{
		Logger: mgr.GetLogger().WithName("webhook").WithName(providersecretlabels.HandlerName),
		Client: mgr.GetClient(),
	}).AddToManager(mgr); err != nil {
		return fmt.Errorf("failed adding %s webhook handler: %w", providersecretlabels.HandlerName, err)
	}

	if err := (&resourcesize.Handler{
		Logger: mgr.GetLogger().WithName("webhook").WithName(resourcesize.HandlerName),
		Config: cfg.Server.ResourceAdmissionConfiguration,
// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package providersecretlabels

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	// HandlerName is the name of this admission webhook handler.
	HandlerName = "sync-provider-secret-labels"
	// WebhookPath is the HTTP handler path for this admission webhook handler.
	WebhookPath = "/webhooks/sync-provider-secret-labels"
)

// AddToManager adds Handler to the given manager.
func (h *Handler) AddToManager(mgr manager.Manager) error {
	webhook := admission.
		WithCustomDefaulter(mgr.GetScheme(), &corev1.Secret{}, h).
		WithRecoverPanic(true)

	mgr.GetWebhookServer().Register(WebhookPath, webhook)
	return nil
}
// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package providersecretlabels

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	v1beta1helper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"
	securityv1alpha1 "github.com/gardener/gardener/pkg/apis/security/v1alpha1"
)

// Handler syncs the provider labels on Secrets referenced in SecretBindings or CredentialsBindings.
type Handler struct {
	Logger logr.Logger
	Client client.Client
}

// Default syncs the provider labels.
func (h *Handler) Default(ctx context.Context, obj runtime.Object) error {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return fmt.Errorf("expected secret but got %T", obj)
	}

	typesFromSecretBindings, err := h.fetchProviderTypesFromSecretBindings(ctx, secret)
	if err != nil {
		return fmt.Errorf("failed fetching provider types from SecretBindings: %w", err)
	}

	typesFromCredentialsBindings, err := h.fetchProviderTypesFromCredentialsBindings(ctx, secret)
	if err != nil {
		return fmt.Errorf("failed fetching provider types from CredentialsBindings: %w", err)
	}

	if typesFromSecretBindings.Len()+typesFromCredentialsBindings.Len() > 0 {
		maintainLabels(secret, typesFromSecretBindings.Union(typesFromCredentialsBindings).UnsortedList()...)
	}

	return nil
}

func (h *Handler) fetchProviderTypesFromSecretBindings(ctx context.Context, secret *corev1.Secret) (sets.Set[string], error) {
	secretBindingList := &gardencorev1beta1.SecretBindingList{}
	if err := h.Client.List(ctx, secretBindingList); err != nil {
		return nil, fmt.Errorf("failed to list SecretBindings: %w", err)
	}

	providerTypes := sets.New[string]()
	for _, secretBinding := range secretBindingList.Items {
		if secretBinding.SecretRef.Name == secret.Name &&
			secretBinding.SecretRef.Namespace == secret.Namespace {
			providerTypes.Insert(v1beta1helper.GetSecretBindingTypes(&secretBinding)...)
		}
	}
	return providerTypes, nil
}

func (h *Handler) fetchProviderTypesFromCredentialsBindings(ctx context.Context, secret *corev1.Secret) (sets.Set[string], error) {
	credentialsBindingList := &securityv1alpha1.CredentialsBindingList{}
	if err := h.Client.List(ctx, credentialsBindingList); err != nil {
		return nil, fmt.Errorf("failed to list CredentialsBindings: %w", err)
	}

	providerTypes := sets.New[string]()
	for _, credentialsBinding := range credentialsBindingList.Items {
		if credentialsBinding.CredentialsRef.APIVersion == corev1.SchemeGroupVersion.String() &&
			credentialsBinding.CredentialsRef.Kind == "Secret" &&
			credentialsBinding.CredentialsRef.Name == secret.Name &&
			credentialsBinding.CredentialsRef.Namespace == secret.Namespace {
			providerTypes.Insert(credentialsBinding.Provider.Type)
		}
	}
	return providerTypes, nil
}

func maintainLabels(secret *corev1.Secret, providerTypes ...string) {
	for k := range secret.Labels {
		if strings.HasPrefix(k, v1beta1constants.LabelShootProviderPrefix) {
			delete(secret.Labels, k)
		}
	}

	for _, providerType := range providerTypes {
		metav1.SetMetaDataLabel(&secret.ObjectMeta, v1beta1constants.LabelShootProviderPrefix+providerType, "true")
	}
}
// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package providersecretlabels_test

import (
	"context"

	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	logzap "sigs.k8s.io/controller-runtime/pkg/log/zap"

	. "github.com/gardener/gardener/pkg/admissioncontroller/webhook/admission/providersecretlabels"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	securityv1alpha1 "github.com/gardener/gardener/pkg/apis/security/v1alpha1"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	"github.com/gardener/gardener/pkg/logger"
)

var _ = Describe("handler", func() {
	var (
		ctx context.Context

		log        logr.Logger
		fakeClient client.Client
		handler    *Handler

		namespace            string
		provider1, provider2 string
		secret               *corev1.Secret
		secretBinding        *gardencorev1beta1.SecretBinding
		credentialsBinding   *securityv1alpha1.CredentialsBinding
	)

	BeforeEach(func() {
		ctx = context.Background()
		log = logger.MustNewZapLogger(logger.DebugLevel, logger.FormatJSON, logzap.WriteTo(GinkgoWriter))

		fakeClient = fakeclient.NewClientBuilder().WithScheme(kubernetes.GardenScheme).Build()
		handler = &Handler{
			Logger: log,
			Client: fakeClient,
		}

		namespace = "test"
		secret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-secret",
				Namespace: namespace,
			},
		}

		provider1 = "provider1"
		provider2 = "provider2"

		secretBinding = &gardencorev1beta1.SecretBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-secret-binding",
				Namespace: namespace,
			},
			SecretRef: corev1.SecretReference{
				Name:      "test-secret",
				Namespace: namespace,
			},
			Provider: &gardencorev1beta1.SecretBindingProvider{
				Type: provider1,
			},
		}

		credentialsBinding = &securityv1alpha1.CredentialsBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-credentials-binding",
				Namespace: "another-namespace",
			},
			CredentialsRef: corev1.ObjectReference{
				APIVersion: corev1.SchemeGroupVersion.String(),
				Kind:       "Secret",
				Name:       "test-secret",
				Namespace:  namespace,
			},
			Provider: securityv1alpha1.CredentialsBindingProvider{
				Type: provider2,
			},
		}
	})

	It("should set the provider label based on the available credential and secret bindings", func() {
		Expect(fakeClient.Create(ctx, secretBinding)).To(Succeed())
		Expect(fakeClient.Create(ctx, credentialsBinding)).To(Succeed())

		Expect(handler.Default(ctx, secret)).To(Succeed())

		Expect(secret.Labels).To(HaveKeyWithValue("provider.shoot.gardener.cloud/provider1", "true"))
		Expect(secret.Labels).To(HaveKeyWithValue("provider.shoot.gardener.cloud/provider2", "true"))
	})

	It("should remove undesired provider type", func() {
		Expect(fakeClient.Create(ctx, secretBinding)).To(Succeed())
		Expect(fakeClient.Create(ctx, credentialsBinding)).To(Succeed())
		secret.Labels = map[string]string{
			"provider.shoot.gardener.cloud/provider1": "true",
			"provider.shoot.gardener.cloud/provider2": "true",
			"provider.shoot.gardener.cloud/provider3": "true",
		}

		Expect(handler.Default(ctx, secret)).To(Succeed())

		Expect(secret.Labels).To(HaveKeyWithValue("provider.shoot.gardener.cloud/provider1", "true"))
		Expect(secret.Labels).To(HaveKeyWithValue("provider.shoot.gardener.cloud/provider2", "true"))
		Expect(secret.Labels).NotTo(HaveKey("provider.shoot.gardener.cloud/provider3"))
	})

	It("should add the missing provider and delete the wrong one", func() {
		Expect(fakeClient.Create(ctx, secretBinding)).To(Succeed())
		Expect(fakeClient.Create(ctx, credentialsBinding)).To(Succeed())
		secret.Labels = map[string]string{
			"provider.shoot.gardener.cloud/provider1": "true",
			"provider.shoot.gardener.cloud/provider3": "true",
		}

		Expect(handler.Default(ctx, secret)).To(Succeed())

		Expect(secret.Labels).To(HaveKeyWithValue("provider.shoot.gardener.cloud/provider1", "true"))
		Expect(secret.Labels).To(HaveKeyWithValue("provider.shoot.gardener.cloud/provider2", "true"))
		Expect(secret.Labels).NotTo(HaveKey("provider.shoot.gardener.cloud/provider3"))
	})

	It("should not add provider labels when secret is unreferenced", func() {
		Expect(handler.Default(ctx, secret)).To(Succeed())
		Expect(secret.Labels).To(BeEmpty())
	})
})
// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package providersecretlabels_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProviderSecretLabels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AdmissionController Webhook Admission ProviderSecretLabels Suite")
}

// GetSecretBindingTypes returns the SecretBinding provider types.
func GetSecretBindingTypes(secretBinding *core.SecretBinding) []string {
	if secretBinding.Provider == nil {
		return []string{}
	}
	return strings.Split(secretBinding.Provider.Type, ",")
}
			Expect(actual).To(Equal(expected))
		},

		Entry("with nil provider type", &core.SecretBinding{Provider: nil}, []string{}),
		Entry("with single-value provider type", &core.SecretBinding{Provider: &core.SecretBindingProvider{Type: "foo"}}, []string{"foo"}),
		Entry("with multi-value provider type", &core.SecretBinding{Provider: &core.SecretBindingProvider{Type: "foo,bar,baz"}}, []string{"foo", "bar", "baz"}),
	)

// GetSecretBindingTypes returns the SecretBinding provider types.
func GetSecretBindingTypes(secretBinding *gardencorev1beta1.SecretBinding) []string {
	if secretBinding.Provider == nil {
		return []string{}
	}
	return strings.Split(secretBinding.Provider.Type, ",")
}
			Expect(actual).To(Equal(expected))
		},

		Entry("with nil provider type", &gardencorev1beta1.SecretBinding{Provider: nil}, []string{}),
		Entry("with single-value provider type", &gardencorev1beta1.SecretBinding{Provider: &gardencorev1beta1.SecretBindingProvider{Type: "foo"}}, []string{"foo"}),
		Entry("with multi-value provider type", &gardencorev1beta1.SecretBinding{Provider: &gardencorev1beta1.SecretBindingProvider{Type: "foo,bar,baz"}}, []string{"foo", "bar", "baz"}),
	)
	"github.com/gardener/gardener/plugin/pkg/global/deletionconfirmation"
	"github.com/gardener/gardener/plugin/pkg/global/extensionlabels"
	"github.com/gardener/gardener/plugin/pkg/global/extensionvalidation"
	"github.com/gardener/gardener/plugin/pkg/global/finalizerremoval"
	"github.com/gardener/gardener/plugin/pkg/global/resourcereferencemanager"
	managedseedshoot "github.com/gardener/gardener/plugin/pkg/managedseed/shoot"
	managedseedvalidator "github.com/gardener/gardener/plugin/pkg/managedseed/validator"
func RegisterAllAdmissionPlugins(plugins *admission.Plugins) {
	resourcereferencemanager.Register(plugins)
	deletionconfirmation.Register(plugins)
	finalizerremoval.Register(plugins)
	extensionvalidation.Register(plugins)
	extensionlabels.Register(plugins)
	shoottolerationrestriction.Register(plugins)

func TestBackupBucket(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry Core BackupBucket Suite")
}

func TestBackupEntry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry Core BackupEntry Suite")
}

func TestCloudProfile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry Core CloudProfile Suite")
}

func TestControllerInstallation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry ControllerInstallation Suite")
}

var _ = Describe("ToSelectableFields", func() {

func TestNamespacedCloudProfile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry Core NamespacedCloudProfile Suite")
}

func TestProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry Core Project Suite")
}
func TestSecretBinding(t *testing.T) {
	features.RegisterFeatureGates()
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry Core SecretBinding Suite")
}
	return true
}

func (s secretBindingStrategy) PrepareForCreate(_ context.Context, obj runtime.Object) {
	binding := obj.(*core.SecretBinding)

	if binding.GetName() == "" {
		binding.SetName(s.GenerateName(binding.GetGenerateName()))
	}
}

func (secretBindingStrategy) Validate(_ context.Context, obj runtime.Object) field.ErrorList {
		}
	})

	Describe("#PrepareForCreate", func() {
		It("should set the name if not set", func() {
			secretBinding.SetName("")

			secretbindingregistry.Strategy.PrepareForCreate(context.TODO(), secretBinding)

			Expect(secretBinding.GetName()).NotTo(BeEmpty())
		})

		It("should set name with generateName as prefix", func() {
			genName := "prefix-"
			secretBinding.GenerateName = genName
			secretBinding.Name = ""

			secretbindingregistry.Strategy.PrepareForCreate(context.TODO(), secretBinding)

			Expect(secretBinding.GetGenerateName()).To(Equal(genName))
			Expect(secretBinding.GetName()).To(HavePrefix(genName))
		})

		It("should not overwrite already set name", func() {
			secretBinding.SetName("bar")

			secretbindingregistry.Strategy.PrepareForCreate(context.TODO(), secretBinding)

			Expect(secretBinding.GetName()).To(Equal("bar"))
		})
	})

	Describe("#Validate", func() {
		It("should forbid creating SecretBinding when provider is nil or empty", func() {
			secretBinding.Provider = nil

func TestSeed(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry Core Seed Suite")
}
func TestShoot(t *testing.T) {
	features.RegisterFeatureGates()
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry Core Shoot Suite")
}

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry Core Shoot Storage Suite")
}

func TestBastion(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry Operations Bastion Suite")
}
// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package credentialsbinding_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gardener/gardener/pkg/apiserver/features"
)

func TestCredentialsBinding(t *testing.T) {
	features.RegisterFeatureGates()
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry Security CredentialsBinding Suite")
}
	return true
}

func (c credentialsBindingStrategy) PrepareForCreate(_ context.Context, obj runtime.Object) {
	credentialsbinding := obj.(*security.CredentialsBinding)

	if credentialsbinding.GetName() == "" {
		credentialsbinding.SetName(c.GenerateName(credentialsbinding.GetGenerateName()))
	}
}

func (credentialsBindingStrategy) PrepareForUpdate(_ context.Context, _, _ runtime.Object) {
// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package credentialsbinding_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gardener/gardener/pkg/apis/security"
	credentialsbindingregistry "github.com/gardener/gardener/pkg/apiserver/registry/security/credentialsbinding"
)

var _ = Describe("Strategy", func() {
	var credentialsBinding *security.CredentialsBinding

	BeforeEach(func() {
		credentialsBinding = &security.CredentialsBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "profile",
				Namespace: "garden",
			},
		}
	})

	Describe("#PrepareForCreate", func() {
		It("should set the name if not set", func() {
			credentialsBinding.SetName("")

			credentialsbindingregistry.Strategy.PrepareForCreate(context.TODO(), credentialsBinding)

			Expect(credentialsBinding.GetName()).NotTo(BeEmpty())
		})

		It("should set name with generateName as prefix", func() {
			genName := "prefix-"
			credentialsBinding.GenerateName = genName
			credentialsBinding.Name = ""

			credentialsbindingregistry.Strategy.PrepareForCreate(context.TODO(), credentialsBinding)

			Expect(credentialsBinding.GetGenerateName()).To(Equal(genName))
			Expect(credentialsBinding.GetName()).To(HavePrefix(genName))
		})

		It("should not overwrite already set name", func() {
			credentialsBinding.SetName("bar")

			credentialsbindingregistry.Strategy.PrepareForCreate(context.TODO(), credentialsBinding)

			Expect(credentialsBinding.GetName()).To(Equal("bar"))
		})
	})
})

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry Security WorkloadIdentity Storage Suite")
}

func TestWorkloadIdentity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry Security WorkloadIdentity Suite")
}

func TestManagedSeed(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry SeedManagement ManagedSeed Suite")
}

func TestManagedSeedSet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "APIServer Registry SeedManagement ManagedSeedSet Suite")
}
		a.clusterRole(),
		a.clusterRoleBinding(virtualGardenAccessSecret.ServiceAccountName),
		a.validatingWebhookConfiguration(caSecret),
		a.mutatingWebhookConfiguration(caSecret),
	)
	if err != nil {
		return err
		clusterRole(),
		clusterRoleBinding(),
		validatingWebhookConfiguration(namespace, caGardener.Data["bundle.crt"], testValues),
		mutatingWebhookConfiguration(namespace, caGardener.Data["bundle.crt"]),
	))

	virtualManagedResourceSecret := &corev1.Secret{

	return webhookConfig
}

func mutatingWebhookConfiguration(namespace string, caBundle []byte) *admissionregistrationv1.MutatingWebhookConfiguration {
	var (
		failurePolicyFail = admissionregistrationv1.Fail
		sideEffectsNone   = admissionregistrationv1.SideEffectClassNone
	)

	return &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "gardener-admission-controller",
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{
			{
				Name:                    "sync-provider-secret-labels.gardener.cloud",
				AdmissionReviewVersions: []string{"v1", "v1beta1"},
				TimeoutSeconds:          ptr.To[int32](10),
				Rules: []admissionregistrationv1.RuleWithOperations{{
					Operations: []admissionregistrationv1.OperationType{admissionregistrationv1.Create, admissionregistrationv1.Update},
					Rule: admissionregistrationv1.Rule{
						APIGroups:   []string{""},
						APIVersions: []string{"v1"},
						Resources:   []string{"secrets"},
					},
				}},
				FailurePolicy: &failurePolicyFail,
				NamespaceSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"gardener.cloud/role": "project",
					},
				},
				ClientConfig: admissionregistrationv1.WebhookClientConfig{
					URL:      ptr.To("https://gardener-admission-controller." + namespace + "/webhooks/sync-provider-secret-labels"),
					CABundle: caBundle,
				},
				SideEffects: &sideEffectsNone,
			},
		},
	}
}
	return validatingWebhook
}

func (a *gardenerAdmissionController) mutatingWebhookConfiguration(caSecret *corev1.Secret) *admissionregistrationv1.MutatingWebhookConfiguration {
	var (
		failurePolicyFail = admissionregistrationv1.Fail
		sideEffectsNone   = admissionregistrationv1.SideEffectClassNone

		caBundle = caSecret.Data[secrets.DataKeyCertificateBundle]
	)

	return &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: DeploymentName,
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{
			{
				Name:                    "sync-provider-secret-labels.gardener.cloud",
				AdmissionReviewVersions: []string{"v1", "v1beta1"},
				TimeoutSeconds:          ptr.To[int32](10),
				Rules: []admissionregistrationv1.RuleWithOperations{{
					Operations: []admissionregistrationv1.OperationType{admissionregistrationv1.Create, admissionregistrationv1.Update},
					Rule: admissionregistrationv1.Rule{
						APIGroups:   []string{corev1.GroupName},
						APIVersions: []string{"v1"},
						Resources:   []string{"secrets"},
					},
				}},
				FailurePolicy: &failurePolicyFail,
				NamespaceSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						v1beta1constants.GardenRole: v1beta1constants.GardenRoleProject,
					},
				},
				ClientConfig: admissionregistrationv1.WebhookClientConfig{
					URL:      buildClientConfigURL("/webhooks/sync-provider-secret-labels", a.namespace),
					CABundle: caBundle,
				},
				SideEffects: &sideEffectsNone,
			},
		},
	}
}

func buildWebhookConfigRulesForResourceSize(config *admissioncontrollerconfigv1alpha1.ResourceAdmissionConfiguration) []admissionregistrationv1.RuleWithOperations {
	if config == nil || len(config.Limits) == 0 {
		return nil
		}
	}

	types := v1beta1helper.GetSecretBindingTypes(secretBinding)
	for _, t := range types {
		labelKey := v1beta1constants.LabelShootProviderPrefix + t

		if !metav1.HasLabel(secret.ObjectMeta, labelKey) {
			patch := client.MergeFrom(secret.DeepCopy())
			metav1.SetMetaDataLabel(&secret.ObjectMeta, labelKey, "true")
			if err := r.Client.Patch(ctx, secret, patch); err != nil {
				return reconcile.Result{}, fmt.Errorf("failed to add provider type label to Secret referenced in SecretBinding: %w", err)
			}
		}
	}
func GardenWebhookSwitchOptions() *extensionscmdwebhook.SwitchOptions {
	return extensionscmdwebhook.NewSwitchOptions(
		extensionscmdwebhook.Switch(validator.Name, validator.New),
		extensionscmdwebhook.Switch(validator.SecretsValidatorName, validator.NewSecretsWebhook),
		extensionscmdwebhook.Switch(mutator.Name, mutator.New),
	)
}
// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"sigs.k8s.io/controller-runtime/pkg/client"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
)

type secretValidator struct{}

// NewSecretValidator returns a new instance of a secret validator.
func NewSecretValidator() extensionswebhook.Validator {
	return &secretValidator{}
}

// Validate checks whether the data is empty.
func (s *secretValidator) Validate(_ context.Context, newObj, oldObj client.Object) error {
	secret, ok := newObj.(*corev1.Secret)
	if !ok {
		return fmt.Errorf("wrong object type %T", newObj)
	}

	if oldObj != nil {
		oldSecret, ok := oldObj.(*corev1.Secret)
		if !ok {
			return fmt.Errorf("wrong object type %T for old object", oldObj)
		}

		if apiequality.Semantic.DeepEqual(secret.Data, oldSecret.Data) {
			return nil
		}
	}

	if len(secret.Data) != 0 {
		return fmt.Errorf("secret data should be empty")
	}

	return nil
}
package validator

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/pkg/apis/core"
	securityv1alpha1 "github.com/gardener/gardener/pkg/apis/security/v1alpha1"
	"github.com/gardener/gardener/pkg/provider-local/local"
)

const (
	// Name is a name for a validation webhook.
	Name = "validator"
	// SecretsValidatorName is the name of the secrets validator.
	SecretsValidatorName = "secrets." + Name
)

var logger = log.Log.WithName("local-validator-webhook")
		Path:     "/webhooks/validate",
		Validators: map[extensionswebhook.Validator][]extensionswebhook.Type{
			NewNamespacedCloudProfileValidator(mgr): {{Obj: &core.NamespacedCloudProfile{}}},
			NewWorkloadIdentityValidator():          {{Obj: &securityv1alpha1.WorkloadIdentity{}}},
		},
		Target: extensionswebhook.TargetSeed,
		ObjectSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"provider.extensions.gardener.cloud/" + local.Type: "true"},
		},
	})
}

// NewSecretsWebhook creates a new validation webhook for Secrets.
func NewSecretsWebhook(mgr manager.Manager) (*extensionswebhook.Webhook, error) {
	logger.Info("Setting up webhook", "name", SecretsValidatorName)

	return extensionswebhook.New(mgr, extensionswebhook.Args{
		Provider: local.Type,
		Name:     SecretsValidatorName,
		Path:     "/webhooks/validate/secrets",
		Validators: map[extensionswebhook.Validator][]extensionswebhook.Type{
			NewSecretValidator(): {{Obj: &corev1.Secret{}}},
		},
		Target: extensionswebhook.TargetSeed,
		ObjectSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"provider.shoot.gardener.cloud/" + local.Type: "true"},
		},
	})
}
// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"context"
	"errors"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	securityv1alpha1 "github.com/gardener/gardener/pkg/apis/security/v1alpha1"
)

type workloadIdentityValidator struct {
}

// NewWorkloadIdentityValidator returns a new instance of a WorkloadIdentity validator.
func NewWorkloadIdentityValidator() extensionswebhook.Validator {
	return &workloadIdentityValidator{}
}

// Validate checks whether the provider config is empty.
func (wi *workloadIdentityValidator) Validate(_ context.Context, newObj, _ client.Object) error {
	workloadIdentity, ok := newObj.(*securityv1alpha1.WorkloadIdentity)
	if !ok {
		return fmt.Errorf("wrong object type %T", newObj)
	}

	if workloadIdentity.Spec.TargetSystem.ProviderConfig != nil {
		return errors.New("target system provider config must be empty")
	}

	return nil
}
}

func addMetaDataLabelsSecretBinding(secretBinding *core.SecretBinding) {
	types := gardencorehelper.GetSecretBindingTypes(secretBinding)
	for _, t := range types {
		metav1.SetMetaDataLabel(&secretBinding.ObjectMeta, v1beta1constants.LabelExtensionProviderTypePrefix+t, "true")
	}
}

// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package finalizerremoval

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gardener/gardener/pkg/apis/core"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/gardener/pkg/apis/security"
	admissioninitializer "github.com/gardener/gardener/pkg/apiserver/admission/initializer"
	gardencoreinformers "github.com/gardener/gardener/pkg/client/core/informers/externalversions"
	gardencorev1beta1listers "github.com/gardener/gardener/pkg/client/core/listers/core/v1beta1"
	plugin "github.com/gardener/gardener/plugin/pkg"
)

// Register registers a plugin.
func Register(plugins *admission.Plugins) {
	plugins.Register(plugin.PluginNameFinalizerRemoval, func(_ io.Reader) (admission.Interface, error) {
		return New()
	})
}

// FinalizerRemoval contains listers and admission handler.
type FinalizerRemoval struct {
	*admission.Handler
	shootLister gardencorev1beta1listers.ShootLister
	readyFunc   admission.ReadyFunc
}

var (
	_ = admissioninitializer.WantsCoreInformerFactory(&FinalizerRemoval{})

	readyFuncs []admission.ReadyFunc
)

// New creates a new FinalizerRemoval admission plugin.
func New() (*FinalizerRemoval, error) {
	return &FinalizerRemoval{
		Handler: admission.NewHandler(admission.Update),
	}, nil
}

// AssignReadyFunc assigns the ready function to the admission handler.
func (f *FinalizerRemoval) AssignReadyFunc(fn admission.ReadyFunc) {
	f.readyFunc = fn
	f.SetReadyFunc(fn)
}

// SetCoreInformerFactory gets Lister from SharedInformerFactory.
func (f *FinalizerRemoval) SetCoreInformerFactory(g gardencoreinformers.SharedInformerFactory) {
	shootInformer := g.Core().V1beta1().Shoots()
	f.shootLister = shootInformer.Lister()

	readyFuncs = append(readyFuncs,
		shootInformer.Informer().HasSynced,
	)
}

// ValidateInitialization checks whether the plugin was correctly initialized.
func (f *FinalizerRemoval) ValidateInitialization() error {
	if f.shootLister == nil {
		return errors.New("missing shoot lister")
	}
	return nil
}

// Admit ensures that finalizers from objects can only be removed if they are not needed anymore.
func (f *FinalizerRemoval) Admit(_ context.Context, a admission.Attributes, _ admission.ObjectInterfaces) error {
	// Wait until the caches have been synced
	if f.readyFunc == nil {
		f.AssignReadyFunc(func() bool {
			for _, readyFunc := range readyFuncs {
				if !readyFunc() {
					return false
				}
			}
			return true
		})
	}
	if !f.WaitForReady() {
		return admission.NewForbidden(a, errors.New("not yet ready to handle request"))
	}

	var (
		err            error
		newObj, oldObj client.Object
	)

	oldObj, ok := a.GetOldObject().(client.Object)
	if !ok {
		return nil
	}

	newObj, ok = a.GetObject().(client.Object)
	if !ok {
		return nil
	}

	switch a.GetKind().GroupKind() {
	case core.Kind("SecretBinding"):
		binding, ok := a.GetObject().(*core.SecretBinding)
		if !ok {
			return apierrors.NewBadRequest("could not convert resource into SecretBinding object")
		}

		// Allow removal of `gardener` finalizer only if the SecretBinding is not used by any shoot.
		if isFinalizerRemoved(oldObj, newObj, gardencorev1beta1.GardenerName) {
			inUse, err := f.isUsedByShoot(binding.Namespace, func(shoot *gardencorev1beta1.Shoot) bool {
				return ptr.Deref(shoot.Spec.SecretBindingName, "") == binding.Name
			})
			if err != nil {
				return apierrors.NewInternalError(fmt.Errorf("error checking if secret binding is in use: %w", err))
			}
			if inUse {
				return admission.NewForbidden(a, fmt.Errorf("finalizer must not be removed - secret binding %s/%s is still in use by at least one shoot", binding.Namespace, binding.Name))
			}
		}
	case security.Kind("CredentialsBinding"):
		binding, ok := a.GetObject().(*security.CredentialsBinding)
		if !ok {
			return apierrors.NewBadRequest("could not convert resource into CredentialsBinding object")
		}

		// Allow removal of `gardener` finalizer only if the CredentialsBinding is not used by any shoot.
		if isFinalizerRemoved(oldObj, newObj, gardencorev1beta1.GardenerName) {
			inUse, err := f.isUsedByShoot(binding.Namespace, func(shoot *gardencorev1beta1.Shoot) bool {
				return ptr.Deref(shoot.Spec.CredentialsBindingName, "") == binding.Name
			})
			if err != nil {
				return apierrors.NewInternalError(fmt.Errorf("error checking if credentials binding is in use: %w", err))
			}
			if inUse {
				return admission.NewForbidden(a, fmt.Errorf("finalizer must not be removed - credentials binding %s/%s is still in use by at least one shoot", binding.Namespace, binding.Name))
			}
		}
	case core.Kind("Shoot"):
		shoot, ok := a.GetObject().(*core.Shoot)
		if !ok {
			return apierrors.NewBadRequest("could not convert resource into Shoot object")
		}

		// Allow removal of `gardener` finalizer only if the Shoot deletion has completed successfully.
		if isFinalizerRemoved(oldObj, newObj, gardencorev1beta1.GardenerName) && !shootDeletionSucceeded(shoot) {
			return admission.NewForbidden(a, fmt.Errorf("finalizer %q cannot be removed because shoot deletion has not completed successfully yet", core.GardenerName))
		}
	}

	if err != nil {
		return admission.NewForbidden(a, err)
	}
	return nil
}

func (f *FinalizerRemoval) isUsedByShoot(namespace string, inUse func(*gardencorev1beta1.Shoot) bool) (bool, error) {
	shoots, err := f.shootLister.Shoots(namespace).List(labels.Everything())
	if err != nil {
		return false, fmt.Errorf("error retrieving shoots: %w", err)
	}

	return slices.ContainsFunc(shoots, inUse), nil
}

func shootDeletionSucceeded(shoot *core.Shoot) bool {
	if len(shoot.Status.TechnicalID) == 0 || shoot.Status.LastOperation == nil {
		return true
	}

	lastOperation := shoot.Status.LastOperation
	return lastOperation.Type == core.LastOperationTypeDelete &&
		lastOperation.State == core.LastOperationStateSucceeded &&
		lastOperation.Progress == 100
}

func isFinalizerRemoved(old, new metav1.Object, finalizerName string) bool {
	var (
		oldFinalizers = sets.New(old.GetFinalizers()...)
		newFinalizer  = sets.New(new.GetFinalizers()...)
	)

	return oldFinalizers.Has(finalizerName) && !newFinalizer.Has(finalizerName)
}
// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package finalizerremoval_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/utils/ptr"

	"github.com/gardener/gardener/pkg/apis/core"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/gardener/pkg/apis/security"
	gardencoreinformers "github.com/gardener/gardener/pkg/client/core/informers/externalversions"
	. "github.com/gardener/gardener/plugin/pkg/global/finalizerremoval"
)

var _ = Describe("finalizerremoval", func() {
	Describe("#Admit", func() {
		var (
			ctx                       context.Context
			admissionHandler          *FinalizerRemoval
			gardenCoreInformerFactory gardencoreinformers.SharedInformerFactory

			finalizers []string

			namespace              = "default"
			secretBindingName      = "binding-1"
			credentialsBindingName = "credentials-binding-1"
			shootName              = "shoot-1"

			shoot *gardencorev1beta1.Shoot
		)

		BeforeEach(func() {
			ctx = context.Background()
			admissionHandler, _ = New()
			admissionHandler.AssignReadyFunc(func() bool { return true })

			finalizers = []string{core.GardenerName}

			shoot = &gardencorev1beta1.Shoot{
				ObjectMeta: metav1.ObjectMeta{
					Name:      shootName,
					Namespace: namespace,
				},
				Spec: gardencorev1beta1.ShootSpec{
					CredentialsBindingName: ptr.To(credentialsBindingName),
					SecretBindingName:      ptr.To(secretBindingName),
				},
			}

			gardenCoreInformerFactory = gardencoreinformers.NewSharedInformerFactory(nil, 0)
			admissionHandler.SetCoreInformerFactory(gardenCoreInformerFactory)
		})

		Context("SecretBinding", func() {
			var coreSecretBinding *core.SecretBinding

			BeforeEach(func() {
				coreSecretBinding = &core.SecretBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name:       secretBindingName,
						Namespace:  namespace,
						Finalizers: finalizers,
					},
				}
			})

			It("should admit the removal because object is not used by any shoot", func() {
				attrs := admission.NewAttributesRecord(&core.SecretBinding{}, coreSecretBinding, core.Kind("SecretBinding").WithVersion("version"), "", coreSecretBinding.Name, core.Resource("SecretBinding").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, nil)

				Expect(admissionHandler.Admit(ctx, attrs, nil)).NotTo(HaveOccurred())
			})

			It("should admit the removal because finalizer is irrelevant", func() {
				newSecretBinding := coreSecretBinding.DeepCopy()
				coreSecretBinding.Finalizers = append(coreSecretBinding.Finalizers, "irrelevant-finalizer")

				attrs := admission.NewAttributesRecord(newSecretBinding, coreSecretBinding, core.Kind("SecretBinding").WithVersion("version"), "", coreSecretBinding.Name, core.Resource("SecretBinding").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, nil)

				Expect(admissionHandler.Admit(ctx, attrs, nil)).NotTo(HaveOccurred())
			})

			It("should reject the removal because object is not used by any shoot", func() {
				newSecretBinding := coreSecretBinding.DeepCopy()
				newSecretBinding.Finalizers = nil

				secondShoot := shoot.DeepCopy()
				secondShoot.Name = shootName + "-2"
				secondShoot.Spec.SecretBindingName = ptr.To(secretBindingName + "-2")

				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(shoot)).To(Succeed())
				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(secondShoot)).To(Succeed())

				attrs := admission.NewAttributesRecord(newSecretBinding, coreSecretBinding, core.Kind("SecretBinding").WithVersion("version"), "", coreSecretBinding.Name, core.Resource("SecretBinding").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, nil)

				Expect(admissionHandler.Admit(ctx, attrs, nil)).To(MatchError(ContainSubstring("finalizer must not be removed")))
			})
		})

		Context("CredentialsBinding", func() {
			var coreCredentialsBinding *security.CredentialsBinding

			BeforeEach(func() {
				coreCredentialsBinding = &security.CredentialsBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name:       credentialsBindingName,
						Namespace:  namespace,
						Finalizers: finalizers,
					},
				}
			})

			It("should admit the removal because object is not used by any shoot", func() {
				attrs := admission.NewAttributesRecord(&security.CredentialsBinding{}, coreCredentialsBinding, security.Kind("CredentialsBinding").WithVersion("version"), "", coreCredentialsBinding.Name, security.Resource("CredentialsBinding").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, nil)

				Expect(admissionHandler.Admit(ctx, attrs, nil)).NotTo(HaveOccurred())
			})

			It("should admit the removal because finalizer is irrelevant", func() {
				newCredentialsBinding := coreCredentialsBinding.DeepCopy()
				coreCredentialsBinding.Finalizers = append(coreCredentialsBinding.Finalizers, "irrelevant-finalizer")

				attrs := admission.NewAttributesRecord(newCredentialsBinding, coreCredentialsBinding, security.Kind("CredentialsBinding").WithVersion("version"), "", coreCredentialsBinding.Name, security.Resource("CredentialsBinding").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, nil)

				Expect(admissionHandler.Admit(ctx, attrs, nil)).NotTo(HaveOccurred())
			})

			It("should reject the removal because object is not used by any shoot", func() {
				newCredentialsBinding := coreCredentialsBinding.DeepCopy()
				newCredentialsBinding.Finalizers = nil

				secondShoot := shoot.DeepCopy()
				secondShoot.Name = shootName + "-2"
				secondShoot.Spec.CredentialsBindingName = ptr.To(secretBindingName + "-2")

				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(shoot)).To(Succeed())
				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(secondShoot)).To(Succeed())

				attrs := admission.NewAttributesRecord(newCredentialsBinding, coreCredentialsBinding, security.Kind("CredentialsBinding").WithVersion("version"), "", coreCredentialsBinding.Name, security.Resource("CredentialsBinding").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, nil)

				Expect(admissionHandler.Admit(ctx, attrs, nil)).To(MatchError(ContainSubstring("finalizer must not be removed")))
			})
		})

		Context("shoot", func() {
			var coreShoot *core.Shoot

			BeforeEach(func() {
				coreShoot = &core.Shoot{
					ObjectMeta: metav1.ObjectMeta{
						Finalizers: finalizers,
					},
					Status: core.ShootStatus{
						TechnicalID: "some-id",
						LastOperation: &core.LastOperation{
							Type:     core.LastOperationTypeReconcile,
							State:    core.LastOperationStateSucceeded,
							Progress: 100,
						},
					},
				}
			})

			It("should allow the removal because finalizer is irrelevant", func() {
				newShoot := coreShoot.DeepCopy()
				coreShoot.Finalizers = append(coreShoot.Finalizers, "irrelevant-finalizer")

				attrs := admission.NewAttributesRecord(newShoot, coreShoot, security.Kind("Shoot").WithVersion("version"), "", coreShoot.Name, security.Resource("Shoot").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, nil)

				Expect(admissionHandler.Admit(ctx, attrs, nil)).To(Succeed())
			})

			It("should admit the removal if the shoot deletion succeeded ", func() {
				newShoot := coreShoot.DeepCopy()
				newShoot.Finalizers = nil
				newShoot.Status.LastOperation.Type = core.LastOperationTypeDelete

				attrs := admission.NewAttributesRecord(newShoot, coreShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, nil)
				Expect(admissionHandler.Admit(ctx, attrs, nil)).To(Succeed())
			})

			It("should reject the removal if the shoot has not yet been deleted successfully", func() {
				newShoot := coreShoot.DeepCopy()
				newShoot.Finalizers = nil

				attrs := admission.NewAttributesRecord(newShoot, coreShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, nil)
				Expect(admissionHandler.Admit(ctx, attrs, nil)).To(MatchError(ContainSubstring("shoot deletion has not completed successfully yet")))
			})

			It("should admit the removal if the shoot has not yet a last operation", func() {
				newShoot := coreShoot.DeepCopy()
				newShoot.Finalizers = nil
				newShoot.Status.LastOperation = nil

				attrs := admission.NewAttributesRecord(newShoot, coreShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, nil)
				Expect(admissionHandler.Admit(ctx, attrs, nil)).To(Succeed())
			})

			It("should admit the removal if the shoot has not yet a technical id", func() {
				newShoot := coreShoot.DeepCopy()
				newShoot.Finalizers = nil
				newShoot.Status.TechnicalID = ""

				attrs := admission.NewAttributesRecord(newShoot, coreShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, nil)
				Expect(admissionHandler.Admit(ctx, attrs, nil)).To(Succeed())
			})
		})
	})
})
// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package finalizerremoval_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFinalizerRemoval(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AdmissionPlugin Global FinalizerRemoval Suite")
}
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	kubecorev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gardener/gardener/pkg/apis/core"
	return nil
}

// Validate ensures that referenced resources do actually exist.
func (r *ReferenceManager) Validate(ctx context.Context, a admission.Attributes, _ admission.ObjectInterfaces) error {
	// Wait until the caches have been synced
	if r.readyFunc == nil {
		r.AssignReadyFunc(func() bool {

		switch a.GetOperation() {
		case admission.Create:
			oldShoot = &core.Shoot{}
		case admission.Update:
			// skip verification if spec wasn't changed
		if utils.SkipVerification(operation, project.ObjectMeta) {
			return nil
		}

		switch a.GetOperation() {
		case admission.Create:
			err = r.ensureProjectNamespace(project)
		case admission.Update:
			oldProject, ok := a.GetOldObject().(*core.Project)
			}
		}

	case core.Kind("BackupBucket"):
		if operation == admission.Delete {
			// The "delete endpoint" handler of the k8s.io/apiserver library calls the admission controllers
		credentialsNamespace  string
		credentialsName       string
		credentialsKind       string
		providerTypes         []string
		credentialsReferenced func(shoot *gardencorev1beta1.Shoot) bool
	)

	switch attributes.GetKind().GroupKind() {
	case core.Kind("SecretBinding"):
		b, ok := binding.(*core.SecretBinding)
		credentialsNamespace = b.SecretRef.Namespace
		credentialsName = b.SecretRef.Name
		credentialsKind = "Secret"
		providerTypes = helper.GetSecretBindingTypes(b)
		credentialsReferenced = func(shoot *gardencorev1beta1.Shoot) bool {
			return ptr.Deref(shoot.Spec.SecretBindingName, "") == b.Name
		}

	case security.Kind("CredentialsBinding"):
		b, ok := binding.(*security.CredentialsBinding)
		if !ok {
		}
		credentialsNamespace = b.CredentialsRef.Namespace
		credentialsName = b.CredentialsRef.Name
		providerTypes = []string{b.Provider.Type}
		credentialsReferenced = func(shoot *gardencorev1beta1.Shoot) bool {
			return ptr.Deref(shoot.Spec.CredentialsBindingName, "") == b.Name
		}

	default:
		return fmt.Errorf("%s is neither of kind SecretBinding nor CredentialsBinding", attributes.GetKind().GroupKind())
	}

	shoots, err := r.shootLister.Shoots(attributes.GetNamespace()).List(labels.Everything())
	if err != nil {
		return fmt.Errorf("failed listing shoots: %w", err)
	}

	for _, shoot := range shoots {
		if !credentialsReferenced(shoot) {
			continue
		}

		if !slices.Contains(providerTypes, shoot.Spec.Provider.Type) {
			return fmt.Errorf("%s is referenced by shoot %q, but provider types (%+v) do not match with the shoot provider type %q", attributes.GetKind().Kind, shoot.Name, providerTypes, shoot.Spec.Provider.Type)
		}
	}

	readAttributes := authorizer.AttributesRecord{
		User:            attributes.GetUserInfo(),
		Verb:            "get",
		if err := r.lookupSecret(ctx, credentialsNamespace, credentialsName); err != nil {
			return err
		}
		if err := r.sanityCheckProviderSecret(ctx, credentialsNamespace, credentialsName, providerTypes); err != nil {
			return err
		}

	case "WorkloadIdentity":
		workloadIdentity, err := r.lookupWorkloadIdentity(ctx, credentialsNamespace, credentialsName)
		if err != nil {
			return err
		}

		if !slices.Contains(providerTypes, workloadIdentity.Spec.TargetSystem.Type) {
			return fmt.Errorf("CredentialsBinding provider type (%+v) does not match with WorkloadIdentity provider type %s", providerTypes, workloadIdentity.Spec.TargetSystem.Type)
		}

	default:
		return fmt.Errorf("unknown credentials kind: %s", credentialsKind)
	}

type getFn func(context.Context, string, string) (runtime.Object, error)

func lookupResource(ctx context.Context, namespace, name string, get getFn, fallbackGet getFn) (runtime.Object, error) {
	// First try to detect the resource in the cache.
	obj, err := get(ctx, namespace, name)
	if err == nil {
		return obj, nil
	}
	if !apierrors.IsNotFound(err) {
		return nil, err
	}

	// Second try to detect the resource in the cache after the first try failed.
	// Give the cache time to observe the resource before rejecting a create.
	// This helps when creating a resource and immediately creating a binding referencing it.
	time.Sleep(MissingResourceWait)
	obj, err = get(ctx, namespace, name)

	switch {
	case apierrors.IsNotFound(err):
		// no-op
	case err != nil:
		return nil, err
	default:
		return obj, nil
	}

	// Third try to detect the secret, now by doing a live lookup instead of relying on the cache.
	return fallbackGet(ctx, namespace, name)
}

func (r *ReferenceManager) lookupWorkloadIdentity(ctx context.Context, namespace, name string) (*securityv1alpha1.WorkloadIdentity, error) {
	workloadIdentityFromLister := func(_ context.Context, namespace, name string) (runtime.Object, error) {
		return r.workloadIdentityLister.WorkloadIdentities(namespace).Get(name)
	}
		return r.gardenSecurityClient.SecurityV1alpha1().WorkloadIdentities(namespace).Get(ctx, name, kubernetesclient.DefaultGetOptions())
	}

	obj, err := lookupResource(ctx, namespace, name, workloadIdentityFromLister, workloadIdentityFromClient)
	if err != nil {
		return nil, err
	}
	return obj.(*securityv1alpha1.WorkloadIdentity), nil
}

func (r *ReferenceManager) lookupSecret(ctx context.Context, namespace, name string) error {
		return r.kubeClient.CoreV1().Secrets(namespace).Get(ctx, name, kubernetesclient.DefaultGetOptions())
	}

	_, err := lookupResource(ctx, namespace, name, secretFromLister, secretFromClient)
	return err
}

func (r *ReferenceManager) lookupConfigMap(ctx context.Context, namespace, name string) error {
		return r.kubeClient.CoreV1().ConfigMaps(namespace).Get(ctx, name, kubernetesclient.DefaultGetOptions())
	}

	_, err := lookupResource(ctx, namespace, name, configMapFromLister, configMapFromClient)
	return err
}

func (r *ReferenceManager) lookupControllerDeployment(ctx context.Context, name string) error {
		return r.gardenCoreClient.CoreV1beta1().ControllerDeployments().Get(ctx, name, kubernetesclient.DefaultGetOptions())
	}

	_, err := lookupResource(ctx, "", name, deploymentFromLister, deploymentFromClient)
	return err
}

func (r *ReferenceManager) getAPIResource(groupVersion, kind string) (*metav1.APIResource, error) {
	}
	return nil
}

func (r *ReferenceManager) sanityCheckProviderSecret(ctx context.Context, namespace, name string, providerTypes []string) error {
	secret, err := r.kubeClient.CoreV1().Secrets(namespace).Get(ctx, name, kubernetesclient.DefaultGetOptions())
	if err != nil {
		return err
	}

	for _, providerType := range providerTypes {
		dummySecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: name,
				Namespace:    namespace,
				Annotations:  secret.Annotations,
				Labels:       gardenerutils.MergeStringMaps(secret.Labels, map[string]string{v1beta1constants.LabelShootProviderPrefix + providerType: "true"}),
			},
			Type: secret.Type,
			Data: secret.Data,
		}

		if _, err := r.kubeClient.CoreV1().Secrets(dummySecret.Namespace).Create(ctx, dummySecret, metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}}); err != nil {
			return fmt.Errorf("%s provider secret sanity check failed: %w", providerType, err)
		}
	}

	return nil
}
	. "github.com/onsi/gomega/gstruct"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/admission"

	"github.com/gardener/gardener/pkg/apis/core"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/gardener/pkg/apis/security"
	securityv1alpha1 "github.com/gardener/gardener/pkg/apis/security/v1alpha1"
	"github.com/gardener/gardener/pkg/apis/seedmanagement"
						Namespace: namespace,
					},
				},
				Provider: &core.SecretBindingProvider{
					Type: "test",
				},
			}
			secretBinding = gardencorev1beta1.SecretBinding{
				ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
					},
				},
				Provider: &gardencorev1beta1.SecretBindingProvider{
					Type: "test",
				},
			}
			securityCredentialsBindingRefSecret = security.CredentialsBinding{
				ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
					},
				},
				Provider: security.CredentialsBindingProvider{
					Type: "test",
				},
			}
			credentialsBindingRefSecret = securityv1alpha1.CredentialsBinding{
				ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
					},
				},
				Provider: securityv1alpha1.CredentialsBindingProvider{
					Type: "test",
				},
			}
			securityCredentialsBindingRefWorkloadIdentity = security.CredentialsBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:       workloadIdentityName,
					Namespace:  namespace,
				},
				Provider: security.CredentialsBindingProvider{Type: "wiprovider"},
				Quotas: []corev1.ObjectReference{
					{
						Name:      quotaName,

			err = gardencorev1beta1.Convert_core_Project_To_v1beta1_Project(&coreProject, &project, nil)
			Expect(err).To(Succeed())

			workloadIdentity.Spec.TargetSystem = securityv1alpha1.TargetSystem{Type: "wiprovider"}
		})

		It("should return nil because the resource is not BackupBucket and operation is delete", func() {
			attrs := admission.NewAttributesRecord(&controllerRegistration, nil, core.Kind("ControllerRegistration").WithVersion("version"), "", controllerRegistration.Name, core.Resource("controllerregistrations").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, nil)

			err := admissionHandler.Validate(context.TODO(), attrs, nil)

			Expect(err).NotTo(HaveOccurred())

			attrs = admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), "", controllerRegistration.Name, core.Resource("shoots").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, nil)

			err = admissionHandler.Validate(context.TODO(), attrs, nil)

			Expect(err).NotTo(HaveOccurred())
		})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&controllerRegistration, nil, core.Kind("ControllerRegistration").WithVersion("version"), "", controllerRegistration.Name, core.Resource("controllerregistrations").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&controllerRegistration, nil, core.Kind("ControllerRegistration").WithVersion("version"), "", controllerRegistration.Name, core.Resource("controllerregistrations").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
					return true, nil, errors.New("nope, out of luck")
				})

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("nope, out of luck"))
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				kubeClient.AddReactor("create", "secrets", func(_ testing.Action) (bool, runtime.Object, error) {
					return true, nil, nil
				})

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should reject because the sanity check fails", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota)).To(Succeed())
				kubeClient.AddReactor("get", "secrets", func(_ testing.Action) (bool, runtime.Object, error) {
					return true, &corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: secret.Namespace,
							Name:      secret.Name,
						},
					}, nil
				})

				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				kubeClient.AddReactor("create", "secrets", func(_ testing.Action) (bool, runtime.Object, error) {
					return true, nil, fmt.Errorf("sanity check failed")
				})

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(ContainSubstring("test provider secret sanity check failed: sanity check failed")))
			})

			It("should reject because the referenced secret does not exist", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota)).To(Succeed())
				kubeClient.AddReactor("get", "secrets", func(_ testing.Action) (bool, runtime.Object, error) {
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})

			It("should reject because provider types do not match with shoot type", func() {
				coreSecretBinding.Provider.Type = "another-provider"
				coreSecretBinding.Quotas = nil
				shoot.Spec.Provider.Type = "local"

				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(&shoot)).To(Succeed())

				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(ContainSubstring(`SecretBinding is referenced by shoot "shoot-1", but provider types ([another-provider]) do not match with the shoot provider type "local"`)))
			})
		})

		Context("tests for CredentialsBinding objects referencing Secret", func() {
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				kubeClient.AddReactor("create", "secrets", func(_ testing.Action) (bool, runtime.Object, error) {
					return true, nil, nil
				})

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should reject because the sanity check fails", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota)).To(Succeed())
				kubeClient.AddReactor("get", "secrets", func(_ testing.Action) (bool, runtime.Object, error) {
					return true, &corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: secret.Namespace,
							Name:      secret.Name,
						},
					}, nil
				})

				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				kubeClient.AddReactor("create", "secrets", func(_ testing.Action) (bool, runtime.Object, error) {
					return true, nil, fmt.Errorf("sanity check failed")
				})

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(ContainSubstring("test provider secret sanity check failed: sanity check failed")))
			})

			It("should reject because the referenced secret does not exist", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota)).To(Succeed())
				kubeClient.AddReactor("get", "secrets", func(_ testing.Action) (bool, runtime.Object, error) {
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})

			It("should reject because provider types do not match with shoot type", func() {
				securityCredentialsBindingRefSecret.Provider.Type = "another-provider"
				securityCredentialsBindingRefSecret.Quotas = nil
				shoot.Spec.Provider.Type = "local"

				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(&shoot)).To(Succeed())

				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(ContainSubstring(`CredentialsBinding is referenced by shoot "shoot-1", but provider types ([another-provider]) do not match with the shoot provider type "local"`)))
			})
		})

		Context("tests for CredentialsBinding objects referencing WorkloadIdentity", func() {
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
			It("should accept because all referenced objects have been found (workloadidentity looked up live)", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota)).To(Succeed())
				gardenSecurityClient.AddReactor("get", "workloadidentities", func(_ testing.Action) (bool, runtime.Object, error) {
					return true, &workloadIdentity, nil
				})

				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should reject because the provider type does not match in WorkloadIdentity and CredentialsBinding", func() {
				workloadIdentity.Spec.TargetSystem.Type = "foo"
				Expect(gardenSecurityInformerFactory.Security().V1alpha1().WorkloadIdentities().Informer().GetStore().Add(&workloadIdentity)).To(Succeed())
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota)).To(Succeed())

				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(MatchError(ContainSubstring("does not match with WorkloadIdentity provider type")))
			})

			It("should reject because the referenced workload identity does not exist", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota)).To(Succeed())
				gardenSecurityClient.AddReactor("get", "workloadidentities", func(_ testing.Action) (bool, runtime.Object, error) {
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				quotaRefList = append(quotaRefList, quota2Ref)
				securityCredentialsBindingRefWorkloadIdentity.Quotas = quotaRefList

				Expect(gardenSecurityInformerFactory.Security().V1alpha1().WorkloadIdentities().Informer().GetStore().Add(&workloadIdentity)).To(Succeed())
				Expect(kubeInformerFactory.Core().V1().Secrets().Informer().GetStore().Add(&secret)).To(Succeed())
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota)).To(Succeed())
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota2)).To(Succeed())
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("tests for Shoot objects", func() {
			It("should accept because all referenced objects have been found", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().CloudProfiles().Informer().GetStore().Add(&cloudProfile)).To(Succeed())
				Expect(gardenCoreInformerFactory.Core().V1beta1().Seeds().Informer().GetStore().Add(&seed)).To(Succeed())
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				coreShoot.Status.TechnicalID = "should-never-change"
				attrs := admission.NewAttributesRecord(&coreShoot, oldShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				It("should reject because the referenced cloud profile does not exist (create)", func() {
					attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

					err := admissionHandler.Validate(context.TODO(), attrs, nil)

					Expect(err).To(HaveOccurred())
				})

					attrs := admission.NewAttributesRecord(&coreShoot, oldShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

					err := admissionHandler.Validate(context.TODO(), attrs, nil)

					Expect(err).To(HaveOccurred())
				})

					attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

					err := admissionHandler.Validate(context.TODO(), attrs, nil)

					Expect(err).To(HaveOccurred())
				})

					attrs := admission.NewAttributesRecord(&coreShoot, oldShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

					err := admissionHandler.Validate(context.TODO(), attrs, nil)

					Expect(err).To(HaveOccurred())
				})

				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				It("should reject because the referenced exposure class does not exists", func() {
					attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

					err := admissionHandler.Validate(context.TODO(), attrs, nil)
					Expect(err).To(HaveOccurred())
				})

					Expect(gardenCoreInformerFactory.Core().V1beta1().ExposureClasses().Informer().GetStore().Add(&exposureClass)).To(Succeed())
					attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

					err := admissionHandler.Validate(context.TODO(), attrs, nil)
					Expect(err).To(HaveOccurred())
				})
			})

				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(MatchError("shoots.core.gardener.cloud \"shoot-1\" is forbidden: cannot reference a resource you are not allowed to read"))
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&coreShoot, oldShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(MatchError("shoots.core.gardener.cloud \"shoot-1\" is forbidden: cannot reference a resource you are not allowed to read"))
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&coreShoot, oldShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(Not(HaveOccurred()))
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(MatchError(ContainSubstring("failed to resolve resource reference")))
			})
					user := &user.DefaultInfo{Name: allowedUser}
					attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

					Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(ContainSubstring(expectedErrorMessage)))
				})

				It("should reject because the referenced "+description+" does not exist (update)", func() {
					user := &user.DefaultInfo{Name: allowedUser}
					attrs := admission.NewAttributesRecord(&coreShoot, oldShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, user)

					Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(ContainSubstring(expectedErrorMessage)))
				})

				It("should pass because the referenced "+description+" does not exist but shoot has deletion timestamp", func() {
					user := &user.DefaultInfo{Name: allowedUser}
					attrs := admission.NewAttributesRecord(&coreShoot, oldShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, user)

					Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
				})

				It("should pass because the referenced "+description+" exists", func() {
					user := &user.DefaultInfo{Name: allowedUser}
					attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

					Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
				})
			}

				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&coreSeed, oldSeed, core.Kind("Seed").WithVersion("version"), "", coreSeed.Name, core.Resource("seeds").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(MatchError("seeds.core.gardener.cloud \"seed-1\" is forbidden: cannot reference a resource you are not allowed to read"))
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&coreSeed, oldSeed, core.Kind("Seed").WithVersion("version"), "", coreSeed.Name, core.Resource("seeds").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(Not(HaveOccurred()))
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSeed, nil, core.Kind("Seed").WithVersion("version"), "", coreSeed.Name, core.Resource("seeds").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(MatchError(ContainSubstring("failed to resolve resource reference")))
			})
			It("should reject if the referred Seed is not found", func() {
				attrs := admission.NewAttributesRecord(&coreBackupBucket, nil, core.Kind("BackupBucket").WithVersion("version"), "", coreBackupBucket.Name, core.Resource("backupBuckets").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(BeForbiddenError())
				Expect(err).To(MatchError(ContainSubstring("backupBuckets.core.gardener.cloud %q is forbidden: seed.core.gardener.cloud %q not found", coreBackupBucket.Name, seed.Name)))

				attrs := admission.NewAttributesRecord(&coreBackupBucket, nil, core.Kind("BackupBucket").WithVersion("version"), "", coreBackupBucket.Name, core.Resource("backupBuckets").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(BeForbiddenError())
				Expect(err).To(MatchError(ContainSubstring("secret not found")))

				attrs := admission.NewAttributesRecord(&coreBackupBucket, nil, core.Kind("BackupBucket").WithVersion("version"), "", coreBackupBucket.Name, core.Resource("backupBuckets").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&coreBackupBucket, nil, core.Kind("BackupBucket").WithVersion("version"), "", coreBackupBucket.Name, core.Resource("backupBuckets").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(nil, nil, core.Kind("BackupBucket").WithVersion("version"), "", coreBackupBucket.Name, core.Resource("backupBuckets").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(nil, nil, core.Kind("BackupBucket").WithVersion("version"), "", backupBucket.Name, core.Resource("backupBuckets").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(BeForbiddenError())
				Expect(err).To(MatchError(ContainSubstring("backupBuckets.core.gardener.cloud %q is forbidden: cannot delete BackupBucket because BackupEntries are still referencing it, backupEntryNames: %s/%s,%s/%s", backupBucket.Name, backupEntry.Namespace, backupEntry.Name, backupEntry2.Namespace, backupEntry2.Name)))

				attrs := admission.NewAttributesRecord(nil, nil, core.Kind("BackupBucket").WithVersion("version"), "", coreBackupBucket.Name, core.Resource("backupBuckets").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(nil, nil, core.Kind("BackupBucket").WithVersion("version"), "", "", core.Resource("backupBuckets").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(BeForbiddenError())
				Expect(err).To(MatchError(ContainSubstring("backupBuckets.core.gardener.cloud %q is forbidden: cannot delete BackupBucket because BackupEntries are still referencing it, backupEntryNames: %s/%s,%s/%s", backupBucket2.Name, backupEntry.Namespace, backupEntry.Name, backupEntry2.Namespace, backupEntry2.Name)))

				attrs := admission.NewAttributesRecord(nil, nil, core.Kind("BackupBucket").WithVersion("version"), "", "", core.Resource("backupBuckets").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
			It("should reject if the referred Seed is not found", func() {
				attrs := admission.NewAttributesRecord(&coreBackupEntry, nil, core.Kind("BackupEntry").WithVersion("version"), coreBackupEntry.Namespace, coreBackupEntry.Name, core.Resource("backupEntries").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(BeForbiddenError())
				Expect(err).To(MatchError(ContainSubstring("backupEntries.core.gardener.cloud %q is forbidden: seed.core.gardener.cloud %q not found", coreBackupEntry.Name, seed.Name)))
				Expect(gardenCoreInformerFactory.Core().V1beta1().Seeds().Informer().GetStore().Add(&seed)).To(Succeed())
				attrs := admission.NewAttributesRecord(&coreBackupEntry, nil, core.Kind("BackupEntry").WithVersion("version"), coreBackupEntry.Namespace, coreBackupEntry.Name, core.Resource("backupEntries").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(BeForbiddenError())
				Expect(err).To(MatchError(ContainSubstring("backupEntries.core.gardener.cloud %q is forbidden: backupbucket.core.gardener.cloud %q not found", coreBackupEntry.Name, coreBackupBucket.Name)))
				Expect(gardenCoreInformerFactory.Core().V1beta1().BackupBuckets().Informer().GetStore().Add(&backupBucket)).To(Succeed())
				attrs := admission.NewAttributesRecord(&coreBackupEntry, nil, core.Kind("BackupEntry").WithVersion("version"), coreBackupEntry.Namespace, coreBackupEntry.Name, core.Resource("backupEntries").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("tests for Project objects", func() {
			It("should allow specifying a namespace which is not in use (create)", func() {
				project.Spec.Namespace = ptr.To("garden-foo")
				projectCopy := project.DeepCopy()
				coreProject.Spec.Namespace = ptr.To("garden-foo")
				attrs := admission.NewAttributesRecord(&coreProject, nil, core.Kind("Project").WithVersion("version"), coreProject.Namespace, coreProject.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(Not(HaveOccurred()))
			})
				coreProject.Spec.Namespace = ptr.To("garden-foo")
				attrs := admission.NewAttributesRecord(&coreProject, coreProjectOld, core.Kind("Project").WithVersion("version"), coreProject.Namespace, coreProject.Name, core.Resource("projects").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(Not(HaveOccurred()))
			})

				attrs := admission.NewAttributesRecord(&coreProject, nil, core.Kind("Project").WithVersion("version"), coreProject.Namespace, coreProject.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(Not(HaveOccurred()))
			})
				coreProject.Spec.Namespace = ptr.To("garden-foo")
				attrs := admission.NewAttributesRecord(&coreProject, nil, core.Kind("Project").WithVersion("version"), coreProject.Namespace, coreProject.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"ErrStatus": MatchFields(IgnoreExtras, Fields{
				coreProject.Spec.Namespace = ptr.To("garden-foo")
				attrs := admission.NewAttributesRecord(&coreProject, &coreProjectOld, core.Kind("Project").WithVersion("version"), coreProject.Namespace, coreProject.Name, core.Resource("projects").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"ErrStatus": MatchFields(IgnoreExtras, Fields{

				attrs := admission.NewAttributesRecord(&cloudProfile, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("1.24.1"))

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to delete Kubernetes version"),
					ContainSubstring("1.24.1"),
					ContainSubstring("still in use by NamespacedCloudProfile"),

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to delete Kubernetes version"),
					ContainSubstring("1.24.1"),
					ContainSubstring("still in use by shoot '/shoot-Two'"),

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should accept removal of kubernetes versions that are used by shoots using another unrelated NamespacedCloudProfile of same name", func() {

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should accept removal of kubernetes version that is still in use by a shoot that is being deleted", func() {

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&cloudProfile, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("1.17.2"))

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("1.17.2"))
				Expect(err.Error()).To(ContainSubstring(s.Spec.Provider.Workers[1].Name))

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("1.17.2"))

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should fail for removal of a machine version that is used by a NamespacedCloudProfile", func() {

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to delete MachineImage version"),
					ContainSubstring("1.16.0"),
					ContainSubstring("still in use by NamespacedCloudProfile"),

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to delete MachineImage version"),
					ContainSubstring("1.16.0"),
					ContainSubstring("still in use by NamespacedCloudProfile"),

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to delete MachineImage \"coreos\""),
					ContainSubstring("still in use by NamespacedCloudProfile"),
				)))

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to add MachineImage \"gardenlinux\""),
					ContainSubstring("already defined by NamespacedCloudProfile \"project-123/profile-42\""),
				)))
				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(shootTwo)).To(Succeed())
				attrs := admission.NewAttributesRecord(cloudProfile, oldCloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(ctx, attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(shootTwo)).To(Succeed())
				attrs := admission.NewAttributesRecord(cloudProfile, oldCloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(ctx, attrs, nil)

				Expect(err).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"ErrStatus": MatchFields(IgnoreExtras, Fields{
				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(shootOne)).To(Succeed())
				attrs := admission.NewAttributesRecord(cloudProfile, oldCloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(ctx, attrs, nil)

				Expect(err).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"ErrStatus": MatchFields(IgnoreExtras, Fields{

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should succeed for the complete Kubernetes section being removed without usages", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should succeed if a used and already extended kubernetes version expiration is changed to another value still in the future", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should succeed if a used and extended kubernetes version already expired is not modified", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should succeed if an extended and used Kubernetes version is removed with the base version still being valid", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should fail if an extended and used Kubernetes version is being removed with the base version being already expired", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)
				Expect(err).To(MatchError(And(
					ContainSubstring("unable to delete Kubernetes version"),
					ContainSubstring("1.29.0"),

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to delete Kubernetes version"),
					ContainSubstring("1.29.0"),
					ContainSubstring("still in use by shoot"),

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})
		})


				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should succeed if a used and extended MachineImage version already expired is not modified", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should succeed if an extended and used MachineImage version is removed with the base version still being valid", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should fail if an extended and used MachineImage version is being removed with the base version being already expired", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)
				Expect(err).To(MatchError(And(
					ContainSubstring("unable to delete Machine image version"),
					ContainSubstring("1.17.3"),

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)
				Expect(err).To(MatchError(And(
					ContainSubstring("unable to delete Machine image version"),
					ContainSubstring("'coreos/1.1.2'"),

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(context.TODO(), attrs, nil)
				Expect(err).To(MatchError(And(
					ContainSubstring("unable to delete Machine image version"),
					ContainSubstring("'custom-namespaced-image/1.1.2'"),

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to delete Machine image version"),
					ContainSubstring("1.17.3"),
					ContainSubstring("still in use by shoot"),

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should succeed if a new but unused MachineImage version is removed", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})
		})

				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(shootTwo)).To(Succeed())
				attrs := admission.NewAttributesRecord(namespacedCloudProfile, oldNamespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespacedCloudProfile.Namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(ctx, attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(shootTwo)).To(Succeed())
				attrs := admission.NewAttributesRecord(namespacedCloudProfile, oldNamespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespacedCloudProfile.Namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Validate(ctx, attrs, nil)

				Expect(err).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"ErrStatus": MatchFields(IgnoreExtras, Fields{
			It("should accept because there is no managed seed with the same name", func() {
				attrs := admission.NewAttributesRecord(gardenlet, nil, seedmanagement.Kind("Gardenlet").WithVersion("version"), gardenlet.Namespace, gardenlet.Name, seedmanagement.Resource("gardenlets").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, &user.DefaultInfo{Name: allowedUser})

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should forbid because there is a managed seed with the same name", func() {

				attrs := admission.NewAttributesRecord(gardenlet, nil, seedmanagement.Kind("Gardenlet").WithVersion("version"), gardenlet.Namespace, gardenlet.Name, seedmanagement.Resource("gardenlets").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, &user.DefaultInfo{Name: allowedUser})

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(ContainSubstring("there is already a ManagedSeed object with the same name")))
			})
		})

			It("should accept because there is no gardenlet with the same name", func() {
				attrs := admission.NewAttributesRecord(managedSeed, nil, seedmanagement.Kind("ManagedSeed").WithVersion("version"), gardenlet.Namespace, gardenlet.Name, seedmanagement.Resource("gardenlets").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, &user.DefaultInfo{Name: allowedUser})

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should forbid because there is a gardenlet with the same name", func() {

				attrs := admission.NewAttributesRecord(managedSeed, nil, seedmanagement.Kind("ManagedSeed").WithVersion("version"), managedSeed.Namespace, managedSeed.Name, seedmanagement.Resource("managedseeds").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, &user.DefaultInfo{Name: allowedUser})

				Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(ContainSubstring("there is already a Gardenlet object with the same name")))
			})
		})
	})
	PluginNameExtensionLabels = "ExtensionLabels"
	// PluginNameExtensionValidator is the name of the ExtensionValidator admission plugin.
	PluginNameExtensionValidator = "ExtensionValidator"
	// PluginNameFinalizerRemoval is the name of the FinalizerRemoval admission plugin.
	PluginNameFinalizerRemoval = "FinalizerRemoval"
	// PluginNameResourceReferenceManager is the name of the ResourceReferenceManager admission plugin.
	PluginNameResourceReferenceManager = "ResourceReferenceManager"
	// PluginNameManagedSeedShoot is the name of the ManagedSeedShoot admission plugin.
		PluginNameNamespacedCloudProfileValidator,   // NamespacedCloudProfileValidator
		PluginNameProjectValidator,                  // ProjectValidator
		PluginNameDeletionConfirmation,              // DeletionConfirmation
		PluginNameFinalizerRemoval,                  // FinalizerRemoval
		PluginNameOpenIDConnectPreset,               // OpenIDConnectPreset
		PluginNameClusterOpenIDConnectPreset,        // ClusterOpenIDConnectPreset
		PluginNameCustomVerbAuthorizer,              // CustomVerbAuthorizer
		PluginNameNamespacedCloudProfileValidator, // NamespacedCloudProfileValidator
		PluginNameProjectValidator,                // ProjectValidator
		PluginNameDeletionConfirmation,            // DeletionConfirmation
		PluginNameFinalizerRemoval,                // FinalizerRemoval
		PluginNameOpenIDConnectPreset,             // OpenIDConnectPreset
		PluginNameClusterOpenIDConnectPreset,      // ClusterOpenIDConnectPreset
		PluginNameCustomVerbAuthorizer,            // CustomVerbAuthorizer
	"context"
	"fmt"
	"io"
	"slices"
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/admission"

	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	gardenerutils "github.com/gardener/gardener/pkg/utils/gardener"
	plugin "github.com/gardener/gardener/plugin/pkg"
	"github.com/gardener/gardener/plugin/pkg/utils"
)

// Register registers a plugin.
// New creates a new handler admission plugin.
func New() (*handler, error) {
	return &handler{
		Handler: admission.NewHandler(admission.Create, admission.Update),
	}, nil
}

var _ admission.MutationInterface = &handler{}

func (v *handler) Admit(_ context.Context, a admission.Attributes, _ admission.ObjectInterfaces) error {
	// Ignore all kinds other than Project
	if a.GetKind().GroupKind() != gardencore.Kind("Project") {
		return nil
		return apierrors.NewBadRequest("could not convert object to Project")
	}

	// TODO: Remove this check in favor of static validation in a future release, see https://github.com/gardener/gardener/pull/4228.
	if project.Spec.Namespace != nil && *project.Spec.Namespace != v1beta1constants.GardenNamespace && !strings.HasPrefix(*project.Spec.Namespace, gardenerutils.ProjectNamespacePrefix) {
		return admission.NewForbidden(a, fmt.Errorf(".spec.namespace must start with %s", gardenerutils.ProjectNamespacePrefix))
	}

	if utils.SkipVerification(a.GetOperation(), project.ObjectMeta) {
		return nil
	}

	if a.GetOperation() == admission.Create {
		ensureProjectOwner(project, a.GetUserInfo().GetName())
	}

	ensureOwnerIsMember(project)

	return nil
}

func ensureProjectOwner(project *gardencore.Project, userName string) {
	// Set createdBy field in Project
	project.Spec.CreatedBy = &rbacv1.Subject{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     rbacv1.UserKind,
		Name:     userName,
	}

	if project.Spec.Owner == nil {
		project.Spec.Owner = func() *rbacv1.Subject {
			for _, member := range project.Spec.Members {
				for _, role := range member.Roles {
					if role == gardencore.ProjectMemberOwner {
						return member.Subject.DeepCopy()
					}
				}
			}
			return project.Spec.CreatedBy
		}()
	}
}

func ensureOwnerIsMember(project *gardencore.Project) {
	if project.Spec.Owner == nil {
		return
	}

	ownerIsMember := slices.ContainsFunc(project.Spec.Members, func(member gardencore.ProjectMember) bool {
		return member.Subject == *project.Spec.Owner
	})

	if !ownerIsMember {
		project.Spec.Members = append(project.Spec.Members, gardencore.ProjectMember{
			Subject: *project.Spec.Owner,
			Roles: []string{
				gardencore.ProjectMemberAdmin,
				gardencore.ProjectMemberOwner,
			},
		})
	}
}

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/utils/ptr"

	"github.com/gardener/gardener/pkg/apis/core"
		var (
			err              error
			project          core.Project
			admissionHandler admission.MutationInterface
			attrs            admission.Attributes

			namespaceName = "garden-my-project"
			projectName   = "my-project"
					Namespace: namespaceName,
				},
			}

			userInfo user.Info
		)

		BeforeEach(func() {
			Expect(err).NotTo(HaveOccurred())

			project = projectBase

			userInfo = &user.DefaultInfo{Name: "foo"}
		})

		When("project is created", func() {
			BeforeEach(func() {
				attrs = admission.NewAttributesRecord(&project, nil, core.Kind("Project").WithVersion("version"), "", project.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
			})

			It("should allow creating the project (namespace nil)", func() {
				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should allow creating the project(namespace non-nil)", func() {
				project.Spec.Namespace = &namespaceName

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should allow creating the project (namespace is 'garden')", func() {
				project.Spec.Namespace = ptr.To(v1beta1constants.GardenNamespace)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should prevent creating the project because namespace prefix is missing", func() {
				project.Spec.Namespace = ptr.To("foo")

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(MatchError(ContainSubstring(".spec.namespace must start with garden-")))
			})

			It("should maintain createdBy and project owner", func() {
				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())

				Expect(project.Spec.CreatedBy).To(Equal(&rbacv1.Subject{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     "User",
					Name:     userInfo.GetName(),
				}))

				Expect(project.Spec.Owner).To(Equal(&rbacv1.Subject{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     "User",
					Name:     userInfo.GetName(),
				}))

				Expect(project.Spec.Members).To(ConsistOf(core.ProjectMember{
					Subject: rbacv1.Subject{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "User",
						Name:     userInfo.GetName(),
					},
					Roles: []string{
						core.ProjectMemberAdmin,
						core.ProjectMemberOwner,
					},
				}))
			})

			It("should not overwrite project owner", func() {
				project.Spec.Owner = &rbacv1.Subject{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     "User",
					Name:     "bar",
				}

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())

				Expect(project.Spec.Owner).To(Equal(&rbacv1.Subject{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     "User",
					Name:     "bar",
				}))
			})
		})

		When("project is updated", func() {
			BeforeEach(func() {
				attrs = admission.NewAttributesRecord(&project, nil, core.Kind("Project").WithVersion("version"), "", project.Name, core.Resource("projects").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, userInfo)
			})

			It("should add project owner to members", func() {
				projectOwner := core.ProjectMember{
					Subject: rbacv1.Subject{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "User",
						Name:     "foo",
					},
					Roles: []string{
						core.ProjectMemberAdmin,
						core.ProjectMemberOwner,
					},
				}

				projectMemberBar := core.ProjectMember{
					Subject: rbacv1.Subject{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "User",
						Name:     "bar",
					},
					Roles: []string{
						core.ProjectMemberViewer,
					},
				}

				project.Spec.Owner = &projectOwner.Subject
				project.Spec.Members = []core.ProjectMember{projectMemberBar}

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())

				Expect(project.Spec.Members).To(ConsistOf(projectMemberBar, projectOwner))
			})

			It("should not re-add owner as member", func() {
				projectOwner := core.ProjectMember{
					Subject: rbacv1.Subject{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "User",
						Name:     "foo",
					},
					Roles: []string{
						core.ProjectMemberAdmin,
						core.ProjectMemberOwner,
					},
				}

				project.Spec.Owner = &projectOwner.Subject
				project.Spec.Members = []core.ProjectMember{projectOwner}

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())

				Expect(project.Spec.Members).To(ConsistOf(projectOwner))
			})
		})
	})

	})

	Describe("#New", func() {
		It("should handle CREATE and UPDATE operations", func() {
			dr, err := New()
			Expect(err).ToNot(HaveOccurred())
			Expect(dr.Handles(admission.Create)).To(BeTrue())
			Expect(dr.Handles(admission.Update)).To(BeTrue())
			Expect(dr.Handles(admission.Connect)).To(BeFalse())
			Expect(dr.Handles(admission.Delete)).To(BeFalse())
		})
		}
	}

	if a.GetOperation() == admission.Create {
		addCreatedByAnnotation(shoot, a.GetUserInfo().GetName())

		if len(ptr.Deref(shoot.Spec.CloudProfileName, "")) > 0 && shoot.Spec.CloudProfile != nil {
			return fmt.Errorf("new shoot can only specify either cloudProfileName or cloudProfile reference")
		}
	}

	cloudProfileSpec, err := admissionutils.GetCloudProfileSpec(v.cloudProfileLister, v.namespacedCloudProfileLister, shoot)
	if err != nil {
		return apierrors.NewInternalError(fmt.Errorf("could not find referenced cloud profile: %+v", err.Error()))
	}

	if err := admissionutils.ValidateCloudProfileChanges(v.cloudProfileLister, v.namespacedCloudProfileLister, shoot, oldShoot); err != nil {
		return err
	}
		}
	}

	return nil
}


	return allErrs
}

func addCreatedByAnnotation(shoot *core.Shoot, userName string) {
	annotations := shoot.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[v1beta1constants.GardenCreatedBy] = userName
	shoot.Annotations = annotations
}
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
			})
		})

		Context("shoot creation", func() {
			BeforeEach(func() {
				Expect(coreInformerFactory.Core().V1beta1().Projects().Informer().GetStore().Add(&project)).To(Succeed())
				Expect(coreInformerFactory.Core().V1beta1().CloudProfiles().Informer().GetStore().Add(&cloudProfile)).To(Succeed())
				Expect(coreInformerFactory.Core().V1beta1().Seeds().Informer().GetStore().Add(&seed)).To(Succeed())
				Expect(coreInformerFactory.Core().V1beta1().SecretBindings().Informer().GetStore().Add(&secretBinding)).To(Succeed())
				Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())
			})

			Context("with generate name", func() {
				BeforeEach(func() {
					shoot.ObjectMeta = metav1.ObjectMeta{
						GenerateName: "demo-",
						Namespace:    namespaceName,
					}
				})

				It("should admit Shoot resources", func() {
					authorizeAttributes.Name = shoot.Name

					attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
					err := admissionHandler.Admit(ctx, attrs, nil)

					Expect(err).NotTo(HaveOccurred())
				})

				It("should reject Shoot resources with not fulfilling the length constraints", func() {
					tooLongName := "too-long-namespace"
					project.ObjectMeta = metav1.ObjectMeta{
						Name: tooLongName,
					}
					shoot.ObjectMeta = metav1.ObjectMeta{
						GenerateName: "too-long-name",
						Namespace:    namespaceName,
					}

					Expect(coreInformerFactory.Core().V1beta1().Projects().Informer().GetStore().Add(&project)).To(Succeed())

					authorizeAttributes.Name = shoot.Name

					attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
					err := admissionHandler.Admit(ctx, attrs, nil)

					Expect(err).To(BeInvalidError())
					Expect(err.Error()).To(ContainSubstring("name must not exceed"))
				})
			})

			It("should add the created-by annotation", func() {
				Expect(shoot.Annotations).NotTo(HaveKeyWithValue(v1beta1constants.GardenCreatedBy, userInfo.Name))

				attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
				Expect(admissionHandler.Admit(ctx, attrs, nil)).NotTo(HaveOccurred())

				Expect(shoot.Annotations).To(HaveKeyWithValue(v1beta1constants.GardenCreatedBy, userInfo.Name))
			})
		})

				shoot.Spec.SeedName = nil
				shoot.Spec.AccessRestrictions = []core.AccessRestrictionWithOptions{{AccessRestriction: core.AccessRestriction{Name: "foo"}}}

				attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.UpdateOptions{}, false, userInfo)
				err := admissionHandler.Admit(ctx, attrs, nil)

				Expect(err).To(BeForbiddenError())

				shoot.Spec.AccessRestrictions = []core.AccessRestrictionWithOptions{{AccessRestriction: core.AccessRestriction{Name: "foo"}}}

				attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.UpdateOptions{}, false, userInfo)
				err := admissionHandler.Admit(ctx, attrs, nil)

				Expect(err).To(BeForbiddenError())

				shoot.Spec.AccessRestrictions = []core.AccessRestrictionWithOptions{{AccessRestriction: core.AccessRestriction{Name: "foo"}}}

				attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.UpdateOptions{}, false, userInfo)
				err := admissionHandler.Admit(ctx, attrs, nil)

				Expect(err).NotTo(HaveOccurred())
						})

						It("should allow scheduling non-HA shoot", func() {
							attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
							Expect(admissionHandler.Admit(ctx, attrs, nil)).To(Succeed())
						})

							shoot.Annotations = make(map[string]string)
							shoot.Spec.ControlPlane = &core.ControlPlane{HighAvailability: &core.HighAvailability{FailureTolerance: core.FailureTolerance{Type: core.FailureToleranceTypeNode}}}

							attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
							Expect(admissionHandler.Admit(ctx, attrs, nil)).To(Succeed())
						})

							shoot.Annotations = make(map[string]string)
							shoot.Spec.ControlPlane = &core.ControlPlane{HighAvailability: &core.HighAvailability{FailureTolerance: core.FailureTolerance{Type: core.FailureToleranceTypeZone}}}

							attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
							Expect(admissionHandler.Admit(ctx, attrs, nil)).To(BeForbiddenError())
						})
					})
						})

						It("should allow scheduling non-HA shoot", func() {
							attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
							Expect(admissionHandler.Admit(ctx, attrs, nil)).To(Succeed())
						})

							shoot.Annotations = make(map[string]string)
							shoot.Spec.ControlPlane = &core.ControlPlane{HighAvailability: &core.HighAvailability{FailureTolerance: core.FailureTolerance{Type: core.FailureToleranceTypeNode}}}

							attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
							Expect(admissionHandler.Admit(ctx, attrs, nil)).To(Succeed())
						})

							shoot.Annotations = make(map[string]string)
							shoot.Spec.ControlPlane = &core.ControlPlane{HighAvailability: &core.HighAvailability{FailureTolerance: core.FailureTolerance{Type: core.FailureToleranceTypeZone}}}

							attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
							Expect(admissionHandler.Admit(ctx, attrs, nil)).To(Succeed())
						})
					})
						Expect(coreInformerFactory.Core().V1beta1().SecretBindings().Informer().GetStore().Add(&secretBinding)).To(Succeed())
						Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())

						attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
						err := admissionHandler.Admit(ctx, attrs, nil)

						Expect(err).NotTo(HaveOccurred())
						Expect(coreInformerFactory.Core().V1beta1().SecretBindings().Informer().GetStore().Add(&secretBinding)).To(Succeed())
						Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())

						attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
						err := admissionHandler.Admit(ctx, attrs, nil)

						Expect(err).To(BeForbiddenError())
						Expect(coreInformerFactory.Core().V1beta1().SecretBindings().Informer().GetStore().Add(&secretBinding)).To(Succeed())
						Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())

						attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
						err := admissionHandler.Admit(ctx, attrs, nil)

						Expect(err).To(BeForbiddenError())
						Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())
						Expect(kubeInformerFactory.Core().V1().Secrets().Informer().GetStore().Add(&secret)).To(Succeed())

						attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
						err := admissionHandler.Admit(ctx, attrs, nil)

						Expect(err).To(BeForbiddenError())
						Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())
						Expect(kubeInformerFactory.Core().V1().Secrets().Informer().GetStore().Add(&secret)).To(Succeed())

						attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
						err := admissionHandler.Admit(ctx, attrs, nil)

						Expect(err).NotTo(HaveOccurred())
					Expect(coreInformerFactory.Core().V1beta1().SecretBindings().Informer().GetStore().Add(&secretBinding)).To(Succeed())
					Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())

					attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
					err := admissionHandler.Admit(ctx, attrs, nil)

					Expect(err).To(errorMatcher)

	It("should be able to manipulate resource from security.gardener.cloud/v1alpha1", func() {
		credentialsBinding := &securityv1alpha1.CredentialsBinding{ObjectMeta: metav1.ObjectMeta{GenerateName: "test-", Namespace: testNamespace.Name}}
		Expect(testClient.Create(ctx, credentialsBinding)).To(MatchError(MatchRegexp("CredentialsBinding.security.gardener.cloud \"test-.+\" is invalid")))
	})
})
