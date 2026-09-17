package main

		logger.WithError(err).Fatal("could not initialize API client with security provider")
	}

	key := viper.GetString("access_key_id")
	secret := viper.GetString("secret_access_key")
	svc := SetupTestS3Client(key, secret)
	return logger, client, svc
}

func SetupTestS3Client(key, secret string) *s3.S3 {
	s3Endpoint := viper.GetString("s3_endpoint")
	awsSession := session.Must(session.NewSession())
	svc := s3.New(awsSession,
			WithCredentials(credentials.NewCredentials(
				&credentials.StaticProvider{
					Value: credentials.Value{
						AccessKeyID:     key,
						SecretAccessKey: secret,
					}})))
	return svc
}

// ParseEndpointURL parses the given endpoint string
