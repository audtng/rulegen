package main

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
