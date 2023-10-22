package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type BucketPolicyArgs struct {
	BucketName      string
	DistributionArn pulumi.StringInput
}

// E.g <bucket-name>.s3.<region>.amazonaws.com
func BuildBucketRegionalDomainName(bucketName pulumi.StringInput, region pulumi.StringInput) pulumi.StringOutput {
	return pulumi.Sprintf("%s.s3.%s.amazonaws.com", bucketName, region)
}

func LookupBucket(ctx *pulumi.Context, bucketName string) (*s3.BucketV2, error) {
	bucket, err := s3.GetBucketV2(ctx, bucketName, pulumi.ID(bucketName), nil)

	if err != nil {
		return nil, err
	}

	return bucket, nil
}

func CreateBucketPolicy(ctx *pulumi.Context, prefix string, bucketPolicyArgs BucketPolicyArgs) error {
	// Create a bucket policy to allow cloudfront to access the bucket.
	_, err := s3.NewBucketPolicy(ctx, fmt.Sprintf("%s-policy-%s", prefix, bucketPolicyArgs.BucketName), &s3.BucketPolicyArgs{
		Bucket: pulumi.String(bucketPolicyArgs.BucketName),
		Policy: pulumi.All(bucketPolicyArgs.BucketName, bucketPolicyArgs.DistributionArn).ApplyT(func(_args []interface{}) (string, error) {
			bucketName := _args[0].(string)
			distributionArn := _args[1].(string)
			return fmt.Sprintf(`
{
	"Version": "2012-10-17",
	"Id": "PolicyForCloudFrontPrivateContent",
	"Statement": [
		{
			"Sid": "AllowCloudFrontServicePrincipal",
			"Effect": "Allow",
			"Principal": {
				"Service": "cloudfront.amazonaws.com"
			},
			"Action": "s3:GetObject",
			"Resource": "arn:aws:s3:::%v/*",
			"Condition": {
				"StringEquals": {
					"AWS:SourceArn": "%v"
				}
			}
		}
	]
}
`, bucketName, distributionArn), nil
		}).(pulumi.StringOutput),
	})
	if err != nil {
		return err
	}
	return nil
}
