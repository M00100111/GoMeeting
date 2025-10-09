package meeting

import (
	"GoMeeting/pkg/ctxdata"
	"net/http"

	"GoMeeting/api/internal/logic/meeting"
	"GoMeeting/api/internal/svc"
	"GoMeeting/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMeetingInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMeetingInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		//获取Ctx中的用户id
		req.UserId = r.Context().Value(ctxdata.JwtUserId).(uint64)
		l := meeting.NewGetMeetingInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetMeetingInfo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
