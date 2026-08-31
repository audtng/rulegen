package main

	}

	for _, ingress := range service.Status.LoadBalancer.Ingress {
		info.loadBalancerIngressIPs = append(info.loadBalancerIngressIPs, &loadBalancerIngressInfo{ip: ingress.IP})
	}
	return info
}
