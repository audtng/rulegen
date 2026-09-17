package main

	Namespaces      Context = "namespaces"
	Quotas          Context = "quotas"
	Recommendations Context = "recommendations"
	All             Context = "all"
	Plugins         Context = "plugins"
	Volumes         Context = "volumes"
	WriteRequest
}

// JobSummaryRequest is used when we just need to get a specific job summary
type JobSummaryRequest struct {
	JobID string
	nn := new(Node)
	*nn = *n
	nn.Attributes = helper.CopyMapStringString(nn.Attributes)
	nn.Resources = nn.Resources.Copy()
	nn.Reserved = nn.Reserved.Copy()
	nn.NodeResources = nn.NodeResources.Copy()
	nn.ReservedResources = nn.ReservedResources.Copy()
	nn.Links = helper.CopyMapStringString(nn.Links)
	nn.Meta = helper.CopyMapStringString(nn.Meta)
	nn.Events = copyNodeEvents(n.Events)
	nn.DrainStrategy = nn.DrainStrategy.Copy()
	nn.CSIControllerPlugins = copyNodeCSI(nn.CSIControllerPlugins)
	nn.CSINodePlugins = copyNodeCSI(nn.CSINodePlugins)
	nn.Drivers = copyNodeDrivers(n.Drivers)
	nn.HostVolumes = copyNodeHostVolumes(n.HostVolumes)
	return nn
}
// Add adds the resources of the delta to this, potentially
// returning an error if not possible.
// COMPAT(0.10): Remove in 0.10
func (r *Resources) Add(delta *Resources) error {
	if delta == nil {
		return nil
	}
	r.CPU += delta.CPU
	r.MemoryMB += delta.MemoryMB
	r.DiskMB += delta.DiskMB
			r.Networks[idx].Add(n)
		}
	}
	return nil
}

// COMPAT(0.10): Remove in 0.10
	return reflect.DeepEqual(n, o)
}

func (n *NodeNetworkResource) HasAlias(alias string) bool {
	for _, addr := range n.Addresses {
		if addr.Alias == alias {
	}
	newR := new(NetworkResource)
	*newR = *n
	if n.ReservedPorts != nil {
		newR.ReservedPorts = make([]Port, len(n.ReservedPorts))
		copy(newR.ReservedPorts, n.ReservedPorts)
	return labelValues
}

// ConnectPort returns the Connect port for the given service. Returns false if
// no port was found for a service with that name.
func (n *NetworkResource) PortForService(serviceName string) (Port, bool) {
	label := fmt.Sprintf("%s-%s", ConnectProxyPrefix, serviceName)
	for _, port := range n.ReservedPorts {
		if port.Label == label {
			return port, true
		}
	}
	for _, port := range n.DynamicPorts {
		if port.Label == label {
			return port, true
		}
	}

	return Port{}, false
}

// Networks defined for a task on the Resources struct.
type Networks []*NetworkResource

}

// Port assignment and IP for the given label or empty values.
func (ns Networks) Port(label string) (string, int) {
	for _, n := range ns {
		for _, p := range n.ReservedPorts {
			if p.Label == label {
				return n.IP, p.Value
			}
		}
		for _, p := range n.DynamicPorts {
			if p.Label == label {
				return n.IP, p.Value
			}
		}
	}
	return "", 0
}

func (ns Networks) NetIndex(n *NetworkResource) int {

	var mErr multierror.Error
	if r.Name == "" {
		multierror.Append(&mErr, errors.New("device name must be given as one of the following: type, vendor/type, or vendor/type/name"))
	}

	for idx, constr := range r.Constraints {
		switch constr.Operand {
		case ConstraintDistinctHosts, ConstraintDistinctProperty:
			outer := fmt.Errorf("Constraint %d validation failed: using unsupported operand %q", idx+1, constr.Operand)
			multierror.Append(&mErr, outer)
		default:
			if err := constr.Validate(); err != nil {
				outer := fmt.Errorf("Constraint %d validation failed: %s", idx+1, err)
				multierror.Append(&mErr, outer)
			}
		}
	}
	for idx, affinity := range r.Affinities {
		if err := affinity.Validate(); err != nil {
			outer := fmt.Errorf("Affinity %d validation failed: %s", idx+1, err)
			multierror.Append(&mErr, outer)
		}
	}


	newN := new(NodeResources)
	*newN = *n

	// Copy the networks
	newN.Networks = n.Networks.Copy()

	// Copy the devices
	if n.Devices != nil {
		devices := len(n.Devices)
	CpuShares int64
}

func (n *NodeCpuResources) Merge(o *NodeCpuResources) {
	if o == nil {
		return
			},
		},
	}
	if len(a.Networks) > 0 {
		for _, net := range a.Networks {
			ret.Flattened.Networks = append(ret.Flattened.Networks, net)
		}
	}
	return ret
}

}

// NamespacedID returns the namespaced id useful for logging
func (j *Job) NamespacedID() *NamespacedID {
	return &NamespacedID{
		ID:        j.ID,
		Namespace: j.Namespace,
	}
}

// Canonicalize is used to canonicalize fields in the Job. This should be called
// when registering a Job. A set of warnings are returned if the job was changed
// in anyway that the user should be made aware of.
func (j *Job) Canonicalize() (warnings error) {
	if j == nil {
		return nil
	}

	var mErr multierror.Error
	// Ensure that an empty and nil map are treated the same to avoid scheduling
	// problems since we use reflect DeepEquals.
	if len(j.Meta) == 0 {
	if j.Periodic != nil {
		j.Periodic.Canonicalize()
	}

	return mErr.ErrorOrNil()
}

// Copy returns a deep copy of the Job. It is expected that callers use recover.
	var mErr multierror.Error

	// Check the groups
	ap := 0
	for _, tg := range j.TaskGroups {
		if err := tg.Warnings(j); err != nil {
			outer := fmt.Errorf("Group %q has warnings: %v", tg.Name, err)
			mErr.Errors = append(mErr.Errors, outer)
		}
		if tg.Update != nil && tg.Update.AutoPromote {
			ap += 1
		}
	}

	// Check AutoPromote, should be all or none
	if ap > 0 && ap < len(j.TaskGroups) {
		err := fmt.Errorf("auto_promote must be true for all groups to enable automatic promotion")
		mErr.Errors = append(mErr.Errors, err)
	}
	return kinds
}

// ConfigEntries accumulates the Consul Configuration Entries defined in task groups
// of j.
//
// Currently Nomad only supports entries for connect ingress gateways.
func (j *Job) ConfigEntries() map[string]*ConsulIngressConfigEntry {
	igEntries := make(map[string]*ConsulIngressConfigEntry)
	for _, tg := range j.TaskGroups {
		for _, service := range tg.Services {
			if service.Connect.IsGateway() {
				if ig := service.Connect.Gateway.Ingress; ig != nil {
					igEntries[service.Name] = ig
				}
				// imagine also accumulating other entry types in the future
			}
		}
	}
	return igEntries
}

// RequiredSignals returns a mapping of task groups to tasks to their required
// set of signals
func (j *Job) RequiredSignals() map[string]map[string][]string {
	switch u.HealthCheck {
	case UpdateStrategyHealthCheck_Checks, UpdateStrategyHealthCheck_TaskStates, UpdateStrategyHealthCheck_Manual:
	default:
		multierror.Append(&mErr, fmt.Errorf("Invalid health check given: %q", u.HealthCheck))
	}

	if u.MaxParallel < 0 {
		multierror.Append(&mErr, fmt.Errorf("Max parallel can not be less than zero: %d < 0", u.MaxParallel))
	}
	if u.Canary < 0 {
		multierror.Append(&mErr, fmt.Errorf("Canary count can not be less than zero: %d < 0", u.Canary))
	}
	if u.Canary == 0 && u.AutoPromote {
		multierror.Append(&mErr, fmt.Errorf("Auto Promote requires a Canary count greater than zero"))
	}
	if u.MinHealthyTime < 0 {
		multierror.Append(&mErr, fmt.Errorf("Minimum healthy time may not be less than zero: %v", u.MinHealthyTime))
	}
	if u.HealthyDeadline <= 0 {
		multierror.Append(&mErr, fmt.Errorf("Healthy deadline must be greater than zero: %v", u.HealthyDeadline))
	}
	if u.ProgressDeadline < 0 {
		multierror.Append(&mErr, fmt.Errorf("Progress deadline must be zero or greater: %v", u.ProgressDeadline))
	}
	if u.MinHealthyTime >= u.HealthyDeadline {
		multierror.Append(&mErr, fmt.Errorf("Minimum healthy time must be less than healthy deadline: %v > %v", u.MinHealthyTime, u.HealthyDeadline))
	}
	if u.ProgressDeadline != 0 && u.HealthyDeadline >= u.ProgressDeadline {
		multierror.Append(&mErr, fmt.Errorf("Healthy deadline must be less than progress deadline: %v > %v", u.HealthyDeadline, u.ProgressDeadline))
	}
	if u.Stagger <= 0 {
		multierror.Append(&mErr, fmt.Errorf("Stagger must be greater than zero: %v", u.Stagger))
	}

	return mErr.ErrorOrNil()
			Datacenters: []string{},
			Meta:        map[string]string{},
		}
		for _, dc := range region.Datacenters {
			copyRegion.Datacenters = append(copyRegion.Datacenters, dc)
		}
		for k, v := range region.Meta {
			copyRegion.Meta[k] = v
		}
	}

	// Write all the user set fields
	hash.Write([]byte(n.Name))
	hash.Write([]byte(n.Description))
	hash.Write([]byte(n.Quota))

	// Finalize the hash
	hashVal := hash.Sum(nil)

	var mErr multierror.Error
	if p.Spec == "" {
		multierror.Append(&mErr, fmt.Errorf("Must specify a spec"))
	}

	// Check if we got a valid time zone
	if p.TimeZone != "" {
		if _, err := time.LoadLocation(p.TimeZone); err != nil {
			multierror.Append(&mErr, fmt.Errorf("Invalid time zone %q: %v", p.TimeZone, err))
		}
	}

	case PeriodicSpecCron:
		// Validate the cron spec
		if _, err := cronexpr.Parse(p.Spec); err != nil {
			multierror.Append(&mErr, fmt.Errorf("Invalid cron spec %q: %v", p.Spec, err))
		}
	case PeriodicSpecTest:
		// No-op
	default:
		multierror.Append(&mErr, fmt.Errorf("Unknown periodic specification type %q", p.SpecType))
	}

	return mErr.ErrorOrNil()
	switch d.Payload {
	case DispatchPayloadOptional, DispatchPayloadRequired, DispatchPayloadForbidden:
	default:
		multierror.Append(&mErr, fmt.Errorf("Unknown payload requirement: %q", d.Payload))
	}

	// Check that the meta configurations are disjoint sets
	disjoint, offending := helper.SliceSetDisjoint(d.MetaRequired, d.MetaOptional)
	if !disjoint {
		multierror.Append(&mErr, fmt.Errorf("Required and optional meta keys should be disjoint. Following keys exist in both: %v", offending))
	}

	return mErr.ErrorOrNil()
	if p.Max < 0 {
		mErr.Errors = append(mErr.Errors,
			fmt.Errorf("maximum count must be specified and non-negative"))
	} else {
		if p.Max < p.Min {
			mErr.Errors = append(mErr.Errors,
				fmt.Errorf("maximum count must not be less than minimum count"))
		}
	}

	if p.Min < 0 {
	switch r.Mode {
	case RestartPolicyModeDelay, RestartPolicyModeFail:
	default:
		multierror.Append(&mErr, fmt.Errorf("Unsupported restart mode: %q", r.Mode))
	}

	// Check for ambiguous/confusing settings
	if r.Attempts == 0 && r.Mode != RestartPolicyModeFail {
		multierror.Append(&mErr, fmt.Errorf("Restart policy %q with %d attempts is ambiguous", r.Mode, r.Attempts))
	}

	if r.Interval.Nanoseconds() < RestartPolicyMinInterval.Nanoseconds() {
		multierror.Append(&mErr, fmt.Errorf("Interval can not be less than %v (got %v)", RestartPolicyMinInterval, r.Interval))
	}
	if time.Duration(r.Attempts)*r.Delay > r.Interval {
		multierror.Append(&mErr,
			fmt.Errorf("Nomad can't restart the TaskGroup %v times in an interval of %v with a delay of %v", r.Attempts, r.Interval, r.Delay))
	}
	return mErr.ErrorOrNil()
	// Check for ambiguous/confusing settings
	if r.Attempts > 0 {
		if r.Interval <= 0 {
			multierror.Append(&mErr, fmt.Errorf("Interval must be a non zero value if Attempts > 0"))
		}
		if r.Unlimited {
			multierror.Append(&mErr, fmt.Errorf("Reschedule Policy with Attempts = %v, Interval = %v, "+
				"and Unlimited = %v is ambiguous", r.Attempts, r.Interval, r.Unlimited))
			multierror.Append(&mErr, errors.New("If Attempts >0, Unlimited cannot also be set to true"))
		}
	}

	delayPreCheck := true
	// Delay should be bigger than the default
	if r.Delay.Nanoseconds() < ReschedulePolicyMinDelay.Nanoseconds() {
		multierror.Append(&mErr, fmt.Errorf("Delay cannot be less than %v (got %v)", ReschedulePolicyMinDelay, r.Delay))
		delayPreCheck = false
	}

	// Must use a valid delay function
	if !isValidDelayFunction(r.DelayFunction) {
		multierror.Append(&mErr, fmt.Errorf("Invalid delay function %q, must be one of %q", r.DelayFunction, RescheduleDelayFunctions))
		delayPreCheck = false
	}

	// Validate MaxDelay if not using linear delay progression
	if r.DelayFunction != "constant" {
		if r.MaxDelay.Nanoseconds() < ReschedulePolicyMinDelay.Nanoseconds() {
			multierror.Append(&mErr, fmt.Errorf("Max Delay cannot be less than %v (got %v)", ReschedulePolicyMinDelay, r.Delay))
			delayPreCheck = false
		}
		if r.MaxDelay < r.Delay {
			multierror.Append(&mErr, fmt.Errorf("Max Delay cannot be less than Delay %v (got %v)", r.Delay, r.MaxDelay))
			delayPreCheck = false
		}

	// Validate Interval and other delay parameters if attempts are limited
	if !r.Unlimited {
		if r.Interval.Nanoseconds() < ReschedulePolicyMinInterval.Nanoseconds() {
			multierror.Append(&mErr, fmt.Errorf("Interval cannot be less than %v (got %v)", ReschedulePolicyMinInterval, r.Interval))
		}
		if !delayPreCheck {
			// We can't cross validate the rest of the delay params if delayPreCheck fails, so return early
		}
		crossValidationErr := r.validateDelayParams()
		if crossValidationErr != nil {
			multierror.Append(&mErr, crossValidationErr)
		}
	}
	return mErr.ErrorOrNil()
	}
	var mErr multierror.Error
	if r.DelayFunction == "constant" {
		multierror.Append(&mErr, fmt.Errorf("Nomad can only make %v attempts in %v with initial delay %v and "+
			"delay function %q", possibleAttempts, r.Interval, r.Delay, r.DelayFunction))
	} else {
		multierror.Append(&mErr, fmt.Errorf("Nomad can only make %v attempts in %v with initial delay %v, "+
			"delay function %q, and delay ceiling %v", possibleAttempts, r.Interval, r.Delay, r.DelayFunction, r.MaxDelay))
	}
	multierror.Append(&mErr, fmt.Errorf("Set the interval to at least %v to accommodate %v attempts", recommendedInterval.Round(time.Second), r.Attempts))
	return mErr.ErrorOrNil()
}

	var mErr multierror.Error

	if m.MaxParallel < 0 {
		multierror.Append(&mErr, fmt.Errorf("MaxParallel must be >= 0 but found %d", m.MaxParallel))
	}

	switch m.HealthCheck {
		// ok
	case "":
		if m.MaxParallel > 0 {
			multierror.Append(&mErr, fmt.Errorf("Missing HealthCheck"))
		}
	default:
		multierror.Append(&mErr, fmt.Errorf("Invalid HealthCheck: %q", m.HealthCheck))
	}

	if m.MinHealthyTime < 0 {
		multierror.Append(&mErr, fmt.Errorf("MinHealthyTime is %s and must be >= 0", m.MinHealthyTime))
	}

	if m.HealthyDeadline < 0 {
		multierror.Append(&mErr, fmt.Errorf("HealthyDeadline is %s and must be >= 0", m.HealthyDeadline))
	}

	if m.MinHealthyTime > m.HealthyDeadline {
		multierror.Append(&mErr, fmt.Errorf("MinHealthyTime must be less than HealthyDeadline"))
	}

	return mErr.ErrorOrNil()
func (tg *TaskGroup) validateNetworks() error {
	var mErr multierror.Error
	portLabels := make(map[string]string)
	staticPorts := make(map[int]string)

	for _, net := range tg.Networks {
		for _, port := range append(net.ReservedPorts, net.DynamicPorts...) {
			}

			if port.Value != 0 {
				// static port
				if other, ok := staticPorts[port.Value]; ok {
					err := fmt.Errorf("Static port %d already reserved by %s", port.Value, other)
					mErr.Errors = append(mErr.Errors, err)
				} else {
					staticPorts[port.Value] = fmt.Sprintf("taskgroup network:%s", port.Label)
				}
			}

				}

				if port.Value != 0 {
					if other, ok := staticPorts[port.Value]; ok {
						err := fmt.Errorf("Static port %d already reserved by %s", port.Value, other)
						mErr.Errors = append(mErr.Errors, err)
						mErr.Errors = append(mErr.Errors, err)
					} else {
						staticPorts[port.Value] = fmt.Sprintf("%s:%s", task.Name, port.Label)
					}
				}
			}
	return mErr.ErrorOrNil()
}

// validateServices runs Service.Validate() on group-level services,
// checks that group services do not conflict with task services and that
// group service checks that refer to tasks only refer to tasks that exist.
func (tg *TaskGroup) validateServices() error {
	var mErr multierror.Error
	knownTasks := make(map[string]struct{})
	knownServices := make(map[string]struct{})

	// Create a map of known tasks and their services so we can compare
	// vs the group-level services and checks
			continue
		}
		for _, service := range task.Services {
			if _, ok := knownServices[service.Name+service.PortLabel]; ok {
				mErr.Errors = append(mErr.Errors, fmt.Errorf("Service %s is duplicate", service.Name))
			}
			for _, check := range service.Checks {
				if check.TaskName != "" {
					mErr.Errors = append(mErr.Errors, fmt.Errorf("Check %s is invalid: only task group service checks can be assigned tasks", check.Name))
				}
			}
			knownServices[service.Name+service.PortLabel] = struct{}{}
		}
	}
	for i, service := range tg.Services {
		if service.AddressMode == AddressModeDriver {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("service %q cannot use address_mode=\"driver\", only services defined in a \"task\" block can use this mode", service.Name))
		}
		if _, ok := knownServices[service.Name+service.PortLabel]; ok {
			mErr.Errors = append(mErr.Errors, fmt.Errorf("Service %s is duplicate", service.Name))
		}
		knownServices[service.Name+service.PortLabel] = struct{}{}
		for _, check := range service.Checks {
			if check.TaskName != "" {
				if check.Type != ServiceCheckScript && check.Type != ServiceCheckGRPC {
			}
		}

		// Ensure that check names are unique and have valid ports
		knownChecks := make(map[string]struct{})
		for _, check := range service.Checks {
	return k.hasPrefix(ConnectIngressPrefix)
}

func (k TaskKind) IsAnyConnectGateway() bool {
	switch {
	case k.IsConnectIngress():
		return true
	default:
		return false
	}
	// ConnectTerminatingPrefix is the prefix used for fields referencing a Consul
	// Connect Terminating Gateway Proxy.
	//
	// Not yet supported.
	// ConnectTerminatingPrefix = "connect-terminating"

	// ConnectMeshPrefix is the prefix used for fields referencing a Consul Connect
	// Mesh Gateway Proxy.

	// Verify we have something to render
	if t.SourcePath == "" && t.EmbeddedTmpl == "" {
		multierror.Append(&mErr, fmt.Errorf("Must specify a source path or have an embedded template"))
	}

	// Verify we can render somewhere
	if t.DestPath == "" {
		multierror.Append(&mErr, fmt.Errorf("Must specify a destination for the template"))
	}

	// Verify the destination doesn't escape
	case TemplateChangeModeNoop, TemplateChangeModeRestart:
	case TemplateChangeModeSignal:
		if t.ChangeSignal == "" {
			multierror.Append(&mErr, fmt.Errorf("Must specify signal value when change mode is signal"))
		}
		if t.Envvars {
			multierror.Append(&mErr, fmt.Errorf("cannot use signals with env var templates"))
		}
	default:
		multierror.Append(&mErr, TemplateChangeModeInvalidError)
	}

	// Verify the splay is positive
	if t.Splay < 0 {
		multierror.Append(&mErr, fmt.Errorf("Must specify positive splay value"))
	}

	// Verify the permissions
	if t.Perms != "" {
		if _, err := strconv.ParseUint(t.Perms, 8, 12); err != nil {
			multierror.Append(&mErr, fmt.Errorf("Failed to parse %q as octal: %v", t.Perms, err))
		}
	}

	}
	sort.Strings(keys)
	for _, k := range keys {
		h.Write([]byte(k))
		h.Write([]byte(m[k]))
	}
}

		panic(err)
	}

	h.Write([]byte(ta.GetterSource))

	hashStringMap(h, ta.GetterOptions)
	hashStringMap(h, ta.GetterHeaders)

	h.Write([]byte(ta.GetterMode))
	h.Write([]byte(ta.RelativeDest))
	return base64.RawStdEncoding.EncodeToString(h.Sum(nil))
}

	LTarget string // Left-hand target
	RTarget string // Right-hand target
	Operand string // Constraint operand (<=, <, =, !=, >, >=), contains, near
	str     string // Memoized string
}

// Equal checks if two constraints are equal
func (c *Constraint) Equals(o *Constraint) bool {
	return c == o ||
		c.LTarget == o.LTarget &&
			c.Operand == o.Operand
}

func (c *Constraint) Equal(o *Constraint) bool {
	return c.Equals(o)
}
	if c == nil {
		return nil
	}
	nc := new(Constraint)
	*nc = *c
	return nc
}

func (c *Constraint) String() string {
	if c.str != "" {
		return c.str
	}
	c.str = fmt.Sprintf("%s %s %s", c.LTarget, c.Operand, c.RTarget)
	return c.str
}

func (c *Constraint) Validate() error {
	RTarget string // Right-hand target
	Operand string // Affinity operand (<=, <, =, !=, >, >=), set_contains_all, set_contains_any
	Weight  int8   // Weight applied to nodes that match the affinity. Can be negative
	str     string // Memoized string
}

// Equal checks if two affinities are equal
func (a *Affinity) Equals(o *Affinity) bool {
	return a == o ||
		a.LTarget == o.LTarget &&
	if a == nil {
		return nil
	}
	na := new(Affinity)
	*na = *a
	return na
}

func (a *Affinity) String() string {
	if a.str != "" {
		return a.str
	}
	a.str = fmt.Sprintf("%s %s %s %v", a.LTarget, a.Operand, a.RTarget, a.Weight)
	return a.str
}

func (a *Affinity) Validate() error {
		if !ok {
			seen[target.Value] = struct{}{}
		} else {
			mErr.Errors = append(mErr.Errors, errors.New(fmt.Sprintf("Spread target value %q already defined", target.Value)))
		}
		if target.Percent < 0 || target.Percent > 100 {
			mErr.Errors = append(mErr.Errors, errors.New(fmt.Sprintf("Spread target percentage for value %q must be between 0 and 100", target.Value)))
		}
		sumPercent += uint32(target.Percent)
	}
	if sumPercent > 100 {
		mErr.Errors = append(mErr.Errors, errors.New(fmt.Sprintf("Sum of spread target percentages must not be greater than 100%%; got %d%%", sumPercent)))
	}
	return mErr.ErrorOrNil()
}

	var mErr multierror.Error
	if len(v.Policies) == 0 {
		multierror.Append(&mErr, fmt.Errorf("Policy list cannot be empty"))
	}

	for _, p := range v.Policies {
		if p == "root" {
			multierror.Append(&mErr, fmt.Errorf("Can not specify \"root\" policy"))
		}
	}

	switch v.ChangeMode {
	case VaultChangeModeSignal:
		if v.ChangeSignal == "" {
			multierror.Append(&mErr, fmt.Errorf("Signal must be specified when using change mode %q", VaultChangeModeSignal))
		}
	case VaultChangeModeNoop, VaultChangeModeRestart:
	default:
		multierror.Append(&mErr, fmt.Errorf("Unknown change mode %q", v.ChangeMode))
	}

	return mErr.ErrorOrNil()
		return false
	}
	for _, group := range d.TaskGroups {
		if !group.AutoPromote {
			return false
		}
	}
	AutoPromote bool

	// ProgressDeadline is the deadline by which an allocation must transition
	// to healthy before the deployment is considered failed.
	ProgressDeadline time.Duration

	// RequireProgressBy is the time by which an allocation must transition
	// to healthy before the deployment is considered failed.
	RequireProgressBy time.Time

	// Promoted marks whether the canaries have been promoted
	ModifyTime int64
}

// Index returns the index of the allocation. If the allocation is from a task
// group with count greater than 1, there will be multiple allocations for it.
func (a *Allocation) Index() uint {
}

func setDisplayMsg(taskStates map[string]*TaskState) {
	if taskStates != nil {
		for _, taskState := range taskStates {
			for _, event := range taskState.Events {
				event.PopulateEventDisplayMessage()
			}
		}
	}
}
	expected := 0
	actual := 0
	for name, allocList := range plan.NodeAllocation {
		didAlloc, _ := p.NodeAllocation[name]
		expected += len(allocList)
		actual += len(didAlloc)
	}
	}

	// Write all the user set fields
	hash.Write([]byte(c.Name))
	hash.Write([]byte(c.Description))
	hash.Write([]byte(c.Rules))

	// Finalize the hash
	hashVal := hash.Sum(nil)
	ModifyIndex uint64
}

var (
	// AnonymousACLToken is used no SecretID is provided, and the
	// request is made anonymously.
	}

	// Write all the user set fields
	hash.Write([]byte(a.Name))
	hash.Write([]byte(a.Type))
	for _, policyName := range a.Policies {
		hash.Write([]byte(policyName))
	}
	if a.Global {
		hash.Write([]byte("global"))
	} else {
		hash.Write([]byte("local"))
	}

	// Finalize the hash
