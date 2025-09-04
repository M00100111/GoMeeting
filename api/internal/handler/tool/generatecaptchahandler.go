package tool

import (
	"net/http"

	"GoMeeting/api/internal/logic/tool"
	"GoMeeting/api/internal/svc"
	"GoMeeting/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 申请生成验证码发送到邮箱并存入Redis
func GenerateCaptchaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CaptchaReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := tool.NewGenerateCaptchaLogic(r.Context(), svcCtx)
		resp, err := l.GenerateCaptcha(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
