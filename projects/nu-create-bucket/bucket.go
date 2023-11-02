package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func BuildBucketArg(args BucketArgsInput) s3.BucketArgs {
	return s3.BucketArgs{
		AccelerationStatus:                pulumi.StringPtrFromPtr(args.AccelerationStatus),
		Acl:                               pulumi.StringPtrFromPtr(args.Acl),
		Arn:                               pulumi.StringPtrFromPtr(args.Arn),
		Bucket:                            pulumi.StringPtr(args.Bucket),
		BucketPrefix:                      pulumi.StringPtrFromPtr(args.BucketPrefix),
		CorsRules:                         bucketCorsRule(args.CorsRules),
		ForceDestroy:                      pulumi.BoolPtrFromPtr(args.ForceDestroy),
		Grants:                            bucketGrants(args.Grants),
		HostedZoneId:                      pulumi.StringPtrFromPtr(args.HostedZoneId),
		LifecycleRules:                    bucketLifecyleRules(args.LifecycleRules),
		Loggings:                          bucketLoggins(args.Loggings),
		ObjectLockConfiguration:           bucketObjectLock(args.ObjectLockConfiguration),
		Policy:                            pulumi.Any(args.Policy),
		ReplicationConfiguration:          bucketReplication(args.ReplicationConfiguration),
		RequestPayer:                      pulumi.StringPtrFromPtr(args.RequestPayer),
		ServerSideEncryptionConfiguration: bucketServerSide(args.ServerSideEncryptionConfiguration),
		Tags:                              pulumi.ToStringMap(args.Tags),
		Versioning:                        bucketVersioning(args.Versioning),
		Website:                           bucketWebsite(args.Website),
		WebsiteDomain:                     pulumi.StringPtrFromPtr(args.WebsiteDomain),
		WebsiteEndpoint:                   pulumi.StringPtrFromPtr(args.WebsiteEndpoint),
	}
}

func bucketVersioning(versioning *s3.BucketVersioning) s3.BucketVersioningPtrInput {
	if versioning == nil {
		return nil
	}
	return &s3.BucketVersioningArgs{
		Enabled:   pulumi.BoolPtrFromPtr(versioning.Enabled),
		MfaDelete: pulumi.BoolPtrFromPtr(versioning.MfaDelete),
	}
}

func bucketServerSide(serverSide *s3.BucketServerSideEncryptionConfiguration) s3.BucketServerSideEncryptionConfigurationPtrInput {
	if serverSide == nil {
		return nil
	}
	return &s3.BucketServerSideEncryptionConfigurationArgs{
		Rule: s3.BucketServerSideEncryptionConfigurationRuleArgs{
			ApplyServerSideEncryptionByDefault: s3.BucketServerSideEncryptionConfigurationRuleApplyServerSideEncryptionByDefaultArgs{
				KmsMasterKeyId: pulumi.StringPtrFromPtr(serverSide.Rule.ApplyServerSideEncryptionByDefault.KmsMasterKeyId),
				SseAlgorithm:   pulumi.String(serverSide.Rule.ApplyServerSideEncryptionByDefault.SseAlgorithm),
			},
			BucketKeyEnabled: pulumi.BoolPtrFromPtr(serverSide.Rule.BucketKeyEnabled),
		},
	}
}

func replicationRuleConfiguration(rules []s3.BucketReplicationConfigurationRule) s3.BucketReplicationConfigurationRuleArray {
	array := make([]s3.BucketReplicationConfigurationRuleInput, len(rules))
	for i, rule := range rules {
		ruleArgs := s3.BucketReplicationConfigurationRuleArgs{
			DeleteMarkerReplicationStatus: pulumi.StringPtrFromPtr(rule.DeleteMarkerReplicationStatus),
			Prefix:                        pulumi.StringPtrFromPtr(rule.Prefix),
			Id:                            pulumi.StringPtrFromPtr(rule.Id),
			Priority:                      pulumi.IntPtrFromPtr(rule.Priority),
			Destination: s3.BucketReplicationConfigurationRuleDestinationArgs{
				Bucket: pulumi.String(rule.Destination.Bucket),
				AccessControlTranslation: s3.BucketReplicationConfigurationRuleDestinationAccessControlTranslationArgs{
					Owner: pulumi.String(rule.Destination.AccessControlTranslation.Owner),
				},
				AccountId: pulumi.StringPtrFromPtr(rule.Destination.AccountId),
				Metrics: s3.BucketReplicationConfigurationRuleDestinationMetricsArgs{
					Minutes: pulumi.IntPtrFromPtr(rule.Destination.Metrics.Minutes),
					Status:  pulumi.StringPtrFromPtr(rule.Destination.Metrics.Status),
				},
				ReplicationTime: s3.BucketReplicationConfigurationRuleDestinationReplicationTimeArgs{
					Status:  pulumi.StringPtrFromPtr(rule.Destination.ReplicationTime.Status),
					Minutes: pulumi.IntPtrFromPtr(rule.Destination.ReplicationTime.Minutes),
				},

				StorageClass: pulumi.StringPtrFromPtr(rule.Destination.StorageClass),
			},
		}
		array[i] = ruleArgs
	}
	return s3.BucketReplicationConfigurationRuleArray(array)
}

func bucketReplication(replication *s3.BucketReplicationConfiguration) s3.BucketReplicationConfigurationPtrInput {
	if replication == nil {
		return nil
	}
	return s3.BucketReplicationConfigurationArgs{
		Role:  pulumi.String(replication.Role),
		Rules: replicationRuleConfiguration(replication.Rules),
	}
}

func bucketWebsite(website *s3.BucketWebsite) s3.BucketWebsitePtrInput {
	if website == nil {
		return nil
	}
	return s3.BucketWebsiteArgs{
		ErrorDocument:         pulumi.StringPtrFromPtr(website.ErrorDocument),
		IndexDocument:         pulumi.StringPtrFromPtr(website.IndexDocument),
		RedirectAllRequestsTo: pulumi.StringPtrFromPtr(website.RedirectAllRequestsTo),
		RoutingRules:          pulumi.Any(website.RoutingRules),
	}
}

func objectLockRule(rule *s3.BucketObjectLockConfigurationRule) s3.BucketObjectLockConfigurationRulePtrInput {
	if rule == nil {
		return nil
	}
	return s3.BucketObjectLockConfigurationRuleArgs{
		DefaultRetention: s3.BucketObjectLockConfigurationRuleDefaultRetentionArgs{
			Days:  pulumi.IntPtrFromPtr(rule.DefaultRetention.Days),
			Mode:  pulumi.String(rule.DefaultRetention.Mode),
			Years: pulumi.IntPtrFromPtr(rule.DefaultRetention.Years),
		},
	}
}

func bucketObjectLock(objectLock *s3.BucketObjectLockConfiguration) s3.BucketObjectLockConfigurationPtrInput {
	if objectLock == nil {
		return nil
	}
	return &s3.BucketObjectLockConfigurationArgs{
		ObjectLockEnabled: pulumi.String(objectLock.ObjectLockEnabled),
		Rule:              objectLockRule(objectLock.Rule),
	}
}

func bucketLoggins(loggings []s3.BucketLogging) s3.BucketLoggingArray {
	array := make([]s3.BucketLoggingInput, len(loggings))
	for i, logging := range loggings {
		loggingArgs := s3.BucketLoggingArgs{
			TargetBucket: pulumi.String(logging.TargetBucket),
			TargetPrefix: pulumi.StringPtrFromPtr(logging.TargetPrefix),
		}
		array[i] = loggingArgs
	}
	return s3.BucketLoggingArray(array)
}

func bucketCorsRule(corsRule []s3.BucketCorsRule) s3.BucketCorsRuleArray {
	array := make([]s3.BucketCorsRuleInput, len(corsRule))
	for i, rule := range corsRule {
		ruleArgs := s3.BucketCorsRuleArgs{
			AllowedHeaders: pulumi.ToStringArray(rule.AllowedHeaders),
			AllowedMethods: pulumi.ToStringArray(rule.AllowedMethods),
			AllowedOrigins: pulumi.ToStringArray(rule.AllowedOrigins),
			ExposeHeaders:  pulumi.ToStringArray(rule.ExposeHeaders),
			MaxAgeSeconds:  pulumi.IntPtrFromPtr(rule.MaxAgeSeconds),
		}
		array[i] = ruleArgs
	}
	return s3.BucketCorsRuleArray(array)
}

func bucketGrants(grants []s3.BucketGrant) s3.BucketGrantArray {
	array := make([]s3.BucketGrantInput, len(grants))
	for i, grant := range grants {
		grantArgs := s3.BucketGrantArgs{
			Permissions: pulumi.ToStringArray(grant.Permissions),
			Id:          pulumi.StringPtrFromPtr(grant.Id),
			Type:        pulumi.String(grant.Type),
			Uri:         pulumi.StringPtrFromPtr(grant.Uri),
		}
		array[i] = grantArgs
	}
	return s3.BucketGrantArray(array)
}

func transition(transitions []s3.BucketLifecycleRuleTransition) s3.BucketLifecycleRuleTransitionArray {
	array := make([]s3.BucketLifecycleRuleTransitionInput, len(transitions))
	for i, rule := range transitions {
		ruleArgs := s3.BucketLifecycleRuleTransitionArgs{
			Date:         pulumi.StringPtrFromPtr(rule.Date),
			Days:         pulumi.IntPtrFromPtr(rule.Days),
			StorageClass: pulumi.String(rule.StorageClass),
		}
		array[i] = ruleArgs
	}
	return s3.BucketLifecycleRuleTransitionArray(array)
}

func nonCurrentVersionTransitions(transitions []s3.BucketLifecycleRuleNoncurrentVersionTransition) s3.BucketLifecycleRuleNoncurrentVersionTransitionArray {
	array := make([]s3.BucketLifecycleRuleNoncurrentVersionTransitionInput, len(transitions))
	for i, rule := range transitions {
		ruleArgs := s3.BucketLifecycleRuleNoncurrentVersionTransitionArgs{
			Days:         pulumi.IntPtrFromPtr(rule.Days),
			StorageClass: pulumi.String(rule.StorageClass),
		}
		array[i] = ruleArgs
	}
	return s3.BucketLifecycleRuleNoncurrentVersionTransitionArray(array)
}

func bucketLifecyleRules(lifecycleRules []s3.BucketLifecycleRule) s3.BucketLifecycleRuleArray {
	array := make([]s3.BucketLifecycleRuleInput, len(lifecycleRules))
	for i, rule := range lifecycleRules {
		ruleArgs := s3.BucketLifecycleRuleArgs{
			AbortIncompleteMultipartUploadDays: pulumi.IntPtrFromPtr(rule.AbortIncompleteMultipartUploadDays),
			Enabled:                            pulumi.Bool(rule.Enabled),
			Id:                                 pulumi.StringPtrFromPtr(rule.Id),
			Expiration: s3.BucketLifecycleRuleExpirationArgs{
				Date: pulumi.StringPtrFromPtr(rule.Expiration.Date),
				Days: pulumi.IntPtrFromPtr(rule.Expiration.Days),
			},
			NoncurrentVersionExpiration: s3.BucketLifecycleRuleNoncurrentVersionExpirationArgs{
				Days: pulumi.IntPtrFromPtr(rule.NoncurrentVersionExpiration.Days),
			},
			NoncurrentVersionTransitions: nonCurrentVersionTransitions(rule.NoncurrentVersionTransitions),
			Prefix:                       pulumi.StringPtrFromPtr(rule.Prefix),
			Tags:                         pulumi.ToStringMap(rule.Tags),
			Transitions:                  transition(rule.Transitions),
		}
		array[i] = ruleArgs
	}
	return s3.BucketLifecycleRuleArray(array)
}

func CreateBucket(context *ProgramContext, id string, args s3.BucketArgs) (*s3.Bucket, error) {
	bucket, err := s3.NewBucket(context.pulumiCtx, fmt.Sprintf("%s:bucket:%s", context.prefix, id), &args)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}
