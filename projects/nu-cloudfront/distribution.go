package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudfront"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CreateDistributionArgs struct {
	BucketName string
	Bucket     *s3.BucketV2
}

func CreateDistribution(ctx *pulumi.Context, prefix string, args *CreateDistributionArgs) (*cloudfront.Distribution, error) {

	oac, err := createOriginAccessControl(ctx, prefix, args)
	if err != nil {
		return nil, err
	}

	bucketOrigin := []cloudfront.DistributionOriginInput{
		cloudfront.DistributionOriginArgs{
			DomainName:            BuildBucketRegionalDomainName(args.Bucket.Bucket, args.Bucket.Region),
			OriginId:              pulumi.Sprintf("%s-origin", args.BucketName),
			OriginAccessControlId: oac.ID(),
		},
	}

	// Load Distribution config
	var distConfig = CreateDistributionConfig(ctx, bucketOrigin)

	// Create a cloudfront distribution to serve the content from the bucket.
	distribution, err := cloudfront.NewDistribution(ctx, fmt.Sprintf("%s-distribution-%s", prefix, args.BucketName), &distConfig.DistributionArgs, pulumi.DependsOn([]pulumi.Resource{oac}))
	if err != nil {
		return nil, err
	}
	return distribution, nil
}

// Generate OriginAccessControl to access the private s3 bucket.
func createOriginAccessControl(ctx *pulumi.Context, prefix string, args *CreateDistributionArgs) (*cloudfront.OriginAccessControl, error) {
	oac, err := cloudfront.NewOriginAccessControl(ctx, fmt.Sprintf("%s-oac-%s", prefix, args.BucketName), &cloudfront.OriginAccessControlArgs{
		Description:                   pulumi.Sprintf("Origin access control for %s", args.BucketName),
		OriginAccessControlOriginType: pulumi.String("s3"),
		SigningBehavior:               pulumi.String("always"),
		SigningProtocol:               pulumi.String("sigv4"),
	})
	if err != nil {
		return nil, err
	}
	return oac, nil
}
