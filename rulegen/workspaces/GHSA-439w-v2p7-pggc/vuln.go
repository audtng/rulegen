package main

	"fmt"
	"net"
	"os"
	"strings"
	"time"


	k8scloud "github.com/juju/juju/caas/kubernetes/cloud"
	"github.com/juju/juju/core/model"
	"github.com/juju/juju/environs/bootstrap"
	"github.com/juju/juju/environs/cloudspec"
	k8sprovider "github.com/juju/juju/internal/provider/kubernetes"
	if err != nil {
		return errors.Trace(err)
	}
	_, err = broker.ensureSecretAccessToken(context.TODO(), tag, nil, nil, removed.RevisionIDs())
	return errors.Trace(err)
}

func cloudSpecToBackendConfig(spec cloudspec.CloudSpec) (*provider.BackendConfig, error) {
	return true
}

// RestrictedConfig returns the config needed to create a
// secrets backend client restricted to manage the specified
// owned secrets and read shared secrets for the given entity tag.
func (p k8sProvider) RestrictedConfig(
	adminCfg *provider.ModelBackendConfig, sameController, forDrain bool, consumer names.Tag, owned provider.SecretRevisions, read provider.SecretRevisions,
) (*provider.BackendConfig, error) {
	logger.Tracef("getting k8s backend config for %q, owned %v, read %v", consumer, owned, read)

	if consumer == nil {
		return &adminCfg.BackendConfig, nil
		return nil, errors.Trace(err)
	}
	ctx := context.TODO()
	token, err := broker.ensureSecretAccessToken(ctx, consumer, owned.RevisionIDs(), read.RevisionIDs(), nil)
	if err != nil {
		return nil, errors.Trace(err)
	}
}

// TODO: make this configurable.
var expiresInSeconds = int64(60 * 10)

func (k *kubernetesClient) createServiceAccount(ctx context.Context, sa *core.ServiceAccount) (*core.ServiceAccount, error) {
	if k.namespace == "" {
	return out, errors.Trace(err)
}

func (k *kubernetesClient) updateServiceAccount(ctx context.Context, sa *core.ServiceAccount) (*core.ServiceAccount, error) {
	if k.namespace == "" {
		return nil, errNoNamespace
	}
	out, err := k.client.CoreV1().ServiceAccounts(k.namespace).Update(ctx, sa, v1.UpdateOptions{})
	if k8serrors.IsNotFound(err) {
		return nil, errors.NotFoundf("service account %q", sa.GetName())
	}
	return out, errors.Trace(err)
}

func (k *kubernetesClient) deleteServiceAccount(ctx context.Context, name string, uid types.UID) error {
	if k.namespace == "" {
		return errNoNamespace
	return nil
}

// ensureServiceAccount creates or updates a service account, disambiguating the name if necessary.
// If a new service account is created, cleanups contain funcs than can be run to delete any new
// resources on error.
func (k *kubernetesClient) ensureServiceAccount(
	ctx context.Context, serviceAccountName string, labels labels.Set, annotations map[string]string, disambiguateName bool,
) (out *core.ServiceAccount, cleanups []func(), err error) {
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
	if disambiguateName {
		out, err = k.createDisambiguatedServiceAccount(ctx, sa)
	} else {
		out, err = k.createServiceAccount(ctx, sa)
		if err != nil && !errors.Is(err, errors.AlreadyExists) {
			return nil, nil, errors.Trace(err)
		}
	}
	if err == nil {
		logger.Debugf("service account %q created", out.GetName())
		cleanups = append(cleanups, func() { _ = k.deleteServiceAccount(ctx, out.GetName(), out.GetUID()) })
		return out, cleanups, nil
	}

	// Service account already exists so update it.
	out, err = k.updateServiceAccount(ctx, sa)
	logger.Debugf("updating service account %q", sa.GetName())
	return out, cleanups, errors.Trace(err)
}

func (k *kubernetesClient) deleteSecrets(ctx context.Context) error {
	if k.namespace == "" {
		return errNoNamespace
	return nil
}

func (k *kubernetesClient) createRole(ctx context.Context, role *rbacv1.Role) (*rbacv1.Role, error) {
	if k.namespace == "" {
		return nil, errNoNamespace
	}
	out, err := k.client.RbacV1().Roles(k.namespace).Create(ctx, role, v1.CreateOptions{FieldManager: resources.JujuFieldManager})
	if k8serrors.IsAlreadyExists(err) {
		return nil, errors.AlreadyExistsf("role %q", role.GetName())
	}
	return out, errors.Trace(err)
}

// updateRole fetches the latest version of the specified role,
// replaces its Rules with those from the provided role, and updates it
// in the cluster. This method retries on conflicts using exponential backoff
// to handle concurrent modifications by other controllers.
// Note that only the Rules field is updated, all other fields from the latest role are preserved.
func (k *kubernetesClient) updateRole(ctx context.Context, role *rbacv1.Role) (*rbacv1.Role, error) {
	if k.namespace == "" {
		return nil, errNoNamespace
	}

	api := k.client.RbacV1().Roles(k.namespace)
	var out *rbacv1.Role
	err := retry.Call(retry.CallArgs{
		Func: func() error {
			patch := map[string]interface{}{
				"rules": role.Rules,
			}
			data, err := json.Marshal(patch)
			if err != nil {
				return errors.Trace(err)
			}
			out, err = api.Patch(ctx, role.GetName(), types.StrategicMergePatchType, data, v1.PatchOptions{
				FieldManager: resources.JujuFieldManager,
			})
			if k8serrors.IsNotFound(err) {
				return errors.NotFoundf("role %q", role.GetName())
			}
			return errors.Trace(err)
		},
		IsFatalError: func(err error) bool {
			return !k8serrors.IsConflict(err)
		},
		Clock:       jujuclock.WallClock,
		Attempts:    5,
		Delay:       time.Second,
		BackoffFunc: retry.ExpBackoff(time.Second, 5*time.Second, 1.5, true),
	})
	return out, errors.Trace(err)
}

func (k *kubernetesClient) getRole(ctx context.Context, name string) (*rbacv1.Role, error) {
	if k.namespace == "" {
		return nil, errNoNamespace
	}
	out, err := k.client.RbacV1().Roles(k.namespace).Get(ctx, name, v1.GetOptions{})
	if k8serrors.IsNotFound(err) {
		return nil, errors.NotFoundf("role %q", name)
	}
	return out, errors.Trace(err)
}

func (k *kubernetesClient) deleteRoles(ctx context.Context) error {
	if k.namespace == "" {
		return errNoNamespace
	return nil
}

func (k *kubernetesClient) deleteRole(ctx context.Context, name string, uid types.UID) error {
	if k.namespace == "" {
		return errNoNamespace
	return errors.Trace(err)
}

func (k *kubernetesClient) ensureRoleBinding(
	ctx context.Context, rb *rbacv1.RoleBinding,
) (_ *rbacv1.RoleBinding, cleanups []func(), err error) {
	if k.namespace == "" {
		return nil, cleanups, errNoNamespace
	}

	api := k.client.RbacV1().RoleBindings(k.namespace)

	out, err := api.Create(ctx, rb, v1.CreateOptions{
		FieldManager: resources.JujuFieldManager,
	})
	if k8serrors.IsAlreadyExists(err) {
		// we need to ensure that the rb is not empty for callers
		// by getting rb from api again eg cases like resource name empty for
		// attempting to get rb name in caller
		out, err = api.Get(ctx, rb.Name, v1.GetOptions{})
		return out, cleanups, err
	}
	if err == nil {
		cleanups = append(cleanups, func() { _ = k.deleteRoleBinding(ctx, out.GetName(), out.GetUID()) })
	}

	return out, cleanups, errors.Trace(err)
}

func (k *kubernetesClient) deleteRoleBinding(ctx context.Context, name string, uid types.UID) error {
	if k.namespace == "" {
		return errNoNamespace
	}
	err := k.client.RbacV1().RoleBindings(k.namespace).Delete(ctx, name, utils.NewPreconditionDeleteOptions(uid))
	if k8serrors.IsNotFound(err) {
		return nil
	}
	return errors.Trace(err)
}

func cleanRules(existing []rbacv1.PolicyRule, shouldRemove func(string) bool) []rbacv1.PolicyRule {
	if len(existing) == 0 {
		return nil
	}

	i := 0
	for _, r := range existing {
		if len(r.ResourceNames) == 1 && shouldRemove(r.ResourceNames[0]) {
			continue
		}
		existing[i] = r
		i++
	}
	return existing[:i]
}

func rulesForSecretAccess(
	namespace string, isControllerModel bool,
	existing []rbacv1.PolicyRule, owned, read, removed []string,
) []rbacv1.PolicyRule {
	if len(existing) == 0 {
		existing = []rbacv1.PolicyRule{
			{
				APIGroups: []string{rbacv1.APIGroupAll},
				Resources: []string{"secrets"},
				Verbs: []string{
					"create",
					"patch", // TODO: we really should only allow "create" but not patch  but currently we uses .Apply() which requres patch!!!
				},
			},
		}
		if isControllerModel {
			// We need to be able to list/get all namespaces for units in controller model.
			existing = append(existing, rbacv1.PolicyRule{
				APIGroups: []string{rbacv1.APIGroupAll},
				Resources: []string{"namespaces"},
				Verbs:     []string{"get", "list"},
			})
		} else {
			// We just need to be able to list/get our own namespace for units in other models.
			existing = append(existing, rbacv1.PolicyRule{
				APIGroups:     []string{rbacv1.APIGroupAll},
				Resources:     []string{"namespaces"},
				Verbs:         []string{"get", "list"},
				ResourceNames: []string{namespace},
			})
		}
	}

	ownedIDs := set.NewStrings(owned...)
	readIDs := set.NewStrings(read...)
	removedIDs := set.NewStrings(removed...)

	existing = cleanRules(existing,
		func(s string) bool {
			return ownedIDs.Contains(s) || readIDs.Contains(s) || removedIDs.Contains(s)
		},
	)

	for _, rName := range owned {
		if removedIDs.Contains(rName) {
			continue
		}
		existing = append(existing, rbacv1.PolicyRule{
			APIGroups:     []string{rbacv1.APIGroupAll},
			Resources:     []string{"secrets"},
			Verbs:         []string{rbacv1.VerbAll},
			ResourceNames: []string{rName},
		})
	}
	for _, rName := range read {
		if removedIDs.Contains(rName) {
			continue
		}
		existing = append(existing, rbacv1.PolicyRule{
			APIGroups:     []string{rbacv1.APIGroupAll},
			Resources:     []string{"secrets"},
			Verbs:         []string{"get"},
			ResourceNames: []string{rName},
		})
	}
	return existing
}

// ensureBindingForSecretAccessToken creates the role and role binding needed to access the supplied secrets.
// If a new role is created, cleanups contain funcs than can be run to delete any new
// resources on error.
func (k *kubernetesClient) ensureBindingForSecretAccessToken(
	ctx context.Context, sa *core.ServiceAccount, owned, read, removed []string,
) (cleanups []func(), _ error) {
	role, err := k.createRole(ctx,
		&rbacv1.Role{
				Labels:      sa.Labels,
				Annotations: sa.Annotations,
			},
			Rules: rulesForSecretAccess(k.namespace, false, nil, owned, read, removed),
		},
	)
	if errors.Is(err, errors.AlreadyExists) {
		role, err = k.getRole(ctx, sa.Name)
		if err != nil {
			return cleanups, errors.Annotatef(err, "getting role %q", sa.Name)
		}
		role.Rules = rulesForSecretAccess(k.namespace, false, role.Rules, owned, read, removed)
		_, err = k.updateRole(ctx, role)
		if err != nil {
			return cleanups, errors.Annotatef(err, "updating role %q", sa.Name)
		}
	} else if err != nil {
		return cleanups, errors.Annotatef(err, "creating role %q", sa.Name)
	} else {
		cleanups = append(cleanups, func() { _ = k.deleteRole(ctx, role.GetName(), role.GetUID()) })
	}

	rb := &rbacv1.RoleBinding{
		ObjectMeta: v1.ObjectMeta{
			},
		},
	}
	out, rbCleanups, err := k.ensureRoleBinding(ctx, rb)
	if err != nil {
		return cleanups, errors.Trace(err)
	}
	cleanups = append(cleanups, rbCleanups...)

	// Ensure role binding exists before we return to avoid a race where a client
	// attempts to perform an operation before the role is allowed.
	return cleanups, errors.Trace(retry.Call(retry.CallArgs{
		Func: func() error {
			api := k.client.RbacV1().RoleBindings(k.namespace)
			_, err := api.Get(ctx, out.Name, v1.GetOptions{ResourceVersion: out.ResourceVersion})
			if k8serrors.IsNotFound(err) {
				return errors.NewNotFound(err, "k8s")
			}
	return nil
}

// ensureClusterBindingForSecretAccessToken creates the cluster role and role binding needed
// to access the supplied secrets.
// If a new cluster role is created, cleanups contain funcs than can be run to delete any new
// resources on error.
func (k *kubernetesClient) ensureClusterBindingForSecretAccessToken(
	ctx context.Context, saName, baseName string, labels labels.Set, annotations map[string]string, owned, read, removed []string,
) (cleanups []func(), _ error) {
	createRules := func(existing []rbacv1.PolicyRule) []rbacv1.PolicyRule {
		return rulesForSecretAccess(k.namespace, true, existing, owned, read, removed)
	}
	clusterRole, crCleanups, err := k.ensureDisambiguatedClusterRole(ctx, baseName, labels, annotations, createRules)
	if err == nil {
		cleanups = append(cleanups, crCleanups...)
	} else {
		return cleanups, errors.Annotatef(err, "disambiguating cluster role name %q", baseName)
	}

	clusterRoleBinding, crbCleanups, err := k.ensureDisambiguatedClusterRoleBinding(ctx, saName, baseName, clusterRole.Name, labels, annotations)
	if err == nil {
		cleanups = append(cleanups, crbCleanups...)
	} else {
		return cleanups, errors.Annotatef(err, "disambiguating cluster role binding name %q", baseName)
	}

	// Ensure role binding exists before we return to avoid a race where a client
	// attempts to perform an operation before the role is allowed.
	return cleanups, errors.Trace(retry.Call(retry.CallArgs{
		Func: func() error {
			api := k.client.RbacV1().ClusterRoleBindings()
			_, err := api.Get(ctx, clusterRoleBinding.Name, v1.GetOptions{ResourceVersion: clusterRoleBinding.ResourceVersion})
			if k8serrors.IsNotFound(err) {
				return errors.NewNotFound(err, "k8s")
			}
	clusterResourcePrefix = "juju-secrets-"
)

func (k *kubernetesClient) ensureSecretAccessToken(
	ctx context.Context, consumer names.Tag, owned, read, removed []string,
) (_ string, err error) {
	var cleanups []func()
	defer func() {
		}
	}()

	labels := labelsForServiceAccount(k.modelName, k.modelUUID)
	annotations := map[string]string{
		controllerIdKey: k.controllerUUID,
		modelIdKey:      k.modelUUID,
	}

	appName := consumer.Id()
			constants.LabelKubernetesAppName: appName,
		})

	// Compose the name of the service account and role and role binding.
	// We'll use the tag string, but for models we'll use the model name, since
	// the UUID will be used to disambiguate anyway if needed.
	baseResourceName := consumer.String()
	if consumer.Kind() == names.ModelTagKind {
		baseResourceName = fmt.Sprintf("model-%s", k.modelName)
	}
	serviceAccountName := baseResourceName
	// For the controller model, the resources are cluster scoped so
	// given them a meaningful prefix.
	if k.isControllerModel {
		baseResourceName = clusterResourcePrefix + baseResourceName
	}
	// If the resources are going to a namespace other than that of the host model,
	// disambiguate the name.
	disambiguateName, err := k.isExternalNamespace()
	if err != nil {
		return "", errors.Annotate(err, "checking if namespace is external")
	}

	sa, saCleanups, err := k.ensureServiceAccount(ctx, serviceAccountName, labels, annotations, disambiguateName)
	cleanups = append(cleanups, saCleanups...)
	if err != nil {
		return "", errors.Annotatef(err, "cannot ensure service account %q", serviceAccountName)
	}

	if k.isControllerModel {
		cbCleanups, err := k.ensureClusterBindingForSecretAccessToken(ctx, sa.Name, baseResourceName, labels, annotations, owned, read, removed)
		cleanups = append(cleanups, cbCleanups...)
		if err != nil {
			return "", errors.Annotatef(err, "cannot ensure cluster binding for secret access token for %q", sa.Name)
		}
	} else {
		// For roles and role bindings created in the namespace set up to hold the secrets,
		// we assume that the service account, role, role binding all share the same disambiguated
		// name as the service account. This is reasonable since it's not expected that anything
		// other than Juju will be messing with such artefacts in that namespace.
		rCleanups, err := k.ensureBindingForSecretAccessToken(ctx, sa, owned, read, removed)
		cleanups = append(cleanups, rCleanups...)
		if err != nil {
			return "", errors.Annotatef(err, "cannot ensure role binding for secret access token for %q", sa.Name)
		}
	}

	treq := &authenticationv1.TokenRequest{
		Spec: authenticationv1.TokenRequestSpec{
			ExpirationSeconds: &expiresInSeconds,
		},
	}
	tr, err := k.client.CoreV1().ServiceAccounts(k.namespace).CreateToken(
	return tr.Status.Token, nil
}

// createDisambiguatedServiceAccount creates a service account with a disambiguated name.
func (k *kubernetesClient) createDisambiguatedServiceAccount(
	ctx context.Context, sa *core.ServiceAccount,
) (*core.ServiceAccount, error) {
	if k.namespace == "" {
		return nil, errNoNamespace
	}

	listOps := v1.ListOptions{
		LabelSelector: modelLabelSelector(k.modelName).String(),
	}
	existing, err := k.client.CoreV1().ServiceAccounts(k.namespace).List(ctx, listOps)
	if err != nil {
		return nil, errors.Trace(err)
	}

	for _, existingServiceAccount := range existing.Items {
		if existingServiceAccount.Annotations[modelIdKey] == k.modelUUID {
			return &existingServiceAccount, nil
		}
	}

	suffixLength := model.DefaultSuffixDigits
	var proposedName string

	for {
		if proposedName, err = model.DisambiguateResourceNameWithSuffixLength(
			k.modelUUID, sa.Name, maxResourceNameLength, suffixLength); err != nil {
			return nil, errors.Annotatef(err, "disambiguating service account name %q", sa.Name)
		}
		_, err = k.client.CoreV1().ServiceAccounts(k.namespace).Get(ctx, proposedName, v1.GetOptions{})
		if err == nil {
			suffixLength = suffixLength + 1
			continue
		} else if !k8serrors.IsNotFound(err) {
			return nil, errors.Annotatef(err, "getting existing service account %q", proposedName)
		}
		sa.Name = proposedName
		return k.createServiceAccount(ctx, sa)
	}
}

// ensureDisambiguatedClusterRole creates a cluster role with a disambiguated name.
// cleanups contain funcs than can be run to delete any new resources on error.
func (k *kubernetesClient) ensureDisambiguatedClusterRole(
	ctx context.Context, baseName string, labels labels.Set, annotations map[string]string, createRules func(existing []rbacv1.PolicyRule) []rbacv1.PolicyRule,
) (_ *rbacv1.ClusterRole, cleanups []func(), _ error) {
	listOps := v1.ListOptions{
		LabelSelector: modelLabelSelector(k.modelName).String(),
		if clusterRole.Annotations[modelIdKey] != k.modelUUID {
			continue
		}
		clusterRole.Rules = createRules(clusterRole.Rules)
		result, err := k.updateClusterRole(ctx, &clusterRole)
		if err != nil {
			return nil, cleanups, errors.Trace(err)
		} else if !k8serrors.IsNotFound(err) {
			return nil, cleanups, errors.Annotatef(err, "getting existing cluster role %q", proposedName)
		}
		result, err := k.createClusterRole(ctx,
			&rbacv1.ClusterRole{
				ObjectMeta: v1.ObjectMeta{
					Name:        proposedName,
					Labels:      labels,
					Annotations: annotations,
				},
				Rules: createRules(nil),
			},
		)
		if errors.Is(err, errors.AlreadyExists) {
			suffixLength++
			continue
		}
		if err == nil {
			cleanups = append(cleanups, func() { _ = k.deleteClusterRole(ctx, result.GetName(), result.GetUID()) })
		}
		return result, cleanups, nil
	}
}

	validForSeconds := int64(validFor.Truncate(time.Second).Seconds())

	treq := &authenticationv1.TokenRequest{
		Spec: authenticationv1.TokenRequestSpec{
			ExpirationSeconds: &validForSeconds,
		},

import (
	"context"
	"fmt"
	"strings"
	"time"
	return nil
}

func (k *kubernetesClient) listCustomResources(selectorGetter func(apiextensionsv1.CustomResourceDefinition) k8slabels.Selector) (out []unstructured.Unstructured, err error) {
	crds, err := k.extendedClient().ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{
		// CRDs might be provisioned by another application/charm from a different model.
	})
	if err != nil {
		return nil, errors.Trace(err)
	}
	for _, crd := range crds.Items {
		selector := selectorGetter(crd)
		if selector.Empty() {
			continue
		}
		for _, version := range crd.Spec.Versions {
			crdClient, err := k.getCustomResourceDefinitionClient(&crd, version.Name)
			if err != nil {
				return nil, errors.Trace(err)
			}
			list, err := crdClient.List(context.TODO(), metav1.ListOptions{
				LabelSelector: selector.String(),
			})
			if err != nil && !k8serrors.IsNotFound(err) {
				return nil, errors.Trace(err)
			}
			out = append(out, list.Items...)
		}
	}
	if len(out) == 0 {

// SetSnapConfig sets a snap's key to value.
func SetSnapConfig(snap string, key string, value string) error {
	if key == "" {
		return errors.NotValidf("key must not be empty")
	}

	cmd := exec.Command(Command, "set", snap, fmt.Sprintf("%s=%s", key, value))
	_, err := cmd.Output()
	if err != nil {
		return errors.Annotate(err, fmt.Sprintf("setting snap %s config %s to %s", snap, key, value))
	}

	return nil
}

// If no BackgroundServices are provided, Service will wrap all of the snap's
// background services.
func NewService(config ServiceConfig) (Service, error) {
	if config.ServiceName == "" {
		return Service{}, errors.New("ServiceName must be provided")
	}
	app := &App{
	}
	err := app.Validate()
	if err != nil {
		return Service{}, errors.Trace(err)
	}

	isLocal := config.SnapPath != ""

	return Service{
		runnable:       defaultRunner{},
// Validate validates that snap.Service has been correctly configured.
// Validate returns nil when successful and an error when successful.
func (s Service) Validate() error {
	if err := s.app.Validate(); err != nil {
		return errors.Trace(err)
	}

	for _, prerequisite := range s.app.Prerequisites() {
		if err := prerequisite.Validate(); err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}


// Running returns (true, nil) when snap indicates that service is currently active.
func (s Service) Running() (bool, error) {
	_, _, running, err := s.status()
	if err != nil {
		return false, errors.Trace(err)
	}
	return running, nil
}


// Install installs the snap and its background services.
func (s Service) Install() error {
	for _, app := range s.app.Prerequisites() {
		logger.Infof("command: %v", app)

		out, err := s.installAppWithRetry(app)
		if err != nil {
			return errors.Annotatef(err, "output: %v", out)
		}
	}

	out, err := s.installAppWithRetry(s.app)
	if err != nil {
		return errors.Annotatef(err, "output: %v", out)
	}
	return nil
}

func (s Service) installAppWithRetry(app Installable) (string, error) {
	ackAsserts := app.AcknowledgeAssertsArgs()
	if ackAsserts != nil {
		_, err := s.runCommandWithRetry(ackAsserts...)
		if err != nil {
			return "", errors.Trace(err)
		}
	}

	return s.runCommandWithRetry(app.InstallArgs()...)
}

// Installed returns true if the service has been successfully installed.
func (s Service) Installed() (bool, error) {
	installed, _, _, err := s.status()
	if err != nil {
		return false, errors.Trace(err)
	}
	return installed, nil
}

// ConfigOverride writes a systemd override to enable the
// specified limits to be used by the snap.
func (s Service) ConfigOverride() error {
	if len(s.conf.Limit) == 0 {
		return nil
	}

	unitOptions := systemd.ServiceLimits(s.conf)
	data, err := io.ReadAll(systemd.UnitSerialize(unitOptions))
	if err != nil {
		return errors.Trace(err)
	}

	for _, backgroundService := range s.app.BackgroundServices() {
		overridesDir := fmt.Sprintf("%s/snap.%s.%s.service.d", s.configDir, s.name, backgroundService.Name)
		if err := os.MkdirAll(overridesDir, 0755); err != nil {
			return errors.Trace(err)
		}
		if err := os.WriteFile(filepath.Join(overridesDir, "overrides.conf"), data, 0644); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

// shell commands to be executed by a shell which start the service.
func (s Service) StartCommands() ([]string, error) {
	deps := s.app.Prerequisites()
	commands := make([]string, 0, 1+len(deps))
	for _, prerequisite := range deps {
		commands = append(commands, prerequisite.StartCommands(s.executable)...)
	}
	return append(commands, s.app.StartCommands(s.executable)...), nil
}

// status returns an interpreted output from the `snap services` command.
//
//	(true, true, false, nil)
func (s *Service) status() (isInstalled, enabledAtStartup, isCurrentlyActive bool, err error) {
	out, err := s.runCommand("services", s.Name())
	if err != nil {
		return false, false, false, errors.Trace(err)
	}
	for _, line := range strings.Split(out, "\n") {
		if !strings.HasPrefix(line, s.Name()) {
			continue
		}

		fields := strings.Fields(line)
		return true, fields[1] == "enabled", fields[2] == "active", nil
	}

	return false, false, false, nil
}

// Start starts the service, returning nil when successful.
// If the service is already running, Start does not restart it.
func (s Service) Start() error {
	running, err := s.Running()
	if err != nil {
		return errors.Trace(err)
	}
	if running {
		return nil
	}

	commands, err := s.StartCommands()
	if err != nil {
		return errors.Trace(err)
	}
	for _, command := range commands {
		commandParts := strings.Fields(command)
		out, err := utils.RunCommand(commandParts[0], commandParts[1:]...)
		if err != nil {
			if strings.Contains(out, "has no services") {
				continue
			}
			return errors.Annotatef(err, "%v -> %v", command, out)
		}
	}

	return nil
}

// Stop stops a running service. Returns nil when the underlying
// call to `snap stop <service-name>` exits with error code 0.
func (s Service) Stop() error {
	running, err := s.Running()
	if err != nil {
		return errors.Trace(err)
	}
	if !running {
		return nil
	}

	args := []string{"stop", s.Name()}
	return s.execThenExpect(args, "Stopped.")
}

// Restart restarts the service, or starts if it's not currently
//
// Restart is part of the service.RestartableService interface
func (s Service) Restart() error {
	args := []string{"restart", s.Name()}
	return s.execThenExpect(args, "Restarted.")
}

// execThenExpect calls `snap <commandArgs>...` and then checks
// stdout against expectation and snap's exit code. When there's a
// mismatch or non-0 exit code, execThenExpect returns an error.
func (s Service) execThenExpect(commandArgs []string, expectation string) error {
	out, err := s.runCommand(commandArgs...)
	if err != nil {
		return errors.Trace(err)
	}
	if !strings.Contains(out, expectation) {
		return errors.Annotatef(err, `expected "%s", got "%s"`, expectation, out)
	}
	return nil
}

}

func (s Service) runCommandWithRetry(args ...string) (res string, err error) {
	if resErr := retry.Call(retry.CallArgs{
		Clock: s.clock,
		Func: func() error {
			res, err = s.runCommand(args...)
			return errors.Trace(err)
		},
		Delay:    5 * time.Second,
		Attempts: 2,
	}); resErr != nil {
		return "", errors.Trace(resErr)
	}

	// Named args are set via the retry.
	return
}
mongo/mongo.go | 8 ++++----
mongo/open.go  | 4 ++--
2 files changed, 6 insertions(+), 6 deletions(-)
