package main

			}
			sources = appSpec.GetSources()
		} else {
			// For sourceHydrator applications, use the dry source to generate manifests
			var source v1alpha1.ApplicationSource
			if a.Spec.SourceHydrator != nil {
				source = a.Spec.SourceHydrator.GetDrySource()
			} else {
				source = a.Spec.GetSource()
			}

			if q.GetRevision() != "" {
				source.TargetRevision = q.GetRevision()
			}
		if policyFinalizer == "" {
			return nil, status.Errorf(codes.InvalidArgument, "invalid propagation policy: %s", *q.PropagationPolicy)
		}
		// Kubernetes forbids adding finalizers to an object that is already being deleted,
		// so skip the patch if the app is mid-deletion. The underlying Delete call below is
		// still issued to keep the RPC idempotent for callers retrying a cascade delete.
		if !a.IsFinalizerPresent(policyFinalizer) && a.DeletionTimestamp == nil {
			a.SetCascadedDeletion(policyFinalizer)
			patchFinalizer = true
		}
}

func (s *Server) RevisionMetadata(ctx context.Context, q *application.RevisionMetadataQuery) (*v1alpha1.RevisionMetadata, error) {
	// Read via the client instead of the informer cache to avoid "revision history not found" errors due to stale informer cache
	a, proj, err := s.getApplicationEnforceRBACClient(ctx, rbac.ActionGet, q.GetProject(), q.GetAppNamespace(), q.GetName(), "")
	if err != nil {
		return nil, err
	}
		}
		log.Warnf("failed to set operation for app %q due to update conflict. retrying again...", *termOpReq.Name)
		time.Sleep(100 * time.Millisecond)
		a, err = s.appclientset.ArgoprojV1alpha1().Applications(appNs).Get(ctx, appName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("error getting application by name: %w", err)
		}
			if err != nil {
				return nil, fmt.Errorf("error unmarshaling live state for %s/%s: %w", liveResource.Kind, liveResource.Name, err)
			}
			if liveObj.GetName() != liveResource.Name {
				return nil, fmt.Errorf("name mismatch: expected %s, got %s", liveResource.Name, liveObj.GetName())
			}
			if liveObj.GetNamespace() != liveResource.Namespace {
				return nil, fmt.Errorf("namespace mismatch: expected %s, got %s", liveResource.Namespace, liveObj.GetNamespace())
			}
			if liveObj.GroupVersionKind().Group != liveResource.Group {
				return nil, fmt.Errorf("group mismatch: expected %s, got %s", liveResource.Group, liveObj.GroupVersionKind().Group)
			}
			if liveObj.GroupVersionKind().Kind != liveResource.Kind {
				return nil, fmt.Errorf("kind mismatch: expected %s, got %s", liveResource.Kind, liveObj.GroupVersionKind().Kind)
			}
			liveObjs = append(liveObjs, liveObj)
		} else {
			liveObjs = append(liveObjs, nil)
	if err != nil {
		return nil, fmt.Errorf("error performing state diffs: %w", err)
	}
	managedResources := make([]*v1alpha1.ResourceDiff, 0)
	err = s.getCachedAppState(ctx, a, func() error {
		return s.cache.GetAppManagedResources(a.InstanceName(s.ns), &managedResources)
	})
	if err != nil {
		return nil, fmt.Errorf("error getting managed resources: %w", err)
	}

	// Convert StateDiffs results to ResourceDiff format for API response
	responseDiffs := make([]*v1alpha1.ResourceDiff, 0, len(diffResults.Diffs))
				i, len(q.GetLiveResources()), len(targetObjs))
		}

		found := false
		for _, item := range managedResources {
			if item.Kind == kind && item.Group == group && item.Namespace == namespace && item.Name == name {
				found = true
				break
			}
		}
		if !found {
			return nil, status.Errorf(codes.PermissionDenied, "%s %s %s not found as part of application %s", kind, group, name, a.Name)
		}

		// Create ResourceDiff with StateDiffs results
		// TargetState = PredictedLive (what the target should be after applying)
		// LiveState = NormalizedLive (current normalized live state)
		targetState := string(diffRes.PredictedLive)
		liveState := string(diffRes.NormalizedLive)

		if kind == kube.SecretKind && group == "" {
			var targetObj, liveObj *unstructured.Unstructured
			if len(diffRes.PredictedLive) > 0 && string(diffRes.PredictedLive) != "null" {
				targetObj = &unstructured.Unstructured{}
				if err := json.Unmarshal(diffRes.PredictedLive, targetObj); err != nil {
					return nil, fmt.Errorf("error unmarshaling predicted live for secret masking: %w", err)
				}
			}
			if len(diffRes.NormalizedLive) > 0 && string(diffRes.NormalizedLive) != "null" {
				liveObj = &unstructured.Unstructured{}
				if err := json.Unmarshal(diffRes.NormalizedLive, liveObj); err != nil {
					return nil, fmt.Errorf("error unmarshaling normalized live for secret masking: %w", err)
				}
			}
			maskedTarget, maskedLive, err := diff.HideSecretData(targetObj, liveObj, s.settingsMgr.GetSensitiveAnnotations())
			if err != nil {
				return nil, fmt.Errorf("error masking secret data: %w", err)
			}
			if maskedTarget != nil {
				data, err := json.Marshal(maskedTarget)
				if err != nil {
					return nil, fmt.Errorf("error marshaling masked target state: %w", err)
				}
				targetState = string(data)
			}
			if maskedLive != nil {
				data, err := json.Marshal(maskedLive)
				if err != nil {
					return nil, fmt.Errorf("error marshaling masked live state: %w", err)
				}
				liveState = string(data)
			}
		}

		responseDiffs = append(responseDiffs, &v1alpha1.ResourceDiff{
			Group:           group,
			Kind:            kind,
			Namespace:       namespace,
			Name:            name,
			TargetState:     targetState,
			LiveState:       liveState,
			Diff:            "", // Diff string is generated client-side
			Hook:            hook,
			Modified:        diffRes.Modified,

import (
	"context"
	stderrors "errors"
	"fmt"
	"os"
	// resources which in this case applies the live values in the configured
	// ignore differences fields.
	if syncOp.SyncOptions.HasOption("RespectIgnoreDifferences=true") {
		patchedTargets, err := normalizeTargetResources(compareResult)
		if err != nil {
			state.Phase = common.OperationError
			state.Message = fmt.Sprintf("Failed to normalize target resources: %s", err)
//   - applies normalization to the target resources based on the live resources
//   - copies ignored fields from the matching live resources: apply normalizer to the live resource,
//     calculates the patch performed by normalizer and applies the patch to the target resource
func normalizeTargetResources(cr *comparisonResult) ([]*unstructured.Unstructured, error) {
	// normalize live and target resources
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
		originalTarget := cr.reconciliationResult.Target[idx]
		if live == nil {
			patchedTargets = append(patchedTargets, originalTarget)
			continue
		}

		var lookupPatchMeta *strategicpatch.PatchMetaFromStruct
		versionedObject, err := scheme.Scheme.New(normalizedTarget.GroupVersionKind())
		if err == nil {
			meta, err := strategicpatch.NewPatchMetaFromStruct(versionedObject)
			if err != nil {
				return nil, err
			}
			lookupPatchMeta = &meta
		}

		livePatch, err := getMergePatch(normalized.Lives[idx], live, lookupPatchMeta)
		if err != nil {
			return nil, err
		}

		normalizedTarget, err = applyMergePatch(normalizedTarget, livePatch, versionedObject)
		if err != nil {
			return nil, err
		}

		patchedTargets = append(patchedTargets, normalizedTarget)
	}
	return patchedTargets, nil
}

// getMergePatch calculates and returns the patch between the original and the
// modified unstructures.
func getMergePatch(original, modified *unstructured.Unstructured, lookupPatchMeta *strategicpatch.PatchMetaFromStruct) ([]byte, error) {
	originalJSON, err := original.MarshalJSON()
	if err != nil {
		return nil, err
	return jsonpatch.CreateMergePatch(originalJSON, modifiedJSON)
}

// applyMergePatch will apply the given patch in the obj and return the patched
// unstructure.
func applyMergePatch(obj *unstructured.Unstructured, patch []byte, versionedObject any) (*unstructured.Unstructured, error) {
	originalJSON, err := obj.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var patchedJSON []byte
	if versionedObject == nil {
		patchedJSON, err = jsonpatch.MergePatch(originalJSON, patch)
	} else {
		patchedJSON, err = strategicpatch.StrategicMergePatch(originalJSON, patch, versionedObject)
	}
	if err != nil {
		return nil, err
	appIf := a.appsLister.Applications(nsFilter)
	apps, err := appIf.List(labels.Everything())
	if err != nil {
		log.Errorf("Failed to list applications: %v", err)
		return
	}

	installationID, err := a.settingsSrc.GetInstallationID()
	if err != nil {
		log.Errorf("Failed to get installation ID: %v", err)
		return
	}
	trackingMethod, err := a.settingsSrc.GetTrackingMethod()
	if err != nil {
		log.Errorf("Failed to get trackingMethod: %v", err)
		return
	}
	appInstanceLabelKey, err := a.settingsSrc.GetAppInstanceLabelKey()
	if err != nil {
		log.Errorf("Failed to get appInstanceLabelKey: %v", err)
		return
	}

	for _, webURL := range webURLs {
		repoRegexp, err := GetWebURLRegex(webURL)
		if err != nil {
			log.Errorf("Failed to get repoRegexp: %s", err)
			continue
		}

		// iterate over apps and check if any files specified in their sources have changed
		for _, app := range filteredApps {
			// get all sources, including sync source and dry source if source hydrator is configured
			sources := app.Spec.GetSources()
			if app.Spec.SourceHydrator != nil {
				// we already have sync source, so add dry source if source hydrator is configured
				sources = append(sources, app.Spec.SourceHydrator.GetDrySource())
			}

			// iterate over all sources and check if any files specified in refresh paths have changed
			for _, source := range sources {
				if sourceRevisionHasChanged(source, revision, touchedHead) && sourceUsesURL(source, webURL, repoRegexp) {
					refreshPaths := path.GetSourceRefreshPaths(&app, source)
					if path.AppFilesHaveChanged(refreshPaths, changedFiles) {
						hydrate := false
						if app.Spec.SourceHydrator != nil {
							drySource := app.Spec.SourceHydrator.GetDrySource()
							if (&source).Equals(&drySource) {
								hydrate = true
							}
						}

						// refresh paths have changed, so we need to refresh the app
						log.Infof("refreshing app '%s' from webhook", app.Name)
						if hydrate {
							// log if we need to hydrate the app
							log.Infof("webhook trigger refresh app to hydrate '%s'", app.Name)
						}
						namespacedAppInterface := a.appClientset.ArgoprojV1alpha1().Applications(app.Namespace)
						if _, err := argo.RefreshApp(namespacedAppInterface, app.Name, v1alpha1.RefreshTypeNormal, hydrate); err != nil {
							log.Errorf("Failed to refresh app '%s': %v", app.Name, err)
						}
						break // we don't need to check other sources
					} else if change.shaBefore != "" && change.shaAfter != "" {
						// update the cached manifests with the new revision cache key
						if err := a.storePreviouslyCachedManifests(&app, change, trackingMethod, appInstanceLabelKey, installationID, source); err != nil {
							log.Errorf("Failed to store cached manifests of previous revision for app '%s': %v", app.Name, err)
						}
					}
				}
	return repoRegexp, nil
}

func (a *ArgoCDWebhookHandler) storePreviouslyCachedManifests(app *v1alpha1.Application, change changeInfo, trackingMethod string, appInstanceLabelKey string, installationID string, source v1alpha1.ApplicationSource) error {
	destCluster, err := argo.GetDestinationCluster(context.Background(), app.Spec.Destination, a.db)
	if err != nil {
		return fmt.Errorf("error validating destination: %w", err)
	if err != nil {
		return fmt.Errorf("error getting ref sources: %w", err)
	}

	cache.LogDebugManifestCacheKeyFields("moving manifests cache", "webhook app revision changed", change.shaBefore, &source, refSources, &clusterInfo, app.Spec.Destination.Namespace, trackingMethod, appInstanceLabelKey, app.Name, nil)

	if err := a.repoCache.SetNewRevisionManifests(change.shaAfter, change.shaBefore, &source, refSources, &clusterInfo, app.Spec.Destination.Namespace, trackingMethod, appInstanceLabelKey, app.Name, nil, installationID); err != nil {
