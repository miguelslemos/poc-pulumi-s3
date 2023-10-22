package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudfront"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CreateDistributionArgs struct {
	BucketName               string
	BucketRegionalDomainName pulumi.StringInput
	DistributionConfig       NuDistributionArgs
}

func CreateDistribution(ctx *pulumi.Context, prefix string, args *CreateDistributionArgs) (*cloudfront.Distribution, error) {

	// Generate OriginAccessControl to access the private s3 bucket.
	oac, err := cloudfront.NewOriginAccessControl(ctx, fmt.Sprintf("%s-oac-%s", prefix, args.BucketName), &cloudfront.OriginAccessControlArgs{
		Description:                   pulumi.String(fmt.Sprintf("Origin access control for %s", args.BucketName)),
		OriginAccessControlOriginType: pulumi.String("s3"),
		SigningBehavior:               pulumi.String("always"),
		SigningProtocol:               pulumi.String("sigv4"),
	})
	if err != nil {
		return nil, err
	}

	// Override config to use the information from the bucket.
	args.DistributionConfig.DistributionArgs.Origins = cloudfront.DistributionOriginArray([]cloudfront.DistributionOriginInput{
		cloudfront.DistributionOriginArgs{
			DomainName: pulumi.String("poc-cloudfront-static-4213.s3.us-east-1.amazonaws.com"),
			// DomainName:            args.BucketRegionalDomainName,
			OriginId:              pulumi.String("bucket-origin"),
			OriginAccessControlId: oac.ID(),
		},
	})

	distribution, err := cloudfront.NewDistribution(ctx, fmt.Sprintf("%s-distribution-%s", prefix, args.BucketName), &args.DistributionConfig.DistributionArgs, pulumi.DependsOn([]pulumi.Resource{oac}))
	if err != nil {
		return nil, err
	}
	return distribution, nil
}
