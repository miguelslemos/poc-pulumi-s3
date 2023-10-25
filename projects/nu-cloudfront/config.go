package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/cloudfront"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const (
	Distribution = "distribution"
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
	DistributionArgs             cloudfront.DistributionArgs
}

func decodeConfig(encodedJson string, distributionArgs *NuDistributionArgs) error {
	bytes, err := b64.StdEncoding.DecodeString(encodedJson)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &distributionArgs)
	if err != nil {
		return err
	}

	return nil
}

func readConfig(ctx *pulumi.Context) NuDistributionArgs {
	var distributionArgs NuDistributionArgs
	configData, err := config.Try(ctx, "configData")
	if err == nil {
		fmt.Println("configData found")
		err := decodeConfig(configData, &distributionArgs)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("configData not found")
		config.RequireObject(ctx, Distribution, &distributionArgs)
	}
	return distributionArgs
}

func CreateDistributionConfig(ctx *pulumi.Context, additionalOrigin []cloudfront.DistributionOriginInput) NuDistributionArgs {

	distributionArgs := readConfig(ctx)
	dist := cloudfront.DistributionArgs{
		Enabled: pulumi.Bool(distributionArgs.Enabled),
	}

	setAliases(&dist, distributionArgs.Aliases)
	setStringPtrFromPtr(&dist.Comment, distributionArgs.Comment)
	setStringPtrFromPtr(&dist.ContinuousDeploymentPolicyId, distributionArgs.ContinuousDeploymentPolicyId)
	setStringPtrFromPtr(&dist.DefaultRootObject, distributionArgs.DefaultRootObject)
	setStringPtrFromPtr(&dist.HttpVersion, distributionArgs.HttpVersion)
	setBoolPtrFromPtr(&dist.IsIpv6Enabled, distributionArgs.IsIpv6Enabled)

	setLoggingConfig(&dist, distributionArgs.LoggingConfig)
	setStringPtrFromPtr(&dist.PriceClass, distributionArgs.PriceClass)
	setBoolPtrFromPtr(&dist.Staging, distributionArgs.Staging)
	dist.Tags = pulumi.ToStringMap(distributionArgs.Tags)
	setBoolPtrFromPtr(&dist.WaitForDeployment, distributionArgs.WaitForDeployment)
	setStringPtrFromPtr(&dist.WebAclId, distributionArgs.WebAclId)

	setOrigins(&dist, distributionArgs.Origins, additionalOrigin)
	setDefaultCacheBehavior(&dist, distributionArgs.DefaultCacheBehavior, additionalOrigin)
	setViewerCertificate(&dist, distributionArgs.ViewerCertificate)
	setGeoRestriction(&dist, distributionArgs.Restrictions.GeoRestriction)
	setCustomErrorResponses(&dist, distributionArgs.CustomErrorResponses)
	setOriginGroups(&dist, distributionArgs.OriginGroups)
	setOrderedCacheBehaviors(&dist, distributionArgs.OrderedCacheBehaviors)

	distributionArgs.DistributionArgs = dist

	return distributionArgs
}

func setAliases(dist *cloudfront.DistributionArgs, aliases []string) {
	if aliases != nil {
		dist.Aliases = pulumi.ToStringArray(aliases)
	}
}

func setLoggingConfig(dist *cloudfront.DistributionArgs, loggingConfig *cloudfront.DistributionLoggingConfig) {
	if loggingConfig != nil {
		logging := cloudfront.DistributionLoggingConfigArgs{
			Bucket: pulumi.String(loggingConfig.Bucket),
		}
		setBoolPtrFromPtr(&logging.IncludeCookies, loggingConfig.IncludeCookies)
		setStringPtrFromPtr(&logging.Prefix, loggingConfig.Prefix)
		dist.LoggingConfig = cloudfront.DistributionLoggingConfigPtr(&logging)
	}
}

func setDefaultCacheBehavior(dist *cloudfront.DistributionArgs, defaultCacheBehavior cloudfront.DistributionDefaultCacheBehavior, additionalOrigin []cloudfront.DistributionOriginInput) {

	cacheBehavior := cloudfront.DistributionDefaultCacheBehaviorArgs{
		ViewerProtocolPolicy: pulumi.String(defaultCacheBehavior.ViewerProtocolPolicy),
		TargetOriginId:       pulumi.String(defaultCacheBehavior.TargetOriginId),
	}

	if len(additionalOrigin) > 0 {
		cacheBehavior.TargetOriginId = additionalOrigin[0].ToDistributionOriginOutput().OriginId().ToStringOutput()
	}
	setStringArrayFrom(&cacheBehavior.AllowedMethods, defaultCacheBehavior.AllowedMethods)
	setStringArrayFrom(&cacheBehavior.CachedMethods, defaultCacheBehavior.CachedMethods)
	setBoolPtrFromPtr(&cacheBehavior.Compress, defaultCacheBehavior.Compress)
	setIntPtrFromPtr(&cacheBehavior.DefaultTtl, defaultCacheBehavior.DefaultTtl)
	setStringPtrFromPtr(&cacheBehavior.CachePolicyId, defaultCacheBehavior.CachePolicyId)
	setIntPtrFromPtr(&cacheBehavior.MinTtl, defaultCacheBehavior.MinTtl)
	setStringArrayFrom(&cacheBehavior.TrustedSigners, defaultCacheBehavior.TrustedSigners)
	setIntPtrFromPtr(&cacheBehavior.MaxTtl, defaultCacheBehavior.MaxTtl)
	setStringPtrFromPtr(&cacheBehavior.RealtimeLogConfigArn, defaultCacheBehavior.RealtimeLogConfigArn)
	setStringPtrFromPtr(&cacheBehavior.ResponseHeadersPolicyId, defaultCacheBehavior.ResponseHeadersPolicyId)
	setBoolPtrFromPtr(&cacheBehavior.SmoothStreaming, defaultCacheBehavior.SmoothStreaming)
	setStringArrayFrom(&cacheBehavior.TrustedKeyGroups, defaultCacheBehavior.TrustedKeyGroups)

	if defaultCacheBehavior.LambdaFunctionAssociations != nil {
		cacheBehavior.LambdaFunctionAssociations = getCacheBehaviorLambdaFunctionAssociations(defaultCacheBehavior.LambdaFunctionAssociations)
	}

	if defaultCacheBehavior.FunctionAssociations != nil {
		cacheBehavior.FunctionAssociations = getCacheBehaviorFunctionAssociations(defaultCacheBehavior.FunctionAssociations)
	}

	setStringPtrFromPtr(&cacheBehavior.OriginRequestPolicyId, defaultCacheBehavior.OriginRequestPolicyId)
	setStringPtrFromPtr(&cacheBehavior.FieldLevelEncryptionId, defaultCacheBehavior.FieldLevelEncryptionId)

	dist.DefaultCacheBehavior = cacheBehavior
}

func setViewerCertificate(dist *cloudfront.DistributionArgs, viewerCertificate cloudfront.DistributionViewerCertificate) {
	certificate := cloudfront.DistributionViewerCertificateArgs{}

	setBoolPtrFromPtr(&certificate.CloudfrontDefaultCertificate, viewerCertificate.CloudfrontDefaultCertificate)
	setStringPtrFromPtr(&certificate.AcmCertificateArn, viewerCertificate.AcmCertificateArn)
	setStringPtrFromPtr(&certificate.IamCertificateId, viewerCertificate.IamCertificateId)
	setStringPtrFromPtr(&certificate.MinimumProtocolVersion, viewerCertificate.MinimumProtocolVersion)
	setStringPtrFromPtr(&certificate.SslSupportMethod, viewerCertificate.SslSupportMethod)

	dist.ViewerCertificate = certificate
}

func setGeoRestriction(dist *cloudfront.DistributionArgs, geoRestriction cloudfront.DistributionRestrictionsGeoRestriction) {
	geo := cloudfront.DistributionRestrictionsGeoRestrictionArgs{
		RestrictionType: pulumi.String(geoRestriction.RestrictionType),
	}

	if geoRestriction.Locations != nil {
		geo.Locations = pulumi.ToStringArray(geoRestriction.Locations)
	}

	dist.Restrictions = cloudfront.DistributionRestrictionsArgs{
		GeoRestriction: geo,
	}
}

func setCustomErrorResponses(dist *cloudfront.DistributionArgs, customErrorResponses []cloudfront.DistributionCustomErrorResponse) {
	if customErrorResponses == nil {
		return
	}

	responses := make([]cloudfront.DistributionCustomErrorResponseInput, len(customErrorResponses))
	for i, customErrorResponse := range customErrorResponses {
		response := cloudfront.DistributionCustomErrorResponseArgs{
			ErrorCode: pulumi.Int(customErrorResponse.ErrorCode),
		}
		setIntPtrFromPtr(&response.ErrorCachingMinTtl, customErrorResponse.ErrorCachingMinTtl)
		setIntPtrFromPtr(&response.ResponseCode, customErrorResponse.ResponseCode)
		setStringPtrFromPtr(&response.ResponsePagePath, customErrorResponse.ResponsePagePath)
		responses[i] = response
	}

	dist.CustomErrorResponses = cloudfront.DistributionCustomErrorResponseArray(responses)
}

func setOrigins(dist *cloudfront.DistributionArgs, origins []cloudfront.DistributionOrigin, additionalOrigin []cloudfront.DistributionOriginInput) {
	if origins == nil && additionalOrigin == nil {
		return
	}

	allOrigins := make([]cloudfront.DistributionOriginInput, 0)
	for _, origin := range origins {
		originArgs := cloudfront.DistributionOriginArgs{
			DomainName:         pulumi.String(origin.DomainName),
			OriginId:           pulumi.String(origin.OriginId),
			ConnectionAttempts: pulumi.IntPtrFromPtr(origin.ConnectionAttempts),
			ConnectionTimeout:  pulumi.IntPtrFromPtr(origin.ConnectionTimeout),
		}
		setStringPtrFromPtr(&originArgs.OriginAccessControlId, origin.OriginAccessControlId)
		setStringPtrFromPtr(&originArgs.OriginPath, origin.OriginPath)

		if origin.CustomOriginConfig != nil {
			config := cloudfront.DistributionOriginCustomOriginConfigArgs{
				HttpPort:             pulumi.Int(origin.CustomOriginConfig.HttpPort),
				HttpsPort:            pulumi.Int(origin.CustomOriginConfig.HttpsPort),
				OriginProtocolPolicy: pulumi.String(origin.CustomOriginConfig.OriginProtocolPolicy),
			}
			setStringArrayFrom(&config.OriginSslProtocols, origin.CustomOriginConfig.OriginSslProtocols)
			setIntPtrFromPtr(&config.OriginReadTimeout, origin.CustomOriginConfig.OriginReadTimeout)
			setIntPtrFromPtr(&config.OriginKeepaliveTimeout, origin.CustomOriginConfig.OriginKeepaliveTimeout)
			originArgs.CustomOriginConfig = cloudfront.DistributionOriginCustomOriginConfigPtr(&config)
		}

		allOrigins = append(allOrigins, originArgs)
	}

	if additionalOrigin != nil {
		allOrigins = append(allOrigins, additionalOrigin...)
	}

	dist.Origins = cloudfront.DistributionOriginArray(allOrigins)

}

func setOriginGroups(dist *cloudfront.DistributionArgs, originGroups []cloudfront.DistributionOriginGroup) {
	if originGroups == nil {
		return
	}

	groupArray := make([]cloudfront.DistributionOriginGroupInput, len(originGroups))
	for i, originGroup := range originGroups {
		groupArgs := cloudfront.DistributionOriginGroupArgs{
			OriginId: pulumi.String(originGroup.OriginId),
			FailoverCriteria: cloudfront.DistributionOriginGroupFailoverCriteriaArgs{
				StatusCodes: pulumi.ToIntArray(originGroup.FailoverCriteria.StatusCodes),
			},
		}

		if originGroup.Members != nil {
			memberArray := make([]cloudfront.DistributionOriginGroupMemberInput, len(originGroup.Members))
			for j, member := range originGroup.Members {
				memberArgs := cloudfront.DistributionOriginGroupMemberArgs{
					OriginId: pulumi.String(member.OriginId),
				}
				memberArray[j] = memberArgs
			}
			groupArgs.Members = cloudfront.DistributionOriginGroupMemberArray(memberArray)
		}

		groupArray[i] = groupArgs
	}

	dist.OriginGroups = cloudfront.DistributionOriginGroupArray(groupArray)
}

func setOrderedCacheBehaviors(dist *cloudfront.DistributionArgs, orderedCacheBehaviors []cloudfront.DistributionOrderedCacheBehavior) {
	if orderedCacheBehaviors == nil {
		return
	}

	behaviorArray := make([]cloudfront.DistributionOrderedCacheBehaviorInput, len(orderedCacheBehaviors))
	for i, orderedCacheBehavior := range orderedCacheBehaviors {
		behaviorArgs := cloudfront.DistributionOrderedCacheBehaviorArgs{
			TargetOriginId:       pulumi.String(orderedCacheBehavior.TargetOriginId),
			ViewerProtocolPolicy: pulumi.String(orderedCacheBehavior.ViewerProtocolPolicy),
		}

		setStringArrayFrom(&behaviorArgs.AllowedMethods, orderedCacheBehavior.AllowedMethods)
		setStringArrayFrom(&behaviorArgs.CachedMethods, orderedCacheBehavior.CachedMethods)
		setBoolPtrFromPtr(&behaviorArgs.Compress, orderedCacheBehavior.Compress)
		setIntPtrFromPtr(&behaviorArgs.DefaultTtl, orderedCacheBehavior.DefaultTtl)
		setIntPtrFromPtr(&behaviorArgs.MaxTtl, orderedCacheBehavior.MaxTtl)
		setIntPtrFromPtr(&behaviorArgs.MinTtl, orderedCacheBehavior.MinTtl)
		setBoolPtrFromPtr(&behaviorArgs.SmoothStreaming, orderedCacheBehavior.SmoothStreaming)
		setStringArrayFrom(&behaviorArgs.TrustedSigners, orderedCacheBehavior.TrustedSigners)
		setStringArrayFrom(&behaviorArgs.TrustedKeyGroups, orderedCacheBehavior.TrustedKeyGroups)

		if orderedCacheBehavior.LambdaFunctionAssociations != nil {
			behaviorArgs.LambdaFunctionAssociations = getOrderedCacheLambdaFunctionAssociations(orderedCacheBehavior.LambdaFunctionAssociations)
		}

		if orderedCacheBehavior.FunctionAssociations != nil {
			behaviorArgs.FunctionAssociations = getOrderedFunctionAssociations(orderedCacheBehavior.FunctionAssociations)
		}

		setStringPtrFromPtr(&behaviorArgs.OriginRequestPolicyId, orderedCacheBehavior.OriginRequestPolicyId)
		setStringPtrFromPtr(&behaviorArgs.FieldLevelEncryptionId, orderedCacheBehavior.FieldLevelEncryptionId)

		behaviorArray[i] = behaviorArgs
	}

	dist.OrderedCacheBehaviors = cloudfront.DistributionOrderedCacheBehaviorArray(behaviorArray)
}

func getOrderedCacheLambdaFunctionAssociations(associations []cloudfront.DistributionOrderedCacheBehaviorLambdaFunctionAssociation) cloudfront.DistributionOrderedCacheBehaviorLambdaFunctionAssociationArray {
	lambdaFunctionAssociations := make([]cloudfront.DistributionOrderedCacheBehaviorLambdaFunctionAssociationInput, len(associations))
	for i, association := range associations {
		lambdaFunctionAssociation := cloudfront.DistributionOrderedCacheBehaviorLambdaFunctionAssociationArgs{
			EventType: pulumi.String(association.EventType),
			LambdaArn: pulumi.String(association.LambdaArn),
		}
		setBoolPtrFromPtr(&lambdaFunctionAssociation.IncludeBody, association.IncludeBody)

		lambdaFunctionAssociations[i] = lambdaFunctionAssociation
	}
	return cloudfront.DistributionOrderedCacheBehaviorLambdaFunctionAssociationArray(lambdaFunctionAssociations)
}

func getOrderedFunctionAssociations(associations []cloudfront.DistributionOrderedCacheBehaviorFunctionAssociation) cloudfront.DistributionOrderedCacheBehaviorFunctionAssociationArray {
	functionAssociations := make([]cloudfront.DistributionOrderedCacheBehaviorFunctionAssociationInput, len(associations))
	for i, association := range associations {
		functionAssociation := cloudfront.DistributionOrderedCacheBehaviorFunctionAssociationArgs{
			EventType:   pulumi.String(association.EventType),
			FunctionArn: pulumi.String(association.FunctionArn),
		}
		functionAssociations[i] = functionAssociation
	}
	return cloudfront.DistributionOrderedCacheBehaviorFunctionAssociationArray(functionAssociations)
}

func getCacheBehaviorLambdaFunctionAssociations(associations []cloudfront.DistributionDefaultCacheBehaviorLambdaFunctionAssociation) cloudfront.DistributionDefaultCacheBehaviorLambdaFunctionAssociationArray {
	lambdaFunctionAssociations := make([]cloudfront.DistributionDefaultCacheBehaviorLambdaFunctionAssociationInput, len(associations))
	for i, association := range associations {
		lambdaFunctionAssociation := cloudfront.DistributionDefaultCacheBehaviorLambdaFunctionAssociationArgs{
			EventType: pulumi.String(association.EventType),
			LambdaArn: pulumi.String(association.LambdaArn),
		}
		setBoolPtrFromPtr(&lambdaFunctionAssociation.IncludeBody, association.IncludeBody)
		lambdaFunctionAssociations[i] = lambdaFunctionAssociation
	}
	return cloudfront.DistributionDefaultCacheBehaviorLambdaFunctionAssociationArray(lambdaFunctionAssociations)
}

func getCacheBehaviorFunctionAssociations(associations []cloudfront.DistributionDefaultCacheBehaviorFunctionAssociation) cloudfront.DistributionDefaultCacheBehaviorFunctionAssociationArray {
	functionAssociations := make([]cloudfront.DistributionDefaultCacheBehaviorFunctionAssociationInput, len(associations))
	for i, association := range associations {
		functionAssociation := cloudfront.DistributionDefaultCacheBehaviorFunctionAssociationArgs{
			EventType:   pulumi.String(association.EventType),
			FunctionArn: pulumi.String(association.FunctionArn),
		}
		functionAssociations[i] = functionAssociation
	}
	return cloudfront.DistributionDefaultCacheBehaviorFunctionAssociationArray(functionAssociations)
}
