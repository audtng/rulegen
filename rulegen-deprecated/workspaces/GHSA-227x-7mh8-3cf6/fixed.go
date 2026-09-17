package main

// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

var (
	// k8s resource names have max 253 characters, consist of lower case alphanumeric characters, '-' or '.'
	k8sResourceNameRegex = `^[a-z0-9.-]+$`
	// VpcIDRegex matches e.g. vpc-064b5b7771f6331aa
	VpcIDRegex = `^vpc-[a-z0-9]+$`
	// EipAllocationIDRegex matches e.g. eipalloc-0676786f3e288044c
	EipAllocationIDRegex = `^eipalloc-[a-z0-9]+$`
	// SnapshotIDRegex matches e.g. snap-0676786f3e288044c
	SnapshotIDRegex = `^snap-[a-z0-9]+$`
	// IamInstanceProfileNameRegex matches https://docs.aws.amazon.com/AWSCloudFormation/latest/TemplateReference/aws-resource-iam-instanceprofile.html#:~:text=Properties-,InstanceProfileName,-The%20name%20of
	IamInstanceProfileNameRegex = `^[\w+=,.@-]+$`
	// IamInstanceProfileArnRegex matches arn:aws:iam::<account-id>:instance-profile/<path>/<profile-name>
	// Note: for china landscapes it's arn:aws-cn:iam::<account-id>:instance-profile/<path>/<profile-name>
	IamInstanceProfileArnRegex = `^arn:[\w +=,.@\-/:]+$`
	// ZoneNameRegex matches e.g. us-east-1a
	ZoneNameRegex = `^[a-z0-9-]+$`
	// TagKeyRegex matches Letters (a–z, A–Z), numbers (0–9), spaces, and the following symbols: + - = . _ : / @
	TagKeyRegex = `^[\w +\-=\.:/@]+$`
	// GatewayEndpointRegex matches one or more word characters, optionally followed by dot-separated word segments
	GatewayEndpointRegex = `^\w+(\.\w+)*$`

	validateK8sResourceName        = combineValidationFuncs(regex(k8sResourceNameRegex), notEmpty, maxLength(253))
	validateVpcID                  = combineValidationFuncs(regex(VpcIDRegex), notEmpty, maxLength(255))
	validateEipAllocationID        = combineValidationFuncs(regex(EipAllocationIDRegex), maxLength(255))
	validateSnapshotID             = combineValidationFuncs(regex(SnapshotIDRegex), maxLength(255))
	validateIamInstanceProfileName = combineValidationFuncs(regex(IamInstanceProfileNameRegex), notEmpty, maxLength(128))
	validateIamInstanceProfileArn  = combineValidationFuncs(regex(IamInstanceProfileArnRegex), maxLength(255))
	validateZoneName               = combineValidationFuncs(regex(ZoneNameRegex), maxLength(255))
	validateTagKey                 = combineValidationFuncs(regex(TagKeyRegex), notEmpty, maxLength(128))
	validateGatewayEndpointName    = combineValidationFuncs(regex(GatewayEndpointRegex), maxLength(255))
)

type validateFunc[T any] func(T, *field.Path) field.ErrorList

// combineValidationFuncs validates a value against a list of filters.
func combineValidationFuncs[T any](filters ...validateFunc[T]) validateFunc[T] {
	return func(t T, fld *field.Path) field.ErrorList {
		var allErrs field.ErrorList
		for _, f := range filters {
			allErrs = append(allErrs, f(t, fld)...)
		}
		return allErrs
	}
}

// regex returns a filterFunc that validates a string against a regular expression.
func regex(regex string) validateFunc[string] {
	compiled := regexp.MustCompile(regex)
	return func(name string, fld *field.Path) field.ErrorList {
		var allErrs field.ErrorList
		if name == "" {
			return allErrs // Allow empty strings to pass through
		}
		if !compiled.MatchString(name) {
			allErrs = append(allErrs, field.Invalid(fld, name, fmt.Sprintf("does not match expected regex %s", compiled.String())))
		}
		return allErrs
	}
}

func notEmpty(name string, fld *field.Path) field.ErrorList {
	if utf8.RuneCountInString(name) == 0 {
		return field.ErrorList{field.Required(fld, "cannot be empty")}
	}
	return nil
}

func maxLength(max int) validateFunc[string] {
	return func(name string, fld *field.Path) field.ErrorList {
		var allErrs field.ErrorList
		if l := utf8.RuneCountInString(name); l > max {
			return field.ErrorList{field.Invalid(fld, name, fmt.Sprintf("must not be more than %d characters, got %d", max, l))}
		}
		return allErrs
	}
}

import (
	"fmt"
	"slices"
	"strings"

	apisaws "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws"
)

// ValidateInfrastructureConfigAgainstCloudProfile validates the given `InfrastructureConfig` against the given `CloudProfile`.
func ValidateInfrastructureConfigAgainstCloudProfile(oldInfra, infra *apisaws.InfrastructureConfig, shoot *core.Shoot, cloudProfileSpec *gardencorev1beta1.CloudProfileSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(infra.Networks.VPC.GatewayEndpoints) > 0 {
		epsPath := networksPath.Child("vpc", "gatewayEndpoints")
		for i, svc := range infra.Networks.VPC.GatewayEndpoints {
			allErrs = append(allErrs, validateGatewayEndpointName(svc, epsPath.Index(i))...)
		}
	}

	for i, zone := range infra.Networks.Zones {
		zonePath := networksPath.Child("zones").Index(i)

		allErrs = append(allErrs, validateZoneName(zone.Name, zonePath.Child("name"))...)

		publicPath := zonePath.Child("public")
		cidrs = append(cidrs, cidrvalidation.NewCIDR(zone.Public, publicPath))
		allErrs = append(allErrs, cidrvalidation.ValidateCIDRIsCanonical(publicPath, zone.Public)...)
			}
			referencedElasticIPAllocationIDs = append(referencedElasticIPAllocationIDs, *zone.ElasticIPAllocationID)

			allErrs = append(allErrs, validateEipAllocationID(*zone.ElasticIPAllocationID, zonePath.Child("elasticIPAllocationID"))...)
		}
	}

		allErrs = append(allErrs, nodes.ValidateSubset(workerCIDRs...)...)
	}

	idProvided := infra.Networks.VPC.ID != nil
	cidrProvided := infra.Networks.VPC.CIDR != nil
	switch {
	case !idProvided && !cidrProvided:
		allErrs = append(allErrs, field.Invalid(networksPath.Child("vpc"), infra.Networks.VPC, "must specify either a vpc id or a cidr"))
	case idProvided && cidrProvided:
		allErrs = append(allErrs, field.Invalid(networksPath.Child("vpc"), infra.Networks.VPC, "cannot specify both vpc id and cidr"))
	case cidrProvided && !idProvided && !slices.Contains(ipFamilies, core.IPFamilyIPv6):
		cidrPath := networksPath.Child("vpc", "cidr")
		vpcCIDR := cidrvalidation.NewCIDR(*infra.Networks.VPC.CIDR, cidrPath)
		allErrs = append(allErrs, cidrvalidation.ValidateCIDRIsCanonical(cidrPath, *infra.Networks.VPC.CIDR)...)
		allErrs = append(allErrs, vpcCIDR.ValidateSubset(nodes)...)
		allErrs = append(allErrs, vpcCIDR.ValidateSubset(cidrs...)...)
		allErrs = append(allErrs, vpcCIDR.ValidateNotOverlap(pods, services)...)
	case idProvided && !cidrProvided:
		allErrs = append(allErrs, validateVpcID(*infra.Networks.VPC.ID, networksPath.Child("vpc", "id"))...)
	}

	// make sure that VPC cidrs don't overlap with each other
	keysPath := fldPath.Child("keys")
	for i, key := range ignoreTags.Keys {
		idxPath := keysPath.Index(i)
		if errs := validateTagKey(key, idxPath); errs != nil {
			allErrs = append(allErrs, errs...)
			continue
		}
		allErrs = append(allErrs, validateKeyIsReserved(idxPath, key)...)
	prefixesPath := fldPath.Child("keyPrefixes")
	for i, prefix := range ignoreTags.KeyPrefixes {
		idxPath := prefixesPath.Index(i)
		if errs := validateTagKey(prefix, idxPath); errs != nil {
			allErrs = append(allErrs, errs...)
			continue
		}
		allErrs = append(allErrs, validatePrefixIncludesReservedKey(idxPath, prefix)...)
		} else {
			dataVolumeConfigNames.Insert(dv.Name)
		}

		if id := dv.SnapshotID; id != nil {
			allErrs = append(allErrs, validateSnapshotID(*id, idxPath.Child("snapshotID"))...)
		}
	}

	if iam := workerConfig.IAMInstanceProfile; iam != nil {
		if (iam.Name == nil && iam.ARN == nil) || (iam.Name != nil && iam.ARN != nil) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("iamInstanceProfile"), iam,
				"exactly one of 'name' or 'arn' must be specified"))
		}
		if iam.Name != nil {
			allErrs = append(allErrs, validateIamInstanceProfileName(*iam.Name, fldPath.Child("iamInstanceProfile", "name"))...)
		}
		if iam.ARN != nil {
			allErrs = append(allErrs, validateIamInstanceProfileArn(*iam.ARN, fldPath.Child("iamInstanceProfile", "arn"))...)
		}
	}

