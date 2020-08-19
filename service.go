package oauthsocialite

import (
	goservice "github.com/baozhenglab/go-sdk"
)

const (
	KeyFacebook = "facebook"
	KeyGoogle   = "google"
)

func NewOAuthSocialite(key string) goservice.PrefixConfigure {
	switch key {
	case KeyFacebook:
		return &facebookSocialite{
			name: key,
		}
	case KeyGoogle:
		return &googleSocialite{
			name: key,
		}
	default:
		panic("Not found driver")
	}
}
