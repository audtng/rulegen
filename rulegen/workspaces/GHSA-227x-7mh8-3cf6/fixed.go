package main

)

// ValidateControlPlaneConfig validates a ControlPlaneConfig object.
func ValidateControlPlaneConfig(cpConfig *apisaws.ControlPlaneConfig, version string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if cpConfig.CloudControllerManager != nil {
		allErrs = append(allErrs, featurevalidation.ValidateFeatureGates(cpConfig.CloudControllerManager.FeatureGates, version, fldPath.Child("cloudControllerManager", "featureGates"))...)
	}

	if cpConfig.LoadBalancerController != nil && cpConfig.LoadBalancerController.IngressClassName != nil {
		ingressClassName := *cpConfig.LoadBalancerController.IngressClassName
		ingressPath := fldPath.Child("loadBalancerController", "ingressClassName")
		allErrs = append(allErrs, validateK8sResourceName(ingressClassName, ingressPath)...)
	}

	return allErrs
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"

	apisaws "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws"
	. "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws/validation"
				})),
			))
		})

		Context("LoadBalancerController", func() {
			It("should pass for valid ingress class name", func() {
				controlPlane.LoadBalancerController = &apisaws.LoadBalancerControllerConfig{
					IngressClassName: ptr.To("valid-ingress-class"),
				}
				Expect(ValidateControlPlaneConfig(controlPlane, "", fldPath)).To(BeEmpty())
			})

			It("should fail for invalid ingress class name", func() {
				controlPlane.LoadBalancerController = &apisaws.LoadBalancerControllerConfig{
					IngressClassName: ptr.To("NoUpperCaseAllowed"),
				}
				errorList := ValidateControlPlaneConfig(controlPlane, "", fldPath)
				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal("loadBalancerController.ingressClassName"),
					"Detail": Equal("does not match expected regex ^[a-z0-9.-]+$"),
				}))))
			})
		})
	})
})
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
package validation_test

import (
	"fmt"

	"github.com/gardener/gardener/pkg/apis/core"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	. "github.com/gardener/gardener/pkg/utils/test/matchers"
		pods        = "100.96.0.0/11"
		services    = "100.64.0.0/13"
		nodes       = "10.250.0.0/16"
		vpcCIDR     = "10.0.0.0/8"
		invalidCIDR = "invalid-cidr"
		zone        = "eu-central-1c"
		zone2       = "us-east-1a"

		awsZone2 = apisaws.Zone{
			Name:     zone2,
		infrastructureConfig = &apisaws.InfrastructureConfig{
			Networks: apisaws.Networks{
				VPC: apisaws.VPC{
					CIDR: &vpcCIDR,
				},
				Zones: []apisaws.Zone{
					{
	})

	Describe("#ValidateInfrastructureConfig", func() {
		Context("VPC", func() {
			Context("ID", func() {
				It("should pass with ID", func() {
					infrastructureConfig.Networks.VPC.ID = ptr.To("vpc-064b5b7771f63317c")
					infrastructureConfig.Networks.VPC.CIDR = nil
					errorList := ValidateInfrastructureConfig(infrastructureConfig, familyIPv4, &nodes, &pods, &services)
					Expect(errorList).To(BeEmpty())
				})

				It("should reject empty ID string", func() {
					infrastructureConfig.Networks.VPC.ID = ptr.To("")
					infrastructureConfig.Networks.VPC.CIDR = nil
					errorList := ValidateInfrastructureConfig(infrastructureConfig, familyIPv4, &nodes, &pods, &services)
					Expect(errorList).To(ConsistOfFields(Fields{
						"Type":   Equal(field.ErrorTypeRequired),
						"Field":  Equal("networks.vpc.id"),
						"Detail": Equal("cannot be empty"),
					}))
				})

				It("should reject setting both ID and CIDR", func() {
					infrastructureConfig.Networks.VPC.ID = ptr.To("vpc-123456")
					infrastructureConfig.Networks.VPC.CIDR = &vpcCIDR
					errorList := ValidateInfrastructureConfig(infrastructureConfig, familyIPv4, &nodes, &pods, &services)
					Expect(errorList).To(ConsistOfFields(Fields{
						"Type":   Equal(field.ErrorTypeInvalid),
						"Field":  Equal("networks.vpc"),
						"Detail": Equal("cannot specify both vpc id and cidr"),
					}))
				})

				It("should reject invalid format", func() {
					infrastructureConfig.Networks.VPC.ID = ptr.To("no-vpc-prefix-1234")
					infrastructureConfig.Networks.VPC.CIDR = nil
					errorList := ValidateInfrastructureConfig(infrastructureConfig, familyIPv4, &nodes, &pods, &services)
					Expect(errorList).To(ConsistOfFields(Fields{
						"Type":   Equal(field.ErrorTypeInvalid),
						"Field":  Equal("networks.vpc.id"),
						"Detail": Equal(fmt.Sprintf("does not match expected regex %s", VpcIDRegex)),
					}))
				})
			})

			Context("gatewayEndpoints", func() {
				It("should accept empty list", func() {
					errorList := ValidateInfrastructureConfig(infrastructureConfig, familyIPv4, &nodes, &pods, &services)
					Expect(errorList).To(BeEmpty())
				})

				It("should reject non-alphanumeric endpoints", func() {
					infrastructureConfig.Networks.VPC.GatewayEndpoints = []string{"s3", "my-endpoint"}
					errorList := ValidateInfrastructureConfig(infrastructureConfig, familyIPv4, &nodes, &pods, &services)
					Expect(errorList).To(ConsistOfFields(Fields{
						"Type":     Equal(field.ErrorTypeInvalid),
						"Field":    Equal("networks.vpc.gatewayEndpoints[1]"),
						"BadValue": Equal("my-endpoint"),
						"Detail":   Equal(fmt.Sprintf("does not match expected regex %s", GatewayEndpointRegex)),
					}))
				})

				It("should accept all-valid lists", func() {
					infrastructureConfig.Networks.VPC.GatewayEndpoints = []string{"myservice", "s3", "my.other.service"}
					errorList := ValidateInfrastructureConfig(infrastructureConfig, familyIPv4, &nodes, &pods, &services)
					Expect(errorList).To(BeEmpty())
				})
			})
		})

		Context("Zones", func() {
			It("should reject invalid zone name", func() {
				infrastructureConfig.Networks.Zones[0].Name = "US-East-1a"

				errorList := ValidateInfrastructureConfig(infrastructureConfig, familyIPv4, &nodes, &pods, &services)

				Expect(errorList).To(ConsistOfFields(Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal("networks.zones[0].name"),
					"Detail": Equal(fmt.Sprintf("does not match expected regex %s", ZoneNameRegex)),
				}))
			})

			It("should forbid empty zones", func() {
				infrastructureConfig.Networks.Zones = nil

				infrastructureConfig.Networks.Zones[0].ElasticIPAllocationID = ptr.To("foo")
				errorList := ValidateInfrastructureConfig(infrastructureConfig, familyIPv4, &nodes, &pods, &services)
				Expect(errorList).To(ConsistOfFields(Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal("networks.zones[0].elasticIPAllocationID"),
					"Detail": Equal(fmt.Sprintf("does not match expected regex %s", EipAllocationIDRegex)),
				}))

				infrastructureConfig.Networks.Zones[0].ElasticIPAllocationID = ptr.To("eipalloc-123456")
			})
		})

		Context("ignoreTags", func() {
			It("should forbid ignoring reserved tags", func() {
				infrastructureConfig.IgnoreTags = &apisaws.IgnoreTags{
			})
			Expect(errorList).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("ignoreTags.keys[1]"),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("ignoreTags.keyPrefixes[1]"),
				})),
			))
		})

		It("should forbid invalid values", func() {
			errorList := ValidateIgnoreTags(fldPath, &apisaws.IgnoreTags{
				Keys: []string{"notAllowedChar{}"},
			})
			Expect(errorList).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal("ignoreTags.keys[0]"),
					"Detail": Equal(fmt.Sprintf("does not match expected regex %s", TagKeyRegex)),
				})),
			))
		})

		It("should forbid ignoring Name tag", func() {
			errorList := ValidateIgnoreTags(fldPath, &apisaws.IgnoreTags{
				Keys:        []string{"Name"},
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

package validation_test

import (
	"fmt"

	"github.com/gardener/gardener/pkg/apis/core"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
			nodeTemplate    *extensionsv1alpha1.NodeTemplate

			iamInstanceProfileName = "name"
			iamInstanceProfileARN  = "arn:aws:iam::123456789012:instance-profile/path/to/profile-name"

			worker  *apisaws.WorkerConfig
			fldPath = field.NewPath("config")
				})),
			))
		})

		It("should enforce that the IOPS is positive", func() {
			var negative int64 = -100
			worker.Volume.IOPS = &negative
				})),
			))
		})

		It("should prevent duplicate entries for data volumes in workerconfig", func() {
			worker.DataVolumes = append(worker.DataVolumes, apisaws.DataVolume{Name: dataVolume1Name})

				"Field": Equal("config.dataVolumes[1].name"),
			}))))
		})

		It("should enforce that the throughput is positive", func() {
			var negative int64 = -100
			worker.Volume.Throughput = &negative
				})),
			))
		})

		It("should prevent data volume entries in workerconfig for non-existing data volumes shoot", func() {
			worker.DataVolumes = append(worker.DataVolumes, apisaws.DataVolume{Name: "broken"})

			}))))
		})

		It("should reject invalid snapshot ID", func() {
			worker.DataVolumes[0].SnapshotID = ptr.To("must-start-with-snap")

			errorList := ValidateWorkerConfig(worker, rootVolumeIO1, dataVolumes, fldPath)

			Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
				"Type":   Equal(field.ErrorTypeInvalid),
				"Field":  Equal("config.dataVolumes[0].snapshotID"),
				"Detail": Equal(fmt.Sprintf("does not match expected regex %s", SnapshotIDRegex)),
			}))))
		})

		Context("iamInstanceProfile", func() {
			It("should prevent not specifying both IAM name and arn", func() {
				worker.IAMInstanceProfile = &apisaws.IAMInstanceProfile{}
				errorList := ValidateWorkerConfig(worker, rootVolumeIO1, dataVolumes, fldPath)

				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal("config.iamInstanceProfile"),
					"Detail": Equal("exactly one of 'name' or 'arn' must be specified"),
				}))))
			})

				errorList := ValidateWorkerConfig(worker, rootVolumeIO1, dataVolumes, fldPath)

				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal("config.iamInstanceProfile"),
					"Detail": Equal("exactly one of 'name' or 'arn' must be specified"),
				}))))
			})

			It("should forbid specifying an invalid IAM name", func() {
				worker.IAMInstanceProfile = &apisaws.IAMInstanceProfile{
					Name: ptr.To("invalidChar{"),
				}

				errorList := ValidateWorkerConfig(worker, rootVolumeIO1, dataVolumes, fldPath)

				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal("config.iamInstanceProfile.name"),
					"Detail": Equal(fmt.Sprintf("does not match expected regex %s", IamInstanceProfileNameRegex)),
				}))))
			})

			It("should forbid specifying an invalid IAM arn", func() {
				worker.IAMInstanceProfile = &apisaws.IAMInstanceProfile{
					ARN: ptr.To("must-start-with-arn:aws:iam::"),
				}

				errorList := ValidateWorkerConfig(worker, rootVolumeIO1, dataVolumes, fldPath)

				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":   Equal(field.ErrorTypeInvalid),
					"Field":  Equal("config.iamInstanceProfile.arn"),
					"Detail": Equal(fmt.Sprintf("does not match expected regex %s", IamInstanceProfileArnRegex)),
				}))))
			})

