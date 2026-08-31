package main

			return v1.AggregateValidationErrors("Function", err)
		}
	}
	// Cross-namespace EnvironmentRef closes GHSA-cvw6-gfvv-953q. An empty
	// namespace remains accepted — the Fission CLI populates it with the
	// function's own namespace at creation time (pkg/fission-cli/cmd/function/
	// create.go), and downstream controllers tolerate empty via
	// DefaultNSResolver. Rejecting only the explicit cross-namespace value is
	// sufficient for this advisory; defaulting an empty namespace at admission
	// is a separate hardening track tracked outside this fix.
	if envRef := new.Spec.Environment; envRef.Namespace != "" && envRef.Namespace != new.Namespace {
		err := fmt.Errorf("environment's namespace [%s] and function's namespace [%s] are different; cross-namespace Environment reference is not allowed",
			envRef.Namespace, new.Namespace)
		return v1.AggregateValidationErrors("Function", err)
	}
	// Cross-namespace PackageRef closes GHSA-3r8v-2xmj-5c39. Same shape as
	// the EnvironmentRef check above, including the empty-is-accepted rule.
	if pkgRef := new.Spec.Package.PackageRef; pkgRef.Namespace != "" && pkgRef.Namespace != new.Namespace {
		err := fmt.Errorf("package's namespace [%s] and function's namespace [%s] are different; cross-namespace Package reference is not allowed",
			pkgRef.Namespace, new.Namespace)
		return v1.AggregateValidationErrors("Function", err)
	}

	if err := new.Validate(); err != nil {
		return v1.AggregateValidationErrors("Function", err)

// RefreshFuncPods deletes pods related to the function so that new pods are replenished
func (deploy *NewDeploy) RefreshFuncPods(ctx context.Context, logger logr.Logger, f fv1.Function) error {
	// Defence in depth for GHSA-cvw6-gfvv-953q — see fnCreate for context.
	if envNs := f.Spec.Environment.Namespace; envNs != "" && envNs != f.Namespace {
		return fmt.Errorf("cross-namespace environment reference is not allowed: fn.namespace=%s env.namespace=%s",
			f.Namespace, envNs)
	}

	env, err := deploy.fissionClient.CoreV1().Environments(f.Spec.Environment.Namespace).Get(ctx, f.Spec.Environment.Name, metav1.GetOptions{})
	if err != nil {
}

func (deploy *NewDeploy) fnCreate(ctx context.Context, fn *fv1.Function) (*fscache.FuncSvc, error) {
	// Defence in depth for GHSA-cvw6-gfvv-953q — primary defence is the
	// admission webhook in pkg/webhook/function.go, but a stale Function
	// from a pre-webhook upgrade window (or failurePolicy=ignore) could
	// still reach this path.
	if envNs := fn.Spec.Environment.Namespace; envNs != "" && envNs != fn.Namespace {
		return nil, fmt.Errorf("cross-namespace environment reference is not allowed: fn.namespace=%s env.namespace=%s",
			fn.Namespace, envNs)
	}
	cleanupFunc := func(ctx context.Context, ns string, name string) {
		err := deploy.cleanupNewdeploy(ctx, ns, name)
		if err != nil {
	var env *fv1.Environment
	otelUtils.SpanTrackEvent(ctx, "getFunctionEnv", otelUtils.GetAttributesForFunction(fn)...)

	// Defence in depth for GHSA-cvw6-gfvv-953q — the admission webhook
	// already rejects this at submit time, but a stale Function object
	// from an upgrade-before-webhook-restart window (or a cluster running
	// with failurePolicy=ignore) could still reach this path.
	if envNs := fn.Spec.Environment.Namespace; envNs != "" && envNs != fn.Namespace {
		return nil, fmt.Errorf("cross-namespace environment reference is not allowed: fn.namespace=%s env.namespace=%s",
			fn.Namespace, envNs)
	}

	// Cached ?
	// TODO: the cache should be able to search by <env name, fn namespace> instead of function metadata.
	result, err := gpm.functionEnv.Get(crd.CacheKeyURFromMeta(&fn.ObjectMeta))
