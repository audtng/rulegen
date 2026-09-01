package main

			violations = append(violations, fieldName)
		}
	}
	// ArtifactGC is allow-listed so that its benign fields (Strategy,
	// ForceFinalizerRemoval) may be set, but its nested ServiceAccountName,
	// PodSpecPatch and PodMetadata reach the artifact-GC Pod and would otherwise
	// re-open the privilege escalation that the top-level ServiceAccountName /
	// PodSpecPatch / PodMetadata blocks are meant to close, so reject them here.
	violations = append(violations, artifactGCOverrideViolations(userSpec.ArtifactGC)...)
	if len(violations) > 0 {
		sort.Strings(violations)
		return fmt.Errorf("fields %v are not permitted when using workflowTemplateRef with templateReferencing restriction", violations)
			dst.Field(i).Set(src.Field(i))
		}
	}
	// ArtifactGC is allow-listed wholesale above, which copies the pointer and
	// so would carry through the user's ServiceAccountName/PodSpecPatch/PodMetadata.
	// Strip those security-sensitive fields (on a copy, so the caller's spec is
	// left untouched) while preserving the benign ones.
	sanitized.ArtifactGC = sanitizeArtifactGC(sanitized.ArtifactGC)
	return sanitized
}

// artifactGCOverrideViolations reports the security-sensitive fields set within a
// user-supplied workflow-level ArtifactGC. These reach the artifact-GC Pod and
// must not be user-controllable when a hardened WorkflowTemplate is referenced
// under Strict/Secure mode.
func artifactGCOverrideViolations(agc *wfv1.WorkflowLevelArtifactGC) []string {
	if agc == nil {
		return nil
	}
	var violations []string
	if agc.PodSpecPatch != "" {
		violations = append(violations, "ArtifactGC.PodSpecPatch")
	}
	if agc.ServiceAccountName != "" {
		violations = append(violations, "ArtifactGC.ServiceAccountName")
	}
	if agc.PodMetadata != nil {
		violations = append(violations, "ArtifactGC.PodMetadata")
	}
	return violations
}

// sanitizeArtifactGC returns a deep copy of the workflow-level ArtifactGC with the
// security-sensitive override fields removed, leaving the benign fields (Strategy,
// ForceFinalizerRemoval) intact. A copy is returned so the caller's original spec
// is never mutated.
func sanitizeArtifactGC(agc *wfv1.WorkflowLevelArtifactGC) *wfv1.WorkflowLevelArtifactGC {
	if agc == nil {
		return nil
	}
	clean := agc.DeepCopy()
	clean.PodSpecPatch = ""
	clean.ServiceAccountName = ""
	clean.PodMetadata = nil
	return clean
}

// MergeTo will merge one workflow (the "patch" workflow) into another (the "target" workflow.
// If the target workflow defines a field, this take precedence over the patch.
func MergeTo(patch, target *wfv1.Workflow) error {
