package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		config := config.New(ctx, "")
		bucketName := config.Require("bucketName")
		var domains []string
		if err := config.TryObject("domains", &domains); err != nil {
			return err
		}

		bucketContent, err := LookupBucket(ctx, bucketName)

		if err != nil {
			return err
		}

		// Create a cloudfront distribution to serve the content from the bucket.
		distribution, err := CreateDistribution(ctx, "poc", &CreateDistributionArgs{
			Bucket:     bucketContent,
			BucketName: bucketName,
		})
		if err != nil {
			return err
		}

		if len(domains) > 0 {
			for _, domain := range domains {
				// Create a record to point to the cloudfront distribution.
				hostedZone, err := LookupHostedZone(ctx, extractDomain(domain))
				if err != nil {
					return err
				}
				_, err = CreateRecord(ctx, "poc", CreateRecordArgs{
					HostedZoneId:                     hostedZone.ZoneId,
					DnsName:                          domain,
					DistributionHostedZoneDomainName: distribution.DomainName,
					DistributionHostedZoneId:         distribution.HostedZoneId,
				}, pulumi.DependsOn([]pulumi.Resource{distribution}))
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
