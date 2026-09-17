package main

			}
			sources = appSpec.GetSources()
		} else {
			source := a.Spec.GetSource()
			if q.GetRevision() != "" {
				source.TargetRevision = q.GetRevision()
			}
		if policyFinalizer == "" {
			return nil, status.Errorf(codes.InvalidArgument, "invalid propagation policy: %s", *q.PropagationPolicy)
		}
		if !a.IsFinalizerPresent(policyFinalizer) {
			a.SetCascadedDeletion(policyFinalizer)
			patchFinalizer = true
		}
}

func (s *Server) RevisionMetadata(ctx context.Context, q *application.RevisionMetadataQuery) (*v1alpha1.RevisionMetadata, error) {
	a, proj, err := s.getApplicationEnforceRBACInformer(ctx, rbac.ActionGet, q.GetProject(), q.GetAppNamespace(), q.GetName())
	if err != nil {
		return nil, err
	}
		}
		log.Warnf("failed to set operation for app %q due to update conflict. retrying again...", *termOpReq.Name)
		time.Sleep(100 * time.Millisecond)
		_, err = s.appclientset.ArgoprojV1alpha1().Applications(appNs).Get(ctx, appName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("error getting application by name: %w", err)
		}
			if err != nil {
				return nil, fmt.Errorf("error unmarshaling live state for %s/%s: %w", liveResource.Kind, liveResource.Name, err)
			}
			liveObjs = append(liveObjs, liveObj)
		} else {
			liveObjs = append(liveObjs, nil)
	if err != nil {
		return nil, fmt.Errorf("error performing state diffs: %w", err)
	}

	// Convert StateDiffs results to ResourceDiff format for API response
	responseDiffs := make([]*v1alpha1.ResourceDiff, 0, len(diffResults.Diffs))
				i, len(q.GetLiveResources()), len(targetObjs))
		}

		// Create ResourceDiff with StateDiffs results
		// TargetState = PredictedLive (what the target should be after applying)
		// LiveState = NormalizedLive (current normalized live state)
		responseDiffs = append(responseDiffs, &v1alpha1.ResourceDiff{
			Group:           group,
			Kind:            kind,
			Namespace:       namespace,
			Name:            name,
			TargetState:     string(diffRes.PredictedLive),
			LiveState:       string(diffRes.NormalizedLive),
			Diff:            "", // Diff string is generated client-side
			Hook:            hook,
			Modified:        diffRes.Modified,

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"os"
	// resources which in this case applies the live values in the configured
	// ignore differences fields.
	if syncOp.SyncOptions.HasOption("RespectIgnoreDifferences=true") {
		patchedTargets, err := normalizeTargetResources(openAPISchema, compareResult)
		if err != nil {
			state.Phase = common.OperationError
			state.Message = fmt.Sprintf("Failed to normalize target resources: %s", err)
//   - applies normalization to the target resources based on the live resources
//   - copies ignored fields from the matching live resources: apply normalizer to the live resource,
//     calculates the patch performed by normalizer and applies the patch to the target resource
func normalizeTargetResources(openAPISchema openapi.Resources, cr *comparisonResult) ([]*unstructured.Unstructured, error) {
	// Normalize live and target resources (cleaning or aligning them)
	normalized, err := diff.Normalize(cr.reconciliationResult.Live, cr.reconciliationResult.Target, cr.diffConfig)
	if err != nil {
		return nil, err
	}

	patchedTargets := []*unstructured.Unstructured{}

	for idx, live := range cr.reconciliationResult.Live {
		normalizedTarget := normalized.Targets[idx]
		if normalizedTarget == nil {
			patchedTargets = append(patchedTargets, nil)
			continue
		}
		gvk := normalizedTarget.GroupVersionKind()

		originalTarget := cr.reconciliationResult.Target[idx]
		if live == nil {
			// No live resource, just use target
			patchedTargets = append(patchedTargets, originalTarget)
			continue
		}

		var (
			lookupPatchMeta strategicpatch.LookupPatchMeta
			versionedObject any
		)

		// Load patch meta struct or OpenAPI schema for CRDs
		if versionedObject, err = scheme.Scheme.New(gvk); err == nil {
			if lookupPatchMeta, err = strategicpatch.NewPatchMetaFromStruct(versionedObject); err != nil {
				return nil, err
			}
		} else if crdSchema := openAPISchema.LookupResource(gvk); crdSchema != nil {
			lookupPatchMeta = strategicpatch.NewPatchMetaFromOpenAPI(crdSchema)
		}

		// Calculate live patch
		livePatch, err := getMergePatch(normalized.Lives[idx], live, lookupPatchMeta)
		if err != nil {
			return nil, err
		}

		// Apply the patch to the normalized target
		// This ensures ignored fields in live are restored into the target before syncing
		normalizedTarget, err = applyMergePatch(normalizedTarget, livePatch, versionedObject, lookupPatchMeta)
		if err != nil {
			return nil, err
		}
		patchedTargets = append(patchedTargets, normalizedTarget)
	}

	return patchedTargets, nil
}

// getMergePatch calculates and returns the patch between the original and the
// modified unstructures.
func getMergePatch(original, modified *unstructured.Unstructured, lookupPatchMeta strategicpatch.LookupPatchMeta) ([]byte, error) {
	originalJSON, err := original.MarshalJSON()
	if err != nil {
		return nil, err
	return jsonpatch.CreateMergePatch(originalJSON, modifiedJSON)
}

// applyMergePatch will apply the given patch in the obj and return the patched unstructure.
func applyMergePatch(obj *unstructured.Unstructured, patch []byte, versionedObject any, meta strategicpatch.LookupPatchMeta) (*unstructured.Unstructured, error) {
	originalJSON, err := obj.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var patchedJSON []byte
	switch {
	case versionedObject != nil:
		patchedJSON, err = strategicpatch.StrategicMergePatch(originalJSON, patch, versionedObject)
	case meta != nil:
		var originalMap, patchMap map[string]any
		if err := json.Unmarshal(originalJSON, &originalMap); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(patch, &patchMap); err != nil {
			return nil, err
		}

		patchedMap, err := strategicpatch.StrategicMergeMapPatchUsingLookupPatchMeta(originalMap, patchMap, meta)
		if err != nil {
			return nil, err
		}
		patchedJSON, err = json.Marshal(patchedMap)
		if err != nil {
			return nil, err
		}
	default:
		patchedJSON, err = jsonpatch.MergePatch(originalJSON, patch)
	}
	if err != nil {
		return nil, err
	appIf := a.appsLister.Applications(nsFilter)
	apps, err := appIf.List(labels.Everything())
	if err != nil {
		log.Warnf("Failed to list applications: %v", err)
		return
	}

	installationID, err := a.settingsSrc.GetInstallationID()
	if err != nil {
		log.Warnf("Failed to get installation ID: %v", err)
		return
	}
	trackingMethod, err := a.settingsSrc.GetTrackingMethod()
	if err != nil {
		log.Warnf("Failed to get trackingMethod: %v", err)
		return
	}
	appInstanceLabelKey, err := a.settingsSrc.GetAppInstanceLabelKey()
	if err != nil {
		log.Warnf("Failed to get appInstanceLabelKey: %v", err)
		return
	}

	for _, webURL := range webURLs {
		repoRegexp, err := GetWebURLRegex(webURL)
		if err != nil {
			log.Warnf("Failed to get repoRegexp: %s", err)
			continue
		}
		for _, app := range filteredApps {
			if app.Spec.SourceHydrator != nil {
				drySource := app.Spec.SourceHydrator.GetDrySource()
				if sourceRevisionHasChanged(drySource, revision, touchedHead) && sourceUsesURL(drySource, webURL, repoRegexp) {
					refreshPaths := path.GetAppRefreshPaths(&app)
					if path.AppFilesHaveChanged(refreshPaths, changedFiles) {
						namespacedAppInterface := a.appClientset.ArgoprojV1alpha1().Applications(app.Namespace)
						log.Infof("webhook trigger refresh app to hydrate '%s'", app.Name)
						_, err = argo.RefreshApp(namespacedAppInterface, app.Name, v1alpha1.RefreshTypeNormal, true)
						if err != nil {
							log.Warnf("Failed to hydrate app '%s' for controller reprocessing: %v", app.Name, err)
							continue
						}
					}
				}
			}

			for _, source := range app.Spec.GetSources() {
				if sourceRevisionHasChanged(source, revision, touchedHead) && sourceUsesURL(source, webURL, repoRegexp) {
					refreshPaths := path.GetAppRefreshPaths(&app)
					if path.AppFilesHaveChanged(refreshPaths, changedFiles) {
						namespacedAppInterface := a.appClientset.ArgoprojV1alpha1().Applications(app.Namespace)
						_, err = argo.RefreshApp(namespacedAppInterface, app.Name, v1alpha1.RefreshTypeNormal, true)
						if err != nil {
							log.Warnf("Failed to refresh app '%s' for controller reprocessing: %v", app.Name, err)
							continue
						}
						// No need to refresh multiple times if multiple sources match.
						break
					} else if change.shaBefore != "" && change.shaAfter != "" {
						if err := a.storePreviouslyCachedManifests(&app, change, trackingMethod, appInstanceLabelKey, installationID); err != nil {
							log.Warnf("Failed to store cached manifests of previous revision for app '%s': %v", app.Name, err)
						}
					}
				}
	return repoRegexp, nil
}

func (a *ArgoCDWebhookHandler) storePreviouslyCachedManifests(app *v1alpha1.Application, change changeInfo, trackingMethod string, appInstanceLabelKey string, installationID string) error {
	destCluster, err := argo.GetDestinationCluster(context.Background(), app.Spec.Destination, a.db)
	if err != nil {
		return fmt.Errorf("error validating destination: %w", err)
	if err != nil {
		return fmt.Errorf("error getting ref sources: %w", err)
	}
	source := app.Spec.GetSource()
	cache.LogDebugManifestCacheKeyFields("moving manifests cache", "webhook app revision changed", change.shaBefore, &source, refSources, &clusterInfo, app.Spec.Destination.Namespace, trackingMethod, appInstanceLabelKey, app.Name, nil)

	if err := a.repoCache.SetNewRevisionManifests(change.shaAfter, change.shaBefore, &source, refSources, &clusterInfo, app.Spec.Destination.Namespace, trackingMethod, appInstanceLabelKey, app.Name, nil, installationID); err != nil {
