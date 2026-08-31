package main

		LeaderElectionResourceLock: opt.LeaderElection.ResourceLock,
		HealthProbeBindAddress:     net.JoinHostPort(opt.BindAddress, strconv.Itoa(opt.SecurePort)),
		LivenessEndpointName:       "/healthz",
		Namespace:                  util.GetCurrentNSOrDefault(),
	})
	if err != nil {
		klog.Errorf("Failed to build controllerManager ,%s", err)
	flags.IntVar(&o.SecurePort, "secure-port", defaultPort,
		"The secure port on which to serve HTTPS.")
	flags.BoolVar(&o.LeaderElection.LeaderElect, "leader-elect", true, "Start a leader election client and gain leadership before executing the main loop. Enable this when running replicated components for high availability.")
	flags.Float32Var(&o.KubeAPIQPS, "kube-api-qps", 100.0, "QPS to use while talking with kubean-apiserver. Doesn't cover events and node heartbeat apis which rate limiting is controlled by a different set of flags.")
	flags.IntVar(&o.KubeAPIBurst, "kube-api-burst", 100, "Burst to use while talking with kubean-apiserver. Doesn't cover events and node heartbeat apis which rate limiting is controlled by a different set of flags.")
}
