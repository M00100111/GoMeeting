package middleware

import (
	"GoMeeting/api/internal/types"
	"GoMeeting/pkg/ctxdata"
	code "GoMeeting/pkg/result"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"strings"
)

type JwtAuthMiddleware struct {
	AccessSecret string
}

func NewJwtAuthMiddleware(accessSecret string) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		AccessSecret: accessSecret,
	}
}

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		//获取Authorization头
		authHeader := r.Header.Get("Authorization")
		//fmt.Println("authHeader:")
		fmt.Println(authHeader)
		if authHeader == "" {
			// 正常返回结果，提示认证失败
			httpx.OkJson(w, types.NewErrorResultWithCodef(code.TokenErrorCode, "未传递Token"))
			return
		}
		//解析Token
		// 解析 Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// 正常返回结果，提示认证失败
			httpx.OkJson(w, types.NewErrorResultWithCodef(code.TokenErrorCode, "Token格式错误"))
			return
		}
		//验证Token
		if !ctxdata.ValidateToken(tokenString, m.AccessSecret) {
			// 正常返回结果，提示认证失败
			httpx.OkJson(w, types.NewErrorResultWithCodef(code.TokenErrorCode, "验证Token失败"))
			return
		}

		uid, err := ctxdata.GetUidFromToken(tokenString, m.AccessSecret)
		if err != nil {
			// 正常返回结果，提示认证失败
			httpx.OkJson(w, types.NewErrorResultWithCodef(code.TokenErrorCode, "解析Token参数失败"+err.Error()))
			return
		}
		// 将用户id添加到 context
		ctx := context.WithValue(r.Context(), ctxdata.JwtUserId, uid)

		// Passthrough to next handler if need
		next(w, r.WithContext(ctx))
	}
}
