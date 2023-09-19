package main

import (
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/s3outposts"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")
		bucketName := conf.Require("bucket-name")
		if _, err := s3outposts.NewBucketPolicy(ctx, "bucketPolicy", &s3outposts.BucketPolicyArgs{
			Bucket: pulumi.String(bucketName),
			PolicyDocument: pulumi.Any(map[string]interface{}{
				"Version": "2012-10-17",
				"Statement": []map[string]interface{}{
					{
						"Effect":    "Allow",
						"Principal": "*",
						"Action": []interface{}{
							"s3:GetObject",
						},
						"Resource": []interface{}{
							pulumi.Sprintf("arn:aws:s3:::%s/lala/*", pulumi.String(bucketName)),
						},
					},
				},
			}),
		}); err != nil {
			return err
		}
		return nil
	})
}
