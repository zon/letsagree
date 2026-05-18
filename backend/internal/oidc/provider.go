package oidc

import (
	"context"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	IssuerURL    string `yaml:"issuerUrl"`
	RedirectURL  string `yaml:"redirectUrl"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}
	return &cfg, nil
}

type Provider struct {
	oauth2Cfg *oauth2.Config
	verifier  *oidc.IDTokenVerifier
}

func NewProvider(ctx context.Context, cfg *Config) (*Provider, error) {
	oidcProvider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("create OIDC provider: %w", err)
	}

	p := &Provider{
		oauth2Cfg: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Endpoint:     oidcProvider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "identity:read"},
		},
		verifier:  oidcProvider.Verifier(&oidc.Config{ClientID: cfg.ClientID}),
	}
	return p, nil
}

func (p *Provider) AuthURL(state string) string {
	return p.oauth2Cfg.AuthCodeURL(state)
}

func (p *Provider) Exchange(ctx context.Context, code string) (*IDToken, error) {
	oauth2Token, err := p.oauth2Cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token in oauth2 token response")
	}

	idToken, err := p.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("verify ID token: %w", err)
	}

	var claims struct {
		Sub     string `json:"sub"`
		IsHuman bool   `json:"is_human"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("decode claims: %w", err)
	}

	return &IDToken{
		sub:     claims.Sub,
		isHuman: claims.IsHuman,
		raw:     idToken,
	}, nil
}