package jwtc

import (
	"kob-kratos/internal/pkg/ctxdata"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	Uid int64
	Iat int64 // issued at，unix timestamp，token颁发时间
	Exp int64 // token 的到期时间
}

func GenJwtToken(secretKey string, payload *Payload) (string, error) {
	iat := payload.Iat
	exp := payload.Exp
	userID := payload.Uid
	claims := make(jwt.MapClaims)
	claims["exp"] = exp
	claims["iat"] = iat
	claims[ctxdata.CtxKeyUid] = userID
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
