package main

	// initiate connections to 10.2.3.0/24 except from IPs in subnet 10.2.3.0/28.
	//
	// +kubebuilder:validation:Optional
	ToCIDRSet CIDRRuleSlice `json:"toCIDRSet,omitzero"`

	// ToEntities is a list of special entities to which the endpoint subject
	// to the rule is allowed to initiate connections. Supported entities are
	if err != nil {
		return &EgressRule{}, err
	}
	newRule.ToCIDRSet = append(newRule.ToCIDRSet, cidrSet...)
	newRule.ToGroups = nil
	return newRule, nil
}
	if err != nil {
		return &EgressDenyRule{}, err
	}
	newRule.ToCIDRSet = append(newRule.ToCIDRSet, cidrSet...)
	newRule.ToGroups = nil
	return newRule, nil
}

import (
	"context"
	"encoding/json"
	"fmt"
	"net/netip"
	"testing"

	newRule, err := eg.CreateDerivative(context.TODO())
	require.NoError(t, err)
	require.Equal(t, &EgressRule{
		EgressCommonRule: EgressCommonRule{
			ToCIDRSet: CIDRRuleSlice{},
		},
	}, newRule)
}

func TestEgressCommonRuleDeepEqual(t *testing.T) {
		})
	}
}

func TestEgressCommonRuleMarshalling(t *testing.T) {
	testCases := []struct {
		name     string
		in       *EgressCommonRule
		expected string
	}{
		{
			name: "ToCIDRSet is nil",
			in: &EgressCommonRule{
				ToCIDRSet: nil,
			},
			expected: `{}`,
		},
		{
			name: "ToCIDRSet is empty",
			in: &EgressCommonRule{
				ToCIDRSet: []CIDRRule{},
			},
			expected: `{"toCIDRSet":[]}`,
		},
		{
			name: "ToCIDRSet has CIDR",
			in: &EgressCommonRule{
				ToCIDRSet: []CIDRRule{
					{
						Cidr: "192.168.1.0/24",
					},
				},
			},
			expected: `{"toCIDRSet":[{"cidr":"192.168.1.0/24"}]}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.in)
			require.NoError(t, err)
			require.Equal(t, tc.expected, string(data))

			rule := EgressCommonRule{}
			err = json.Unmarshal(data, &rule)
			require.NoError(t, err)
			require.True(t, tc.in.DeepEqual(&rule))
		})
	}
}
	// connections from 10.0.0.0/8 except from IPs in subnet 10.96.0.0/12.
	//
	// +kubebuilder:validation:Optional
	FromCIDRSet CIDRRuleSlice `json:"fromCIDRSet,omitzero"`

	// FromEntities is a list of special entities which the endpoint subject
	// to the rule is allowed to receive connections from. Supported entities are
	if err != nil {
		return &IngressRule{}, err
	}
	newRule.FromCIDRSet = append(newRule.FromCIDRSet, cidrSet...)
	newRule.FromGroups = nil
	return newRule, nil
}
	if err != nil {
		return &IngressDenyRule{}, err
	}
	newRule.FromCIDRSet = append(newRule.FromCIDRSet, cidrSet...)
	newRule.FromGroups = nil
	return newRule, nil
}

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
		})
	}
}

func TestIngressCommonRuleMarshalling(t *testing.T) {
	testCases := []struct {
		name     string
		in       *IngressCommonRule
		expected string
	}{
		{
			name: "ToCIDRSet is nil",
			in: &IngressCommonRule{
				FromCIDRSet: nil,
			},
			expected: `{}`,
		},
		{
			name: "ToCIDRSet is empty",
			in: &IngressCommonRule{
				FromCIDRSet: []CIDRRule{},
			},
			expected: `{"fromCIDRSet":[]}`,
		},
		{
			name: "ToCIDRSet has CIDR",
			in: &IngressCommonRule{
				FromCIDRSet: []CIDRRule{
					{
						Cidr: "192.168.1.0/24",
					},
				},
			},
			expected: `{"fromCIDRSet":[{"cidr":"192.168.1.0/24"}]}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.in)
			require.NoError(t, err)
			require.Equal(t, tc.expected, string(data))

			rule := IngressCommonRule{}
			err = json.Unmarshal(data, &rule)
			require.NoError(t, err)
			require.True(t, tc.in.DeepEqual(&rule))
		})
	}
}
