package main

import (
	"encoding/json"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudfront"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		config := config.New(ctx, "")
		bucketName := config.Require("bucketName")

		// Generate Origin Access Identity to access the private s3 bucket.
		originAccessIdentity, err := cloudfront.NewOriginAccessIdentity(ctx, "originAccessIdentity", &cloudfront.OriginAccessIdentityArgs{
			Comment: pulumi.String("this is needed to setup s3 polices and make s3 not public."),
		})
		if err != nil {
			return err
		}

		// Lookup the bucket by name.
		bucketContent, err := s3.LookupBucket(ctx, &s3.LookupBucketArgs{
			Bucket: bucketName,
		}, nil)

		if err != nil {
			return err
		}

		// Create a bucket policy to allow cloudfront to access the bucket.
		_, err = s3.NewBucketPolicy(ctx, "cloudfront-bucket-policy", &s3.BucketPolicyArgs{
			Bucket: pulumi.String(bucketContent.Id),
			Policy: pulumi.All(bucketContent.Arn, originAccessIdentity.IamArn).ApplyT(
				func(args []any) (string, error) {
					bucketArn := args[0].(string)
					iamArn := args[1].(string)
					policy, err := json.Marshal(map[string]any{
						"Version": "2012-10-17",
						"Statement": []map[string]any{
							{
								"Sid":    "CloudfrontAllow",
								"Effect": "Allow",
								"Principal": map[string]any{
									"AWS": iamArn,
								},
								"Action":   "s3:GetObject",
								"Resource": bucketArn + "/*",
							},
						},
					})
					if err != nil {
						return "", err
					}
					return string(policy), nil
				}).(pulumi.StringOutput),
		})

		if err != nil {
			return err
		}

		// Load Distribution config
		var distConfig = DistributionConfig(ctx)

		// Create a cloudfront distribution to serve the content from the bucket.
		distribution, err := cloudfront.NewDistribution(ctx, "distribution", &distConfig)
		if err != nil {
			return err
		}

		ctx.Export("contentBucketUri", pulumi.String(bucketContent.Bucket))
		ctx.Export("contentBucketWebsiteEndpoint", pulumi.String(bucketContent.WebsiteEndpoint))
		ctx.Export("cloudFrontDomain", distribution.DomainName)
		return nil
	})
}
