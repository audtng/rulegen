package main

		log.Info("Error loading config")
		return nil, err
	}
	// Do not log cfg.Spec directly (e.g. %#v / %+v / %v).
	// It embeds K8sAPIToken, KubeconfigInline, EtcdPassword, EtcdKey, and EtcdCert.
	log.WithFields(log.Fields{
		"datastoreType":       cfg.Spec.DatastoreType,
		"etcdEndpoints":       cfg.Spec.EtcdEndpoints,
		"k8sAPIEndpoint":      cfg.Spec.K8sAPIEndpoint,
		"kubeconfig":          cfg.Spec.Kubeconfig,
		"kubeconfigInlineSet": cfg.Spec.KubeconfigInline != "",
		"etcdUsernameSet":     cfg.Spec.EtcdUsername != "",
		"etcdPasswordSet":     cfg.Spec.EtcdPassword != "",
		"etcdKeyInlineSet":    cfg.Spec.EtcdKey != "",
		"etcdCertInlineSet":   cfg.Spec.EtcdCert != "",
		"k8sAPITokenSet":      cfg.Spec.K8sAPIToken != "",
	}).Info("Loaded client config")
	return NewClientFromConfig(cfg)
}

