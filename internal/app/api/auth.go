package api

import (
	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/pkg/config/env"
	"github.com/Januadrym/seennit/internal/pkg/jwt"
)

//
func newAuthHandler(signer jwt.Signer, authenticator auth.UserAuthen) *auth.Handler {
	srv := auth.NewService(signer, authenticator)
	return auth.NewHandler(srv)
}

//
func newJWTSignVerifier() jwt.SignVerifier {
	var conf jwt.Config
	env.Load(&conf)
	return jwt.New(conf)
}
