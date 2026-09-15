package main

	"k8s.io/apimachinery/pkg/util/validation"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/envoyproxy/gateway/internal/gatewayapi/resource"
				unsupportedFilters = true
				continue
			}

			// BackendRef URLRewrite only supports hostname rewrites.
			// Path rewrites are not supported because Envoy weighted clusters
			// do not support path rewrite actions.
			if filter.Type == gwapiv1.HTTPRouteFilterURLRewrite &&
				filter.URLRewrite != nil &&
				filter.URLRewrite.Path != nil {
				return status.NewRouteStatusError(
					errors.New("URLRewrite path modifier is not supported within BackendRef"),
					status.RouteReasonUnsupportedRefValue,
				)
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
		if !isCrossNamespaceReferencePermitted(
			crossNamespaceFrom{
				group:     gwapiv1.GroupName,
				kind:      string(routeKind),
		return
	}

	onlyResolvedRefFailure := len(lConditions) == 1 && lConditions[0].Type == string(gwapiv1.ListenerConditionResolvedRefs)
	if onlyResolvedRefFailure {
		switch lConditions[0].Reason {
		case string(status.ListenerReasonPartiallyInvalidCertificateRef):
			// The listener is ready because we program it using only the valid certificates.
			listener.SetCondition(gwapiv1.ListenerConditionAccepted, metav1.ConditionTrue, gwapiv1.ListenerReasonAccepted,
				"Listener has been successfully translated")
			listener.SetCondition(gwapiv1.ListenerConditionProgrammed, metav1.ConditionTrue, gwapiv1.ListenerReasonProgrammed,
				"Sending translated listener configuration to the data plane")
			return
		case string(gwapiv1.ListenerReasonInvalidCertificateRef):
			// The listener configuration is semantically valid, but the listener cannot serve traffic with an invalid certificate.
			listener.SetCondition(gwapiv1.ListenerConditionAccepted, metav1.ConditionTrue, gwapiv1.ListenerReasonAccepted,
				"Listener has been successfully translated")
			listener.SetCondition(gwapiv1.ListenerConditionProgrammed, metav1.ConditionFalse, gwapiv1.ListenerReasonInvalid,
				"Listener is invalid, see other Conditions for details.")
			return
		}
	}

	// Any condition on the listener apart from Programmed=true indicates an error.
				"Listener references have been resolved",
			)
		}
	}
}

// hasInvalidCondition checks if a listener has been marked as invalid during per-listener validation.
// A listener is considered invalid if it has Programmed=False, Accepted=False, or ResolvedRefs=False
// (except for the special case of PartiallyInvalidCertificateRef which is allowed).
// This is used during conflict resolution to skip invalid listeners so they don't block valid ones.
func hasInvalidCondition(listener *ListenerContext) bool {
	conditions := listener.GetConditions()
	for _, cond := range conditions {
		if cond.Type == string(gwapiv1.ListenerConditionProgrammed) && cond.Status == metav1.ConditionFalse {
			return true
		}
		if cond.Type == string(gwapiv1.ListenerConditionAccepted) && cond.Status == metav1.ConditionFalse {
			return true
		}
		// ResolvedRefs=False is invalid except for PartiallyInvalidCertificateRef which allows
		// the listener to still be programmed with valid certificates
		if cond.Type == string(gwapiv1.ListenerConditionResolvedRefs) &&
			cond.Status == metav1.ConditionFalse &&
			cond.Reason != string(status.ListenerReasonPartiallyInvalidCertificateRef) {
			return true
		}
	}
	return false
}

// isSpecValidForConflictChecks returns whether a listener should participate in
// conflict detection. In the normal translation flow this is driven by
// listener.specValid. The fallback to hasInvalidCondition exists only for unit
// tests that invoke conflict checks directly without running per-listener spec
// validation (Phase 1) first. Production code paths always run validateListenerSpec
// before conflict detection.
func isSpecValidForConflictChecks(listener *ListenerContext) bool {
	if listener.specValid {
		return true
	}
	return !hasInvalidCondition(listener)
}

// validateAllowedNamespaces validates namespace selector configuration.
// Returns true if the namespace spec is valid, false otherwise.
func (t *Translator) validateAllowedNamespaces(listener *ListenerContext) bool {
	if listener.AllowedRoutes != nil &&
		listener.AllowedRoutes.Namespaces != nil &&
		listener.AllowedRoutes.Namespaces.From != nil &&
				gwapiv1.ListenerReasonInvalid,
				"The allowedRoutes.namespaces.selector field must be specified when allowedRoutes.namespaces.from is set to \"Selector\".",
			)
			return false
		}
		selector, err := metav1.LabelSelectorAsSelector(listener.AllowedRoutes.Namespaces.Selector)
		if err != nil {
			listener.SetCondition(
				gwapiv1.ListenerConditionProgrammed,
				metav1.ConditionFalse,
				gwapiv1.ListenerReasonInvalid,
				fmt.Sprintf("The allowedRoutes.namespaces.selector could not be parsed: %v.", err),
			)
			return false
		}

		listener.namespaceSelector = selector
	}
	return true
}

func (t *Translator) validateTerminateModeAndGetTLSSecrets(
	listener *ListenerContext,
	resources *resource.Resources,
) ([]*corev1.Secret, []*x509.Certificate, bool) {
	if len(listener.TLS.CertificateRefs) == 0 {
		listener.SetCondition(
			gwapiv1.ListenerConditionProgrammed,
			gwapiv1.ListenerReasonInvalid,
			"Listener must have at least 1 TLS certificate ref",
		)
		return nil, nil, false
	}

	var errs []status.ListenerError
				fromKind = resource.KindListenerSet
			}

			if !isCrossNamespaceReferencePermitted(
				crossNamespaceFrom{
					group:     fromGroup,
					kind:      fromKind,
			fmt.Sprintf("No valid secrets exist: %v", errors.Join(errList...)),
		)

		return nil, nil, false
	}

	validSecrets, certs, err := parseCertsFromTLSSecretsData(secrets)
				err.Reason(),
				fmt.Sprintf("No valid secrets exist: %v.", err.Error()),
			)
			return nil, nil, false
		} else {
			errs = append(errs, err)
		}
			fmt.Sprintf("Some secrets are invalid: %v", errors.Join(errList...)),
		)
	}
	return validSecrets, certs, true
}

// validateTLSConfiguration validates TLS configuration per protocol.
// Returns true if the TLS spec is valid, false otherwise.
func (t *Translator) validateTLSConfiguration(
	listener *ListenerContext,
	resources *resource.Resources,
) bool {
	specValid := true

	switch listener.Protocol {
	case gwapiv1.HTTPProtocolType, gwapiv1.UDPProtocolType, gwapiv1.TCPProtocolType:
		if listener.TLS != nil {
				gwapiv1.ListenerReasonInvalid,
				fmt.Sprintf("Listener must not have TLS set when protocol is %s.", listener.Protocol),
			)
			specValid = false
		}
	case gwapiv1.HTTPSProtocolType:
		if listener.TLS == nil {
				gwapiv1.ListenerReasonInvalid,
				fmt.Sprintf("Listener must have TLS set when protocol is %s.", listener.Protocol),
			)
			specValid = false
		} else {
			if listener.TLS.Mode != nil && *listener.TLS.Mode != gwapiv1.TLSModeTerminate {
				listener.SetCondition(
					gwapiv1.ListenerConditionProgrammed,
					metav1.ConditionFalse,
					"UnsupportedTLSMode",
					fmt.Sprintf("TLS %s mode is not supported, TLS mode must be Terminate.", *listener.TLS.Mode),
				)
				specValid = false
			} else {
				secrets, certs, ok := t.validateTerminateModeAndGetTLSSecrets(listener, resources)
				listener.SetTLSSecrets(secrets)

				if !ok {
					specValid = false
				}

				listener.tls.certDNSNames = make([]string, 0)
				for _, cert := range certs {
					listener.tls.certDNSNames = append(listener.tls.certDNSNames, cert.DNSNames...)
				}
			}
		}
	case gwapiv1.TLSProtocolType:
		if listener.TLS == nil {
				gwapiv1.ListenerReasonInvalid,
				fmt.Sprintf("Listener must have TLS set when protocol is %s.", listener.Protocol),
			)
			specValid = false
		} else {
			if listener.TLS.Mode != nil && *listener.TLS.Mode == gwapiv1.TLSModePassthrough {
				if len(listener.TLS.CertificateRefs) > 0 {
					listener.SetCondition(
						gwapiv1.ListenerConditionProgrammed,
						metav1.ConditionFalse,
						gwapiv1.ListenerReasonInvalid,
						"Listener must not have TLS certificate refs set for TLS mode Passthrough.",
					)
					specValid = false
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
					specValid = false
				} else {
					secrets, _, ok := t.validateTerminateModeAndGetTLSSecrets(listener, resources)
					listener.SetTLSSecrets(secrets)

					if !ok {
						specValid = false
					}
				}
			}
		}
	}

			gwapiv1.ListenerReasonNoValidCACertificate,
			message,
		)
		specValid = false
	}

	return specValid
}

// validateHostName validates hostname configuration per protocol.
// Returns true if the hostname spec is valid, false otherwise.
func (t *Translator) validateHostName(listener *ListenerContext) bool {
	if listener.Protocol == gwapiv1.UDPProtocolType || listener.Protocol == gwapiv1.TCPProtocolType {
		if listener.Hostname != nil {
			listener.SetCondition(
				gwapiv1.ListenerReasonInvalid,
				fmt.Sprintf("Listener must not have hostname set when protocol is %s.", listener.Protocol),
			)
			return false
		}
	}
	return true
}

// validateAllowedRoutes validates allowed route kinds configuration.
// Returns true if the allowed routes spec is valid, false otherwise.
func (t *Translator) validateAllowedRoutes(listener *ListenerContext, routeKinds ...gwapiv1.Kind) bool {
	canSupportKinds := make([]gwapiv1.RouteGroupKind, len(routeKinds))
	for i, routeKind := range routeKinds {
		canSupportKinds[i] = gwapiv1.RouteGroupKind{Group: GroupPtr(gwapiv1.GroupName), Kind: routeKind}
	}
	if listener.AllowedRoutes == nil || len(listener.AllowedRoutes.Kinds) == 0 {
		listener.SetSupportedKinds(canSupportKinds...)
		return true
	}

	specValid := true
	supportedRouteKinds := make([]gwapiv1.Kind, 0)
	supportedKinds := make([]gwapiv1.RouteGroupKind, 0)
	unSupportedKinds := make([]gwapiv1.RouteGroupKind, 0)
				gwapiv1.ListenerReasonInvalidRouteKinds,
				fmt.Sprintf("Group is not supported, group must be %s", gwapiv1.GroupName),
			)
			specValid = false
			continue
		}

			gwapiv1.ListenerReasonInvalidRouteKinds,
			fmt.Sprintf("%s is not supported, kind must be one of %v", string(kind.Kind), printRouteKinds),
		)
		specValid = false
	}

	listener.SetSupportedKinds(supportedKinds...)
	return specValid
}

type portListeners struct {
	listenerSets := sets.Set[string]{}
	for _, gateway := range gateways {
		for _, listener := range gateway.listeners {
			// Skip listeners that are already marked as invalid from per-listener validation.
			// This prevents an invalid first listener from blocking valid subsequent listeners.
			if !isSpecValidForConflictChecks(listener) {
				continue
			}
			hostname := new(gwapiv1.Hostname)
			if listener.Hostname != nil {
				hostname = listener.Hostname
	}
}

// validateConflictedProtocolsListeners checks for listeners that have conflicting protocols on the same port.
// UDP can coexist with any protocol. HTTPS and TLS are treated as compatible via getProtocolForListener.
func (t *Translator) validateConflictedProtocolsListeners(gateways []*GatewayContext) {
	validateByPort := func(listeners []*ListenerContext) {
		portListenerInfo := map[gwapiv1.PortNumber][]*ListenerContext{}
		for _, listener := range listeners {
			if !isSpecValidForConflictChecks(listener) || !isSupportedListenerProtocol(listener.Protocol) {
				continue
			}
			portListenerInfo[listener.Port] = append(portListenerInfo[listener.Port], listener)
		}

		for _, listenersOnPort := range portListenerInfo {
			nonUDPProtocols := sets.New[string]()
			nonListenerSetCount := 0
			for _, listener := range listenersOnPort {
				protocol := getProtocolForListener(listener)
				if protocol == string(gwapiv1.UDPProtocolType) {
					continue
				}
				nonUDPProtocols.Insert(protocol)
				if !listener.isFromListenerSet() {
					nonListenerSetCount++
				}
			}

			// No protocol conflict when all non-UDP listeners are compatible.
			if len(nonUDPProtocols) <= 1 {
				continue
			}

			// If there are more than 1 non-UDP protocols and more than 1 listener not from ListenerSet,
			// we cannot determine a clear winner and all listeners on this port are in conflict.
			if nonListenerSetCount > 1 {
				// If any conflicted listener is not from ListenerSet, do not pick a winner.
				for _, listener := range listenersOnPort {
					if getProtocolForListener(listener) == string(gwapiv1.UDPProtocolType) {
						continue
					}
					listener.SetCondition(
						gwapiv1.ListenerConditionConflicted,
						metav1.ConditionTrue,
						gwapiv1.ListenerReasonProtocolConflict,
						"All listeners for a given port must use a compatible protocol",
					)
				}
				continue
			}

			// When nonListenerSetCount == 1, explicitly pick the Gateway-owned listener as winner.
			// When nonListenerSetCount == 0, pick the first ListenerSet listener as winner.
			// Note: UDP conflicts are handled by validateConflictedLayer4Listeners, so we skip
			// UDP listeners here (this branch is only reached when len(nonUDPProtocols) > 1).
			var winnerProtocol string
			if nonListenerSetCount == 1 {
				// Find and use the non-ListenerSet listener's protocol as the winner
				for _, listener := range listenersOnPort {
					protocol := getProtocolForListener(listener)
					if !listener.isFromListenerSet() && protocol != string(gwapiv1.UDPProtocolType) {
						winnerProtocol = protocol
						break
					}
				}
			}

			for _, listener := range listenersOnPort {
				protocol := getProtocolForListener(listener)
				// Skip UDP listeners as they are handled by validateConflictedLayer4Listeners
				if protocol == string(gwapiv1.UDPProtocolType) {
					continue
				}

				// If we have an explicit winner protocol, use it; otherwise first one wins
				if winnerProtocol != "" {
					if protocol != winnerProtocol {
						listener.SetCondition(
							gwapiv1.ListenerConditionConflicted,
							metav1.ConditionTrue,
							gwapiv1.ListenerReasonProtocolConflict,
							"All listeners for a given port must use a compatible protocol",
						)
					}
				} else {
					// All conflicted listeners are from ListenerSet, first one wins
					if winnerProtocol == "" {
						winnerProtocol = protocol
					} else if protocol != winnerProtocol {
						listener.SetCondition(
							gwapiv1.ListenerConditionConflicted,
							metav1.ConditionTrue,
							gwapiv1.ListenerReasonProtocolConflict,
							"All listeners for a given port must use a compatible protocol",
						)
					}
				}
			}
		}
	}

	for _, gateway := range gateways {
		validateByPort(gateway.listeners)
	}

	if t.MergeGateways {
		allListeners := make([]*ListenerContext, 0)
		for _, gateway := range gateways {
			allListeners = append(allListeners, gateway.listeners...)
		}
		validateByPort(allListeners)
	}
}

func (t *Translator) validateConflictedLayer7Listeners(gateways []*GatewayContext) {
	// Iterate through all layer-7 (HTTP, HTTPS, TLS) listeners and collect info about protocols
	// and hostnames per port.
			if listener.Protocol == gwapiv1.UDPProtocolType || listener.Protocol == gwapiv1.TCPProtocolType {
				continue
			}
			// Skip listeners that are already marked as invalid from per-listener validation.
			// This prevents an invalid first listener from blocking valid subsequent listeners.
			if !isSpecValidForConflictChecks(listener) {
				continue
			}
			if portListenerInfo[listener.Port] == nil {
				portListenerInfo[listener.Port] = &portListeners{
					protocols: sets.Set[string]{},
	for _, gateway := range gateways {
		portListenerInfo := map[gwapiv1.PortNumber]*portListeners{}
		for _, listener := range gateway.listeners {
			// Skip listeners that are already marked as invalid from per-listener validation.
			// This prevents an invalid first listener from blocking valid subsequent listeners.
			if !isSpecValidForConflictChecks(listener) {
				continue
			}
			for _, protocol := range protocols {
				if listener.Protocol == protocol {
					if portListenerInfo[listener.Port] == nil {
	}
}

func getProtocolForListener(listener *ListenerContext) string {
	switch listener.Protocol {
	// HTTPS and TLS can co-exist on the same port.
	case gwapiv1.HTTPSProtocolType, gwapiv1.TLSProtocolType:
		return "https/tls"
	default:
		return string(listener.Protocol)
	}
}

func isSupportedListenerProtocol(protocol gwapiv1.ProtocolType) bool {
	switch protocol {
	case gwapiv1.HTTPProtocolType, gwapiv1.HTTPSProtocolType, gwapiv1.TLSProtocolType,
		gwapiv1.TCPProtocolType, gwapiv1.UDPProtocolType:
		return true
	default:
		return false
	}
}

// Checks if a hostname is valid according to RFC 1123 and gateway API's requirement that it not be an IP address
				from.namespace)
		}

		if !isCrossNamespaceReferencePermitted(
			from,
			crossNamespaceTo{
				group:     "",
	// check if the cross-namespace reference is permitted
	if backendRef.Namespace != nil && string(*backendRef.Namespace) != "" &&
		string(*backendRef.Namespace) != ownerNamespace {
		if !isCrossNamespaceReferencePermitted(
			crossNamespaceFrom{
				group:     egv1a1.GroupName,
				kind:      policyKind,
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/utils/ptr"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/envoyproxy/gateway/internal/gatewayapi/resource"
		SectionIndex: make(map[types.NamespacedName]sets.Set[string], gatewayMapSize),
	}

	policyCopies := securityPolicyCopiesWithStatusDeepCopy(securityPolicies)

	handledPolicies := make(map[types.NamespacedName]*egv1a1.SecurityPolicy, policyMapSize)

	// Map of attached Policy to Gateway. Used for policy merge process.
	// 4. Finally, the policies targeting Gateways

	// Build gateway policy maps, which are needed when processing the policies targeting xRoutes.
	t.buildGatewayPolicyMapForSecurity(securityPolicies, gateways, gatewayMap, gatewayPolicyMap, resources.ReferenceGrants)

	// Process the policies targeting RouteRules (HTTP + TCP)
	for i, currPolicy := range securityPolicies {
		policyName := utils.NamespacedName(currPolicy)
		// Only resolve TargetRefs from targetRefs field since TargetSelectors can't specify sectionName.
		targetRefs := resolvePolicyTargetsFromReferences(currPolicy.Spec.PolicyTargetReferences, currPolicy.Namespace)
		for _, currTarget := range targetRefs {
			if isRouteRule(currTarget) {
				policy, found := handledPolicies[policyName]
				if !found {
					policy = policyCopies[i]
					handledPolicies[policyName] = policy
					res = append(res, policy)
				}
		}
	}
	// Process the policies targeting xRoutes (HTTP + TCP)
	for i, currPolicy := range securityPolicies {
		policyName := utils.NamespacedName(currPolicy)
		targetRefs := resolvePolicyTargets(
			currPolicy.Spec.PolicyTargetReferences,
			routes,
			resources.ReferenceGrants,
			egv1a1.GroupName,
			egv1a1.KindSecurityPolicy,
			currPolicy.Namespace,
			t.GetNamespace)
		for _, currTarget := range targetRefs {
			if isRoute(currTarget) {
				policy, found := handledPolicies[policyName]
				if !found {
					policy = policyCopies[i]
					handledPolicies[policyName] = policy
					res = append(res, policy)
				}
		}
	}
	// Process the policies targeting Listeners
	for i, currPolicy := range securityPolicies {
		policyName := utils.NamespacedName(currPolicy)
		// Only resolve TargetRefs from targetRefs field since TargetSelectors can't specify sectionName.
		targetRefs := resolvePolicyTargetsFromReferences(currPolicy.Spec.PolicyTargetReferences, currPolicy.Namespace)
		for _, currTarget := range targetRefs {
			if isListener(currTarget) {
				policy, found := handledPolicies[policyName]
				if !found {
					policy = policyCopies[i]
					handledPolicies[policyName] = policy
					res = append(res, policy)
				}
		}
	}
	// Process the policies targeting Gateways
	for i, currPolicy := range securityPolicies {
		policyName := utils.NamespacedName(currPolicy)
		targetRefs := resolvePolicyTargets(
			currPolicy.Spec.PolicyTargetReferences,
			gateways,
			resources.ReferenceGrants,
			egv1a1.GroupName,
			egv1a1.KindSecurityPolicy,
			currPolicy.Namespace,
			t.GetNamespace)

		for _, currTarget := range targetRefs {
			if isGateway(currTarget) {
				policy, found := handledPolicies[policyName]
				if !found {
					policy = policyCopies[i]
					handledPolicies[policyName] = policy
					res = append(res, policy)
				}
	gateways []*GatewayContext,
	gatewayMap map[types.NamespacedName]*policyGatewayTargetContext,
	gatewayPolicyMap map[NamespacedNameWithSection]*egv1a1.SecurityPolicy,
	referenceGrants []*gwapiv1b1.ReferenceGrant,
) {
	for _, currPolicy := range securityPolicies {
		targetRefs := resolvePolicyTargets(
			currPolicy.Spec.PolicyTargetReferences,
			gateways,
			referenceGrants,
			egv1a1.GroupName,
			egv1a1.KindSecurityPolicy,
			currPolicy.Namespace,
			t.GetNamespace)
		for _, currTarget := range targetRefs {
			if currTarget.Kind == resource.KindGateway {
				// Check if the gateway exists
				key := types.NamespacedName{
					Name:      string(currTarget.Name),
					Namespace: string(currTarget.Namespace),
				}
				gateway, ok := gatewayMap[key]
				if !ok {
	gatewayPolicyMerged *GatewayPolicyRouteMap,
	gatewayPolicyMap map[NamespacedNameWithSection]*egv1a1.SecurityPolicy,
	policy *egv1a1.SecurityPolicy,
	currTarget policyTargetReferenceWithSectionName,
) {
	var (
		targetedRoute RouteContext
		resolveErr    *status.PolicyResolveError
	)

	targetedRoute, resolveErr = resolveSecurityPolicyRouteTargetRef(currTarget, routeMap)
	// Skip if the route is not found
	// It's not necessarily an error because the SecurityPolicy may be
	// reconciled by multiple controllers. And the other controller may
	// Check if merging is enabled
	if policy.Spec.MergeType == nil {
		// No merging - use existing translation logic
		if err := t.translateSecurityPolicyForRoute(policy, &securityPolicyOwners{}, targetedRoute, currTarget, resources, xdsIR, nil, nil); err != nil {
			status.SetTranslationErrorForPolicyAncestors(&policy.Status,
				ancestorRefs,
				t.GatewayControllerName,

				if gwPolicy == nil && listenerPolicy == nil {
					// No parent policy found, fall back to current policy
					if err := t.translateSecurityPolicyForRoute(policy, &securityPolicyOwners{}, targetedRoute, currTarget, resources, xdsIR, &gwNN, &listener.Name); err != nil {
						status.SetConditionForPolicyAncestor(&policy.Status,
							&ancestorRef,
							t.GatewayControllerName,
				}

				// Merge with parent policy
				mergedPolicy, owners, err := mergeSecurityPolicy(policy, parentPolicy)
				if err != nil {
					status.SetConditionForPolicyAncestor(&policy.Status,
						&ancestorRef,
				}

				// Apply merged policy
				if err := t.translateSecurityPolicyForRoute(mergedPolicy, owners, targetedRoute, currTarget, resources, xdsIR, &gwNN, &listener.Name); err != nil {
					status.SetConditionForPolicyAncestor(&policy.Status,
						&ancestorRef,
						t.GatewayControllerName,
	key := policyTargetRouteKey{
		Kind:      string(currTarget.Kind),
		Name:      string(currTarget.Name),
		Namespace: string(currTarget.Namespace),
	}
	overriddenTargetsMessage := getOverriddenTargetsMessageForRoute(routeMap[key], currTarget.SectionName)
	if overriddenTargetsMessage != "" {
	gatewayRouteMap *GatewayPolicyRouteMap,
	gatewayPolicyMergedMap *GatewayPolicyRouteMap,
	policy *egv1a1.SecurityPolicy,
	currTarget policyTargetReferenceWithSectionName,
) {
	var (
		targetedGateway *GatewayContext
		resolveErr      *status.PolicyResolveError
	)

	targetedGateway, resolveErr = resolveSecurityPolicyGatewayTargetRef(currTarget, gatewayMap)
	// Skip if the gateway is not found
	// It's not necessarily an error because the SecurityPolicy may be
	// reconciled by multiple controllers. And the other controller may
}

func resolveSecurityPolicyGatewayTargetRef(
	target policyTargetReferenceWithSectionName,
	gateways map[types.NamespacedName]*policyGatewayTargetContext,
) (*GatewayContext, *status.PolicyResolveError) {
	// Find the Gateway
	key := types.NamespacedName{
		Name:      string(target.Name),
		Namespace: string(target.Namespace),
	}
	gateway, ok := gateways[key]

}

func resolveSecurityPolicyRouteTargetRef(
	target policyTargetReferenceWithSectionName,
	routes map[policyTargetRouteKey]*policyRouteTargetContext,
) (RouteContext, *status.PolicyResolveError) {
	// Check if the route exists
	key := policyTargetRouteKey{
		Kind:      string(target.Kind),
		Name:      string(target.Name),
		Namespace: string(target.Namespace),
	}
	route, ok := routes[key]


func (t *Translator) translateSecurityPolicyForRoute(
	policy *egv1a1.SecurityPolicy,
	owners *securityPolicyOwners,
	route RouteContext,
	target policyTargetReferenceWithSectionName,
	resources *resource.Resources,
	xdsIR resource.XdsIRMap,
	policyTargetGateway *types.NamespacedName,
	if policy.Spec.BasicAuth != nil {
		if basicAuth, err = t.buildBasicAuth(
			policy,
			owners,
			resources,
		); err != nil {
			err = perr.WithMessage(err, "BasicAuth")
			errs = errors.Join(errs, err)
		}
	if policy.Spec.APIKeyAuth != nil {
		if apiKeyAuth, err = t.buildAPIKeyAuth(
			policy,
			owners,
			resources,
		); err != nil {
			err = perr.WithMessage(err, "APIKeyAuth")
			errs = errors.Join(errs, err)
		}
	}

	if policy.Spec.Authorization != nil {
		if authorization, err = t.buildAuthorization(policy, owners); err != nil {
			err = perr.WithMessage(err, "Authorization")
			errs = errors.Join(errs, err)
		}
		if policy.Spec.ExtAuth != nil {
			if extAuth, extAuthErr = t.buildExtAuth(
				policy,
				owners,
				resources,
				gtwCtx,
			); extAuthErr != nil {
				extAuthErr = perr.WithMessage(extAuthErr, "ExtAuth")
				errs = errors.Join(errs, extAuthErr)
			}
		if policy.Spec.OIDC != nil {
			if oidc, err = t.buildOIDC(
				policy,
				owners,
				resources,
				gtwCtx,
			); err != nil {
				err = perr.WithMessage(err, "OIDC")
				errs = errors.Join(errs, err)
				hasNonExtAuthError = true
		if policy.Spec.JWT != nil {
			if jwt, err = t.buildJWT(
				policy,
				owners,
				resources,
				gtwCtx,
			); err != nil {
				err = perr.WithMessage(err, "JWT")
				errs = errors.Join(errs, err)
				hasNonExtAuthError = true
						continue
					}
					// Only authorization for TCP
					if authorization != nil {
						authCopy := *authorization
						r.Authorization = &authCopy
					}
				}
			}
		case resource.KindHTTPRoute, resource.KindGRPCRoute:
func (t *Translator) translateSecurityPolicyForGateway(
	policy *egv1a1.SecurityPolicy,
	gtwCtx *GatewayContext,
	target policyTargetReferenceWithSectionName,
	resources *resource.Resources,
	xdsIR resource.XdsIRMap,
) error {
	// Build IR
	noOwners := &securityPolicyOwners{}
	var (
		cors                  *ir.CORS
		jwt                   *ir.JWT
	if policy.Spec.JWT != nil {
		if jwt, err = t.buildJWT(
			policy,
			noOwners,
			resources,
			gtwCtx,
		); err != nil {
			err = perr.WithMessage(err, "JWT")
			errs = errors.Join(errs, err)
		}
	if policy.Spec.OIDC != nil {
		if oidc, err = t.buildOIDC(
			policy,
			noOwners,
			resources,
			gtwCtx,
		); err != nil {
			err = perr.WithMessage(err, "OIDC")
			errs = errors.Join(errs, err)
		}
	if policy.Spec.BasicAuth != nil {
		if basicAuth, err = t.buildBasicAuth(
			policy,
			noOwners,
			resources,
		); err != nil {
			err = perr.WithMessage(err, "BasicAuth")
			errs = errors.Join(errs, err)
		}
	if policy.Spec.APIKeyAuth != nil {
		if apiKeyAuth, err = t.buildAPIKeyAuth(
			policy,
			noOwners,
			resources,
		); err != nil {
			err = perr.WithMessage(err, "APIKeyAuth")
			errs = errors.Join(errs, err)
		}
	}

	if policy.Spec.Authorization != nil {
		if authorization, err = t.buildAuthorization(policy, noOwners); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	if policy.Spec.ExtAuth != nil {
		if extAuth, extAuthErr = t.buildExtAuth(
			policy,
			noOwners,
			resources,
			gtwCtx,
		); extAuthErr != nil {
			extAuthErr = perr.WithMessage(extAuthErr, "ExtAuth")
			errs = errors.Join(errs, extAuthErr)
		}

func (t *Translator) buildJWT(
	policy *egv1a1.SecurityPolicy,
	owners *securityPolicyOwners,
	resources *resource.Resources,
	gtwCtx *GatewayContext,
) (*ir.JWT, error) {
		return nil, err
	}

	jwtOwnerPolicy := policyOwnerOr(owners.jwtProviders, policy)
	providers := make([]ir.JWTProvider, 0, len(policy.Spec.JWT.Providers))
	for i, p := range policy.Spec.JWT.Providers {
		provider := ir.JWTProvider{
			ExtractFrom:    p.ExtractFrom,
		}
		if p.RemoteJWKS != nil {
			remoteJWKS, err := t.buildRemoteJWKS(jwtOwnerPolicy, p.RemoteJWKS, i, resources, gtwCtx)
			if err != nil {
				return nil, err
			}
			provider.RemoteJWKS = remoteJWKS
		} else {
			localJWKS, err := t.buildLocalJWKS(jwtOwnerPolicy, p.LocalJWKS)
			if err != nil {
				return nil, err
			}

func (t *Translator) buildOIDC(
	policy *egv1a1.SecurityPolicy,
	owners *securityPolicyOwners,
	resources *resource.Resources,
	gtwCtx *GatewayContext,
) (*ir.OIDC, error) {
		err                    error
	)

	if provider, err = t.buildOIDCProvider(policy, owners, resources, gtwCtx); err != nil {
		return nil, err
	}

	// Client ID can be specified either as a string or as a reference to a secret.
	switch {
	case oidc.ClientID != nil:
		clientID = *oidc.ClientID
	case oidc.ClientIDRef != nil:
		ownerPolicy := policyOwnerOr(owners.oidcClientIDRef, policy)
		from := crossNamespaceFrom{
			group:     egv1a1.GroupName,
			kind:      resource.KindSecurityPolicy,
			namespace: ownerPolicy.Namespace,
		}

		var clientIDSecret *corev1.Secret
		if clientIDSecret, err = t.validateSecretRef(true, from, *oidc.ClientIDRef, resources); err != nil {
			return nil, err
		return nil, fmt.Errorf("client ID must be specified in OIDC policy %s/%s", policy.Namespace, policy.Name)
	}

	clientSecretOwner := policyOwnerOr(owners.oidcClientSecret, policy)
	from := crossNamespaceFrom{
		group:     egv1a1.GroupName,
		kind:      resource.KindSecurityPolicy,
		namespace: clientSecretOwner.Namespace,
	}
	if clientSecret, err = t.validateSecretRef(true, from, oidc.ClientSecret, resources); err != nil {
		return nil, err
	}
		disableTokenEncryption = *oidc.DisableTokenEncryption
	}

	oidcOwner := policyOwnerOr(owners.oidc, policy)

	// Generate a unique cookie suffix for oauth filters.
	// This is to avoid cookie name collision when multiple security policies are applied
	// to the same route.
	suffix := utils.Digest32(string(oidcOwner.UID))

	// Get the HMAC secret.
	// HMAC secret is generated by the CertGen job and stored in a secret
	}

	irOIDC := &ir.OIDC{
		Name:                   irConfigName(oidcOwner),
		Provider:               *provider,
		ClientID:               clientID,
		ClientSecret:           clientSecretBytes,

func (t *Translator) buildOIDCProvider(
	policy *egv1a1.SecurityPolicy,
	owners *securityPolicyOwners,
	resources *resource.Resources,
	gtwCtx *GatewayContext,
) (*ir.OIDCProvider, error) {
		protocol = ir.HTTP
	}

	oidcProviderOwner := policyOwnerOr(owners.oidcProviderBackendRefs, policy)
	if len(provider.BackendRefs) > 0 {
		if rd, err = t.translateExtServiceBackendRefs(
			oidcProviderOwner, provider.BackendRefs, protocol, resources, gtwCtx, "oidc", 0); err != nil {
			return nil, err
		}
	}

func (t *Translator) buildAPIKeyAuth(
	policy *egv1a1.SecurityPolicy,
	owners *securityPolicyOwners,
	resources *resource.Resources,
) (*ir.APIKeyAuth, error) {
	ownerPolicy := policyOwnerOr(owners.apiKeyAuthCredentialRefs, policy)
	from := crossNamespaceFrom{
		group:     egv1a1.GroupName,
		kind:      resource.KindSecurityPolicy,
		namespace: ownerPolicy.Namespace,
	}

	expected := len(policy.Spec.APIKeyAuth.CredentialRefs)

func (t *Translator) buildBasicAuth(
	policy *egv1a1.SecurityPolicy,
	owners *securityPolicyOwners,
	resources *resource.Resources,
) (*ir.BasicAuth, error) {
	var (
		err         error
	)

	ownerPolicy := policyOwnerOr(owners.basicAuth, policy)
	from := crossNamespaceFrom{
		group:     egv1a1.GroupName,
		kind:      resource.KindSecurityPolicy,
		namespace: ownerPolicy.Namespace,
	}
	if usersSecret, err = t.validateSecretRef(true, from, basicAuth.Users, resources); err != nil {
		return nil, err
	}

	return &ir.BasicAuth{
		Name:                  irConfigName(ownerPolicy),
		Users:                 usersSecretBytes,
		ForwardUsernameHeader: basicAuth.ForwardUsernameHeader,
	}, nil

func (t *Translator) buildExtAuth(
	policy *egv1a1.SecurityPolicy,
	owners *securityPolicyOwners,
	resources *resource.Resources,
	gtwCtx *GatewayContext,
) (*ir.ExtAuth, error) {
		contextExtensions []*ir.ContextExtention
	)

	backendRefsOwnerPolicy := policyOwnerOr(owners.extAuthBackendRefs, policy)

	// These are sanity checks, they should never happen because the API server
	// should have caught them
	if http == nil && grpc == nil {
	}

	if rd, err = t.translateExtServiceBackendRefs(
		backendRefsOwnerPolicy, backendRefs, protocol, resources, gtwCtx, "extauth", 0); err != nil {
		return nil, err
	}

		// When translated to XDS, the authority is used on the filter level not on the cluster level.
		// There's no way to translate to XDS and use a different authority for each backendref
		if authority == "" {
			authority = t.backendRefAuthority(&backendRef.BackendObjectReference, backendRefsOwnerPolicy)
		}
	}

		return nil, err
	}

	if contextExtensions, err = t.buildContextExtensions(policy.Spec.ExtAuth.ContextExtensions, owners, policy); err != nil {
		return nil, err
	}

	extAuthOwner := policyOwnerOr(owners.extAuth, policy)
	extAuth := &ir.ExtAuth{
		Name:                 irConfigName(extAuthOwner),
		HeadersToExtAuth:     policy.Spec.ExtAuth.HeadersToExtAuth,
		ContextExtensions:    contextExtensions,
		FailOpen:             policy.Spec.ExtAuth.FailOpen,
			Destination:      *rd,
			Authority:        authority,
			Path:             ptr.Deref(http.Path, ""),
			PathOverride:     ptr.Deref(http.PathOverride, ""),
			HeadersToBackend: http.HeadersToBackend,
		}
	} else {

func (t *Translator) buildContextExtensions(
	contextExtensions []*egv1a1.ContextExtension,
	owners *securityPolicyOwners,
	defaultOwner *egv1a1.SecurityPolicy,
) ([]*ir.ContextExtention, error) {
	if len(contextExtensions) == 0 {
		return nil, nil
	for _, ext := range contextExtensions {
		var value ir.PrivateBytes
		if ext.Type == egv1a1.ContextExtensionValueTypeValueRef {
			ownerPolicy := policyOwnerOr(owners.extAuthContextExtensions[ext.Name], defaultOwner)
			var err error
			if value, err = t.getContextExtensionValueFromRef(ext.ValueRef, ownerPolicy.Namespace); err != nil {
				return nil, err
			}
		} else if ext.Value != nil {
	return fmt.Sprintf("%s.%s", backendRef.Name, backendNamespace)
}

func (t *Translator) buildAuthorization(
	policy *egv1a1.SecurityPolicy,
	owners *securityPolicyOwners,
) (*ir.Authorization, error) {
	var (
		authorization = policy.Spec.Authorization
		irAuth        = &ir.Authorization{}
		defaultAction = egv1a1.AuthorizationActionDeny
	)

	ownerPolicy := policyOwnerOr(owners.authorizationRules, policy)

	if authorization.DefaultAction != nil {
		defaultAction = *authorization.DefaultAction
	}
		if rule.Name != nil && *rule.Name != "" {
			name = *rule.Name
		} else {
			name = defaultAuthorizationRuleName(ownerPolicy, i)
		}
		irAuth.Rules = append(irAuth.Rules, &ir.AuthorizationRule{
			Name:      name,
		strconv.Itoa(index))
}

type securityPolicyOwners struct {
	basicAuth                *egv1a1.SecurityPolicy
	apiKeyAuthCredentialRefs *egv1a1.SecurityPolicy
	authorizationRules       *egv1a1.SecurityPolicy
	extAuth                  *egv1a1.SecurityPolicy
	extAuthBackendRefs       *egv1a1.SecurityPolicy
	extAuthContextExtensions map[string]*egv1a1.SecurityPolicy
	oidc                     *egv1a1.SecurityPolicy
	oidcProviderBackendRefs  *egv1a1.SecurityPolicy
	oidcClientIDRef          *egv1a1.SecurityPolicy
	oidcClientSecret         *egv1a1.SecurityPolicy
	jwtProviders             *egv1a1.SecurityPolicy
}

// policyOwnerOr returns owner if non-nil, otherwise fallback.
// Used to resolve per-field owners from securityPolicyOwners: the owner is the policy
// that contributed the field (route overrides parent), falling back to the active policy
// when no merge occurred or the field was not set by either side.
func policyOwnerOr(owner, fallback *egv1a1.SecurityPolicy) *egv1a1.SecurityPolicy {
	if owner != nil {
		return owner
	}
	return fallback
}

// mergeSecurityPolicy merges a route-level SecurityPolicy with a parent (Gateway/Listener) SecurityPolicy.
func mergeSecurityPolicy(routePolicy, parentPolicy *egv1a1.SecurityPolicy) (*egv1a1.SecurityPolicy, *securityPolicyOwners, error) {
	if routePolicy.Spec.MergeType == nil || parentPolicy == nil {
		return routePolicy, nil, nil
	}
	mergedPolicy, err := utils.Merge[*egv1a1.SecurityPolicy](parentPolicy, routePolicy, *routePolicy.Spec.MergeType)
	if err != nil {
		return nil, nil, err
	}
	return mergedPolicy, buildSecurityPolicyOwners(routePolicy, parentPolicy), nil
}

// ownerOf returns route if routeOwns(route) is true, otherwise parent.
// Use this when ownership of a merged field is determined by a single predicate.
func ownerOf(
	route, parent *egv1a1.SecurityPolicy,
	routeOwns func(*egv1a1.SecurityPolicy) bool,
) *egv1a1.SecurityPolicy {
	if routeOwns(route) {
		return route
	}
	return parent
}

// buildSecurityPolicyOwners determines, for each merged field, which policy
// (route or parent) is considered the owner. The owner is used later to resolve
// references (e.g. Secrets, BackendRefs) scoped to the owning policy's namespace,
// and to derive IR resource names tied to the owning policy.
func buildSecurityPolicyOwners(route, parent *egv1a1.SecurityPolicy) *securityPolicyOwners {
	return &securityPolicyOwners{
		basicAuth: ownerOf(route, parent, func(p *egv1a1.SecurityPolicy) bool {
			return p.Spec.BasicAuth != nil
		}),
		apiKeyAuthCredentialRefs: ownerOf(route, parent, func(p *egv1a1.SecurityPolicy) bool {
			return p.Spec.APIKeyAuth != nil && len(p.Spec.APIKeyAuth.CredentialRefs) > 0
		}),
		authorizationRules: ownerOf(route, parent, func(p *egv1a1.SecurityPolicy) bool {
			return p.Spec.Authorization != nil && len(p.Spec.Authorization.Rules) > 0
		}),
		extAuth: ownerOf(route, parent, func(p *egv1a1.SecurityPolicy) bool {
			return p.Spec.ExtAuth != nil
		}),
		extAuthBackendRefs: ownerOf(route, parent, func(p *egv1a1.SecurityPolicy) bool {
			ea := p.Spec.ExtAuth
			if ea == nil {
				return false
			}
			if ea.HTTP != nil && (len(ea.HTTP.BackendRefs) > 0 || ea.HTTP.BackendRef != nil) {
				return true
			}
			if ea.GRPC != nil && (len(ea.GRPC.BackendRefs) > 0 || ea.GRPC.BackendRef != nil) {
				return true
			}
			return false
		}),
		extAuthContextExtensions: buildExtAuthContextExtensionOwners(route, parent),
		oidc: ownerOf(route, parent, func(p *egv1a1.SecurityPolicy) bool {
			return p.Spec.OIDC != nil
		}),
		oidcProviderBackendRefs: ownerOf(route, parent, func(p *egv1a1.SecurityPolicy) bool {
			return p.Spec.OIDC != nil && len(p.Spec.OIDC.Provider.BackendRefs) > 0
		}),
		oidcClientIDRef: ownerOf(route, parent, func(p *egv1a1.SecurityPolicy) bool {
			return p.Spec.OIDC != nil && p.Spec.OIDC.ClientIDRef != nil
		}),
		oidcClientSecret: ownerOf(route, parent, func(p *egv1a1.SecurityPolicy) bool {
			return p.Spec.OIDC != nil
		}),
		jwtProviders: ownerOf(route, parent, func(p *egv1a1.SecurityPolicy) bool {
			return p.Spec.JWT != nil && len(p.Spec.JWT.Providers) > 0
		}),
	}
}

// buildExtAuthContextExtensionOwners returns a per-key owner map for ExtAuth ContextExtensions.
// Parent keys are added first so that route-level extensions take precedence on conflict.
func buildExtAuthContextExtensionOwners(route, parent *egv1a1.SecurityPolicy) map[string]*egv1a1.SecurityPolicy {
	owners := make(map[string]*egv1a1.SecurityPolicy)
	if parent.Spec.ExtAuth != nil {
		for _, ext := range parent.Spec.ExtAuth.ContextExtensions {
			owners[ext.Name] = parent
		}
	}
	if route.Spec.ExtAuth != nil {
		for _, ext := range route.Spec.ExtAuth.ContextExtensions {
			owners[ext.Name] = route
		}
	}
	return owners
}

// securityPolicyCopiesWithStatusDeepCopy returns shallow copies with deep-copied Status fields.
// Status is mutated during translation and shares a pointer with the watchable coalesce goroutine.
func securityPolicyCopiesWithStatusDeepCopy(policies []*egv1a1.SecurityPolicy) []*egv1a1.SecurityPolicy {
	copies := make([]*egv1a1.SecurityPolicy, len(policies))
	for i, p := range policies {
		out := *p
		p.Status.DeepCopyInto(&out.Status)
		copies[i] = &out
	}
	return copies
}
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	"github.com/envoyproxy/gateway/internal/gatewayapi/resource"
}

type targetRefWithTimestamp struct {
	policyTargetReferenceWithSectionName
	CreationTimestamp metav1.Time
}

// policyTargetReferenceWithSectionName extends the Gateway API's LocalPolicyTargetReference to include a Namespace field.
// This is necessary because policies may reference targets in other namespaces.
type policyTargetReferenceWithSectionName struct {
	// Group is the group of the target resource.
	// +required
	Group gwapiv1.Group `json:"group"`

	// Kind is kind of the target resource.
	// +required
	Kind gwapiv1.Kind `json:"kind"`

	// Name is the name of the target resource.
	// +required
	Name gwapiv1.ObjectName `json:"name"`

	// Namespace is the namespace of the target resource. When unspecified, it is assumed to be in the same namespace as the policy.
	Namespace gwapiv1.Namespace `json:"namespace"`

	// SectionName is the name of a section within the target resource. When
	// unspecified, this targetRef targets the entire resource. In the following
	// resources, SectionName is interpreted as the following:
	//
	// * Gateway: Listener name
	// * HTTPRoute: HTTPRouteRule name
	// * Service: Port name
	//
	// If a SectionName is specified, but does not exist on the targeted object,
	// the Policy must fail to attach, and the policy implementation should record
	// a `ResolvedRefs` or similar Condition in the Policy's status.
	//
	// +optional
	SectionName *gwapiv1.SectionName `json:"sectionName,omitempty"`
}

func isRouteRule(target policyTargetReferenceWithSectionName) bool {
	// If the target is not a gateway and the section name is not nil, then it's a route rule.
	return target.Kind != resource.KindGateway && target.SectionName != nil
}

func isRoute(target policyTargetReferenceWithSectionName) bool {
	// If the target is not a gateway and the section name is nil, then it's a route.
	return target.Kind != resource.KindGateway && target.SectionName == nil
}

func isGateway(target policyTargetReferenceWithSectionName) bool {
	// If the target is a gateway and the section name is nil, then it's a gateway.
	return target.Kind == resource.KindGateway && target.SectionName == nil
}

func isListener(target policyTargetReferenceWithSectionName) bool {
	// If the target is a gateway and the section name is not nil, then it's a listener.
	return target.Kind == resource.KindGateway && target.SectionName != nil
}

func selectorFromTargetSelector(selector egv1a1.TargetSelector) labels.Selector {
	l, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels:      selector.MatchLabels,
	return l
}

func targetNamespaceLabelSelector(namespaces *egv1a1.TargetSelectorNamespaces) labels.Selector {
	if namespaces == nil || namespaces.Selector == nil {
		return labels.Nothing()
	}

	selector, err := metav1.LabelSelectorAsSelector(namespaces.Selector)
	if err != nil {
		return labels.Nothing()
	}

	return selector
}

func targetSelectorNamespacesMatch(
	namespaces *egv1a1.TargetSelectorNamespaces,
	policyNamespace,
	targetNamespace string,
	targetNamespaceLabels map[string]string,
) bool {
	if namespaces == nil {
		return targetNamespace == policyNamespace
	}

	switch namespaces.From {
	case "", egv1a1.TargetNamespaceFromSame:
		return targetNamespace == policyNamespace
	case egv1a1.TargetNamespaceFromAll:
		return true
	case egv1a1.TargetNamespaceFromSelector:
		if targetNamespaceLabels == nil {
			return false
		}
		return targetNamespaceLabelSelector(namespaces).Matches(labels.Set(targetNamespaceLabels))
	default:
		return false
	}
}

func targetNamespaceMatches(
	selector egv1a1.TargetSelector,
	policyNamespace,
	targetNamespace string,
	namespaceLookup func(string) *corev1.Namespace,
) bool {
	var targetNamespaceLabels map[string]string
	if namespaceLookup != nil {
		if ns := namespaceLookup(targetNamespace); ns != nil {
			targetNamespaceLabels = ns.GetLabels()
		}
	}

	return targetSelectorNamespacesMatch(selector.Namespaces, policyNamespace, targetNamespace, targetNamespaceLabels)
}

// isCrossNamespaceReferencePermitted checks if a cross-namespace reference from a policy in one namespace to a target
// in another namespace is allowed by a ReferenceGrant in the target namespace.
func isCrossNamespaceReferencePermitted(
	from crossNamespaceFrom,
	to crossNamespaceTo,
	referenceGrants []*gwapiv1b1.ReferenceGrant,
) bool {
	if from.namespace == to.namespace {
		return true
	}

	for _, referenceGrant := range referenceGrants {
		if referenceGrant.Namespace != to.namespace {
			continue
		}

		var fromAllowed bool
		for _, refGrantFrom := range referenceGrant.Spec.From {
			if string(refGrantFrom.Namespace) == from.namespace &&
				string(refGrantFrom.Group) == from.group &&
				string(refGrantFrom.Kind) == from.kind {
				fromAllowed = true
				break
			}
		}
		if !fromAllowed {
			continue
		}

		for _, refGrantTo := range referenceGrant.Spec.To {
			if string(refGrantTo.Group) == to.group &&
				string(refGrantTo.Kind) == to.kind &&
				(refGrantTo.Name == nil || *refGrantTo.Name == "" || string(*refGrantTo.Name) == to.name) {
				return true
			}
		}
	}

	return false
}

// resolvePolicyTargetsFromSelectors returns policy target refs allowed by the policy's TargetSelectors.
func resolvePolicyTargetsFromSelectors[T client.Object](
	targetSelectors []egv1a1.TargetSelector,
	potentialTargets []T,
	referenceGrants []*gwapiv1b1.ReferenceGrant,
	policyGroup string,
	policyKind string,
	policyNamespace string,
	namespaceLookup func(string) *corev1.Namespace,
) []targetRefWithTimestamp {
	allowedDedup := sets.New[targetRefWithTimestamp]()
	targetRefs := make([]targetRefWithTimestamp, 0)
	for _, currSelector := range targetSelectors {
		labelSelector := selectorFromTargetSelector(currSelector)
		for _, obj := range potentialTargets {
			gvk := obj.GetObjectKind().GroupVersionKind()
			if gvk.Kind != string(currSelector.Kind) ||
				gvk.Group != string(ptr.Deref(currSelector.Group, gwapiv1.GroupName)) {
				continue
			}

			// Check if the target object's namespace matches the selector's namespace criteria.
			if !targetNamespaceMatches(currSelector, policyNamespace, obj.GetNamespace(), namespaceLookup) {
				continue
			}

			ref := policyTargetReferenceWithSectionName{
				Group:     gwapiv1.Group(gvk.Group),
				Kind:      gwapiv1.Kind(gvk.Kind),
				Name:      gwapiv1.ObjectName(obj.GetName()),
				Namespace: gwapiv1.Namespace(obj.GetNamespace()),
			}

			// Check if the target object's labels match the selector's label criteria.
			if !labelSelector.Matches(labels.Set(obj.GetLabels())) {
				continue
			}

			// Check if cross-namespace reference is allowed if the policy and target are in different namespaces.
			if !isCrossNamespaceReferencePermitted(
				crossNamespaceFrom{
					group:     policyGroup,
					kind:      policyKind,
					namespace: policyNamespace,
				},
				crossNamespaceTo{
					group:     gvk.Group,
					kind:      gvk.Kind,
					namespace: obj.GetNamespace(),
					name:      obj.GetName(),
				},
				referenceGrants,
			) {
				continue
			}

			targetRef := targetRefWithTimestamp{
				CreationTimestamp:                    obj.GetCreationTimestamp(),
				policyTargetReferenceWithSectionName: ref,
			}
			if allowedDedup.Has(targetRef) {
				continue
			}
			allowedDedup.Insert(targetRef)
			targetRefs = append(targetRefs, targetRef)
		}
	}

	return targetRefs
}

// resolvePolicyTargetsFromReferences returns a list of policy target refs specified in the policy's TargetRefs, with the namespace field populated.
func resolvePolicyTargetsFromReferences(
	targetRefs egv1a1.PolicyTargetReferences,
	policyNamespace string,
) []policyTargetReferenceWithSectionName {
	refs := targetRefs.GetTargetRefs()
	ret := make([]policyTargetReferenceWithSectionName, 0, len(refs))
	var emptyTargetRef gwapiv1.LocalPolicyTargetReferenceWithSectionName
	for _, v := range refs {
		if v == emptyTargetRef {
			// This can happen when the targetRef structure is read from extension server policies
			continue
		}
		ret = append(ret, policyTargetReferenceWithSectionName{
			Group:       v.Group,
			Kind:        v.Kind,
			Name:        v.Name,
			Namespace:   gwapiv1.Namespace(policyNamespace),
			SectionName: v.SectionName,
		})
	}

	return ret
}

// composePolicyTargetRefs combines the allowed target refs derived from the selectors and the plain target refs specified in the policy.
func composePolicyTargetRefs(
	selectorTargetRefs []targetRefWithTimestamp,
	plainTargetRefs []policyTargetReferenceWithSectionName,
) []policyTargetReferenceWithSectionName {
	// First add the target refs derived from the selectors, sorted by the creation timestamp of the matched objects.
	slices.SortFunc(selectorTargetRefs, func(i, j targetRefWithTimestamp) int {
		return i.CreationTimestamp.Compare(j.CreationTimestamp.Time)
	})
	ret := make([]policyTargetReferenceWithSectionName, len(selectorTargetRefs))
	for i, v := range selectorTargetRefs {
		ret[i] = v.policyTargetReferenceWithSectionName
	}

	// Plain targetRefs in the policy don't have an associated creation timestamp, but can still refer
	// to targets that were already found via the selectors. Only add them to the returned list if
	// they are not yet there. Always add them at the end.
	fastLookup := sets.New(ret...)
	for _, targetRef := range plainTargetRefs {
		if !fastLookup.Has(targetRef) {
			ret = append(ret, targetRef)
		}
	}

	return ret
}

// resolvePolicyTargets returns a list of policy target refs that are allowed by the policy's TargetSelectors.
// The list includes both target refs derived from the selectors and plain target refs specified in the policy.
func resolvePolicyTargets[T client.Object](
	targetRefs egv1a1.PolicyTargetReferences,
	potentialTargets []T,
	referenceGrants []*gwapiv1b1.ReferenceGrant,
	policyGroup string,
	policyKind string,
	policyNamespace string,
	namespaceLookup func(string) *corev1.Namespace,
) []policyTargetReferenceWithSectionName {
	selectorTargetRefs := resolvePolicyTargetsFromSelectors(
		targetRefs.TargetSelectors,
		potentialTargets,
		referenceGrants,
		policyGroup,
		policyKind,
		policyNamespace,
		namespaceLookup)
	plainTargetRefs := resolvePolicyTargetsFromReferences(targetRefs, policyNamespace)
	return composePolicyTargetRefs(selectorTargetRefs, plainTargetRefs)
}

// legacy function to get policy target refs without considering cross-namespace policy attachment.
// This is only used for extension server policies.
// TODO: add cross-namesapce policy attachment to extension server if needed, and remove this function.
func getPolicyTargetRefs[T client.Object](policy egv1a1.PolicyTargetReferences, potentialTargets []T, policyNamespace string) []gwapiv1.LocalPolicyTargetReferenceWithSectionName {
	dedup := sets.New[targetRefWithTimestamp]()
	for _, currSelector := range policy.TargetSelectors {
			if labelSelector.Matches(labels.Set(obj.GetLabels())) {
				dedup.Insert(targetRefWithTimestamp{
					CreationTimestamp: obj.GetCreationTimestamp(),
					policyTargetReferenceWithSectionName: policyTargetReferenceWithSectionName{
						Group: gwapiv1.Group(gvk.Group),
						Kind:  gwapiv1.Kind(gvk.Kind),
						Name:  gwapiv1.ObjectName(obj.GetName()),
					},
				})
			}
	})
	ret := make([]gwapiv1.LocalPolicyTargetReferenceWithSectionName, len(selectorsList))
	for i, v := range selectorsList {
		ret[i] = gwapiv1.LocalPolicyTargetReferenceWithSectionName{
			LocalPolicyTargetReference: gwapiv1.LocalPolicyTargetReference{
				Group: v.Group,
				Kind:  v.Kind,
				Name:  v.Name,
			},
			SectionName: v.SectionName,
		}
	}
	// Plain targetRefs in the policy don't have an associated creation timestamp, but can still refer
	// to targets that were already found via the selectors. Only add them to the returned list if
