package main

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

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

// Extract domain from a given url
// e.g. https://my-website.example.com -> example.com
func extractDomain(domainUrl string) string {
	domainUrl = strings.TrimSpace(domainUrl)
	if !regexp.MustCompile(`^https?`).MatchString(domainUrl) {
		domainUrl = "https://" + domainUrl
	}
	url, _ := url.Parse(domainUrl)
	parts := strings.Split(url.Hostname(), ".")
	return parts[len(parts)-2] + "." + parts[len(parts)-1]
}
