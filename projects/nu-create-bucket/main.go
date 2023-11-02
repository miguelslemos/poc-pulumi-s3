package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		context := &ProgramContext{
			pulumiCtx: ctx,
			stack:     ctx.Stack(),
			prefix:    fmt.Sprintf("nu:%s", ctx.Stack()),
		}
		bucketInput := ReadConfig(context)
		bucketArgs := BuildBucketArg(bucketInput)

		bucket, err := CreateBucket(context, bucketInput.Bucket, bucketArgs)
		if err != nil {
			return err
		}

		context.pulumiCtx.Export("BucketName", bucket.Bucket)
		context.pulumiCtx.Export("BucketDomainName", bucket.BucketDomainName)
		context.pulumiCtx.Export("BucketArn", bucket.Arn)
		return nil
	})
}
