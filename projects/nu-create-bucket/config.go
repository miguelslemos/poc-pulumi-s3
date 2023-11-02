package main

import (
	b64 "encoding/base64"
	"encoding/json"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const (
	ConfigKey = "bucket"
)

type BucketArgsInput struct {
	AccelerationStatus                *string                                     `pulumi:"accelerationStatus"`
	Acl                               *string                                     `pulumi:"acl"`
	Arn                               *string                                     `pulumi:"arn"`
	Bucket                            string                                      `pulumi:"bucket"`
	BucketPrefix                      *string                                     `pulumi:"bucketPrefix"`
	CorsRules                         []s3.BucketCorsRule                         `pulumi:"corsRules"`
	ForceDestroy                      *bool                                       `pulumi:"forceDestroy"`
	Grants                            []s3.BucketGrant                            `pulumi:"grants"`
	HostedZoneId                      *string                                     `pulumi:"hostedZoneId"`
	LifecycleRules                    []s3.BucketLifecycleRule                    `pulumi:"lifecycleRules"`
	Loggings                          []s3.BucketLogging                          `pulumi:"loggings"`
	ObjectLockConfiguration           *s3.BucketObjectLockConfiguration           `pulumi:"objectLockConfiguration"`
	Policy                            interface{}                                 `pulumi:"policy"`
	ReplicationConfiguration          *s3.BucketReplicationConfiguration          `pulumi:"replicationConfiguration"`
	RequestPayer                      *string                                     `pulumi:"requestPayer"`
	ServerSideEncryptionConfiguration *s3.BucketServerSideEncryptionConfiguration `pulumi:"serverSideEncryptionConfiguration"`
	Tags                              map[string]string                           `pulumi:"tags"`
	Versioning                        *s3.BucketVersioning                        `pulumi:"versioning"`
	Website                           *s3.BucketWebsite                           `pulumi:"website"`
	WebsiteDomain                     *string                                     `pulumi:"websiteDomain"`
	WebsiteEndpoint                   *string                                     `pulumi:"websiteEndpoint"`
}

type ProgramContext struct {
	pulumiCtx *pulumi.Context
	prefix    string
	stack     string
}

func decodeConfig(encodedJson string, args *BucketArgsInput) error {
	bytes, err := b64.StdEncoding.DecodeString(encodedJson)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &args)
	if err != nil {
		return err
	}

	return nil
}

func ReadConfig(context *ProgramContext) BucketArgsInput {
	var args BucketArgsInput
	configData, err := config.Try(context.pulumiCtx, "configData")
	if err == nil {
		err := decodeConfig(configData, &args)
		if err != nil {
			panic(err)
		}
	} else {
		config.RequireObject(context.pulumiCtx, ConfigKey, &args)
	}
	return args
}
