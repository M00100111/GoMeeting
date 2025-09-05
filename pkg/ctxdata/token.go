package ctxdata

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

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

// ParseJwtToken 验证并解析JWT令牌
func parseJwtToken(tokenString string, secretKey string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
}

// GetUidFromToken 从token中提取用户ID
func GetUidFromToken(tokenString string, secretKey string) (string, error) {
	token, err := parseJwtToken(tokenString, secretKey)
	if err != nil {
		return "", err
	}

	// 验证token是否有效
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 提取用户ID
		if uid, exists := claims[Identify]; exists {
			if uidStr, ok := uid.(string); ok {
				return uidStr, nil
			}
			return "", fmt.Errorf("user id is not a string")
		}
		return "", fmt.Errorf("user id not found in token")
	}

	return "", fmt.Errorf("invalid token claims")
}

// ValidateToken 验证token是否有效
func ValidateToken(tokenString string, secretKey string) bool {
	token, err := parseJwtToken(tokenString, secretKey)
	if err != nil {
		return false
	}
	return token.Valid
}
