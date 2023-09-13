package main

import (
	"github.com/miguelslemos/pulumi-nu-packages/sdk/go/nu-packages/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		S3Bucket, err := storage.NewS3Bucket(ctx, "my-first-bucket-lib", nil)
		if err != nil {
			return err
		}
		ctx.Export("bucketName", S3Bucket.BucketId.Name())
		ctx.Export("bucketArn", S3Bucket.BucketId.Arn())
		return nil
	})
}
