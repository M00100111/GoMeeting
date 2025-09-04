package ctxdata

import "github.com/golang-jwt/jwt/v4"

const Identify = "meeting"

func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["iat"] = iat
	claims["exp"] = iat + seconds
	claims[Identify] = uid

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	return token.SignedString([]byte(secretKey))
}
