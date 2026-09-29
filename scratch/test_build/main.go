package rules

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

func Test(ctx context.Context, conf *oauth2.Config, verifier *oidc.IDTokenVerifier) {
	_, _ = conf.Exchange(ctx, "code")
	_, _ = verifier.Verify(ctx, "raw")
	_ = oidc.Config{SkipClientIDCheck: true}
}
