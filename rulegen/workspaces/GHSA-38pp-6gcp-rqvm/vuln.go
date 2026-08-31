package main

	// initiate connections to 10.2.3.0/24 except from IPs in subnet 10.2.3.0/28.
	//
	// +kubebuilder:validation:Optional
	ToCIDRSet CIDRRuleSlice `json:"toCIDRSet,omitempty"`

	// ToEntities is a list of special entities to which the endpoint subject
	// to the rule is allowed to initiate connections. Supported entities are
	if err != nil {
		return &EgressRule{}, err
	}
	newRule.ToCIDRSet = append(e.ToCIDRSet, cidrSet...)
	newRule.ToGroups = nil
	return newRule, nil
}
	if err != nil {
		return &EgressDenyRule{}, err
	}
	newRule.ToCIDRSet = append(e.ToCIDRSet, cidrSet...)
	newRule.ToGroups = nil
	return newRule, nil
}

import (
	"context"
	"fmt"
	"net/netip"
	"testing"

	newRule, err := eg.CreateDerivative(context.TODO())
	require.NoError(t, err)
	require.Equal(t, &EgressRule{}, newRule)
}

func TestEgressCommonRuleDeepEqual(t *testing.T) {
		})
	}
}
	// connections from 10.0.0.0/8 except from IPs in subnet 10.96.0.0/12.
	//
	// +kubebuilder:validation:Optional
	FromCIDRSet CIDRRuleSlice `json:"fromCIDRSet,omitempty"`

	// FromEntities is a list of special entities which the endpoint subject
	// to the rule is allowed to receive connections from. Supported entities are
	if err != nil {
		return &IngressRule{}, err
	}
	newRule.FromCIDRSet = append(e.FromCIDRSet, cidrSet...)
	newRule.FromGroups = nil
	return newRule, nil
}
	if err != nil {
		return &IngressDenyRule{}, err
	}
	newRule.FromCIDRSet = append(e.FromCIDRSet, cidrSet...)
	newRule.FromGroups = nil
	return newRule, nil
}

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
		})
	}
}
