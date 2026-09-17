package main

		log.Info("Error loading config")
		return nil, err
	}
	log.Infof("Loaded client config: %#v", cfg.Spec)
	return NewClientFromConfig(cfg)
}

