package main

package esti

import (
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestDeleteObjects(t *testing.T) {
	assert.NoError(t, err)
	assert.Len(t, listOut.Contents, 0)
}
		logger.WithError(err).Fatal("could not initialize API client with security provider")
	}

	s3Endpoint := viper.GetString("s3_endpoint")
	awsSession := session.Must(session.NewSession())
	svc := s3.New(awsSession,
			WithCredentials(credentials.NewCredentials(
				&credentials.StaticProvider{
					Value: credentials.Value{
						AccessKeyID:     viper.GetString("access_key_id"),
						SecretAccessKey: viper.GetString("secret_access_key"),
					}})))

	return logger, client, svc
}

// ParseEndpointURL parses the given endpoint string
