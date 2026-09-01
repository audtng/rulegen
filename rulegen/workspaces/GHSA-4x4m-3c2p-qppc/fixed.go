package main

			return admission.NewForbidden(a, fmt.Errorf("node %q is not allowed to modify taints", nodeName))
		}

		// Don't allow a node to update its own ownerReferences.
		if !apiequality.Semantic.DeepEqual(node.OwnerReferences, oldNode.OwnerReferences) {
			return admission.NewForbidden(a, fmt.Errorf("node %q is not allowed to modify ownerReferences", nodeName))
		}

		// Don't allow a node to update labels outside the allowed set.
		// This would allow a node to add or modify its labels in a way that would let it steer privileged workloads to itself.
		modifiedLabels := getModifiedLabels(node.Labels, oldNode.Labels)
