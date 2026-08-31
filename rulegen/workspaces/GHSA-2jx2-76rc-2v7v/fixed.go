package main

			return nil, errors.NewInternalError(fmt.Errorf("error converting object from store to a pod security policy: %v", c))
		}

		// if no user info exists then the API is being hit via the unsecured port. In this case authorize the request.
		if user == nil || authorizedForPolicy(user, constraint, authz) || authorizedForPolicy(sa, constraint, authz) {
			matchedPolicies = append(matchedPolicies, constraint)
		}
	}

// authorizedForPolicy returns true if info is authorized to perform a "get" on policy.
func authorizedForPolicy(info user.Info, policy *extensions.PodSecurityPolicy, authz authorizer.Authorizer) bool {
	if info == nil {
		return false
	}
	attr := buildAttributes(info, policy)
	allowed, _, _ := authz.Authorize(attr)
			// (ie. a request hitting the unsecure port)
			expectedPolicies: sets.NewString("policy1", "policy2", "policy3"),
		},
		"policies are not allowed for nil sa info": {
			user: &user.DefaultInfo{Name: "user"},
			sa:   nil,
			disallowedPolicies: map[string][]string{
				policyWithName("policy2"),
				policyWithName("policy3"),
			},
			// only the policies for the user are allowed when sa info is nil
			expectedPolicies: sets.NewString("policy2"),
		},
	}
	for k, v := range tests {
