package oidc

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
)

type IDToken struct {
	sub     string
	isHuman bool
	raw     *oidc.IDToken
}

func (t *IDToken) Sub() string {
	return t.sub
}

func (t *IDToken) IsHuman() bool {
	return t.isHuman
}

func (t *IDToken) Raw() *oidc.IDToken {
	return t.raw
}

func NewIDToken(sub string, isHuman bool, raw *oidc.IDToken) *IDToken {
	return &IDToken{sub: sub, isHuman: isHuman, raw: raw}
}

type StubProvider struct {
	token *IDToken
}

func NewStubProvider(token *IDToken) *StubProvider {
	return &StubProvider{token: token}
}

func (p *StubProvider) AuthURL(state string) string {
	return "https://auth.humanityprotocol.com/auth?state=" + state
}

func (p *StubProvider) Exchange(ctx context.Context, code string) (*IDToken, error) {
	return p.token, nil
}

func AnyIDToken() *IDToken {
	return &IDToken{sub: "hp_test_sub", isHuman: true}
}

func IDTokenWithIsHuman(isHuman bool) *IDToken {
	return &IDToken{sub: "hp_test_sub", isHuman: isHuman}
}

func IDTokenWithSub(sub string) *IDToken {
	return &IDToken{sub: sub, isHuman: true}
}