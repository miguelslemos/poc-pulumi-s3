package main

import (
	"encoding/json"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudfront"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type NuDistributionArgs struct {
	Aliases                      []string
	Comment                      *string
	ContinuousDeploymentPolicyId *string
	CustomErrorResponses         []cloudfront.DistributionCustomErrorResponse
	DefaultCacheBehavior         cloudfront.DistributionDefaultCacheBehavior   `pulumi:"defaultCacheBehavior"`
	DefaultRootObject            *string                                       `pulumi:"defaultRootObject"`
	Enabled                      bool                                          `pulumi:"enabled"`
	HttpVersion                  *string                                       `pulumi:"httpVersion"`
	IsIpv6Enabled                *bool                                         `pulumi:"isIpv6Enabled"`
	LoggingConfig                *cloudfront.DistributionLoggingConfig         `pulumi:"loggingConfig"`
	OrderedCacheBehaviors        []cloudfront.DistributionOrderedCacheBehavior `pulumi:"orderedCacheBehaviors"`
	OriginGroups                 []cloudfront.DistributionOriginGroup          `pulumi:"originGroups"`
	Origins                      []cloudfront.DistributionOrigin               `pulumi:"origins"`
	PriceClass                   *string                                       `pulumi:"priceClass"`
	Restrictions                 cloudfront.DistributionRestrictions
	Staging                      *bool             `pulumi:"staging"`
	Tags                         map[string]string `pulumi:"tags"`
	ViewerCertificate            cloudfront.DistributionViewerCertificate
	WaitForDeployment            *bool   `pulumi:"waitForDeployment"`
	WebAclId                     *string `pulumi:"webAclId"`
}

func toDistributionOrderedCacheBehavior(a []cloudfront.DistributionOrderedCacheBehavior) cloudfront.DistributionOrderedCacheBehaviorArrayInput {
	var res []cloudfront.DistributionOrderedCacheBehaviorInput
	for _, s := range a {
		res = append(res, cloudfront.DistributionOrderedCacheBehaviorArgs{
			AllowedMethods: pulumi.ToStringArray(s.AllowedMethods),
			CachedMethods:  pulumi.ToStringArray(s.CachedMethods),
			Compress:       pulumi.BoolPtrFromPtr(s.Compress),
			DefaultTtl:     pulumi.IntPtrFromPtr(s.DefaultTtl),
			ForwardedValues: cloudfront.DistributionOrderedCacheBehaviorForwardedValuesArgs{
				Cookies: cloudfront.DistributionOrderedCacheBehaviorForwardedValuesCookiesArgs{
					Forward: pulumi.String(s.ForwardedValues.Cookies.Forward),
				},
				QueryString: pulumi.Bool(s.ForwardedValues.QueryString),
				Headers:     pulumi.ToStringArray(s.ForwardedValues.Headers),
			},
		})
	}
	return cloudfront.DistributionOrderedCacheBehaviorArray(res)
}

func toDistributionOriginGroupMember(a []cloudfront.DistributionOriginGroupMember) cloudfront.DistributionOriginGroupMemberArrayInput {
	var res []cloudfront.DistributionOriginGroupMemberInput
	for _, s := range a {
		res = append(res, cloudfront.DistributionOriginGroupMemberArgs{
			OriginId: pulumi.String(s.OriginId),
		})
	}
	return cloudfront.DistributionOriginGroupMemberArray(res)
}
func toDistributionOriginGroup(a []cloudfront.DistributionOriginGroup) cloudfront.DistributionOriginGroupArrayInput {
	var res []cloudfront.DistributionOriginGroupInput
	for _, s := range a {
		res = append(res, cloudfront.DistributionOriginGroupArgs{
			OriginId: pulumi.String(s.OriginId),
			FailoverCriteria: cloudfront.DistributionOriginGroupFailoverCriteriaArgs{
				StatusCodes: pulumi.ToIntArray(s.FailoverCriteria.StatusCodes),
			},
			Members: toDistributionOriginGroupMember(s.Members),
		})
	}
	return cloudfront.DistributionOriginGroupArray(res)
}

// func toDistributionOriginCustomHeader(a []cloudfront.DistributionOriginCustomHeader) cloudfront.DistributionOriginCustomHeaderArrayInput {
// 	var res []cloudfront.DistributionOriginCustomHeaderInput
// 	for _, s := range a {
// 		res = append(res, cloudfront.DistributionOriginCustomHeaderArgs{
// 			Name:  pulumi.String(s.Name),
// 			Value: pulumi.String(s.Value),
// 		})
// 	}
// 	return cloudfront.DistributionOriginCustomHeaderArray(res)
// }

// func toDistributionOriginArray(a []cloudfront.DistributionOrigin) cloudfront.DistributionOriginArrayInput {
// 	var res []cloudfront.DistributionOriginInput
// 	for _, s := range a {
// 		res = append(res, cloudfront.DistributionOriginArgs{
// 			DomainName: pulumi.String(s.DomainName),
// 			OriginId:   pulumi.String(s.OriginId),
// 			CustomOriginConfig: cloudfront.DistributionOriginCustomOriginConfigArgs{
// 				HttpPort:             pulumi.Int(s.CustomOriginConfig.HttpPort),
// 				HttpsPort:            pulumi.Int(s.CustomOriginConfig.HttpsPort),
// 				OriginProtocolPolicy: pulumi.String(s.CustomOriginConfig.OriginProtocolPolicy),
// 				OriginSslProtocols:   pulumi.ToStringArray(s.CustomOriginConfig.OriginSslProtocols),
// 			},
// 			ConnectionAttempts:    pulumi.IntPtrFromPtr(s.ConnectionAttempts),
// 			ConnectionTimeout:     pulumi.IntPtrFromPtr(s.ConnectionTimeout),
// 			OriginPath:            pulumi.StringPtrFromPtr(s.OriginPath),
// 			OriginAccessControlId: pulumi.StringPtrFromPtr(s.OriginAccessControlId),
// 			OriginShield: cloudfront.DistributionOriginOriginShieldArgs{
// 				Enabled:            pulumi.Bool(s.OriginShield.Enabled),
// 				OriginShieldRegion: pulumi.StringPtrFromPtr(s.OriginShield.OriginShieldRegion),
// 			},
// 			S3OriginConfig: cloudfront.DistributionOriginS3OriginConfigArgs{
// 				OriginAccessIdentity: pulumi.String(s.S3OriginConfig.OriginAccessIdentity),
// 			},
// 			CustomHeaders: toDistributionOriginCustomHeader(s.CustomHeaders),
// 		})
// 	}
// 	return cloudfront.DistributionOriginArray(res)
// }

func toDistributionCustomErrorResponse(a []cloudfront.DistributionCustomErrorResponse) cloudfront.DistributionCustomErrorResponseArrayInput {
	var res []cloudfront.DistributionCustomErrorResponseInput
	for _, s := range a {
		res = append(res, cloudfront.DistributionCustomErrorResponseArgs{
			ErrorCode:          pulumi.Int(s.ErrorCode),
			ErrorCachingMinTtl: pulumi.IntPtrFromPtr(s.ErrorCachingMinTtl),
			ResponseCode:       pulumi.IntPtrFromPtr(s.ResponseCode),
			ResponsePagePath:   pulumi.StringPtrFromPtr(s.ResponsePagePath),
		})
	}
	return cloudfront.DistributionCustomErrorResponseArray(res)
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		config := config.New(ctx, "")
		bucketName := config.Require("bucketName")
		var distributionArgs NuDistributionArgs
		config.RequireObject("distribution", &distributionArgs)

		// Generate Origin Access Identity to access the private s3 bucket.
		originAccessIdentity, err := cloudfront.NewOriginAccessIdentity(ctx, "originAccessIdentity", &cloudfront.OriginAccessIdentityArgs{
			Comment: pulumi.String("this is needed to setup s3 polices and make s3 not public."),
		})
		if err != nil {
			return err
		}

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

		// Create a cloudfront distribution to serve the content from the bucket.
		distribution, err := cloudfront.NewDistribution(ctx, "distribution", &cloudfront.DistributionArgs{
			Aliases: pulumi.ToStringArray(distributionArgs.Aliases),
			Comment: pulumi.StringPtrFromPtr(distributionArgs.Comment),
			// Origins: toDistributionOriginArray(distributionArgs.Origins),
			DefaultCacheBehavior: cloudfront.DistributionDefaultCacheBehaviorArgs{
				AllowedMethods: pulumi.ToStringArray(distributionArgs.DefaultCacheBehavior.AllowedMethods),
				CachedMethods:  pulumi.ToStringArray(distributionArgs.DefaultCacheBehavior.CachedMethods),
				Compress:       pulumi.BoolPtrFromPtr(distributionArgs.DefaultCacheBehavior.Compress),
				DefaultTtl:     pulumi.IntPtrFromPtr(distributionArgs.DefaultCacheBehavior.DefaultTtl),
				ForwardedValues: cloudfront.DistributionDefaultCacheBehaviorForwardedValuesArgs{
					Cookies: cloudfront.DistributionDefaultCacheBehaviorForwardedValuesCookiesArgs{
						Forward: pulumi.String(distributionArgs.DefaultCacheBehavior.ForwardedValues.Cookies.Forward),
					},
					QueryString: pulumi.Bool(distributionArgs.DefaultCacheBehavior.ForwardedValues.QueryString),
					Headers:     pulumi.ToStringArray(distributionArgs.DefaultCacheBehavior.ForwardedValues.Headers),
				},
			},
			OriginGroups:                 toDistributionOriginGroup(distributionArgs.OriginGroups),
			ContinuousDeploymentPolicyId: pulumi.StringPtrFromPtr(distributionArgs.ContinuousDeploymentPolicyId),
			CustomErrorResponses:         toDistributionCustomErrorResponse(distributionArgs.CustomErrorResponses),
			DefaultRootObject:            pulumi.StringPtrFromPtr(distributionArgs.DefaultRootObject),
			Enabled:                      pulumi.Bool(distributionArgs.Enabled),
			HttpVersion:                  pulumi.StringPtrFromPtr(distributionArgs.HttpVersion),
			IsIpv6Enabled:                pulumi.BoolPtrFromPtr(distributionArgs.IsIpv6Enabled),
			LoggingConfig: cloudfront.DistributionLoggingConfigPtr(&cloudfront.DistributionLoggingConfigArgs{
				Bucket:         pulumi.String(distributionArgs.LoggingConfig.Bucket),
				IncludeCookies: pulumi.BoolPtrFromPtr(distributionArgs.LoggingConfig.IncludeCookies),
				Prefix:         pulumi.StringPtrFromPtr(distributionArgs.LoggingConfig.Prefix),
			}),
			OrderedCacheBehaviors: toDistributionOrderedCacheBehavior(distributionArgs.OrderedCacheBehaviors),
			PriceClass:            pulumi.StringPtrFromPtr(distributionArgs.PriceClass),
			Restrictions: cloudfront.DistributionRestrictionsArgs{
				GeoRestriction: cloudfront.DistributionRestrictionsGeoRestrictionArgs{
					Locations:       pulumi.ToStringArray(distributionArgs.Restrictions.GeoRestriction.Locations),
					RestrictionType: pulumi.String(distributionArgs.Restrictions.GeoRestriction.RestrictionType),
				},
			},
			Staging: pulumi.BoolPtrFromPtr(distributionArgs.Staging),
			Tags:    pulumi.ToStringMap(distributionArgs.Tags),
			ViewerCertificate: cloudfront.DistributionViewerCertificateArgs{
				AcmCertificateArn:      pulumi.StringPtrFromPtr(distributionArgs.ViewerCertificate.AcmCertificateArn),
				MinimumProtocolVersion: pulumi.StringPtrFromPtr(distributionArgs.ViewerCertificate.MinimumProtocolVersion),
				SslSupportMethod:       pulumi.StringPtrFromPtr(distributionArgs.ViewerCertificate.SslSupportMethod),
			},
			WaitForDeployment: pulumi.BoolPtrFromPtr(distributionArgs.WaitForDeployment),
			WebAclId:          pulumi.StringPtrFromPtr(distributionArgs.WebAclId),
		})
		if err != nil {
			return err
		}

		ctx.Export("contentBucketUri", pulumi.String(bucketContent.Bucket))
		ctx.Export("contentBucketWebsiteEndpoint", pulumi.String(bucketContent.WebsiteEndpoint))
		ctx.Export("cloudFrontDomain", distribution.DomainName)
		return nil
	})
}
