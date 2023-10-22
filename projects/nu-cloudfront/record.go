package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CreateRecordArgs struct {
	HostedZoneId                     string
	DnsName                          string
	DistributionHostedZoneDomainName pulumi.StringInput
	DistributionHostedZoneId         pulumi.StringInput
}

func LookupHostedZone(ctx *pulumi.Context, hostedZoneName string) (*route53.LookupZoneResult, error) {
	zone, err := route53.LookupZone(ctx, &route53.LookupZoneArgs{
		Name: &hostedZoneName,
	})
	if err != nil {
		return nil, err
	}
	return zone, nil
}

func CreateRecord(ctx *pulumi.Context, prefix string, args CreateRecordArgs) (*route53.Record, error) {
	record, err := route53.NewRecord(ctx, fmt.Sprintf("%s-record-%s", prefix, args.DnsName), &route53.RecordArgs{
		ZoneId: pulumi.String(args.HostedZoneId),
		Name:   pulumi.String(args.DnsName),
		Type:   pulumi.String("A"),
		Aliases: route53.RecordAliasArray{
			&route53.RecordAliasArgs{
				Name:                 args.DistributionHostedZoneDomainName,
				ZoneId:               args.DistributionHostedZoneId,
				EvaluateTargetHealth: pulumi.Bool(false),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return record, nil
}
