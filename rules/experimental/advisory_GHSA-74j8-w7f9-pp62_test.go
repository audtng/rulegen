package rules

import (
	"context"
	"fmt"
)

type TokenReviewStatus struct {
	Authenticated bool
	User          string
}

type TokenReview struct {
	Status TokenReviewStatus
}

type SubjectAccessReviewStatus struct {
	Allowed bool
	Denied  bool
	Reason  string
}

type SubjectAccessReview struct {
	Status SubjectAccessReviewStatus
}

type TokenReviewInterface interface {
	Create(ctx context.Context, tr *TokenReview) (*TokenReview, error)
}

type SubjectAccessReviewInterface interface {
	Create(ctx context.Context, sar *SubjectAccessReview) (*SubjectAccessReview, error)
}

type AuthClient struct{}

func (c *AuthClient) TokenReviews() TokenReviewInterface {
	return nil
}

func (c *AuthClient) SubjectAccessReviews() SubjectAccessReviewInterface {
	return nil
}

type PolicyRule struct {
	Verbs     []string
	Resources []string
}

// 1. Insecure: TokenReview error discarded via blank identifier
func testVulnTokenReviewBlankError(ctx context.Context, client *AuthClient) bool {
	// ruleid: kubernetes-rbac-improper-authentication
	tr, _ := client.TokenReviews().Create(ctx, &TokenReview{})
	_ = tr
	return true
}

// 2. Insecure: TokenReview error unchecked before proceeding
func testVulnTokenReviewUncheckedError(ctx context.Context, client *AuthClient) bool {
	// ruleid: kubernetes-rbac-improper-authentication
	tr, err := client.TokenReviews().Create(ctx, &TokenReview{})
	fmt.Println("Created token review without checking error:", err)
	_ = tr
	return true
}

// 3. Insecure: SubjectAccessReview error discarded via blank identifier
func testVulnSARBlankError(ctx context.Context, client *AuthClient) bool {
	// ruleid: kubernetes-rbac-improper-authentication
	sar, _ := client.SubjectAccessReviews().Create(ctx, &SubjectAccessReview{})
	_ = sar
	return true
}

// 4. Insecure: SubjectAccessReview error unchecked before proceeding
func testVulnSARUncheckedError(ctx context.Context, client *AuthClient) bool {
	// ruleid: kubernetes-rbac-improper-authentication
	sar, err := client.SubjectAccessReviews().Create(ctx, &SubjectAccessReview{})
	fmt.Println("Created SAR without checking error:", err)
	_ = sar
	return true
}

// 5. Insecure: Overly permissive wildcard PolicyRule
func testVulnWildcardPolicyRule() PolicyRule {
	// ruleid: kubernetes-rbac-improper-authentication
	return PolicyRule{
		Verbs:     []string{"*"},
		Resources: []string{"*"},
	}
}

// 6. Safe: TokenReview properly checks error and authentication status
func testSafeTokenReview(ctx context.Context, client *AuthClient) (bool, error) {
	// ok: kubernetes-rbac-improper-authentication
	tr, err := client.TokenReviews().Create(ctx, &TokenReview{})
	if err != nil {
		return false, err
	}
	if !tr.Status.Authenticated {
		return false, fmt.Errorf("token not authenticated")
	}
	return true, nil
}

// 7. Safe: SubjectAccessReview properly checks error and allowed status
func testSafeSAR(ctx context.Context, client *AuthClient) (bool, error) {
	// ok: kubernetes-rbac-improper-authentication
	sar, err := client.SubjectAccessReviews().Create(ctx, &SubjectAccessReview{})
	if err != nil {
		return false, err
	}
	if !sar.Status.Allowed {
		return false, fmt.Errorf("access forbidden")
	}
	return true, nil
}

// 8. Safe: Restricted RBAC PolicyRule scoped to specific verbs and resources
func testSafePolicyRule() PolicyRule {
	// ok: kubernetes-rbac-improper-authentication
	return PolicyRule{
		Verbs:     []string{"get", "list"},
		Resources: []string{"configmaps"},
	}
}

// 9. Safe: TokenReview inside if initializer
func testSafeTokenReviewIfInit(ctx context.Context, client *AuthClient) bool {
	// ok: kubernetes-rbac-improper-authentication
	if tr, err := client.TokenReviews().Create(ctx, &TokenReview{}); err != nil {
		return false
	} else {
		_ = tr
		return true
	}
}
