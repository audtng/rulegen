package main

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
