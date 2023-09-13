// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package aws

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Use this data source to get the default tags configured on the provider.
//
// With this data source, you can apply default tags to resources not _directly_ managed by a resource, such as the instances underneath an Auto Scaling group or the volumes created for an EC2 instance.
//
// ## Example Usage
// ### Basic Usage
//
// ```go
// package main
//
// import (
//
//	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//
// )
//
//	func main() {
//		pulumi.Run(func(ctx *pulumi.Context) error {
//			_, err := aws.GetDefaultTags(ctx, nil, nil)
//			if err != nil {
//				return err
//			}
//			return nil
//		})
//	}
//
// ```
func GetDefaultTags(ctx *pulumi.Context, args *GetDefaultTagsArgs, opts ...pulumi.InvokeOption) (*GetDefaultTagsResult, error) {
	var rv GetDefaultTagsResult
	err := ctx.Invoke("aws:index/getDefaultTags:getDefaultTags", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getDefaultTags.
type GetDefaultTagsArgs struct {
	// Blocks of default tags set on the provider. See details below.
	Tags map[string]string `pulumi:"tags"`
}

// A collection of values returned by getDefaultTags.
type GetDefaultTagsResult struct {
	// The provider-assigned unique ID for this managed resource.
	Id string `pulumi:"id"`
	// Blocks of default tags set on the provider. See details below.
	Tags map[string]string `pulumi:"tags"`
}

func GetDefaultTagsOutput(ctx *pulumi.Context, args GetDefaultTagsOutputArgs, opts ...pulumi.InvokeOption) GetDefaultTagsResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (GetDefaultTagsResult, error) {
			args := v.(GetDefaultTagsArgs)
			r, err := GetDefaultTags(ctx, &args, opts...)
			var s GetDefaultTagsResult
			if r != nil {
				s = *r
			}
			return s, err
		}).(GetDefaultTagsResultOutput)
}

// A collection of arguments for invoking getDefaultTags.
type GetDefaultTagsOutputArgs struct {
	// Blocks of default tags set on the provider. See details below.
	Tags pulumi.StringMapInput `pulumi:"tags"`
}

func (GetDefaultTagsOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*GetDefaultTagsArgs)(nil)).Elem()
}

// A collection of values returned by getDefaultTags.
type GetDefaultTagsResultOutput struct{ *pulumi.OutputState }

func (GetDefaultTagsResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*GetDefaultTagsResult)(nil)).Elem()
}

func (o GetDefaultTagsResultOutput) ToGetDefaultTagsResultOutput() GetDefaultTagsResultOutput {
	return o
}

func (o GetDefaultTagsResultOutput) ToGetDefaultTagsResultOutputWithContext(ctx context.Context) GetDefaultTagsResultOutput {
	return o
}

// The provider-assigned unique ID for this managed resource.
func (o GetDefaultTagsResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v GetDefaultTagsResult) string { return v.Id }).(pulumi.StringOutput)
}

// Blocks of default tags set on the provider. See details below.
func (o GetDefaultTagsResultOutput) Tags() pulumi.StringMapOutput {
	return o.ApplyT(func(v GetDefaultTagsResult) map[string]string { return v.Tags }).(pulumi.StringMapOutput)
}

func init() {
	pulumi.RegisterOutputType(GetDefaultTagsResultOutput{})
}
