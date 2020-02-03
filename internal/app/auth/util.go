package auth

import (
	"time"

	"github.com/seennit/internal/app/types"
	"github.com/seennit/internal/pkg/jwt"
)

func userToClaims(user *types.User, lifeTime time.Duration) jwt.Claims {
	return jwt.Claims{
		AvatarURL: user.AvatarURL,
		Role:      user.Roles,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserID:    user.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(lifeTime).Unix(),
			Id:        user.UserID,
			IssuedAt:  time.Now().Unix(),
			Issuer:    jwt.DefaultIssuer,
			Subject:   user.UserID,
		},
	}
}

func claimsToUser(claims *jwt.Claims) *types.User {
	return &types.User{
		AvatarURL: claims.AvatarURL,
		Roles:     claims.Role,
		FirstName: claims.FirstName,
		LastName:  claims.LastName,
		UserID:    claims.UserID,
	}
}