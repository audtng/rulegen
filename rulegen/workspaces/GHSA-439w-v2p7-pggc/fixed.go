package main

	"fmt"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"


	k8scloud "github.com/juju/juju/caas/kubernetes/cloud"
	"github.com/juju/juju/core/model"
	coresecrets "github.com/juju/juju/core/secrets"
	"github.com/juju/juju/environs/bootstrap"
	"github.com/juju/juju/environs/cloudspec"
	k8sprovider "github.com/juju/juju/internal/provider/kubernetes"
	if err != nil {
		return errors.Trace(err)
	}

	ctx := context.TODO()
	err = broker.dropSecretAccess(ctx, removed.RevisionIDs())
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func cloudSpecToBackendConfig(spec cloudspec.CloudSpec) (*provider.BackendConfig, error) {
	return true
}

// CleanupIssuedTokens removes all ACLs/tokens related to the given issued
// token UUIDs. It returns, even during error, the list of tokens it revoked
// so far.
func (p k8sProvider) CleanupIssuedTokens(
	adminCfg *provider.ModelBackendConfig, issuedTokenUUIDs []string,
) ([]string, error) {
	broker, err := p.getBroker(adminCfg)
	if err != nil {
		return nil, errors.Trace(err)
	}

	ctx := context.TODO()

	for i, uuid := range issuedTokenUUIDs {
		err = broker.revokeSecretAccessToken(ctx, uuid)
		if err != nil {
			// return the tokens deleted so far.
			return issuedTokenUUIDs[:i], errors.New(
				"removing k8s secret backend issued tokens",
			)
		}
	}

	return issuedTokenUUIDs, nil
}

// RestrictedConfig returns the config needed to create a
// secrets backend client restricted to manage the specified
// owned secrets and read shared secrets for the given entity tag.
func (p k8sProvider) RestrictedConfig(
	adminCfg *provider.ModelBackendConfig,
	sameController, forDrain bool,
	issuedTokenUUID string,
	consumer names.Tag,
	owned []string,
	ownedRevs provider.SecretRevisions,
	readRevs provider.SecretRevisions,
) (*provider.BackendConfig, error) {
	logger.Tracef("getting k8s backend config for %q, owned %v, readRevs %v",
		consumer, owned, readRevs)

	if consumer == nil {
		return &adminCfg.BackendConfig, nil
		return nil, errors.Trace(err)
	}
	ctx := context.TODO()

	// Kubernetes secrets cannot restrict create operations by name. To ensure
	// a restricted config cannot create secrets with other names, we must add
	// an extra pre-created secret object for the next revision. For secrets
	// that have not yet been created, we must make the first revision secret
	// object.
	maxOwnedRev := make(map[string]int)
	for _, rev := range ownedRevs.RevisionIDs() {
		id, rev, err := coresecrets.ParseRevisionName(rev)
		if err != nil {
			return nil, errors.Trace(err)
		}
		maxOwnedRev[id] = max(maxOwnedRev[id], rev)
	}
	preCreateRevisions := make([]string, 0, len(owned))
	for _, id := range owned {
		nextRev := maxOwnedRev[id] + 1
		preCreateRevisions = append(preCreateRevisions,
			coresecrets.RevisionName(id, nextRev))
	}
	err = broker.precreateSecretRevs(ctx, preCreateRevisions)
	if err != nil {
		return nil, errors.Trace(err)
	}

	writeRevs := slices.Concat(ownedRevs.RevisionIDs(), preCreateRevisions)
	token, err := broker.createSecretAccessToken(
		ctx, issuedTokenUUID, consumer, writeRevs, readRevs.RevisionIDs(),
	)
	if err != nil {
		return nil, errors.Trace(err)
	}
}

// TODO: make this configurable.
const (
	minExpireSeconds = 600
)

func (k *kubernetesClient) createServiceAccount(ctx context.Context, sa *core.ServiceAccount) (*core.ServiceAccount, error) {
	if k.namespace == "" {
	return out, errors.Trace(err)
}

func (k *kubernetesClient) deleteServiceAccount(ctx context.Context, name string, uid types.UID) error {
	if k.namespace == "" {
		return errNoNamespace
	return nil
}

func (k *kubernetesClient) deleteSecrets(ctx context.Context) error {
	if k.namespace == "" {
		return errNoNamespace
	return nil
}

func (k *kubernetesClient) createRole(
	ctx context.Context, role *rbacv1.Role,
) (*rbacv1.Role, error) {
	if k.namespace == "" {
		return nil, errNoNamespace
	}
	out, err := k.client.RbacV1().Roles(k.namespace).Create(
		ctx, role, v1.CreateOptions{FieldManager: resources.JujuFieldManager})
	if k8serrors.IsAlreadyExists(err) {
		return nil, errors.AlreadyExistsf("role %q", role.GetName())
	}
	return out, errors.Trace(err)
}

func (k *kubernetesClient) deleteRoles(ctx context.Context) error {
	if k.namespace == "" {
		return errNoNamespace
	return nil
}

func (k *kubernetesClient) updateRole(
	ctx context.Context, role *rbacv1.Role,
) (*rbacv1.Role, error) {
	api := k.client.RbacV1().Roles(k.namespace)

	var out *rbacv1.Role
	err := retry.Call(retry.CallArgs{
		Func: func() error {
			patch := map[string]interface{}{
				"rules": role.Rules,
			}
			data, err := json.Marshal(patch)
			if err != nil {
				return errors.Annotatef(err, "marshaling role patch")
			}
			out, err = api.Patch(
				ctx, role.GetName(), types.StrategicMergePatchType, data,
				v1.PatchOptions{
					FieldManager: resources.JujuFieldManager,
				},
			)
			if k8serrors.IsNotFound(err) {
				return errors.NotFoundf("role %q", role.GetName())
			}
			return errors.Annotatef(err, "patching role %q", role.GetName())
		},
		IsFatalError: func(err error) bool {
			return !k8serrors.IsConflict(err)
		},
		Clock:       jujuclock.WallClock,
		Attempts:    5,
		Delay:       time.Second,
		BackoffFunc: retry.ExpBackoff(time.Second, 5*time.Second, 1.5, true),
	})

	return out, errors.Annotatef(err, "updating role %q", role.GetName())
}

func (k *kubernetesClient) deleteRole(ctx context.Context, name string, uid types.UID) error {
	if k.namespace == "" {
		return errNoNamespace
	return errors.Trace(err)
}

func (k *kubernetesClient) createRoleBinding(
	ctx context.Context, rb *rbacv1.RoleBinding,
) (_ *rbacv1.RoleBinding, cleanups []func(), err error) {
	if k.namespace == "" {
		return nil, nil, errNoNamespace
	}

	api := k.client.RbacV1().RoleBindings(k.namespace)
	out, err := api.Create(ctx, rb, v1.CreateOptions{
		FieldManager: resources.JujuFieldManager,
	})
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	cleanups = append(cleanups, func() {
		_ = k.deleteRoleBinding(ctx, out.GetName(), out.GetUID())
	})

	return out, cleanups, nil
}

func (k *kubernetesClient) deleteRoleBinding(
	ctx context.Context, name string, uid types.UID,
) error {
	if k.namespace == "" {
		return errNoNamespace
	}
	err := k.client.RbacV1().RoleBindings(k.namespace).Delete(
		ctx, name, utils.NewPreconditionDeleteOptions(uid))
	if k8serrors.IsNotFound(err) {
		return nil
	}
	return errors.Trace(err)
}

// policyRulesForSecretAccess returns the full policy rules required for
// secrets.
func policyRulesForSecretAccess(
	namespace string, owned, read []string,
) []rbacv1.PolicyRule {
	rules := []rbacv1.PolicyRule{{
		APIGroups:     []string{rbacv1.APIGroupAll},
		Resources:     []string{"namespaces"},
		Verbs:         []string{"get", "list"},
		ResourceNames: []string{namespace},
	}}
	if len(owned) > 0 {
		// owned cannot be empty, otherwise this policy rule grants access to
		// all secrets.
		rules = append(rules, rbacv1.PolicyRule{
			APIGroups: []string{rbacv1.APIGroupAll},
			Resources: []string{"secrets"},
			Verbs: []string{
				// NOTE: create is not given here as it cannot be enforced due
				// to kubernetes rbac limitation.
				"get", "patch", "update", "replace", "delete",
			},
			ResourceNames: owned,
		})
	}
	if len(read) > 0 {
		// read cannot be empty, otherwise this policy rule grants access to
		// all secrets.
		rules = append(rules, rbacv1.PolicyRule{
			APIGroups:     []string{rbacv1.APIGroupAll},
			Resources:     []string{"secrets"},
			Verbs:         []string{"get"},
			ResourceNames: read,
		})
	}
	return rules
}

func (k *kubernetesClient) createRoleAndBinding(
	ctx context.Context, sa *core.ServiceAccount, rules []rbacv1.PolicyRule,
) (cleanups []func(), _ error) {
	role, err := k.createRole(ctx,
		&rbacv1.Role{
				Labels:      sa.Labels,
				Annotations: sa.Annotations,
			},
			Rules: rules,
		},
	)
	if err != nil {
		return cleanups, errors.Annotatef(err, "creating role %q", sa.Name)
	}
	cleanups = append(cleanups, func() {
		_ = k.deleteRole(ctx, role.GetName(), role.GetUID())
	})

	rb := &rbacv1.RoleBinding{
		ObjectMeta: v1.ObjectMeta{
			},
		},
	}
	out, rbCleanups, err := k.createRoleBinding(ctx, rb)
	if err != nil {
		return cleanups, errors.Trace(err)
	}
	cleanups = append(cleanups, rbCleanups...)

	// Ensure role binding exists before we return to avoid a race where a
	// client attempts to perform an operation before the role is allowed.
	return cleanups, errors.Trace(retry.Call(retry.CallArgs{
		Func: func() error {
			api := k.client.RbacV1().RoleBindings(k.namespace)
			_, err := api.Get(ctx, out.Name, v1.GetOptions{
				ResourceVersion: out.ResourceVersion,
			})
			if k8serrors.IsNotFound(err) {
				return errors.NewNotFound(err, "k8s")
			}
	return nil
}

// ensureControllerClusterBindingForSecretAccessToken creates the cluster role
// and role binding needed to access the supplied secrets for the controller.
// If a new cluster role is created, cleanups contain funcs than can be run to
// delete any new resources on error.
func (k *kubernetesClient) createClusterRoleAndBinding(
	ctx context.Context, sa *core.ServiceAccount,
	rules []rbacv1.PolicyRule,
) (cleanups []func(), _ error) {
	cr, err := k.createClusterRole(ctx,
		&rbacv1.ClusterRole{
			ObjectMeta: v1.ObjectMeta{
				Name:        sa.Name,
				Labels:      sa.Labels,
				Annotations: sa.Annotations,
			},
			Rules: rules,
		},
	)
	if err != nil {
		return cleanups, errors.Annotatef(
			err, "creating cluster role %q", sa.Name)
	}
	cleanups = append(cleanups, func() {
		_ = k.deleteClusterRole(ctx, cr.GetName(), cr.GetUID())
	})

	crb, err := k.createClusterRoleBinding(ctx,
		&rbacv1.ClusterRoleBinding{
			ObjectMeta: v1.ObjectMeta{
				Name:        sa.Name,
				Labels:      sa.Labels,
				Annotations: sa.Annotations,
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     sa.Name,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      sa.Name,
					Namespace: sa.Namespace,
				},
			},
		},
	)
	if err != nil {
		return cleanups, errors.Annotatef(
			err, "creating cluster role binding %q", sa.Name)
	}
	cleanups = append(cleanups, func() {
		_ = k.deleteClusterRoleBinding(ctx, crb.GetName(), crb.GetUID())
	})

	// Ensure role binding exists before we return to avoid a race where a
	// client attempts to perform an operation before the role is allowed.
	return cleanups, errors.Trace(retry.Call(retry.CallArgs{
		Func: func() error {
			api := k.client.RbacV1().ClusterRoleBindings()
			_, err := api.Get(ctx, crb.Name, v1.GetOptions{
				ResourceVersion: crb.ResourceVersion,
			})
			if k8serrors.IsNotFound(err) {
				return errors.NewNotFound(err, "k8s")
			}
	clusterResourcePrefix = "juju-secrets-"
)

// precreateSecretRevs ensures that a secret exists for a secret revision.
func (k *kubernetesClient) precreateSecretRevs(
	ctx context.Context, revs []string,
) error {
	labels := labelsForSecretRevision(k.modelName, k.modelUUID)
	client := k.client.CoreV1().Secrets(k.namespace)
	existingSecrets, err := client.List(ctx, v1.ListOptions{
		LabelSelector: labels.AsSelector().String(),
	})
	if err != nil {
		return errors.Trace(err)
	}

	existing := set.NewStrings()
	for _, secret := range existingSecrets.Items {
		existing.Add(secret.Name)
	}

	tmpl := &core.Secret{
		ObjectMeta: v1.ObjectMeta{
			Namespace: k.namespace,
			Labels:    labels,
		},
		Type: core.SecretTypeOpaque,
	}
	for _, name := range revs {
		if existing.Contains(name) {
			continue
		}
		tmpl.Name = name
		_, err := client.Create(ctx, tmpl, v1.CreateOptions{
			FieldManager: resources.JujuFieldManager,
		})
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func (k *kubernetesClient) createSecretAccessToken(
	ctx context.Context,
	issuedTokenUUID string,
	consumer names.Tag,
	ownedRevs []string,
	readRevs []string,
) (_ string, err error) {
	var cleanups []func()
	defer func() {
		}
	}()

	expireAt := time.Now().Add(coresecrets.IssuedTokenValidity)

	labels := labelsForServiceAccount(k.modelName, k.modelUUID, consumer)
	annotations := map[string]string{
		controllerIdKey:              k.controllerUUID,
		modelIdKey:                   k.modelUUID,
		annotationJujuSecretExpireAt: strconv.FormatInt(expireAt.Unix(), 10),
	}

	appName := consumer.Id()
			constants.LabelKubernetesAppName: appName,
		})

	// Service Account name and all the ACLs for this SA are derrived from the
	// issued token UUID. This allows juju to revoke the issued token and
	// perform cleanup of tokens.
	serviceAccountName := fmt.Sprintf(
		"juju-secret-consumer-%s", issuedTokenUUID,
	)

	automountServiceAccountToken := true
	sa := &core.ServiceAccount{
		ObjectMeta: v1.ObjectMeta{
			Name:        serviceAccountName,
			Labels:      labels,
			Annotations: annotations,
			Namespace:   k.namespace,
		},
		AutomountServiceAccountToken: &automountServiceAccountToken,
	}
	sa, err = k.createServiceAccount(ctx, sa)
	if err != nil {
		return "", errors.Annotatef(err, "cannot ensure service account %q", serviceAccountName)
	}
	cleanups = append(cleanups, func() {
		_ = k.deleteServiceAccount(ctx, sa.Name, sa.UID)
	})

	rules := policyRulesForSecretAccess(k.namespace, ownedRevs, readRevs)
	rCleanups, err := k.createRoleAndBinding(ctx, sa, rules)
	cleanups = append(cleanups, rCleanups...)
	if err != nil {
		return "", errors.Annotatef(err, "cannot ensure role binding for secret access token for %q", sa.Name)
	}

	if k.isControllerModel {
		// We need to be able to list/get all namespaces for units in controller
		// model.
		clusterRules := append([]rbacv1.PolicyRule{{
			APIGroups: []string{rbacv1.APIGroupAll},
			Resources: []string{"namespaces"},
			Verbs:     []string{"get", "list"},
		}}, rules...)
		cbCleanups, err := k.createClusterRoleAndBinding(ctx, sa, clusterRules)
		cleanups = append(cleanups, cbCleanups...)
		if err != nil {
			return "", errors.Annotatef(err, "cannot ensure cluster binding for secret access token for %q", sa.Name)
		}
	}

	treq := &authenticationv1.TokenRequest{
		ObjectMeta: v1.ObjectMeta{
			Name: sa.Name,
		},
		Spec: authenticationv1.TokenRequestSpec{
			ExpirationSeconds: func() *int64 {
				until := time.Until(expireAt)
				seconds := int64(max(minExpireSeconds, until.Seconds()))
				return &seconds
			}(),
		},
	}
	tr, err := k.client.CoreV1().ServiceAccounts(k.namespace).CreateToken(
	return tr.Status.Token, nil
}

// revokeSecretAccessTokens removes all the roles, role bindings and service
// accounts related to the named issued token UUID.
func (k *kubernetesClient) revokeSecretAccessToken(
	ctx context.Context, issuedTokenUUID string,
) error {
	if k.namespace == "" {
		return errNoNamespace
	}

	serviceAccountName := fmt.Sprintf(
		"juju-secret-consumer-%s", issuedTokenUUID,
	)

	err := k.client.RbacV1().ClusterRoleBindings().Delete(
		ctx, serviceAccountName, *v1.NewDeleteOptions(0))
	if err != nil && !k8serrors.IsNotFound(err) {
		return errors.Trace(err)
	}

	err = k.client.RbacV1().ClusterRoles().Delete(
		ctx, serviceAccountName, v1.DeleteOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		return errors.Trace(err)
	}

	err = k.client.RbacV1().RoleBindings(k.namespace).Delete(
		ctx, serviceAccountName, *v1.NewDeleteOptions(0))
	if err != nil && !k8serrors.IsNotFound(err) {
		return errors.Trace(err)
	}

	err = k.client.RbacV1().Roles(k.namespace).Delete(
		ctx, serviceAccountName, v1.DeleteOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		return errors.Trace(err)
	}

	err = k.client.CoreV1().ServiceAccounts(k.namespace).Delete(
		ctx, serviceAccountName, v1.DeleteOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		return errors.Trace(err)
	}

	return nil
}

// filterRemovedSecretsPolicyRules removes from the given rules access to the
// specified secret revisions.
func filterRemovedSecretsPolicyRules(
	rules []rbacv1.PolicyRule, removed []string,
) []rbacv1.PolicyRule {
	toRemove := set.NewStrings(removed...)
	var out []rbacv1.PolicyRule
	for _, rule := range rules {
		if slices.Contains(rule.Resources, "secrets") {
			rule.ResourceNames = slices.DeleteFunc(
				rule.ResourceNames, toRemove.Contains)
			if len(rule.ResourceNames) == 0 {
				continue
			}
		}
		out = append(out, rule)
	}
	return out
}

func (k *kubernetesClient) dropSecretAccess(
	ctx context.Context, removed []string,
) error {
	labels := labelsForServiceAccount(k.modelName, k.modelUUID, nil)

	listOps := v1.ListOptions{
		LabelSelector: labels.AsSelector().String(),
	}

	clusterRoles, err := k.client.RbacV1().ClusterRoles().List(ctx, listOps)
	if err != nil {
		return errors.Trace(err)
	}
	for _, clusterRole := range clusterRoles.Items {
		clusterRole.Rules = filterRemovedSecretsPolicyRules(
			clusterRole.Rules, removed)
		_, err := k.updateClusterRole(ctx, &clusterRole)
		if err != nil {
			return errors.Trace(err)
		}
	}

	roles, err := k.client.RbacV1().Roles(k.namespace).List(ctx, listOps)
	if err != nil {
		return errors.Trace(err)
	}
	for _, role := range roles.Items {
		role.Rules = filterRemovedSecretsPolicyRules(role.Rules, removed)
		_, err := k.updateRole(ctx, &role)
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

// ensureDisambiguatedClusterRole creates a cluster role with a disambiguated name.
// cleanups contain funcs than can be run to delete any new resources on error.
func (k *kubernetesClient) ensureDisambiguatedClusterRole(
	ctx context.Context, baseName string, labels labels.Set, annotations map[string]string, rules []rbacv1.PolicyRule,
) (_ *rbacv1.ClusterRole, cleanups []func(), _ error) {
	listOps := v1.ListOptions{
		LabelSelector: modelLabelSelector(k.modelName).String(),
		if clusterRole.Annotations[modelIdKey] != k.modelUUID {
			continue
		}
		clusterRole.Rules = rules
		result, err := k.updateClusterRole(ctx, &clusterRole)
		if err != nil {
			return nil, cleanups, errors.Trace(err)
		} else if !k8serrors.IsNotFound(err) {
			return nil, cleanups, errors.Annotatef(err, "getting existing cluster role %q", proposedName)
		}
	}
}

	validForSeconds := int64(validFor.Truncate(time.Second).Seconds())

	treq := &authenticationv1.TokenRequest{
		ObjectMeta: v1.ObjectMeta{
			Name: broker.serviceAccount,
		},
		Spec: authenticationv1.TokenRequestSpec{
			ExpirationSeconds: &validForSeconds,
		},

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	return nil
}

// getAllNamespacesCustomResourceDefinitionClient returns a dynamic resource
// client for the given CRD and version that operates everywhere. For namespaced
// CRDs this returns the unscoped NamespaceableResourceInterface.
func (k *kubernetesClient) getAllNamespacesCustomResourceDefinitionClient(
	crd *apiextensionsv1.CustomResourceDefinition,
	version string,
) (dynamic.NamespaceableResourceInterface, error) {
	if version == "" {
		return nil, errors.NotValidf(
			"empty version for custom resource definition %q", crd.GetName(),
		)
	}
	found := false
	for _, v := range crd.Spec.Versions {
		if !v.Served {
			continue
		}
		if version == v.Name {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.NotValidf(
			"custom resource definition %s %s is not a supported and served version",
			crd.GetName(), version,
		)
	}
	return k.dynamicClient().Resource(schema.GroupVersionResource{
		Group:    crd.Spec.Group,
		Version:  version,
		Resource: crd.Spec.Names.Plural,
	}), nil
}

// removeAllCustomResourceFinalizers lists all CRs everywhere that matches
// the selector, and patches each one to remove all finalisers. This must be
// done before deletion so that resources with finalisers are not left stuck
// in a terminating state.
func (k *kubernetesClient) removeAllCustomResourceFinalizers(
	ctx context.Context, selector k8slabels.Selector,
) error {
	client := k.extendedClient().ApiextensionsV1().CustomResourceDefinitions()
	crds, err := client.List(ctx, metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		return errors.Trace(err)
	}
	// finalizersPatch is the merge-patch payload that clears all finalizers.
	finalizersPatch, err := json.Marshal(map[string]any{
		"metadata": map[string]any{
			"finalizers": []string{},
		},
	})
	if err != nil {
		return errors.Trace(err)
	}
	patchAll := func(
		crd *apiextensionsv1.CustomResourceDefinition, versionName string,
	) error {
		crdClient, err := k.getAllNamespacesCustomResourceDefinitionClient(
			crd, versionName)
		if err != nil {
			return errors.Trace(err)
		}
		list, err := crdClient.List(ctx, metav1.ListOptions{
			// CRs might be provisioned by another application/charm from a different model.
			LabelSelector: "",
		})
		if err != nil && !k8serrors.IsNotFound(err) {
			return errors.Trace(err)
		}
		if list == nil {
			return nil
		}
		for _, cr := range list.Items {
			if len(cr.GetFinalizers()) == 0 {
				continue
			}
			client := dynamic.ResourceInterface(crdClient)
			if isCRDScopeNamespaced(crd.Spec.Scope) && cr.GetNamespace() != "" {
				client = crdClient.Namespace(cr.GetNamespace())
			}
			_, err = client.Patch(
				context.TODO(),
				cr.GetName(),
				types.MergePatchType,
				finalizersPatch,
				metav1.PatchOptions{},
			)
			if err != nil && !k8serrors.IsNotFound(err) {
				return errors.Annotatef(
					err, "removing finalizers from custom resource %q (namespace %q)",
					cr.GetName(), cr.GetNamespace(),
				)
			}
		}
		return nil
	}
	for _, crd := range crds.Items {
		if selector.Empty() {
			continue
		}
		for _, version := range crd.Spec.Versions {
			if !version.Served {
				continue
			}
			err := patchAll(&crd, version.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// deleteAllCustomResourcesAllNamespaces deletes custom resources matching the
// supplied selector everywhere.
func (k *kubernetesClient) deleteAllCustomResourcesAllNamespaces(
	ctx context.Context, selector k8slabels.Selector,
) error {
	client := k.extendedClient().ApiextensionsV1().CustomResourceDefinitions()
	crds, err := client.List(ctx, metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		return errors.Trace(err)
	}
	for _, crd := range crds.Items {
		if selector.Empty() {
			continue
		}
		for _, version := range crd.Spec.Versions {
			crdClient, err := k.getAllNamespacesCustomResourceDefinitionClient(
				&crd, version.Name)
			if err != nil {
				return errors.Trace(err)
			}
			err = crdClient.DeleteCollection(ctx, metav1.DeleteOptions{
				PropagationPolicy: constants.DefaultPropagationPolicy(),
			}, metav1.ListOptions{
				// CRs might be provisioned by another application/charm from a different model.
				LabelSelector: "",
			})
			if err != nil && !k8serrors.IsNotFound(err) {
				return errors.Trace(err)
			}
		}
	}
	return nil
}

// listAllCustomResourcesAllNamespaces lists custom resources matching the
// selector everywhere.
func (k *kubernetesClient) listAllCustomResourcesAllNamespaces(
	ctx context.Context, selector k8slabels.Selector,
) (out []unstructured.Unstructured, err error) {
	client := k.extendedClient().ApiextensionsV1().CustomResourceDefinitions()
	crds, err := client.List(ctx, metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		return nil, errors.Trace(err)
	}
	for _, crd := range crds.Items {
		if selector.Empty() {
			continue
		}
		for _, version := range crd.Spec.Versions {
			crdClient, err := k.getAllNamespacesCustomResourceDefinitionClient(
				&crd, version.Name)
			if err != nil {
				return nil, errors.Trace(err)
			}
			list, err := crdClient.List(ctx, metav1.ListOptions{
				// CRs might be provisioned by another application/charm from a different model.
				LabelSelector: "",
			})
			if err != nil && !k8serrors.IsNotFound(err) {
				return nil, errors.Trace(err)
			}
			if list != nil {
				out = append(out, list.Items...)
			}
		}
	}
	if len(out) == 0 {

// SetSnapConfig sets a snap's key to value.
func SetSnapConfig(snap string, key string, value string) error {
	logger.Infof("setting snap %q config key %q to value %q", snap, key, value)
	if key == "" {
		logger.Warningf("set snap config called with empty key for snap %q", snap)
		return errors.NotValidf("key must not be empty")
	}

	cmd := exec.Command(Command, "set", snap, fmt.Sprintf("%s=%s", key, value))
	_, err := cmd.Output()
	if err != nil {
		logger.Errorf("failed to set snap %q config %q=%q: %v", snap, key, value, err)
		return errors.Annotate(err, fmt.Sprintf("setting snap %s config %s to %s", snap, key, value))
	}

	logger.Infof("successfully set snap %q config %q", snap, key)
	return nil
}

// If no BackgroundServices are provided, Service will wrap all of the snap's
// background services.
func NewService(config ServiceConfig) (Service, error) {
	logger.Infof("creating new snap service %q (path=%q, channel=%q, confinement=%q)",
		config.ServiceName, config.SnapPath, config.Channel, config.ConfinementPolicy)
	if config.ServiceName == "" {
		logger.Warningf("NewService called with empty ServiceName")
		return Service{}, errors.New("ServiceName must be provided")
	}
	app := &App{
	}
	err := app.Validate()
	if err != nil {
		logger.Warningf("snap app validation failed for %q: %v", config.ServiceName, err)
		return Service{}, errors.Trace(err)
	}

	isLocal := config.SnapPath != ""
	logger.Debugf("snap service %q: isLocal=%v, configDir=%q, executable=%q",
		config.ServiceName, isLocal, config.ConfigDir, config.SnapExecutable)

	return Service{
		runnable:       defaultRunner{},
// Validate validates that snap.Service has been correctly configured.
// Validate returns nil when successful and an error when successful.
func (s Service) Validate() error {
	logger.Debugf("validating snap service %q", s.name)
	if err := s.app.Validate(); err != nil {
		logger.Warningf("snap service %q app validation failed: %v", s.name, err)
		return errors.Trace(err)
	}

	for _, prerequisite := range s.app.Prerequisites() {
		if err := prerequisite.Validate(); err != nil {
			logger.Warningf(
				"snap service %q prerequisite %q validation failed: %v",
				s.name, prerequisite.Name(), err,
			)
			return errors.Trace(err)
		}
	}

	logger.Debugf("snap service %q validation successful", s.name)
	return nil
}


// Running returns (true, nil) when snap indicates that service is currently active.
func (s Service) Running() (bool, error) {
	logger.Debugf("checking if snap service %q is running", s.name)
	_, _, running, err := s.status()
	if err != nil {
		logger.Warningf("failed to check running status for snap service %q: %v", s.name, err)
		return false, errors.Trace(err)
	}
	logger.Debugf("snap service %q running=%v", s.name, running)
	return running, nil
}


// Install installs the snap and its background services.
func (s Service) Install() error {
	logger.Infof("installing snap service %q (isLocal=%v)", s.name, s.isLocal)
	prerequisites := s.app.Prerequisites()
	logger.Infof("snap service %q has %d prereq(s) to install", s.name, len(prerequisites))
	for i, app := range prerequisites {
		logger.Infof("installing prerequisite %d/%d: %q", i+1, len(prerequisites), app.Name())

		out, err := s.installAppWithRetry(app)
		if err != nil {
			logger.Errorf(
				"failed to install prereq %q for snap service %q: %v (output: %v)",
				app.Name(), s.name, err, out,
			)
			return errors.Annotatef(err, "output: %v", out)
		}
		logger.Infof("successfully installed prerequisite %q", app.Name())
	}

	logger.Infof("installing snap app %q with args: %v", s.app.Name(), s.app.InstallArgs())
	out, err := s.installAppWithRetry(s.app)
	if err != nil {
		logger.Errorf("failed to install snap service %q: %v (output: %v)", s.name, err, out)
		return errors.Annotatef(err, "output: %v", out)
	}
	logger.Infof("successfully installed snap service %q", s.name)
	return nil
}

func (s Service) installAppWithRetry(app Installable) (string, error) {
	ackAsserts := app.AcknowledgeAssertsArgs()
	if ackAsserts != nil {
		logger.Infof("acknowledging asserts for snap %q: %v", app.Name(), ackAsserts)
		_, err := s.runCommandWithRetry(ackAsserts...)
		if err != nil {
			logger.Errorf("failed to acknowledge asserts for snap %q: %v", app.Name(), err)
			return "", errors.Trace(err)
		}
		logger.Infof("successfully acknowledged asserts for snap %q", app.Name())
	} else {
		logger.Debugf("no asserts to acknowledge for snap %q", app.Name())
	}

	logger.Infof("running install command for snap %q with args: %v", app.Name(), app.InstallArgs())
	return s.runCommandWithRetry(app.InstallArgs()...)
}

// Installed returns true if the service has been successfully installed.
func (s Service) Installed() (bool, error) {
	logger.Debugf("checking if snap service %q is installed", s.name)
	installed, _, _, err := s.status()
	if err != nil {
		logger.Warningf("failed to check installed status for snap service %q: %v", s.name, err)
		return false, errors.Trace(err)
	}
	logger.Debugf("snap service %q installed=%v", s.name, installed)
	return installed, nil
}

// ConfigOverride writes a systemd override to enable the
// specified limits to be used by the snap.
func (s Service) ConfigOverride() error {
	logger.Debugf(
		"applying config overrides for snap service %q (limits count: %d)",
		s.name, len(s.conf.Limit),
	)
	if len(s.conf.Limit) == 0 {
		logger.Debugf("no config limits defined for snap service %q, skipping override", s.name)
		return nil
	}

	unitOptions := systemd.ServiceLimits(s.conf)
	data, err := io.ReadAll(systemd.UnitSerialize(unitOptions))
	if err != nil {
		logger.Errorf("failed to serialise systemd unit options for snap service %q: %v", s.name, err)
		return errors.Trace(err)
	}

	backgroundServices := s.app.BackgroundServices()
	logger.Infof(
		"writing config overrides for %d background services of snap %q",
		len(backgroundServices), s.name,
	)
	for _, backgroundService := range backgroundServices {
		overridesDir := fmt.Sprintf("%s/snap.%s.%s.service.d", s.configDir, s.name, backgroundService.Name)
		logger.Debugf(
			"creating overrides directory %q for background service %q",
			overridesDir, backgroundService.Name,
		)
		if err := os.MkdirAll(overridesDir, 0755); err != nil {
			logger.Errorf("failed to create overrides directory %q: %v", overridesDir, err)
			return errors.Trace(err)
		}
		overridePath := filepath.Join(overridesDir, "overrides.conf")
		logger.Debugf("writing overrides config to %q", overridePath)
		if err := os.WriteFile(overridePath, data, 0644); err != nil {
			logger.Errorf("failed to write overrides config to %q: %v", overridePath, err)
			return errors.Trace(err)
		}
	}
	logger.Infof("successfully applied config overrides for snap service %q", s.name)
	return nil
}

// shell commands to be executed by a shell which start the service.
func (s Service) StartCommands() ([]string, error) {
	deps := s.app.Prerequisites()
	logger.Debugf(
		"generating start commands for snap service %q (%d prerequisites)",
		s.name, len(deps),
	)
	commands := make([]string, 0, 1+len(deps))
	for _, prerequisite := range deps {
		cmds := prerequisite.StartCommands(s.executable)
		logger.Debugf("prerequisite %q start commands: %v", prerequisite.Name(), cmds)
		commands = append(commands, cmds...)
	}
	appCmds := s.app.StartCommands(s.executable)
	logger.Debugf("snap service %q start commands: %v", s.name, appCmds)
	commands = append(commands, appCmds...)
	logger.Debugf("total start commands for snap service %q: %v", s.name, commands)
	return commands, nil
}

// status returns an interpreted output from the `snap services` command.
//
//	(true, true, false, nil)
func (s *Service) status() (isInstalled, enabledAtStartup, isCurrentlyActive bool, err error) {
	logger.Debugf("querying status for snap service %q", s.Name())
	out, err := s.runCommand("services", s.Name())
	if err != nil {
		logger.Warningf("failed to query snap services for %q: %v", s.Name(), err)
		return false, false, false, errors.Trace(err)
	}
	logger.Debugf("snap services output for %q: %q", s.Name(), out)
	for _, line := range strings.Split(out, "\n") {
		if !strings.HasPrefix(line, s.Name()) {
			continue
		}

		fields := strings.Fields(line)
		installed := true
		enabled := fields[1] == "enabled"
		active := fields[2] == "active"
		logger.Debugf(
			"snap service %q status: installed=%v, enabledAtStartup=%v, active=%v",
			s.Name(), installed, enabled, active,
		)
		return installed, enabled, active, nil
	}

	logger.Debugf("snap service %q not found in services output", s.Name())
	return false, false, false, nil
}

// Start starts the service, returning nil when successful.
// If the service is already running, Start does not restart it.
func (s Service) Start() error {
	logger.Infof("starting snap service %q", s.name)
	running, err := s.Running()
	if err != nil {
		return errors.Trace(err)
	}
	if running {
		logger.Debugf("snap service %q is already running, skipping start", s.name)
		return nil
	}

	commands, err := s.StartCommands()
	if err != nil {
		logger.Errorf("failed to get start commands for snap service %q: %v", s.name, err)
		return errors.Trace(err)
	}
	logger.Infof("executing %d start commands for snap service %q", len(commands), s.name)
	for i, command := range commands {
		logger.Infof(
			"executing start command %d/%d for snap service %q: %q",
			i, len(commands), s.name, command,
		)
		commandParts := strings.Fields(command)
		out, err := utils.RunCommand(commandParts[0], commandParts[1:]...)
		if err != nil {
			if strings.Contains(out, "has no services") {
				logger.Debugf("snap %q has no services, skipping command %q", s.name, command)
				continue
			}
			logger.Errorf(
				"start command failed for snap service %q: %q -> %v (output: %v)",
				s.name, command, err, out,
			)
			return errors.Annotatef(err, "%v -> %v", command, out)
		}
		logger.Debugf(
			"start command %d/%d completed successfully for snap service %q (output: %q)",
			i, len(commands), s.name, out,
		)
	}

	logger.Infof("successfully started snap service %q", s.name)
	return nil
}

// Stop stops a running service. Returns nil when the underlying
// call to `snap stop <service-name>` exits with error code 0.
func (s Service) Stop() error {
	logger.Infof("stopping snap service %q", s.name)
	running, err := s.Running()
	if err != nil {
		return errors.Trace(err)
	}
	if !running {
		logger.Debugf("snap service %q is not running, skipping stop", s.name)
		return nil
	}

	args := []string{"stop", s.Name()}
	if err := s.execThenExpect(args, "Stopped."); err != nil {
		logger.Errorf("failed to stop snap service %q: %v", s.name, err)
		return err
	}
	logger.Infof("successfully stopped snap service %q", s.name)
	return nil
}

// Restart restarts the service, or starts if it's not currently
//
// Restart is part of the service.RestartableService interface
func (s Service) Restart() error {
	logger.Infof("restarting snap service %q", s.name)
	args := []string{"restart", s.Name()}
	if err := s.execThenExpect(args, "Restarted."); err != nil {
		logger.Errorf("failed to restart snap service %q: %v", s.name, err)
		return err
	}
	logger.Infof("successfully restarted snap service %q", s.name)
	return nil
}

// execThenExpect calls `snap <commandArgs>...` and then checks
// stdout against expectation and snap's exit code. When there's a
// mismatch or non-0 exit code, execThenExpect returns an error.
func (s Service) execThenExpect(commandArgs []string, expectation string) error {
	logger.Debugf("executing snap command %v, expecting %q", commandArgs, expectation)
	out, err := s.runCommand(commandArgs...)
	if err != nil {
		logger.Errorf("snap command %v failed: %v", commandArgs, err)
		return errors.Trace(err)
	}
	if !strings.Contains(out, expectation) {
		logger.Errorf(
			"snap command %v: expected %q in output, got %q",
			commandArgs, expectation, out,
		)
		return errors.Annotatef(err, `expected "%s", got "%s"`, expectation, out)
	}
	logger.Debugf("snap command %v output matched expectation %q", commandArgs, expectation)
	return nil
}

}

func (s Service) runCommandWithRetry(args ...string) (res string, err error) {
	const delay = 5 * time.Second
	const attempts = 2
	logger.Debugf(
		"running snap command with retry: %v (delay=%v, attempts=%v)",
		args, delay, attempts,
	)
	attempt := 0
	if resErr := retry.Call(retry.CallArgs{
		Clock: s.clock,
		Func: func() error {
			attempt++
			logger.Debugf("snap command attempt %d: %v", attempt, args)
			res, err = s.runCommand(args...)
			if err != nil {
				logger.Warningf(
					"snap command attempt %d failed: %v (output: %q)",
					attempt, err, res,
				)
			}
			return errors.Trace(err)
		},
		Delay:    delay,
		Attempts: attempts,
	}); resErr != nil {
		logger.Errorf("snap command %v failed after %d attempts: %v", args, attempt, resErr)
		return "", errors.Trace(resErr)
	}

	logger.Debugf("snap command %v succeeded on attempt %d (output: %q)", args, attempt, res)
	// Named args are set via the retry.
	return
}
mongo/mongo.go | 8 ++++----
mongo/open.go  | 4 ++--
2 files changed, 6 insertions(+), 6 deletions(-)
