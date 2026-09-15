package main

	"k8s.io/apimachinery/pkg/util/validation"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/envoyproxy/gateway/internal/gatewayapi/resource"
				unsupportedFilters = true
				continue
			}
			if filter.Type != gwapiv1.HTTPRouteFilterRequestHeaderModifier &&
				filter.Type != gwapiv1.HTTPRouteFilterResponseHeaderModifier &&
				filter.Type != gwapiv1.HTTPRouteFilterExtensionRef &&
				unsupportedFilters = true
			}
		}
	case resource.KindGRPCRoute:
		for _, filter := range filters.([]gwapiv1.GRPCRouteFilter) {
			if filter.Type != gwapiv1.GRPCRouteFilterRequestHeaderModifier &&
				unsupportedFilters = true
			}
		}
	default:
		return nil
	}
	resources *resource.Resources, routeKind gwapiv1.Kind,
) status.Error {
	if backendRef.Namespace != nil && string(*backendRef.Namespace) != "" && string(*backendRef.Namespace) != route.GetNamespace() {
		if !t.validateCrossNamespaceRef(
			crossNamespaceFrom{
				group:     gwapiv1.GroupName,
				kind:      string(routeKind),
		return
	}

	// Edge case: only one condition which is ResolvedRefs=False, Reason=PartiallyInvalidCertificateRef
	// In this case, we can still consider the listener as ready because we only program the listener using only the valid certificates.
	if len(lConditions) == 1 && lConditions[0].Type == string(gwapiv1.ListenerConditionResolvedRefs) &&
		lConditions[0].Reason == string(status.ListenerReasonPartiallyInvalidCertificateRef) {
		listener.SetCondition(gwapiv1.ListenerConditionAccepted, metav1.ConditionTrue, gwapiv1.ListenerReasonAccepted,
			"Listener has been successfully translated")
		listener.SetCondition(gwapiv1.ListenerConditionProgrammed, metav1.ConditionTrue, gwapiv1.ListenerReasonProgrammed,
			"Sending translated listener configuration to the data plane")
		return
	}

	// Any condition on the listener apart from Programmed=true indicates an error.
				"Listener references have been resolved",
			)
		}
		// skip computing IR
		return
	}
}

func (t *Translator) validateAllowedNamespaces(listener *ListenerContext) {
	if listener.AllowedRoutes != nil &&
		listener.AllowedRoutes.Namespaces != nil &&
		listener.AllowedRoutes.Namespaces.From != nil &&
				gwapiv1.ListenerReasonInvalid,
				"The allowedRoutes.namespaces.selector field must be specified when allowedRoutes.namespaces.from is set to \"Selector\".",
			)
		} else {
			selector, err := metav1.LabelSelectorAsSelector(listener.AllowedRoutes.Namespaces.Selector)
			if err != nil {
				listener.SetCondition(
					gwapiv1.ListenerConditionProgrammed,
					metav1.ConditionFalse,
					gwapiv1.ListenerReasonInvalid,
					fmt.Sprintf("The allowedRoutes.namespaces.selector could not be parsed: %v.", err),
				)
			}

			listener.namespaceSelector = selector
		}
	}
}

func (t *Translator) validateTerminateModeAndGetTLSSecrets(
	listener *ListenerContext,
	resources *resource.Resources,
) ([]*corev1.Secret, []*x509.Certificate) {
	if len(listener.TLS.CertificateRefs) == 0 {
		listener.SetCondition(
			gwapiv1.ListenerConditionProgrammed,
			gwapiv1.ListenerReasonInvalid,
			"Listener must have at least 1 TLS certificate ref",
		)
		return nil, nil
	}

	var errs []status.ListenerError
				fromKind = resource.KindListenerSet
			}

			if !t.validateCrossNamespaceRef(
				crossNamespaceFrom{
					group:     fromGroup,
					kind:      fromKind,
			fmt.Sprintf("No valid secrets exist: %v", errors.Join(errList...)),
		)

		return nil, nil
	}

	validSecrets, certs, err := parseCertsFromTLSSecretsData(secrets)
				err.Reason(),
				fmt.Sprintf("No valid secrets exist: %v.", err.Error()),
			)
			return nil, nil
		} else {
			errs = append(errs, err)
		}
			fmt.Sprintf("Some secrets are invalid: %v", errors.Join(errList...)),
		)
	}
	return validSecrets, certs
}

func (t *Translator) validateTLSConfiguration(
	listener *ListenerContext,
	resources *resource.Resources,
) {
	switch listener.Protocol {
	case gwapiv1.HTTPProtocolType, gwapiv1.UDPProtocolType, gwapiv1.TCPProtocolType:
		if listener.TLS != nil {
				gwapiv1.ListenerReasonInvalid,
				fmt.Sprintf("Listener must not have TLS set when protocol is %s.", listener.Protocol),
			)
		}
	case gwapiv1.HTTPSProtocolType:
		if listener.TLS == nil {
				gwapiv1.ListenerReasonInvalid,
				fmt.Sprintf("Listener must have TLS set when protocol is %s.", listener.Protocol),
			)
			break
		}

		if listener.TLS.Mode != nil && *listener.TLS.Mode != gwapiv1.TLSModeTerminate {
			listener.SetCondition(
				gwapiv1.ListenerConditionProgrammed,
				metav1.ConditionFalse,
				"UnsupportedTLSMode",
				fmt.Sprintf("TLS %s mode is not supported, TLS mode must be Terminate.", *listener.TLS.Mode),
			)
			break
		}

		secrets, certs := t.validateTerminateModeAndGetTLSSecrets(listener, resources)
		listener.SetTLSSecrets(secrets)

		listener.tls.certDNSNames = make([]string, 0)
		for _, cert := range certs {
			listener.tls.certDNSNames = append(listener.tls.certDNSNames, cert.DNSNames...)
		}
	case gwapiv1.TLSProtocolType:
		if listener.TLS == nil {
				gwapiv1.ListenerReasonInvalid,
				fmt.Sprintf("Listener must have TLS set when protocol is %s.", listener.Protocol),
			)
			break
		}

		if listener.TLS.Mode != nil && *listener.TLS.Mode == gwapiv1.TLSModePassthrough {
			if len(listener.TLS.CertificateRefs) > 0 {
				listener.SetCondition(
					gwapiv1.ListenerConditionProgrammed,
					metav1.ConditionFalse,
					gwapiv1.ListenerReasonInvalid,
					"Listener must not have TLS certificate refs set for TLS mode Passthrough.",
				)
				break
			}
		}

		if listener.TLS.Mode != nil && *listener.TLS.Mode == gwapiv1.TLSModeTerminate {
			if len(listener.TLS.CertificateRefs) == 0 {
				listener.SetCondition(
					gwapiv1.ListenerConditionProgrammed,
					metav1.ConditionFalse,
					gwapiv1.ListenerReasonInvalid,
					"Listener must have TLS certificate refs set for TLS mode Terminate.",
				)
				break
			}
			secrets, _ := t.validateTerminateModeAndGetTLSSecrets(listener, resources)
			listener.SetTLSSecrets(secrets)
		}
	}

			gwapiv1.ListenerReasonNoValidCACertificate,
			message,
		)
	}
}

func (t *Translator) validateHostName(listener *ListenerContext) {
	if listener.Protocol == gwapiv1.UDPProtocolType || listener.Protocol == gwapiv1.TCPProtocolType {
		if listener.Hostname != nil {
			listener.SetCondition(
				gwapiv1.ListenerReasonInvalid,
				fmt.Sprintf("Listener must not have hostname set when protocol is %s.", listener.Protocol),
			)
		}
	}
}

func (t *Translator) validateAllowedRoutes(listener *ListenerContext, routeKinds ...gwapiv1.Kind) {
	canSupportKinds := make([]gwapiv1.RouteGroupKind, len(routeKinds))
	for i, routeKind := range routeKinds {
		canSupportKinds[i] = gwapiv1.RouteGroupKind{Group: GroupPtr(gwapiv1.GroupName), Kind: routeKind}
	}
	if listener.AllowedRoutes == nil || len(listener.AllowedRoutes.Kinds) == 0 {
		listener.SetSupportedKinds(canSupportKinds...)
		return
	}

	supportedRouteKinds := make([]gwapiv1.Kind, 0)
	supportedKinds := make([]gwapiv1.RouteGroupKind, 0)
	unSupportedKinds := make([]gwapiv1.RouteGroupKind, 0)
				gwapiv1.ListenerReasonInvalidRouteKinds,
				fmt.Sprintf("Group is not supported, group must be %s", gwapiv1.GroupName),
			)
			continue
		}

			gwapiv1.ListenerReasonInvalidRouteKinds,
			fmt.Sprintf("%s is not supported, kind must be one of %v", string(kind.Kind), printRouteKinds),
		)
	}

	listener.SetSupportedKinds(supportedKinds...)
}

type portListeners struct {
	listenerSets := sets.Set[string]{}
	for _, gateway := range gateways {
		for _, listener := range gateway.listeners {
			hostname := new(gwapiv1.Hostname)
			if listener.Hostname != nil {
				hostname = listener.Hostname
	}
}

func (t *Translator) validateConflictedLayer7Listeners(gateways []*GatewayContext) {
	// Iterate through all layer-7 (HTTP, HTTPS, TLS) listeners and collect info about protocols
	// and hostnames per port.
			if listener.Protocol == gwapiv1.UDPProtocolType || listener.Protocol == gwapiv1.TCPProtocolType {
				continue
			}
			if portListenerInfo[listener.Port] == nil {
				portListenerInfo[listener.Port] = &portListeners{
					protocols: sets.Set[string]{},
	for _, gateway := range gateways {
		portListenerInfo := map[gwapiv1.PortNumber]*portListeners{}
		for _, listener := range gateway.listeners {
			for _, protocol := range protocols {
				if listener.Protocol == protocol {
					if portListenerInfo[listener.Port] == nil {
	}
}

func (t *Translator) validateCrossNamespaceRef(from crossNamespaceFrom, to crossNamespaceTo, referenceGrants []*gwapiv1b1.ReferenceGrant) bool {
	for _, referenceGrant := range referenceGrants {
		// The ReferenceGrant must be defined in the namespace of
		// the "to" (the referent).
		if referenceGrant.Namespace != to.namespace {
			continue
		}

		// Check if the ReferenceGrant has a matching "from".
		var fromAllowed bool
		for _, refGrantFrom := range referenceGrant.Spec.From {
			if string(refGrantFrom.Namespace) == from.namespace && string(refGrantFrom.Group) == from.group && string(refGrantFrom.Kind) == from.kind {
				fromAllowed = true
				break
			}
		}
		if !fromAllowed {
			continue
		}

		// Check if the ReferenceGrant has a matching "to".
		var toAllowed bool
		for _, refGrantTo := range referenceGrant.Spec.To {
			if string(refGrantTo.Group) == to.group && string(refGrantTo.Kind) == to.kind && (refGrantTo.Name == nil || *refGrantTo.Name == "" || string(*refGrantTo.Name) == to.name) {
				toAllowed = true
				break
			}
		}
		if !toAllowed {
			continue
		}

		// If we got here, both the "from" and the "to" were allowed by this
		// reference grant.
		return true
	}

	// If we got here, no reference policy or reference grant allowed both the "from" and "to".
	return false
}

// Checks if a hostname is valid according to RFC 1123 and gateway API's requirement that it not be an IP address
				from.namespace)
		}

		if !t.validateCrossNamespaceRef(
			from,
			crossNamespaceTo{
				group:     "",
	// check if the cross-namespace reference is permitted
	if backendRef.Namespace != nil && string(*backendRef.Namespace) != "" &&
		string(*backendRef.Namespace) != ownerNamespace {
		if !t.validateCrossNamespaceRef(
			crossNamespaceFrom{
				group:     egv1a1.GroupName,
				kind:      policyKind,
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/utils/ptr"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/envoyproxy/gateway/internal/gatewayapi/resource"
		SectionIndex: make(map[types.NamespacedName]sets.Set[string], gatewayMapSize),
	}

	handledPolicies := make(map[types.NamespacedName]*egv1a1.SecurityPolicy, policyMapSize)

	// Map of attached Policy to Gateway. Used for policy merge process.
	// 4. Finally, the policies targeting Gateways

	// Build gateway policy maps, which are needed when processing the policies targeting xRoutes.
	t.buildGatewayPolicyMapForSecurity(securityPolicies, gateways, gatewayMap, gatewayPolicyMap)

	// Process the policies targeting RouteRules (HTTP + TCP)
	for _, currPolicy := range securityPolicies {
		policyName := utils.NamespacedName(currPolicy)
		targetRefs := getPolicyTargetRefs(currPolicy.Spec.PolicyTargetReferences, routes, currPolicy.Namespace)
		for _, currTarget := range targetRefs {
			// If the target is not a gateway, then it's an xRoute. If the section name is defined, then it's a route rule.
			if currTarget.Kind != resource.KindGateway && currTarget.SectionName != nil {
				policy, found := handledPolicies[policyName]
				if !found {
					policy = currPolicy
					handledPolicies[policyName] = policy
					res = append(res, policy)
				}
		}
	}
	// Process the policies targeting xRoutes (HTTP + TCP)
	for _, currPolicy := range securityPolicies {
		policyName := utils.NamespacedName(currPolicy)
		targetRefs := getPolicyTargetRefs(currPolicy.Spec.PolicyTargetReferences, routes, currPolicy.Namespace)
		for _, currTarget := range targetRefs {
			// If the target is not a gateway, then it's an xRoute. If the section name is not defined, then it's a route.
			if currTarget.Kind != resource.KindGateway && currTarget.SectionName == nil {
				policy, found := handledPolicies[policyName]
				if !found {
					policy = currPolicy
					handledPolicies[policyName] = policy
					res = append(res, policy)
				}
		}
	}
	// Process the policies targeting Listeners
	for _, currPolicy := range securityPolicies {
		policyName := utils.NamespacedName(currPolicy)
		targetRefs := getPolicyTargetRefs(currPolicy.Spec.PolicyTargetReferences, gateways, currPolicy.Namespace)
		for _, currTarget := range targetRefs {
			// If the target is a gateway and the section name is defined, then it's a listener.
			if currTarget.Kind == resource.KindGateway && currTarget.SectionName != nil {
				policy, found := handledPolicies[policyName]
				if !found {
					policy = currPolicy
					handledPolicies[policyName] = policy
					res = append(res, policy)
				}
		}
	}
	// Process the policies targeting Gateways
	for _, currPolicy := range securityPolicies {
		policyName := utils.NamespacedName(currPolicy)
		targetRefs := getPolicyTargetRefs(currPolicy.Spec.PolicyTargetReferences, gateways, currPolicy.Namespace)
		for _, currTarget := range targetRefs {
			// If the target is a gateway and the section name is not defined, then it's a gateway.
			if currTarget.Kind == resource.KindGateway && currTarget.SectionName == nil {
				policy, found := handledPolicies[policyName]
				if !found {
					policy = currPolicy
					handledPolicies[policyName] = policy
					res = append(res, policy)
				}
	gateways []*GatewayContext,
	gatewayMap map[types.NamespacedName]*policyGatewayTargetContext,
	gatewayPolicyMap map[NamespacedNameWithSection]*egv1a1.SecurityPolicy,
) {
	for _, currPolicy := range securityPolicies {
		targetRefs := getPolicyTargetRefs(currPolicy.Spec.PolicyTargetReferences, gateways, currPolicy.Namespace)
		for _, currTarget := range targetRefs {
			if currTarget.Kind == resource.KindGateway {
				// Check if the gateway exists
				key := types.NamespacedName{
					Name:      string(currTarget.Name),
					Namespace: currPolicy.Namespace,
				}
				gateway, ok := gatewayMap[key]
				if !ok {
	gatewayPolicyMerged *GatewayPolicyRouteMap,
	gatewayPolicyMap map[NamespacedNameWithSection]*egv1a1.SecurityPolicy,
	policy *egv1a1.SecurityPolicy,
	currTarget gwapiv1.LocalPolicyTargetReferenceWithSectionName,
) {
	var (
		targetedRoute RouteContext
		resolveErr    *status.PolicyResolveError
	)

	targetedRoute, resolveErr = resolveSecurityPolicyRouteTargetRef(policy, currTarget, routeMap)
	// Skip if the route is not found
	// It's not necessarily an error because the SecurityPolicy may be
	// reconciled by multiple controllers. And the other controller may
	// Check if merging is enabled
	if policy.Spec.MergeType == nil {
		// No merging - use existing translation logic
		if err := t.translateSecurityPolicyForRoute(policy, targetedRoute, currTarget, resources, xdsIR, nil, nil); err != nil {
			status.SetTranslationErrorForPolicyAncestors(&policy.Status,
				ancestorRefs,
				t.GatewayControllerName,

				if gwPolicy == nil && listenerPolicy == nil {
					// No parent policy found, fall back to current policy
					if err := t.translateSecurityPolicyForRoute(policy, targetedRoute, currTarget, resources, xdsIR, &gwNN, &listener.Name); err != nil {
						status.SetConditionForPolicyAncestor(&policy.Status,
							&ancestorRef,
							t.GatewayControllerName,
				}

				// Merge with parent policy
				mergedPolicy, err := mergeSecurityPolicy(policy, parentPolicy)
				if err != nil {
					status.SetConditionForPolicyAncestor(&policy.Status,
						&ancestorRef,
				}

				// Apply merged policy
				if err := t.translateSecurityPolicyForRoute(mergedPolicy, targetedRoute, currTarget, resources, xdsIR, &gwNN, &listener.Name); err != nil {
					status.SetConditionForPolicyAncestor(&policy.Status,
						&ancestorRef,
						t.GatewayControllerName,
	key := policyTargetRouteKey{
		Kind:      string(currTarget.Kind),
		Name:      string(currTarget.Name),
		Namespace: policy.Namespace,
	}
	overriddenTargetsMessage := getOverriddenTargetsMessageForRoute(routeMap[key], currTarget.SectionName)
	if overriddenTargetsMessage != "" {
	gatewayRouteMap *GatewayPolicyRouteMap,
	gatewayPolicyMergedMap *GatewayPolicyRouteMap,
	policy *egv1a1.SecurityPolicy,
	currTarget gwapiv1.LocalPolicyTargetReferenceWithSectionName,
) {
	var (
		targetedGateway *GatewayContext
		resolveErr      *status.PolicyResolveError
	)

	targetedGateway, resolveErr = resolveSecurityPolicyGatewayTargetRef(policy, currTarget, gatewayMap)
	// Skip if the gateway is not found
	// It's not necessarily an error because the SecurityPolicy may be
	// reconciled by multiple controllers. And the other controller may
}

func resolveSecurityPolicyGatewayTargetRef(
	policy *egv1a1.SecurityPolicy,
	target gwapiv1.LocalPolicyTargetReferenceWithSectionName,
	gateways map[types.NamespacedName]*policyGatewayTargetContext,
) (*GatewayContext, *status.PolicyResolveError) {
	// Find the Gateway
	key := types.NamespacedName{
		Name:      string(target.Name),
		Namespace: policy.Namespace,
	}
	gateway, ok := gateways[key]

}

func resolveSecurityPolicyRouteTargetRef(
	policy *egv1a1.SecurityPolicy,
	target gwapiv1.LocalPolicyTargetReferenceWithSectionName,
	routes map[policyTargetRouteKey]*policyRouteTargetContext,
) (RouteContext, *status.PolicyResolveError) {
	// Check if the route exists
	key := policyTargetRouteKey{
		Kind:      string(target.Kind),
		Name:      string(target.Name),
		Namespace: policy.Namespace,
	}
	route, ok := routes[key]


func (t *Translator) translateSecurityPolicyForRoute(
	policy *egv1a1.SecurityPolicy,
	route RouteContext,
	target gwapiv1.LocalPolicyTargetReferenceWithSectionName,
	resources *resource.Resources,
	xdsIR resource.XdsIRMap,
	policyTargetGateway *types.NamespacedName,
	if policy.Spec.BasicAuth != nil {
		if basicAuth, err = t.buildBasicAuth(
			policy,
			resources); err != nil {
			err = perr.WithMessage(err, "BasicAuth")
			errs = errors.Join(errs, err)
		}
	if policy.Spec.APIKeyAuth != nil {
		if apiKeyAuth, err = t.buildAPIKeyAuth(
			policy,
			resources); err != nil {
			err = perr.WithMessage(err, "APIKeyAuth")
			errs = errors.Join(errs, err)
		}
	}

	if policy.Spec.Authorization != nil {
		if authorization, err = t.buildAuthorization(policy); err != nil {
			err = perr.WithMessage(err, "Authorization")
			errs = errors.Join(errs, err)
		}
		if policy.Spec.ExtAuth != nil {
			if extAuth, extAuthErr = t.buildExtAuth(
				policy,
				resources,
				gtwCtx); extAuthErr != nil {
				extAuthErr = perr.WithMessage(extAuthErr, "ExtAuth")
				errs = errors.Join(errs, extAuthErr)
			}
		if policy.Spec.OIDC != nil {
			if oidc, err = t.buildOIDC(
				policy,
				resources,
				gtwCtx); err != nil {
				err = perr.WithMessage(err, "OIDC")
				errs = errors.Join(errs, err)
				hasNonExtAuthError = true
		if policy.Spec.JWT != nil {
			if jwt, err = t.buildJWT(
				policy,
				resources,
				gtwCtx); err != nil {
				err = perr.WithMessage(err, "JWT")
				errs = errors.Join(errs, err)
				hasNonExtAuthError = true
						continue
					}
					// Only authorization for TCP
					authCopy := *authorization
					r.Authorization = &authCopy
				}
			}
		case resource.KindHTTPRoute, resource.KindGRPCRoute:
func (t *Translator) translateSecurityPolicyForGateway(
	policy *egv1a1.SecurityPolicy,
	gtwCtx *GatewayContext,
	target gwapiv1.LocalPolicyTargetReferenceWithSectionName,
	resources *resource.Resources,
	xdsIR resource.XdsIRMap,
) error {
	// Build IR
	var (
		cors                  *ir.CORS
		jwt                   *ir.JWT
	if policy.Spec.JWT != nil {
		if jwt, err = t.buildJWT(
			policy,
			resources,
			gtwCtx); err != nil {
			err = perr.WithMessage(err, "JWT")
			errs = errors.Join(errs, err)
		}
	if policy.Spec.OIDC != nil {
		if oidc, err = t.buildOIDC(
			policy,
			resources,
			gtwCtx); err != nil {
			err = perr.WithMessage(err, "OIDC")
			errs = errors.Join(errs, err)
		}
	if policy.Spec.BasicAuth != nil {
		if basicAuth, err = t.buildBasicAuth(
			policy,
			resources); err != nil {
			err = perr.WithMessage(err, "BasicAuth")
			errs = errors.Join(errs, err)
		}
	if policy.Spec.APIKeyAuth != nil {
		if apiKeyAuth, err = t.buildAPIKeyAuth(
			policy,
			resources); err != nil {
			err = perr.WithMessage(err, "APIKeyAuth")
			errs = errors.Join(errs, err)
		}
	}

	if policy.Spec.Authorization != nil {
		if authorization, err = t.buildAuthorization(policy); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	if policy.Spec.ExtAuth != nil {
		if extAuth, extAuthErr = t.buildExtAuth(
			policy,
			resources,
			gtwCtx); extAuthErr != nil {
			extAuthErr = perr.WithMessage(extAuthErr, "ExtAuth")
			errs = errors.Join(errs, extAuthErr)
		}

func (t *Translator) buildJWT(
	policy *egv1a1.SecurityPolicy,
	resources *resource.Resources,
	gtwCtx *GatewayContext,
) (*ir.JWT, error) {
		return nil, err
	}

	providers := make([]ir.JWTProvider, 0, len(policy.Spec.JWT.Providers))
	for i, p := range policy.Spec.JWT.Providers {
		provider := ir.JWTProvider{
			ExtractFrom:    p.ExtractFrom,
		}
		if p.RemoteJWKS != nil {
			remoteJWKS, err := t.buildRemoteJWKS(policy, p.RemoteJWKS, i, resources, gtwCtx)
			if err != nil {
				return nil, err
			}
			provider.RemoteJWKS = remoteJWKS
		} else {
			localJWKS, err := t.buildLocalJWKS(policy, p.LocalJWKS)
			if err != nil {
				return nil, err
			}

func (t *Translator) buildOIDC(
	policy *egv1a1.SecurityPolicy,
	resources *resource.Resources,
	gtwCtx *GatewayContext,
) (*ir.OIDC, error) {
		err                    error
	)

	if provider, err = t.buildOIDCProvider(policy, resources, gtwCtx); err != nil {
		return nil, err
	}

	from := crossNamespaceFrom{
		group:     egv1a1.GroupName,
		kind:      resource.KindSecurityPolicy,
		namespace: policy.Namespace,
	}

	// Client ID can be specified either as a string or as a reference to a secret.
	switch {
	case oidc.ClientID != nil:
		clientID = *oidc.ClientID
	case oidc.ClientIDRef != nil:
		var clientIDSecret *corev1.Secret
		if clientIDSecret, err = t.validateSecretRef(true, from, *oidc.ClientIDRef, resources); err != nil {
			return nil, err
		return nil, fmt.Errorf("client ID must be specified in OIDC policy %s/%s", policy.Namespace, policy.Name)
	}

	if clientSecret, err = t.validateSecretRef(true, from, oidc.ClientSecret, resources); err != nil {
		return nil, err
	}
		disableTokenEncryption = *oidc.DisableTokenEncryption
	}

	// Generate a unique cookie suffix for oauth filters.
	// This is to avoid cookie name collision when multiple security policies are applied
	// to the same route.
	suffix := utils.Digest32(string(policy.UID))

	// Get the HMAC secret.
	// HMAC secret is generated by the CertGen job and stored in a secret
	}

	irOIDC := &ir.OIDC{
		Name:                   irConfigName(policy),
		Provider:               *provider,
		ClientID:               clientID,
		ClientSecret:           clientSecretBytes,

func (t *Translator) buildOIDCProvider(
	policy *egv1a1.SecurityPolicy,
	resources *resource.Resources,
	gtwCtx *GatewayContext,
) (*ir.OIDCProvider, error) {
		protocol = ir.HTTP
	}

	if len(provider.BackendRefs) > 0 {
		if rd, err = t.translateExtServiceBackendRefs(
			policy, provider.BackendRefs, protocol, resources, gtwCtx, "oidc", 0); err != nil {
			return nil, err
		}
	}

func (t *Translator) buildAPIKeyAuth(
	policy *egv1a1.SecurityPolicy,
	resources *resource.Resources,
) (*ir.APIKeyAuth, error) {
	from := crossNamespaceFrom{
		group:     egv1a1.GroupName,
		kind:      resource.KindSecurityPolicy,
		namespace: policy.Namespace,
	}

	expected := len(policy.Spec.APIKeyAuth.CredentialRefs)

func (t *Translator) buildBasicAuth(
	policy *egv1a1.SecurityPolicy,
	resources *resource.Resources,
) (*ir.BasicAuth, error) {
	var (
		err         error
	)

	from := crossNamespaceFrom{
		group:     egv1a1.GroupName,
		kind:      resource.KindSecurityPolicy,
		namespace: policy.Namespace,
	}
	if usersSecret, err = t.validateSecretRef(true, from, basicAuth.Users, resources); err != nil {
		return nil, err
	}

	return &ir.BasicAuth{
		Name:                  irConfigName(policy),
		Users:                 usersSecretBytes,
		ForwardUsernameHeader: basicAuth.ForwardUsernameHeader,
	}, nil

func (t *Translator) buildExtAuth(
	policy *egv1a1.SecurityPolicy,
	resources *resource.Resources,
	gtwCtx *GatewayContext,
) (*ir.ExtAuth, error) {
		contextExtensions []*ir.ContextExtention
	)

	// These are sanity checks, they should never happen because the API server
	// should have caught them
	if http == nil && grpc == nil {
	}

	if rd, err = t.translateExtServiceBackendRefs(
		policy, backendRefs, protocol, resources, gtwCtx, "extauth", 0); err != nil {
		return nil, err
	}

		// When translated to XDS, the authority is used on the filter level not on the cluster level.
		// There's no way to translate to XDS and use a different authority for each backendref
		if authority == "" {
			authority = t.backendRefAuthority(&backendRef.BackendObjectReference, policy)
		}
	}

		return nil, err
	}

	if contextExtensions, err = t.buildContextExtensions(policy.Spec.ExtAuth.ContextExtensions, policy.Namespace); err != nil {
		return nil, err
	}

	extAuth := &ir.ExtAuth{
		Name:                 irConfigName(policy),
		HeadersToExtAuth:     policy.Spec.ExtAuth.HeadersToExtAuth,
		ContextExtensions:    contextExtensions,
		FailOpen:             policy.Spec.ExtAuth.FailOpen,
			Destination:      *rd,
			Authority:        authority,
			Path:             ptr.Deref(http.Path, ""),
			HeadersToBackend: http.HeadersToBackend,
		}
	} else {

func (t *Translator) buildContextExtensions(
	contextExtensions []*egv1a1.ContextExtension,
	policyNs string,
) ([]*ir.ContextExtention, error) {
	if len(contextExtensions) == 0 {
		return nil, nil
	for _, ext := range contextExtensions {
		var value ir.PrivateBytes
		if ext.Type == egv1a1.ContextExtensionValueTypeValueRef {
			var err error
			if value, err = t.getContextExtensionValueFromRef(ext.ValueRef, policyNs); err != nil {
				return nil, err
			}
		} else if ext.Value != nil {
	return fmt.Sprintf("%s.%s", backendRef.Name, backendNamespace)
}

func (t *Translator) buildAuthorization(policy *egv1a1.SecurityPolicy) (*ir.Authorization, error) {
	var (
		authorization = policy.Spec.Authorization
		irAuth        = &ir.Authorization{}
		defaultAction = egv1a1.AuthorizationActionDeny
	)

	if authorization.DefaultAction != nil {
		defaultAction = *authorization.DefaultAction
	}
		if rule.Name != nil && *rule.Name != "" {
			name = *rule.Name
		} else {
			name = defaultAuthorizationRuleName(policy, i)
		}
		irAuth.Rules = append(irAuth.Rules, &ir.AuthorizationRule{
			Name:      name,
		strconv.Itoa(index))
}

// mergeSecurityPolicy merges a route-level SecurityPolicy with a parent (Gateway/Listener) SecurityPolicy.
func mergeSecurityPolicy(routePolicy, parentPolicy *egv1a1.SecurityPolicy) (*egv1a1.SecurityPolicy, error) {
	if routePolicy.Spec.MergeType == nil || parentPolicy == nil {
		return routePolicy, nil
	}

	return utils.Merge[*egv1a1.SecurityPolicy](parentPolicy, routePolicy, *routePolicy.Spec.MergeType)
}
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/envoyproxy/gateway/internal/gatewayapi/resource"
}

type targetRefWithTimestamp struct {
	gwapiv1.LocalPolicyTargetReferenceWithSectionName
	CreationTimestamp metav1.Time
}

func selectorFromTargetSelector(selector egv1a1.TargetSelector) labels.Selector {
	l, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels:      selector.MatchLabels,
	return l
}

func getPolicyTargetRefs[T client.Object](policy egv1a1.PolicyTargetReferences, potentialTargets []T, policyNamespace string) []gwapiv1.LocalPolicyTargetReferenceWithSectionName {
	dedup := sets.New[targetRefWithTimestamp]()
	for _, currSelector := range policy.TargetSelectors {
			if labelSelector.Matches(labels.Set(obj.GetLabels())) {
				dedup.Insert(targetRefWithTimestamp{
					CreationTimestamp: obj.GetCreationTimestamp(),
					LocalPolicyTargetReferenceWithSectionName: gwapiv1.LocalPolicyTargetReferenceWithSectionName{
						LocalPolicyTargetReference: gwapiv1.LocalPolicyTargetReference{
							Group: gwapiv1.Group(gvk.Group),
							Kind:  gwapiv1.Kind(gvk.Kind),
							Name:  gwapiv1.ObjectName(obj.GetName()),
						},
					},
				})
			}
	})
	ret := make([]gwapiv1.LocalPolicyTargetReferenceWithSectionName, len(selectorsList))
	for i, v := range selectorsList {
		ret[i] = v.LocalPolicyTargetReferenceWithSectionName
	}
	// Plain targetRefs in the policy don't have an associated creation timestamp, but can still refer
	// to targets that were already found via the selectors. Only add them to the returned list if
