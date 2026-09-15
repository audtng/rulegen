package main

			return nil, errors.NewInternalError(fmt.Errorf("error converting object from store to a pod security policy: %v", c))
		}

		if authorizedForPolicy(user, constraint, authz) || authorizedForPolicy(sa, constraint, authz) {
			matchedPolicies = append(matchedPolicies, constraint)
		}
	}

// authorizedForPolicy returns true if info is authorized to perform a "get" on policy.
func authorizedForPolicy(info user.Info, policy *extensions.PodSecurityPolicy, authz authorizer.Authorizer) bool {
	// if no info exists then the API is being hit via the unsecured port.  In this case
	// authorize the request.
	if info == nil {
		return true
	}
	attr := buildAttributes(info, policy)
	allowed, _, _ := authz.Authorize(attr)
