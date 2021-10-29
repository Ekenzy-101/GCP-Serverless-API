package config

const (
	AccessTokenCookieName    = "sls_acc_token"
	AccessTokenTTLInSeconds  = 60 * 60
	RefreshTokenCookieName   = "sls_ref_token"
	RefreshTokenTTLInSeconds = 60 * 60 * 24 * 7
)
