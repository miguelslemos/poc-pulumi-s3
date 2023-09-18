package main

import (
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")
		bucketName := conf.Require("bucket-name")
		bucket, err := s3.NewBucket(ctx, bucketName, nil)
		if err != nil {
			return err
		}

		// Export the name of the bucket
		ctx.Export("bucketNameID", bucket.ID())
		ctx.Export("bucketArn", bucket.Arn)
		ctx.Export("bucketName", bucket.BucketName)
		return nil
	})
}
