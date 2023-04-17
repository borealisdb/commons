package auth

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/coreos/go-oidc/v3/oidc"
)

type IDTokenClaims struct {
	Email    string   `json:"email"`
	Verified bool     `json:"email_verified"`
	Groups   []string `json:"groups"`
}

func InitializeAuth(ctx context.Context, issuerUrl, clientID string) (*oidc.Provider, *oidc.IDTokenVerifier, error) {
	ctx = oidc.InsecureIssuerURLContext(ctx, issuerUrl)
	var provider *oidc.Provider

	// It could be that DEX is not yet ready, thus we retry
	if err := retry.Do(
		func() error {
			// TODO provider url should not hardcoded
			prvdr, err := oidc.NewProvider(ctx, "http://localhost:5556/borealis/identity")
			if err != nil {
				return err
			}
			provider = prvdr
			return nil
		},
		retry.Attempts(10),
	); err != nil {
		return nil, nil, err
	}

	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: clientID})

	return provider, idTokenVerifier, nil
}
