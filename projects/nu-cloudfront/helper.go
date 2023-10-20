package main

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

func setStringPtrFromPtr(str *pulumi.StringPtrInput, ptr *string) {
	if ptr != nil {
		*str = pulumi.StringPtrFromPtr(ptr)
	}
}

func setBoolPtrFromPtr(val *pulumi.BoolPtrInput, ptr *bool) {
	if ptr != nil {
		*val = pulumi.BoolPtrFromPtr(ptr)
	}
}

func setStringArrayFrom(array *pulumi.StringArrayInput, val []string) {
	if val != nil {
		*array = pulumi.ToStringArray(val)
	}
}

func setIntPtrFromPtr(ptr *pulumi.IntPtrInput, val *int) {
	if val != nil {
		*ptr = pulumi.IntPtrFromPtr(val)
	}
}
