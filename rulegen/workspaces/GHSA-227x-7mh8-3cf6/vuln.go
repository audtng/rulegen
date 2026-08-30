package main

)

// ValidateControlPlaneConfig validates a ControlPlaneConfig object.
func ValidateControlPlaneConfig(controlPlaneConfig *apisaws.ControlPlaneConfig, version string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if controlPlaneConfig.CloudControllerManager != nil {
		allErrs = append(allErrs, featurevalidation.ValidateFeatureGates(controlPlaneConfig.CloudControllerManager.FeatureGates, version, fldPath.Child("cloudControllerManager", "featureGates"))...)
	}

	return allErrs
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"k8s.io/apimachinery/pkg/util/validation/field"

	apisaws "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws"
	. "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws/validation"
				})),
			))
		})
	})
})

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	apisaws "github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws"
)

// valid values for networks.vpc.gatewayEndpoints
var gatewayEndpointPattern = regexp.MustCompile(`^[a-zA-Z0-9\-]+(\.[a-zA-Z0-9\-]+)*$`)

// ValidateInfrastructureConfigAgainstCloudProfile validates the given `InfrastructureConfig` against the given `CloudProfile`.
func ValidateInfrastructureConfigAgainstCloudProfile(oldInfra, infra *apisaws.InfrastructureConfig, shoot *core.Shoot, cloudProfileSpec *gardencorev1beta1.CloudProfileSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(infra.Networks.VPC.GatewayEndpoints) > 0 {
		epsPath := networksPath.Child("vpc", "gatewayEndpoints")
		for i, svc := range infra.Networks.VPC.GatewayEndpoints {
			if !gatewayEndpointPattern.MatchString(svc) {
				allErrs = append(allErrs, field.Invalid(epsPath.Index(i), svc, "must be a valid domain name"))
			}
		}
	}

	for i, zone := range infra.Networks.Zones {
		zonePath := networksPath.Child("zones").Index(i)

		publicPath := zonePath.Child("public")
		cidrs = append(cidrs, cidrvalidation.NewCIDR(zone.Public, publicPath))
		allErrs = append(allErrs, cidrvalidation.ValidateCIDRIsCanonical(publicPath, zone.Public)...)
			}
			referencedElasticIPAllocationIDs = append(referencedElasticIPAllocationIDs, *zone.ElasticIPAllocationID)

			if !strings.HasPrefix(*zone.ElasticIPAllocationID, "eipalloc-") {
				allErrs = append(allErrs, field.Invalid(zonePath.Child("elasticIPAllocationID"), *zone.ElasticIPAllocationID, "must start with eipalloc-"))
			}
		}
	}

		allErrs = append(allErrs, nodes.ValidateSubset(workerCIDRs...)...)
	}

	if (infra.Networks.VPC.ID == nil && infra.Networks.VPC.CIDR == nil) || (infra.Networks.VPC.ID != nil && infra.Networks.VPC.CIDR != nil) {
		allErrs = append(allErrs, field.Invalid(networksPath.Child("vpc"), infra.Networks.VPC, "must specify either a vpc id or a cidr"))
	} else if infra.Networks.VPC.CIDR != nil && infra.Networks.VPC.ID == nil && !slices.Contains(ipFamilies, core.IPFamilyIPv6) {
		cidrPath := networksPath.Child("vpc", "cidr")
		vpcCIDR := cidrvalidation.NewCIDR(*infra.Networks.VPC.CIDR, cidrPath)
		allErrs = append(allErrs, cidrvalidation.ValidateCIDRIsCanonical(cidrPath, *infra.Networks.VPC.CIDR)...)
		allErrs = append(allErrs, vpcCIDR.ValidateSubset(nodes)...)
		allErrs = append(allErrs, vpcCIDR.ValidateSubset(cidrs...)...)
		allErrs = append(allErrs, vpcCIDR.ValidateNotOverlap(pods, services)...)
	}

	// make sure that VPC cidrs don't overlap with each other
	keysPath := fldPath.Child("keys")
	for i, key := range ignoreTags.Keys {
		idxPath := keysPath.Index(i)
		if key == "" {
			allErrs = append(allErrs, field.Invalid(idxPath, key, "ignored key must not be empty"))
			continue
		}
		allErrs = append(allErrs, validateKeyIsReserved(idxPath, key)...)
	prefixesPath := fldPath.Child("keyPrefixes")
	for i, prefix := range ignoreTags.KeyPrefixes {
		idxPath := prefixesPath.Index(i)
		if prefix == "" {
			allErrs = append(allErrs, field.Invalid(idxPath, prefix, "ignored key prefix must not be empty"))
			continue
		}
		allErrs = append(allErrs, validatePrefixIncludesReservedKey(idxPath, prefix)...)
package validation_test

import (
	"github.com/gardener/gardener/pkg/apis/core"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	. "github.com/gardener/gardener/pkg/utils/test/matchers"
		pods        = "100.96.0.0/11"
		services    = "100.64.0.0/13"
		nodes       = "10.250.0.0/16"
		vpc         = "10.0.0.0/8"
		invalidCIDR = "invalid-cidr"
		zone        = "zone1"
		zone2       = "zone2"

		awsZone2 = apisaws.Zone{
			Name:     zone2,
		infrastructureConfig = &apisaws.InfrastructureConfig{
			Networks: apisaws.Networks{
				VPC: apisaws.VPC{
					CIDR: &vpc,
				},
				Zones: []apisaws.Zone{
					{
	})

	Describe("#ValidateInfrastructureConfig", func() {
		Context("Zones", func() {
			It("should forbid empty zones", func() {
				infrastructureConfig.Networks.Zones = nil

				infrastructureConfig.Networks.Zones[0].ElasticIPAllocationID = ptr.To("foo")
				errorList := ValidateInfrastructureConfig(infrastructureConfig, familyIPv4, &nodes, &pods, &services)
				Expect(errorList).To(ConsistOfFields(Fields{
					"Type":  Equal(field.ErrorTypeInvalid),
					"Field": Equal("networks.zones[0].elasticIPAllocationID"),
				}))

				infrastructureConfig.Networks.Zones[0].ElasticIPAllocationID = ptr.To("eipalloc-123456")
			})
		})

		Context("gatewayEndpoints", func() {
			It("should accept empty list", func() {
				errorList := ValidateInfrastructureConfig(infrastructureConfig, familyIPv4, &nodes, &pods, &services)
				Expect(errorList).To(BeEmpty())
			})

			It("should reject non-domain name endpoints", func() {
				infrastructureConfig.Networks.VPC.GatewayEndpoints = []string{"s3", "my_endpoint", "guardduty-data"}
				errorList := ValidateInfrastructureConfig(infrastructureConfig, familyIPv4, &nodes, &pods, &services)
				Expect(errorList).To(ConsistOfFields(Fields{
					"Type":     Equal(field.ErrorTypeInvalid),
					"Field":    Equal("networks.vpc.gatewayEndpoints[1]"),
					"BadValue": Equal("my_endpoint"),
					"Detail":   Equal("must be a valid domain name"),
				}))
			})

			It("should accept all-valid lists", func() {
				infrastructureConfig.Networks.VPC.GatewayEndpoints = []string{"myservice", "s3", "my.other.service"}
				errorList := ValidateInfrastructureConfig(infrastructureConfig, familyIPv4, &nodes, &pods, &services)
				Expect(errorList).To(BeEmpty())
			})
		})

		Context("ignoreTags", func() {
			It("should forbid ignoring reserved tags", func() {
				infrastructureConfig.IgnoreTags = &apisaws.IgnoreTags{
			})
			Expect(errorList).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeInvalid),
					"Field": Equal("ignoreTags.keys[1]"),
				})),
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeInvalid),
					"Field": Equal("ignoreTags.keyPrefixes[1]"),
				})),
			))
		})

		It("should forbid ignoring Name tag", func() {
			errorList := ValidateIgnoreTags(fldPath, &apisaws.IgnoreTags{
				Keys:        []string{"Name"},
		} else {
			dataVolumeConfigNames.Insert(dv.Name)
		}
	}

	if iam := workerConfig.IAMInstanceProfile; iam != nil {
		if (iam.Name == nil && iam.ARN == nil) || (iam.Name != nil && iam.ARN != nil) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("iamInstanceProfile"), iam, "either <name> or <arn> must be provided"))
		}
		if iam.Name != nil && len(*iam.Name) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("iamInstanceProfile", "name"), "name must not be empty"))
		}
		if iam.ARN != nil && len(*iam.ARN) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("iamInstanceProfile", "arn"), "arn must not be empty"))
		}
	}

package validation_test

import (
	"github.com/gardener/gardener/pkg/apis/core"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
			nodeTemplate    *extensionsv1alpha1.NodeTemplate

			iamInstanceProfileName = "name"
			iamInstanceProfileARN  = "arn"

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

		Context("iamInstanceProfile", func() {
			It("should prevent not specifying both IAM name and arn", func() {
				worker.IAMInstanceProfile = &apisaws.IAMInstanceProfile{}
				errorList := ValidateWorkerConfig(worker, rootVolumeIO1, dataVolumes, fldPath)

				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeInvalid),
					"Field": Equal("config.iamInstanceProfile"),
				}))))
			})

				errorList := ValidateWorkerConfig(worker, rootVolumeIO1, dataVolumes, fldPath)

				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeInvalid),
					"Field": Equal("config.iamInstanceProfile"),
				}))))
			})

			It("should forbid specifying an invalid IAM name", func() {
				worker.IAMInstanceProfile = &apisaws.IAMInstanceProfile{
					Name: ptr.To(""),
				}

				errorList := ValidateWorkerConfig(worker, rootVolumeIO1, dataVolumes, fldPath)

				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("config.iamInstanceProfile.name"),
				}))))
			})

			It("should forbid specifying an invalid IAM arn", func() {
				worker.IAMInstanceProfile = &apisaws.IAMInstanceProfile{
					ARN: ptr.To(""),
				}

				errorList := ValidateWorkerConfig(worker, rootVolumeIO1, dataVolumes, fldPath)

				Expect(errorList).To(ConsistOf(PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeRequired),
					"Field": Equal("config.iamInstanceProfile.arn"),
				}))))
			})

