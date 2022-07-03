package util

import "github.com/aws/aws-sdk-go-v2/aws"

// LocalResolver provides resolver which takes in account provided awsEndpoint and routes requests to defined endpoint
func LocalResolver(awsEndpoint, awsRegion string) aws.EndpointResolverWithOptionsFunc {
	return func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if awsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpoint,
				SigningRegion: awsRegion,
			}, nil
		}

		// fallback to default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}
}
