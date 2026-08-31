package main

			return v1.AggregateValidationErrors("Function", err)
		}
	}

	if err := new.Validate(); err != nil {
		return v1.AggregateValidationErrors("Function", err)

// RefreshFuncPods deletes pods related to the function so that new pods are replenished
func (deploy *NewDeploy) RefreshFuncPods(ctx context.Context, logger logr.Logger, f fv1.Function) error {

	env, err := deploy.fissionClient.CoreV1().Environments(f.Spec.Environment.Namespace).Get(ctx, f.Spec.Environment.Name, metav1.GetOptions{})
	if err != nil {
}

func (deploy *NewDeploy) fnCreate(ctx context.Context, fn *fv1.Function) (*fscache.FuncSvc, error) {
	cleanupFunc := func(ctx context.Context, ns string, name string) {
		err := deploy.cleanupNewdeploy(ctx, ns, name)
		if err != nil {
	var env *fv1.Environment
	otelUtils.SpanTrackEvent(ctx, "getFunctionEnv", otelUtils.GetAttributesForFunction(fn)...)

	// Cached ?
	// TODO: the cache should be able to search by <env name, fn namespace> instead of function metadata.
	result, err := gpm.functionEnv.Get(crd.CacheKeyURFromMeta(&fn.ObjectMeta))
