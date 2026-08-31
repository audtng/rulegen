package main

	extensionscmdcontroller "github.com/gardener/gardener/extensions/pkg/controller/cmd"
	"github.com/gardener/gardener/extensions/pkg/util"
	extensionscmdwebhook "github.com/gardener/gardener/extensions/pkg/webhook/cmd"
	"github.com/gardener/gardener/pkg/apis/core/install"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	gardenerhealthz "github.com/gardener/gardener/pkg/healthz"
	admissioncmd "github.com/gardener/gardener/pkg/provider-local/admission/cmd"
	localinstall "github.com/gardener/gardener/pkg/provider-local/apis/local/install"
				return fmt.Errorf("could not instantiate manager: %w", err)
			}

			install.Install(mgr.GetScheme())

			if err := localinstall.AddToScheme(mgr.GetScheme()); err != nil {
				return fmt.Errorf("could not update manager scheme: %w", err)
	"github.com/gardener/gardener/pkg/admissioncontroller/webhook/admission/internaldomainsecret"
	"github.com/gardener/gardener/pkg/admissioncontroller/webhook/admission/kubeconfigsecret"
	"github.com/gardener/gardener/pkg/admissioncontroller/webhook/admission/namespacedeletion"
	"github.com/gardener/gardener/pkg/admissioncontroller/webhook/admission/resourcesize"
	"github.com/gardener/gardener/pkg/admissioncontroller/webhook/admission/seedrestriction"
	"github.com/gardener/gardener/pkg/admissioncontroller/webhook/admission/shootkubeconfigsecretref"
		return fmt.Errorf("failed adding %s webhook handler: %w", namespacedeletion.HandlerName, err)
	}

	if err := (&resourcesize.Handler{
		Logger: mgr.GetLogger().WithName("webhook").WithName(resourcesize.HandlerName),
		Config: cfg.Server.ResourceAdmissionConfiguration,

// GetSecretBindingTypes returns the SecretBinding provider types.
func GetSecretBindingTypes(secretBinding *core.SecretBinding) []string {
	return strings.Split(secretBinding.Provider.Type, ",")
}
			Expect(actual).To(Equal(expected))
		},

		Entry("with single-value provider type", &core.SecretBinding{Provider: &core.SecretBindingProvider{Type: "foo"}}, []string{"foo"}),
		Entry("with multi-value provider type", &core.SecretBinding{Provider: &core.SecretBindingProvider{Type: "foo,bar,baz"}}, []string{"foo", "bar", "baz"}),
	)

// GetSecretBindingTypes returns the SecretBinding provider types.
func GetSecretBindingTypes(secretBinding *gardencorev1beta1.SecretBinding) []string {
	return strings.Split(secretBinding.Provider.Type, ",")
}
			Expect(actual).To(Equal(expected))
		},

		Entry("with single-value provider type", &gardencorev1beta1.SecretBinding{Provider: &gardencorev1beta1.SecretBindingProvider{Type: "foo"}}, []string{"foo"}),
		Entry("with multi-value provider type", &gardencorev1beta1.SecretBinding{Provider: &gardencorev1beta1.SecretBindingProvider{Type: "foo,bar,baz"}}, []string{"foo", "bar", "baz"}),
	)
	"github.com/gardener/gardener/plugin/pkg/global/deletionconfirmation"
	"github.com/gardener/gardener/plugin/pkg/global/extensionlabels"
	"github.com/gardener/gardener/plugin/pkg/global/extensionvalidation"
	"github.com/gardener/gardener/plugin/pkg/global/resourcereferencemanager"
	managedseedshoot "github.com/gardener/gardener/plugin/pkg/managedseed/shoot"
	managedseedvalidator "github.com/gardener/gardener/plugin/pkg/managedseed/validator"
func RegisterAllAdmissionPlugins(plugins *admission.Plugins) {
	resourcereferencemanager.Register(plugins)
	deletionconfirmation.Register(plugins)
	extensionvalidation.Register(plugins)
	extensionlabels.Register(plugins)
	shoottolerationrestriction.Register(plugins)

func TestBackupBucket(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Core BackupBucket Suite")
}

func TestBackupEntry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Core BackupEntry Suite")
}

func TestCloudProfile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Core CloudProfile Suite")
}

func TestControllerInstallation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry ControllerInstallation Suite")
}

var _ = Describe("ToSelectableFields", func() {

func TestNamespacedCloudProfile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Core NamespacedCloudProfile Suite")
}

func TestProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Core Project Suite")
}
func TestSecretBinding(t *testing.T) {
	features.RegisterFeatureGates()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Core SecretBinding Suite")
}
	return true
}

func (secretBindingStrategy) PrepareForCreate(_ context.Context, _ runtime.Object) {
}

func (secretBindingStrategy) Validate(_ context.Context, obj runtime.Object) field.ErrorList {
		}
	})

	Describe("#Validate", func() {
		It("should forbid creating SecretBinding when provider is nil or empty", func() {
			secretBinding.Provider = nil

func TestSeed(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Core Seed Suite")
}
func TestShoot(t *testing.T) {
	features.RegisterFeatureGates()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Core Shoot Suite")
}

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Core Shoot Storage Suite")
}

func TestBastion(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Operations Bastion Suite")
}
	return true
}

func (credentialsBindingStrategy) PrepareForCreate(_ context.Context, _ runtime.Object) {

}

func (credentialsBindingStrategy) PrepareForUpdate(_ context.Context, _, _ runtime.Object) {

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Security WorkloadIdentity Storage Suite")
}

func TestWorkloadIdentity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Security WorkloadIdentity Suite")
}

func TestManagedSeed(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry SeedManagement ManagedSeed Suite")
}

func TestManagedSeedSet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry SeedManagement ManagedSeedSet Suite")
}
		a.clusterRole(),
		a.clusterRoleBinding(virtualGardenAccessSecret.ServiceAccountName),
		a.validatingWebhookConfiguration(caSecret),
	)
	if err != nil {
		return err
		clusterRole(),
		clusterRoleBinding(),
		validatingWebhookConfiguration(namespace, caGardener.Data["bundle.crt"], testValues),
	))

	virtualManagedResourceSecret := &corev1.Secret{

	return webhookConfig
}
	return validatingWebhook
}

func buildWebhookConfigRulesForResourceSize(config *admissioncontrollerconfigv1alpha1.ResourceAdmissionConfiguration) []admissionregistrationv1.RuleWithOperations {
	if config == nil || len(config.Limits) == 0 {
		return nil
		}
	}

	if secretBinding.Provider != nil {
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
	}
func GardenWebhookSwitchOptions() *extensionscmdwebhook.SwitchOptions {
	return extensionscmdwebhook.NewSwitchOptions(
		extensionscmdwebhook.Switch(validator.Name, validator.New),
		extensionscmdwebhook.Switch(mutator.Name, mutator.New),
	)
}
package validator

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/pkg/apis/core"
	"github.com/gardener/gardener/pkg/provider-local/local"
)

const (
	// Name is a name for a validation webhook.
	Name = "validator"
)

var logger = log.Log.WithName("local-validator-webhook")
		Path:     "/webhooks/validate",
		Validators: map[extensionswebhook.Validator][]extensionswebhook.Type{
			NewNamespacedCloudProfileValidator(mgr): {{Obj: &core.NamespacedCloudProfile{}}},
		},
		Target: extensionswebhook.TargetSeed,
		ObjectSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"provider.extensions.gardener.cloud/local": "true"},
		},
	})
}
}

func addMetaDataLabelsSecretBinding(secretBinding *core.SecretBinding) {
	if secretBinding.Provider != nil {
		types := gardencorehelper.GetSecretBindingTypes(secretBinding)
		for _, t := range types {
			metav1.SetMetaDataLabel(&secretBinding.ObjectMeta, v1beta1constants.LabelExtensionProviderTypePrefix+t, "true")
		}
	}
}

	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	kubecorev1listers "k8s.io/client-go/listers/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gardener/gardener/pkg/apis/core"
	return nil
}

// Admit ensures that referenced resources do actually exist.
func (r *ReferenceManager) Admit(ctx context.Context, a admission.Attributes, _ admission.ObjectInterfaces) error {
	// Wait until the caches have been synced
	if r.readyFunc == nil {
		r.AssignReadyFunc(func() bool {

		switch a.GetOperation() {
		case admission.Create:
			// Add createdBy annotation to Shoot
			annotations := shoot.Annotations
			if annotations == nil {
				annotations = map[string]string{}
			}
			annotations[v1beta1constants.GardenCreatedBy] = a.GetUserInfo().GetName()
			shoot.Annotations = annotations

			oldShoot = &core.Shoot{}
		case admission.Update:
			// skip verification if spec wasn't changed
		if utils.SkipVerification(operation, project.ObjectMeta) {
			return nil
		}
		// Set createdBy field in Project
		switch a.GetOperation() {
		case admission.Create:
			project.Spec.CreatedBy = &rbacv1.Subject{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     rbacv1.UserKind,
				Name:     a.GetUserInfo().GetName(),
			}

			if project.Spec.Owner == nil {
				owner := project.Spec.CreatedBy

			outer:
				for _, member := range project.Spec.Members {
					for _, role := range member.Roles {
						if role == core.ProjectMemberOwner {
							owner = member.Subject.DeepCopy()
							break outer
						}
					}
				}

				project.Spec.Owner = owner
			}

			err = r.ensureProjectNamespace(project)
		case admission.Update:
			oldProject, ok := a.GetOldObject().(*core.Project)
			}
		}

		if project.Spec.Owner != nil {
			ownerIsMember := false
			for _, member := range project.Spec.Members {
				if member.Subject == *project.Spec.Owner {
					ownerIsMember = true
				}
			}
			if !ownerIsMember {
				project.Spec.Members = append(project.Spec.Members, core.ProjectMember{
					Subject: *project.Spec.Owner,
					Roles: []string{
						core.ProjectMemberAdmin,
						core.ProjectMemberOwner,
					},
				})
			}
		}

	case core.Kind("BackupBucket"):
		if operation == admission.Delete {
			// The "delete endpoint" handler of the k8s.io/apiserver library calls the admission controllers
		credentialsNamespace  string
		credentialsName       string
		credentialsKind       string
	)
	switch attributes.GetKind().GroupKind() {
	case core.Kind("SecretBinding"):
		b, ok := binding.(*core.SecretBinding)
		credentialsNamespace = b.SecretRef.Namespace
		credentialsName = b.SecretRef.Name
		credentialsKind = "Secret"
	case security.Kind("CredentialsBinding"):
		b, ok := binding.(*security.CredentialsBinding)
		if !ok {
		}
		credentialsNamespace = b.CredentialsRef.Namespace
		credentialsName = b.CredentialsRef.Name
	default:
		return fmt.Errorf("%s is neither of kind SecretBinding nor CredentialsBinding", attributes.GetKind().GroupKind())
	}
	readAttributes := authorizer.AttributesRecord{
		User:            attributes.GetUserInfo(),
		Verb:            "get",
		if err := r.lookupSecret(ctx, credentialsNamespace, credentialsName); err != nil {
			return err
		}
	case "WorkloadIdentity":
		if err := r.lookupWorkloadIdentity(ctx, credentialsNamespace, credentialsName); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown credentials kind: %s", credentialsKind)
	}

type getFn func(context.Context, string, string) (runtime.Object, error)

func lookupResource(ctx context.Context, namespace, name string, get getFn, fallbackGet getFn) error {
	// First try to detect the resource in the cache.
	var err error

	_, err = get(ctx, namespace, name)
	if err == nil {
		return nil
	}
	if !apierrors.IsNotFound(err) {
		return err
	}

	// Second try to detect the resource in the cache after the first try failed.
	// Give the cache time to observe the resource before rejecting a create.
	// This helps when creating a resource and immediately creating a binding referencing it.
	time.Sleep(MissingResourceWait)
	_, err = get(ctx, namespace, name)

	switch {
	case apierrors.IsNotFound(err):
		// no-op
	case err != nil:
		return err
	default:
		return nil
	}

	// Third try to detect the secret, now by doing a live lookup instead of relying on the cache.
	if _, err := fallbackGet(ctx, namespace, name); err != nil {
		return err
	}

	return nil
}

func (r *ReferenceManager) lookupWorkloadIdentity(ctx context.Context, namespace, name string) error {
	workloadIdentityFromLister := func(_ context.Context, namespace, name string) (runtime.Object, error) {
		return r.workloadIdentityLister.WorkloadIdentities(namespace).Get(name)
	}
		return r.gardenSecurityClient.SecurityV1alpha1().WorkloadIdentities(namespace).Get(ctx, name, kubernetesclient.DefaultGetOptions())
	}

	return lookupResource(ctx, namespace, name, workloadIdentityFromLister, workloadIdentityFromClient)
}

func (r *ReferenceManager) lookupSecret(ctx context.Context, namespace, name string) error {
		return r.kubeClient.CoreV1().Secrets(namespace).Get(ctx, name, kubernetesclient.DefaultGetOptions())
	}

	return lookupResource(ctx, namespace, name, secretFromLister, secretFromClient)
}

func (r *ReferenceManager) lookupConfigMap(ctx context.Context, namespace, name string) error {
		return r.kubeClient.CoreV1().ConfigMaps(namespace).Get(ctx, name, kubernetesclient.DefaultGetOptions())
	}

	return lookupResource(ctx, namespace, name, configMapFromLister, configMapFromClient)
}

func (r *ReferenceManager) lookupControllerDeployment(ctx context.Context, name string) error {
		return r.gardenCoreClient.CoreV1beta1().ControllerDeployments().Get(ctx, name, kubernetesclient.DefaultGetOptions())
	}

	return lookupResource(ctx, "", name, deploymentFromLister, deploymentFromClient)
}

func (r *ReferenceManager) getAPIResource(groupVersion, kind string) (*metav1.APIResource, error) {
	}
	return nil
}
	. "github.com/onsi/gomega/gstruct"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/admission"

	"github.com/gardener/gardener/pkg/apis/core"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"github.com/gardener/gardener/pkg/apis/security"
	securityv1alpha1 "github.com/gardener/gardener/pkg/apis/security/v1alpha1"
	"github.com/gardener/gardener/pkg/apis/seedmanagement"
						Namespace: namespace,
					},
				},
			}
			secretBinding = gardencorev1beta1.SecretBinding{
				ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
					},
				},
			}
			securityCredentialsBindingRefSecret = security.CredentialsBinding{
				ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
					},
				},
			}
			credentialsBindingRefSecret = securityv1alpha1.CredentialsBinding{
				ObjectMeta: metav1.ObjectMeta{
						Namespace: namespace,
					},
				},
			}
			securityCredentialsBindingRefWorkloadIdentity = security.CredentialsBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:       workloadIdentityName,
					Namespace:  namespace,
				},
				Quotas: []corev1.ObjectReference{
					{
						Name:      quotaName,

			err = gardencorev1beta1.Convert_core_Project_To_v1beta1_Project(&coreProject, &project, nil)
			Expect(err).To(Succeed())
		})

		It("should return nil because the resource is not BackupBucket and operation is delete", func() {
			attrs := admission.NewAttributesRecord(&controllerRegistration, nil, core.Kind("ControllerRegistration").WithVersion("version"), "", controllerRegistration.Name, core.Resource("controllerregistrations").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, nil)

			err := admissionHandler.Admit(context.TODO(), attrs, nil)

			Expect(err).NotTo(HaveOccurred())

			attrs = admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), "", controllerRegistration.Name, core.Resource("shoots").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, nil)

			err = admissionHandler.Admit(context.TODO(), attrs, nil)

			Expect(err).NotTo(HaveOccurred())
		})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&controllerRegistration, nil, core.Kind("ControllerRegistration").WithVersion("version"), "", controllerRegistration.Name, core.Resource("controllerregistrations").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&controllerRegistration, nil, core.Kind("ControllerRegistration").WithVersion("version"), "", controllerRegistration.Name, core.Resource("controllerregistrations").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
					return true, nil, errors.New("nope, out of luck")
				})

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("nope, out of luck"))
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should reject because the referenced secret does not exist", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota)).To(Succeed())
				kubeClient.AddReactor("get", "secrets", func(_ testing.Action) (bool, runtime.Object, error) {
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSecretBinding, nil, core.Kind("SecretBinding").WithVersion("version"), coreSecretBinding.Namespace, coreSecretBinding.Name, core.Resource("secretbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("tests for CredentialsBinding objects referencing Secret", func() {
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should reject because the referenced secret does not exist", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota)).To(Succeed())
				kubeClient.AddReactor("get", "secrets", func(_ testing.Action) (bool, runtime.Object, error) {
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefSecret, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefSecret.Namespace, securityCredentialsBindingRefSecret.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("tests for CredentialsBinding objects referencing WorkloadIdentity", func() {
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
			It("should accept because all referenced objects have been found (workloadidentity looked up live)", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota)).To(Succeed())
				gardenSecurityClient.AddReactor("get", "workloadidentities", func(_ testing.Action) (bool, runtime.Object, error) {
					return true, &securityv1alpha1.WorkloadIdentity{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: workloadIdentity.Namespace,
							Name:      workloadIdentity.Name,
						},
					}, nil
				})

				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

			It("should reject because the referenced workload identity does not exist", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota)).To(Succeed())
				gardenSecurityClient.AddReactor("get", "workloadidentities", func(_ testing.Action) (bool, runtime.Object, error) {
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				quotaRefList = append(quotaRefList, quota2Ref)
				securityCredentialsBindingRefWorkloadIdentity.Quotas = quotaRefList

				Expect(kubeInformerFactory.Core().V1().Secrets().Informer().GetStore().Add(&secret)).To(Succeed())
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota)).To(Succeed())
				Expect(gardenCoreInformerFactory.Core().V1beta1().Quotas().Informer().GetStore().Add(&quota2)).To(Succeed())
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&securityCredentialsBindingRefWorkloadIdentity, nil, security.Kind("CredentialsBinding").WithVersion("version"), securityCredentialsBindingRefWorkloadIdentity.Namespace, securityCredentialsBindingRefWorkloadIdentity.Name, security.Resource("credentialsbindings").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
		})

		Context("tests for Shoot objects", func() {
			It("should add the created-by annotation", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().CloudProfiles().Informer().GetStore().Add(&cloudProfile)).To(Succeed())
				Expect(gardenCoreInformerFactory.Core().V1beta1().Seeds().Informer().GetStore().Add(&seed)).To(Succeed())
				Expect(gardenCoreInformerFactory.Core().V1beta1().SecretBindings().Informer().GetStore().Add(&secretBinding)).To(Succeed())
				Expect(gardenSecurityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBindingRefSecret)).To(Succeed())
				Expect(kubeInformerFactory.Core().V1().ConfigMaps().Informer().GetStore().Add(&configMap)).To(Succeed())

				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				Expect(coreShoot.Annotations).NotTo(HaveKeyWithValue(v1beta1constants.GardenCreatedBy, user.Name))

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
				Expect(coreShoot.Annotations).To(HaveKeyWithValue(v1beta1constants.GardenCreatedBy, user.Name))
			})

			It("should accept because all referenced objects have been found", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().CloudProfiles().Informer().GetStore().Add(&cloudProfile)).To(Succeed())
				Expect(gardenCoreInformerFactory.Core().V1beta1().Seeds().Informer().GetStore().Add(&seed)).To(Succeed())
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				coreShoot.Status.TechnicalID = "should-never-change"
				attrs := admission.NewAttributesRecord(&coreShoot, oldShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				It("should reject because the referenced cloud profile does not exist (create)", func() {
					attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

					err := admissionHandler.Admit(context.TODO(), attrs, nil)

					Expect(err).To(HaveOccurred())
				})

					attrs := admission.NewAttributesRecord(&coreShoot, oldShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

					err := admissionHandler.Admit(context.TODO(), attrs, nil)

					Expect(err).To(HaveOccurred())
				})

					attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

					err := admissionHandler.Admit(context.TODO(), attrs, nil)

					Expect(err).To(HaveOccurred())
				})

					attrs := admission.NewAttributesRecord(&coreShoot, oldShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

					err := admissionHandler.Admit(context.TODO(), attrs, nil)

					Expect(err).To(HaveOccurred())
				})

				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				It("should reject because the referenced exposure class does not exists", func() {
					attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

					err := admissionHandler.Admit(context.TODO(), attrs, nil)
					Expect(err).To(HaveOccurred())
				})

					Expect(gardenCoreInformerFactory.Core().V1beta1().ExposureClasses().Informer().GetStore().Add(&exposureClass)).To(Succeed())
					attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

					err := admissionHandler.Admit(context.TODO(), attrs, nil)
					Expect(err).To(HaveOccurred())
				})
			})

				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(MatchError("shoots.core.gardener.cloud \"shoot-1\" is forbidden: cannot reference a resource you are not allowed to read"))
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&coreShoot, oldShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(MatchError("shoots.core.gardener.cloud \"shoot-1\" is forbidden: cannot reference a resource you are not allowed to read"))
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&coreShoot, oldShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(Not(HaveOccurred()))
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(MatchError(ContainSubstring("failed to resolve resource reference")))
			})
					user := &user.DefaultInfo{Name: allowedUser}
					attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

					Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(MatchError(ContainSubstring(expectedErrorMessage)))
				})

				It("should reject because the referenced "+description+" does not exist (update)", func() {
					user := &user.DefaultInfo{Name: allowedUser}
					attrs := admission.NewAttributesRecord(&coreShoot, oldShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, user)

					Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(MatchError(ContainSubstring(expectedErrorMessage)))
				})

				It("should pass because the referenced "+description+" does not exist but shoot has deletion timestamp", func() {
					user := &user.DefaultInfo{Name: allowedUser}
					attrs := admission.NewAttributesRecord(&coreShoot, oldShoot, core.Kind("Shoot").WithVersion("version"), coreShoot.Namespace, coreShoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, user)

					Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
				})

				It("should pass because the referenced "+description+" exists", func() {
					user := &user.DefaultInfo{Name: allowedUser}
					attrs := admission.NewAttributesRecord(&coreShoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

					Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
				})
			}

				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&coreSeed, oldSeed, core.Kind("Seed").WithVersion("version"), "", coreSeed.Name, core.Resource("seeds").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(MatchError("seeds.core.gardener.cloud \"seed-1\" is forbidden: cannot reference a resource you are not allowed to read"))
			})
				user := &user.DefaultInfo{Name: "disallowed-user"}
				attrs := admission.NewAttributesRecord(&coreSeed, oldSeed, core.Kind("Seed").WithVersion("version"), "", coreSeed.Name, core.Resource("seeds").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(Not(HaveOccurred()))
			})
				user := &user.DefaultInfo{Name: allowedUser}
				attrs := admission.NewAttributesRecord(&coreSeed, nil, core.Kind("Seed").WithVersion("version"), "", coreSeed.Name, core.Resource("seeds").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, user)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(MatchError(ContainSubstring("failed to resolve resource reference")))
			})
			It("should reject if the referred Seed is not found", func() {
				attrs := admission.NewAttributesRecord(&coreBackupBucket, nil, core.Kind("BackupBucket").WithVersion("version"), "", coreBackupBucket.Name, core.Resource("backupBuckets").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(BeForbiddenError())
				Expect(err).To(MatchError(ContainSubstring("backupBuckets.core.gardener.cloud %q is forbidden: seed.core.gardener.cloud %q not found", coreBackupBucket.Name, seed.Name)))

				attrs := admission.NewAttributesRecord(&coreBackupBucket, nil, core.Kind("BackupBucket").WithVersion("version"), "", coreBackupBucket.Name, core.Resource("backupBuckets").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(BeForbiddenError())
				Expect(err).To(MatchError(ContainSubstring("secret not found")))

				attrs := admission.NewAttributesRecord(&coreBackupBucket, nil, core.Kind("BackupBucket").WithVersion("version"), "", coreBackupBucket.Name, core.Resource("backupBuckets").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&coreBackupBucket, nil, core.Kind("BackupBucket").WithVersion("version"), "", coreBackupBucket.Name, core.Resource("backupBuckets").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(nil, nil, core.Kind("BackupBucket").WithVersion("version"), "", coreBackupBucket.Name, core.Resource("backupBuckets").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(nil, nil, core.Kind("BackupBucket").WithVersion("version"), "", backupBucket.Name, core.Resource("backupBuckets").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(BeForbiddenError())
				Expect(err).To(MatchError(ContainSubstring("backupBuckets.core.gardener.cloud %q is forbidden: cannot delete BackupBucket because BackupEntries are still referencing it, backupEntryNames: %s/%s,%s/%s", backupBucket.Name, backupEntry.Namespace, backupEntry.Name, backupEntry2.Namespace, backupEntry2.Name)))

				attrs := admission.NewAttributesRecord(nil, nil, core.Kind("BackupBucket").WithVersion("version"), "", coreBackupBucket.Name, core.Resource("backupBuckets").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(nil, nil, core.Kind("BackupBucket").WithVersion("version"), "", "", core.Resource("backupBuckets").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(BeForbiddenError())
				Expect(err).To(MatchError(ContainSubstring("backupBuckets.core.gardener.cloud %q is forbidden: cannot delete BackupBucket because BackupEntries are still referencing it, backupEntryNames: %s/%s,%s/%s", backupBucket2.Name, backupEntry.Namespace, backupEntry.Name, backupEntry2.Namespace, backupEntry2.Name)))

				attrs := admission.NewAttributesRecord(nil, nil, core.Kind("BackupBucket").WithVersion("version"), "", "", core.Resource("backupBuckets").WithVersion("version"), "", admission.Delete, &metav1.DeleteOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
			It("should reject if the referred Seed is not found", func() {
				attrs := admission.NewAttributesRecord(&coreBackupEntry, nil, core.Kind("BackupEntry").WithVersion("version"), coreBackupEntry.Namespace, coreBackupEntry.Name, core.Resource("backupEntries").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(BeForbiddenError())
				Expect(err).To(MatchError(ContainSubstring("backupEntries.core.gardener.cloud %q is forbidden: seed.core.gardener.cloud %q not found", coreBackupEntry.Name, seed.Name)))
				Expect(gardenCoreInformerFactory.Core().V1beta1().Seeds().Informer().GetStore().Add(&seed)).To(Succeed())
				attrs := admission.NewAttributesRecord(&coreBackupEntry, nil, core.Kind("BackupEntry").WithVersion("version"), coreBackupEntry.Namespace, coreBackupEntry.Name, core.Resource("backupEntries").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(BeForbiddenError())
				Expect(err).To(MatchError(ContainSubstring("backupEntries.core.gardener.cloud %q is forbidden: backupbucket.core.gardener.cloud %q not found", coreBackupEntry.Name, coreBackupBucket.Name)))
				Expect(gardenCoreInformerFactory.Core().V1beta1().BackupBuckets().Informer().GetStore().Add(&backupBucket)).To(Succeed())
				attrs := admission.NewAttributesRecord(&coreBackupEntry, nil, core.Kind("BackupEntry").WithVersion("version"), coreBackupEntry.Namespace, coreBackupEntry.Name, core.Resource("backupEntries").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("tests for Project objects", func() {
			It("should set the created-by field", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().Projects().Informer().GetStore().Add(&project)).To(Succeed())

				attrs := admission.NewAttributesRecord(&coreProject, nil, core.Kind("Project").WithVersion("version"), coreProject.Namespace, coreProject.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
				Expect(coreProject.Spec.CreatedBy).To(Equal(&rbacv1.Subject{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     rbacv1.UserKind,
					Name:     defaultUserName,
				}))
			})

			It("should set the owner field (member with owner role found)", func() {
				projectCopy := project.DeepCopy()
				coreProjectCopy := coreProject.DeepCopy()
				ownerMember := &rbacv1.Subject{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     rbacv1.UserKind,
					Name:     "owner",
				}
				projectCopy.Name = "foo"
				projectCopy.Spec.Members = []gardencorev1beta1.ProjectMember{
					{
						Subject: *ownerMember,
						Roles:   []string{core.ProjectMemberOwner},
					},
				}
				coreProjectCopy.Name = "foo"
				coreProjectCopy.Spec.Members = []core.ProjectMember{
					{
						Subject: *ownerMember,
						Roles:   []string{core.ProjectMemberOwner},
					},
				}

				Expect(gardenCoreInformerFactory.Core().V1beta1().Projects().Informer().GetStore().Add(projectCopy)).To(Succeed())

				attrs := admission.NewAttributesRecord(coreProjectCopy, nil, core.Kind("Project").WithVersion("version"), coreProjectCopy.Namespace, coreProjectCopy.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
				Expect(coreProjectCopy.Spec.Owner).To(Equal(ownerMember))
				Expect(coreProjectCopy.Spec.CreatedBy).To(Equal(&rbacv1.Subject{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     rbacv1.UserKind,
					Name:     defaultUserName,
				}))
			})

			It("should set the owner field", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().Projects().Informer().GetStore().Add(&project)).To(Succeed())

				attrs := admission.NewAttributesRecord(&coreProject, nil, core.Kind("Project").WithVersion("version"), coreProject.Namespace, coreProject.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
				Expect(coreProject.Spec.Owner).To(Equal(&rbacv1.Subject{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     rbacv1.UserKind,
					Name:     defaultUserName,
				}))
			})

			It("should add the owner to members", func() {
				Expect(gardenCoreInformerFactory.Core().V1beta1().Projects().Informer().GetStore().Add(&project)).To(Succeed())

				attrs := admission.NewAttributesRecord(&coreProject, nil, core.Kind("Project").WithVersion("version"), coreProject.Namespace, coreProject.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
				Expect(coreProject.Spec.Members).To(ContainElement(Equal(core.ProjectMember{
					Subject: rbacv1.Subject{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     rbacv1.UserKind,
						Name:     defaultUserName,
					},
					Roles: []string{
						core.ProjectMemberAdmin,
						core.ProjectMemberOwner,
					},
				})))
			})

			It("should allow specifying a namespace which is not in use (create)", func() {
				project.Spec.Namespace = ptr.To("garden-foo")
				projectCopy := project.DeepCopy()
				coreProject.Spec.Namespace = ptr.To("garden-foo")
				attrs := admission.NewAttributesRecord(&coreProject, nil, core.Kind("Project").WithVersion("version"), coreProject.Namespace, coreProject.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(Not(HaveOccurred()))
			})
				coreProject.Spec.Namespace = ptr.To("garden-foo")
				attrs := admission.NewAttributesRecord(&coreProject, coreProjectOld, core.Kind("Project").WithVersion("version"), coreProject.Namespace, coreProject.Name, core.Resource("projects").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(Not(HaveOccurred()))
			})

				attrs := admission.NewAttributesRecord(&coreProject, nil, core.Kind("Project").WithVersion("version"), coreProject.Namespace, coreProject.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(Not(HaveOccurred()))
			})
				coreProject.Spec.Namespace = ptr.To("garden-foo")
				attrs := admission.NewAttributesRecord(&coreProject, nil, core.Kind("Project").WithVersion("version"), coreProject.Namespace, coreProject.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"ErrStatus": MatchFields(IgnoreExtras, Fields{
				coreProject.Spec.Namespace = ptr.To("garden-foo")
				attrs := admission.NewAttributesRecord(&coreProject, &coreProjectOld, core.Kind("Project").WithVersion("version"), coreProject.Namespace, coreProject.Name, core.Resource("projects").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"ErrStatus": MatchFields(IgnoreExtras, Fields{

				attrs := admission.NewAttributesRecord(&cloudProfile, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("1.24.1"))

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to delete Kubernetes version"),
					ContainSubstring("1.24.1"),
					ContainSubstring("still in use by NamespacedCloudProfile"),

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to delete Kubernetes version"),
					ContainSubstring("1.24.1"),
					ContainSubstring("still in use by shoot '/shoot-Two'"),

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should accept removal of kubernetes versions that are used by shoots using another unrelated NamespacedCloudProfile of same name", func() {

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should accept removal of kubernetes version that is still in use by a shoot that is being deleted", func() {

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&cloudProfile, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("1.17.2"))

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("1.17.2"))
				Expect(err.Error()).To(ContainSubstring(s.Spec.Provider.Workers[1].Name))

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("1.17.2"))

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should fail for removal of a machine version that is used by a NamespacedCloudProfile", func() {

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to delete MachineImage version"),
					ContainSubstring("1.16.0"),
					ContainSubstring("still in use by NamespacedCloudProfile"),

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to delete MachineImage version"),
					ContainSubstring("1.16.0"),
					ContainSubstring("still in use by NamespacedCloudProfile"),

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to delete MachineImage \"coreos\""),
					ContainSubstring("still in use by NamespacedCloudProfile"),
				)))

				attrs := admission.NewAttributesRecord(&cloudProfileNew, &cloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to add MachineImage \"gardenlinux\""),
					ContainSubstring("already defined by NamespacedCloudProfile \"project-123/profile-42\""),
				)))
				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(shootTwo)).To(Succeed())
				attrs := admission.NewAttributesRecord(cloudProfile, oldCloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(ctx, attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(shootTwo)).To(Succeed())
				attrs := admission.NewAttributesRecord(cloudProfile, oldCloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(ctx, attrs, nil)

				Expect(err).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"ErrStatus": MatchFields(IgnoreExtras, Fields{
				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(shootOne)).To(Succeed())
				attrs := admission.NewAttributesRecord(cloudProfile, oldCloudProfile, core.Kind("CloudProfile").WithVersion("version"), "", cloudProfile.Name, core.Resource("CloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(ctx, attrs, nil)

				Expect(err).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"ErrStatus": MatchFields(IgnoreExtras, Fields{

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should succeed for the complete Kubernetes section being removed without usages", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should succeed if a used and already extended kubernetes version expiration is changed to another value still in the future", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should succeed if a used and extended kubernetes version already expired is not modified", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should succeed if an extended and used Kubernetes version is removed with the base version still being valid", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should fail if an extended and used Kubernetes version is being removed with the base version being already expired", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)
				Expect(err).To(MatchError(And(
					ContainSubstring("unable to delete Kubernetes version"),
					ContainSubstring("1.29.0"),

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to delete Kubernetes version"),
					ContainSubstring("1.29.0"),
					ContainSubstring("still in use by shoot"),

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})
		})


				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should succeed if a used and extended MachineImage version already expired is not modified", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should succeed if an extended and used MachineImage version is removed with the base version still being valid", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should fail if an extended and used MachineImage version is being removed with the base version being already expired", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)
				Expect(err).To(MatchError(And(
					ContainSubstring("unable to delete Machine image version"),
					ContainSubstring("1.17.3"),

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)
				Expect(err).To(MatchError(And(
					ContainSubstring("unable to delete Machine image version"),
					ContainSubstring("'coreos/1.1.2'"),

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(context.TODO(), attrs, nil)
				Expect(err).To(MatchError(And(
					ContainSubstring("unable to delete Machine image version"),
					ContainSubstring("'custom-namespaced-image/1.1.2'"),

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(MatchError(And(
					ContainSubstring("unable to delete Machine image version"),
					ContainSubstring("1.17.3"),
					ContainSubstring("still in use by shoot"),

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should succeed if a new but unused MachineImage version is removed", func() {

				attrs := admission.NewAttributesRecord(updatedNamespacedCloudProfile, namespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})
		})

				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(shootTwo)).To(Succeed())
				attrs := admission.NewAttributesRecord(namespacedCloudProfile, oldNamespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespacedCloudProfile.Namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(ctx, attrs, nil)

				Expect(err).NotTo(HaveOccurred())
			})
				Expect(gardenCoreInformerFactory.Core().V1beta1().Shoots().Informer().GetStore().Add(shootTwo)).To(Succeed())
				attrs := admission.NewAttributesRecord(namespacedCloudProfile, oldNamespacedCloudProfile, core.Kind("NamespacedCloudProfile").WithVersion("version"), namespacedCloudProfile.Namespace, namespacedCloudProfile.Name, core.Resource("NamespacedCloudProfile").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, defaultUserInfo)

				err := admissionHandler.Admit(ctx, attrs, nil)

				Expect(err).To(PointTo(MatchFields(IgnoreExtras, Fields{
					"ErrStatus": MatchFields(IgnoreExtras, Fields{
			It("should accept because there is no managed seed with the same name", func() {
				attrs := admission.NewAttributesRecord(gardenlet, nil, seedmanagement.Kind("Gardenlet").WithVersion("version"), gardenlet.Namespace, gardenlet.Name, seedmanagement.Resource("gardenlets").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, &user.DefaultInfo{Name: allowedUser})

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should forbid because there is a managed seed with the same name", func() {

				attrs := admission.NewAttributesRecord(gardenlet, nil, seedmanagement.Kind("Gardenlet").WithVersion("version"), gardenlet.Namespace, gardenlet.Name, seedmanagement.Resource("gardenlets").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, &user.DefaultInfo{Name: allowedUser})

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(MatchError(ContainSubstring("there is already a ManagedSeed object with the same name")))
			})
		})

			It("should accept because there is no gardenlet with the same name", func() {
				attrs := admission.NewAttributesRecord(managedSeed, nil, seedmanagement.Kind("ManagedSeed").WithVersion("version"), gardenlet.Namespace, gardenlet.Name, seedmanagement.Resource("gardenlets").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, &user.DefaultInfo{Name: allowedUser})

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(Succeed())
			})

			It("should forbid because there is a gardenlet with the same name", func() {

				attrs := admission.NewAttributesRecord(managedSeed, nil, seedmanagement.Kind("ManagedSeed").WithVersion("version"), managedSeed.Namespace, managedSeed.Name, seedmanagement.Resource("managedseeds").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, &user.DefaultInfo{Name: allowedUser})

				Expect(admissionHandler.Admit(context.TODO(), attrs, nil)).To(MatchError(ContainSubstring("there is already a Gardenlet object with the same name")))
			})
		})
	})
	PluginNameExtensionLabels = "ExtensionLabels"
	// PluginNameExtensionValidator is the name of the ExtensionValidator admission plugin.
	PluginNameExtensionValidator = "ExtensionValidator"
	// PluginNameResourceReferenceManager is the name of the ResourceReferenceManager admission plugin.
	PluginNameResourceReferenceManager = "ResourceReferenceManager"
	// PluginNameManagedSeedShoot is the name of the ManagedSeedShoot admission plugin.
		PluginNameNamespacedCloudProfileValidator,   // NamespacedCloudProfileValidator
		PluginNameProjectValidator,                  // ProjectValidator
		PluginNameDeletionConfirmation,              // DeletionConfirmation
		PluginNameOpenIDConnectPreset,               // OpenIDConnectPreset
		PluginNameClusterOpenIDConnectPreset,        // ClusterOpenIDConnectPreset
		PluginNameCustomVerbAuthorizer,              // CustomVerbAuthorizer
		PluginNameNamespacedCloudProfileValidator, // NamespacedCloudProfileValidator
		PluginNameProjectValidator,                // ProjectValidator
		PluginNameDeletionConfirmation,            // DeletionConfirmation
		PluginNameOpenIDConnectPreset,             // OpenIDConnectPreset
		PluginNameClusterOpenIDConnectPreset,      // ClusterOpenIDConnectPreset
		PluginNameCustomVerbAuthorizer,            // CustomVerbAuthorizer
	"context"
	"fmt"
	"io"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/admission"

	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	gardenerutils "github.com/gardener/gardener/pkg/utils/gardener"
	plugin "github.com/gardener/gardener/plugin/pkg"
)

// Register registers a plugin.
// New creates a new handler admission plugin.
func New() (*handler, error) {
	return &handler{
		Handler: admission.NewHandler(admission.Create),
	}, nil
}

var _ admission.ValidationInterface = &handler{}

func (v *handler) Validate(_ context.Context, a admission.Attributes, _ admission.ObjectInterfaces) error {
	// Ignore all kinds other than Project
	if a.GetKind().GroupKind() != gardencore.Kind("Project") {
		return nil
		return apierrors.NewBadRequest("could not convert object to Project")
	}

	// TODO: Remove this admission plugin in favor of static validation in a future release, see https://github.com/gardener/gardener/pull/4228.
	if project.Spec.Namespace != nil && *project.Spec.Namespace != v1beta1constants.GardenNamespace && !strings.HasPrefix(*project.Spec.Namespace, gardenerutils.ProjectNamespacePrefix) {
		return admission.NewForbidden(a, fmt.Errorf(".spec.namespace must start with %s", gardenerutils.ProjectNamespacePrefix))
	}

	return nil
}

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/utils/ptr"

	"github.com/gardener/gardener/pkg/apis/core"
		var (
			err              error
			project          core.Project
			admissionHandler admission.ValidationInterface

			namespaceName = "garden-my-project"
			projectName   = "my-project"
					Namespace: namespaceName,
				},
			}
		)

		BeforeEach(func() {
			Expect(err).NotTo(HaveOccurred())

			project = projectBase
		})

		It("should allow creating the project (namespace nil)", func() {
			attrs := admission.NewAttributesRecord(&project, nil, core.Kind("Project").WithVersion("version"), "", project.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
		})

		It("should allow creating the project(namespace non-nil)", func() {
			project.Spec.Namespace = &namespaceName

			attrs := admission.NewAttributesRecord(&project, nil, core.Kind("Project").WithVersion("version"), "", project.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
		})

		It("should allow creating the project (namespace is 'garden')", func() {
			project.Spec.Namespace = ptr.To(v1beta1constants.GardenNamespace)

			attrs := admission.NewAttributesRecord(&project, nil, core.Kind("Project").WithVersion("version"), "", project.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(Succeed())
		})

		It("should prevent creating the project because namespace prefix is missing", func() {
			project.Spec.Namespace = ptr.To("foo")

			attrs := admission.NewAttributesRecord(&project, nil, core.Kind("Project").WithVersion("version"), "", project.Name, core.Resource("projects").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)

			Expect(admissionHandler.Validate(context.TODO(), attrs, nil)).To(MatchError(ContainSubstring(".spec.namespace must start with garden-")))
		})
	})

	})

	Describe("#New", func() {
		It("should only handle CREATE operations", func() {
			dr, err := New()
			Expect(err).ToNot(HaveOccurred())
			Expect(dr.Handles(admission.Create)).To(BeTrue())
			Expect(dr.Handles(admission.Update)).To(BeFalse())
			Expect(dr.Handles(admission.Connect)).To(BeFalse())
			Expect(dr.Handles(admission.Delete)).To(BeFalse())
		})
		}
	}

	cloudProfileSpec, err := admissionutils.GetCloudProfileSpec(v.cloudProfileLister, v.namespacedCloudProfileLister, shoot)
	if err != nil {
		return apierrors.NewInternalError(fmt.Errorf("could not find referenced cloud profile: %+v", err.Error()))
	}

	if a.GetOperation() == admission.Create && len(ptr.Deref(shoot.Spec.CloudProfileName, "")) > 0 && shoot.Spec.CloudProfile != nil {
		return fmt.Errorf("new shoot can only specify either cloudProfileName or cloudProfile reference")
	}

	if err := admissionutils.ValidateCloudProfileChanges(v.cloudProfileLister, v.namespacedCloudProfileLister, shoot, oldShoot); err != nil {
		return err
	}
		}
	}

	// Allow removal of `gardener` finalizer only if the Shoot deletion has completed successfully
	if len(c.shoot.Status.TechnicalID) > 0 && c.shoot.Status.LastOperation != nil {
		oldFinalizers := sets.New(c.oldShoot.Finalizers...)
		newFinalizers := sets.New(c.shoot.Finalizers...)

		if oldFinalizers.Has(core.GardenerName) && !newFinalizers.Has(core.GardenerName) {
			lastOperation := c.shoot.Status.LastOperation
			deletionSucceeded := lastOperation.Type == core.LastOperationTypeDelete && lastOperation.State == core.LastOperationStateSucceeded && lastOperation.Progress == 100

			if !deletionSucceeded {
				return admission.NewForbidden(a, fmt.Errorf("finalizer %q cannot be removed because shoot deletion has not completed successfully yet", core.GardenerName))
			}
		}
	}

	return nil
}


	return allErrs
}
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
			})
		})

		Context("shoot with generate name", func() {
			BeforeEach(func() {
				shoot.ObjectMeta = metav1.ObjectMeta{
					GenerateName: "demo-",
					Namespace:    namespaceName,
				}
			})

			It("should admit Shoot resources", func() {
				Expect(coreInformerFactory.Core().V1beta1().Projects().Informer().GetStore().Add(&project)).To(Succeed())
				Expect(coreInformerFactory.Core().V1beta1().CloudProfiles().Informer().GetStore().Add(&cloudProfile)).To(Succeed())
				Expect(coreInformerFactory.Core().V1beta1().Seeds().Informer().GetStore().Add(&seed)).To(Succeed())
				Expect(coreInformerFactory.Core().V1beta1().SecretBindings().Informer().GetStore().Add(&secretBinding)).To(Succeed())
				Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())

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
				Expect(coreInformerFactory.Core().V1beta1().CloudProfiles().Informer().GetStore().Add(&cloudProfile)).To(Succeed())
				Expect(coreInformerFactory.Core().V1beta1().Seeds().Informer().GetStore().Add(&seed)).To(Succeed())
				Expect(coreInformerFactory.Core().V1beta1().SecretBindings().Informer().GetStore().Add(&secretBinding)).To(Succeed())
				Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())

				authorizeAttributes.Name = shoot.Name

				attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, userInfo)
				err := admissionHandler.Admit(ctx, attrs, nil)

				Expect(err).To(BeInvalidError())
				Expect(err.Error()).To(ContainSubstring("name must not exceed"))
			})
		})

		Context("finalizer removal checks", func() {
			var (
				oldShoot *core.Shoot
			)

			BeforeEach(func() {
				shoot = *shootBase.DeepCopy()

				shoot.Status.TechnicalID = "some-id"
				shoot.Status.LastOperation = &core.LastOperation{
					Type:     core.LastOperationTypeReconcile,
					State:    core.LastOperationStateSucceeded,
					Progress: 100,
				}

				// set old shoot for update and add gardener finalizer to it
				oldShoot = shoot.DeepCopy()
				finalizers := sets.New(oldShoot.GetFinalizers()...)
				finalizers.Insert(core.GardenerName)
				oldShoot.SetFinalizers(finalizers.UnsortedList())
			})

			It("should reject removing the gardener finalizer if the shoot has not yet been deleted successfully", func() {
				Expect(coreInformerFactory.Core().V1beta1().Projects().Informer().GetStore().Add(&project)).To(Succeed())
				Expect(coreInformerFactory.Core().V1beta1().CloudProfiles().Informer().GetStore().Add(&cloudProfile)).To(Succeed())
				Expect(coreInformerFactory.Core().V1beta1().Seeds().Informer().GetStore().Add(&seed)).To(Succeed())
				Expect(coreInformerFactory.Core().V1beta1().SecretBindings().Informer().GetStore().Add(&secretBinding)).To(Succeed())
				Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())

				attrs := admission.NewAttributesRecord(&shoot, oldShoot, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, nil)
				err := admissionHandler.Admit(ctx, attrs, nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("shoot deletion has not completed successfully yet"))
			})

			It("should admit removing the gardener finalizer if the shoot deletion succeeded ", func() {
				Expect(coreInformerFactory.Core().V1beta1().Projects().Informer().GetStore().Add(&project)).To(Succeed())
				Expect(coreInformerFactory.Core().V1beta1().CloudProfiles().Informer().GetStore().Add(&cloudProfile)).To(Succeed())
				Expect(coreInformerFactory.Core().V1beta1().Seeds().Informer().GetStore().Add(&seed)).To(Succeed())
				Expect(coreInformerFactory.Core().V1beta1().SecretBindings().Informer().GetStore().Add(&secretBinding)).To(Succeed())
				Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())

				shoot.Status.LastOperation = &core.LastOperation{
					Type:     core.LastOperationTypeDelete,
					State:    core.LastOperationStateSucceeded,
					Progress: 100,
				}

				attrs := admission.NewAttributesRecord(&shoot, oldShoot, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Update, &metav1.UpdateOptions{}, false, nil)
				err := admissionHandler.Admit(ctx, attrs, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})

				shoot.Spec.SeedName = nil
				shoot.Spec.AccessRestrictions = []core.AccessRestrictionWithOptions{{AccessRestriction: core.AccessRestriction{Name: "foo"}}}

				attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.UpdateOptions{}, false, nil)
				err := admissionHandler.Admit(ctx, attrs, nil)

				Expect(err).To(BeForbiddenError())

				shoot.Spec.AccessRestrictions = []core.AccessRestrictionWithOptions{{AccessRestriction: core.AccessRestriction{Name: "foo"}}}

				attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.UpdateOptions{}, false, nil)
				err := admissionHandler.Admit(ctx, attrs, nil)

				Expect(err).To(BeForbiddenError())

				shoot.Spec.AccessRestrictions = []core.AccessRestrictionWithOptions{{AccessRestriction: core.AccessRestriction{Name: "foo"}}}

				attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.UpdateOptions{}, false, nil)
				err := admissionHandler.Admit(ctx, attrs, nil)

				Expect(err).NotTo(HaveOccurred())
						})

						It("should allow scheduling non-HA shoot", func() {
							attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)
							Expect(admissionHandler.Admit(ctx, attrs, nil)).To(Succeed())
						})

							shoot.Annotations = make(map[string]string)
							shoot.Spec.ControlPlane = &core.ControlPlane{HighAvailability: &core.HighAvailability{FailureTolerance: core.FailureTolerance{Type: core.FailureToleranceTypeNode}}}

							attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)
							Expect(admissionHandler.Admit(ctx, attrs, nil)).To(Succeed())
						})

							shoot.Annotations = make(map[string]string)
							shoot.Spec.ControlPlane = &core.ControlPlane{HighAvailability: &core.HighAvailability{FailureTolerance: core.FailureTolerance{Type: core.FailureToleranceTypeZone}}}

							attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)
							Expect(admissionHandler.Admit(ctx, attrs, nil)).To(BeForbiddenError())
						})
					})
						})

						It("should allow scheduling non-HA shoot", func() {
							attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)
							Expect(admissionHandler.Admit(ctx, attrs, nil)).To(Succeed())
						})

							shoot.Annotations = make(map[string]string)
							shoot.Spec.ControlPlane = &core.ControlPlane{HighAvailability: &core.HighAvailability{FailureTolerance: core.FailureTolerance{Type: core.FailureToleranceTypeNode}}}

							attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)
							Expect(admissionHandler.Admit(ctx, attrs, nil)).To(Succeed())
						})

							shoot.Annotations = make(map[string]string)
							shoot.Spec.ControlPlane = &core.ControlPlane{HighAvailability: &core.HighAvailability{FailureTolerance: core.FailureTolerance{Type: core.FailureToleranceTypeZone}}}

							attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)
							Expect(admissionHandler.Admit(ctx, attrs, nil)).To(Succeed())
						})
					})
						Expect(coreInformerFactory.Core().V1beta1().SecretBindings().Informer().GetStore().Add(&secretBinding)).To(Succeed())
						Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())

						attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)
						err := admissionHandler.Admit(ctx, attrs, nil)

						Expect(err).NotTo(HaveOccurred())
						Expect(coreInformerFactory.Core().V1beta1().SecretBindings().Informer().GetStore().Add(&secretBinding)).To(Succeed())
						Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())

						attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)
						err := admissionHandler.Admit(ctx, attrs, nil)

						Expect(err).To(BeForbiddenError())
						Expect(coreInformerFactory.Core().V1beta1().SecretBindings().Informer().GetStore().Add(&secretBinding)).To(Succeed())
						Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())

						attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)
						err := admissionHandler.Admit(ctx, attrs, nil)

						Expect(err).To(BeForbiddenError())
						Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())
						Expect(kubeInformerFactory.Core().V1().Secrets().Informer().GetStore().Add(&secret)).To(Succeed())

						attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)
						err := admissionHandler.Admit(ctx, attrs, nil)

						Expect(err).To(BeForbiddenError())
						Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())
						Expect(kubeInformerFactory.Core().V1().Secrets().Informer().GetStore().Add(&secret)).To(Succeed())

						attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)
						err := admissionHandler.Admit(ctx, attrs, nil)

						Expect(err).NotTo(HaveOccurred())
					Expect(coreInformerFactory.Core().V1beta1().SecretBindings().Informer().GetStore().Add(&secretBinding)).To(Succeed())
					Expect(securityInformerFactory.Security().V1alpha1().CredentialsBindings().Informer().GetStore().Add(&credentialsBinding)).To(Succeed())

					attrs := admission.NewAttributesRecord(&shoot, nil, core.Kind("Shoot").WithVersion("version"), shoot.Namespace, shoot.Name, core.Resource("shoots").WithVersion("version"), "", admission.Create, &metav1.CreateOptions{}, false, nil)
					err := admissionHandler.Admit(ctx, attrs, nil)

					Expect(err).To(errorMatcher)

	It("should be able to manipulate resource from security.gardener.cloud/v1alpha1", func() {
		credentialsBinding := &securityv1alpha1.CredentialsBinding{ObjectMeta: metav1.ObjectMeta{GenerateName: "test-", Namespace: testNamespace.Name}}
		Expect(testClient.Create(ctx, credentialsBinding)).To(MatchError(ContainSubstring("credentialsbindings.security.gardener.cloud \"test-\" is forbidden")))
	})
})
