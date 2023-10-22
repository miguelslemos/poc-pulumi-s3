package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		config := config.New(ctx, "")
		bucketName := config.Require("bucketName")

		bucketContent, err := LookupBucket(ctx, bucketName)

		if err != nil {
			return err
		}

		// Load Distribution config
		var distConfig = CreateDistributionConfig(ctx)

		// Create a cloudfront distribution to serve the content from the bucket.
		distribution, err := CreateDistribution(ctx, "poc", &CreateDistributionArgs{
			BucketName:               bucketName,
			BucketRegionalDomainName: bucketContent.BucketRegionalDomainName,
			DistributionConfig:       distConfig,
		})
		if err != nil {
			return err
		}

		if distConfig.Aliases != nil && len(distConfig.Aliases) > 0 {
			for _, alias := range distConfig.Aliases {
				// Create a record to point to the cloudfront distribution.
				hostedZone, err := LookupHostedZone(ctx, extractDomain(alias))
				if err != nil {
					return err
				}
				_, err = CreateRecord(ctx, "poc", CreateRecordArgs{
					HostedZoneId:                     hostedZone.ZoneId,
					DnsName:                          alias,
					DistributionHostedZoneDomainName: distribution.DomainName,
					DistributionHostedZoneId:         distribution.HostedZoneId,
				})
				if err != nil {
					return err
				}
			}
		}

		err = CreateBucketPolicy(ctx, "poc", BucketPolicyArgs{
			BucketName:      bucketName,
			DistributionArn: distribution.Arn,
		})

		if err != nil {
			return err
		}

		ctx.Export("BucketDomainName", bucketContent.BucketDomainName)
		ctx.Export("BucketRegionalDomainName", bucketContent.BucketRegionalDomainName)
		ctx.Export("Region", bucketContent.Region)
		ctx.Export("contentBucketWebsiteEndpoint", bucketContent.WebsiteEndpoint)
		ctx.Export("cloudFrontDomain", distribution.DomainName)
		return nil
	})
}
