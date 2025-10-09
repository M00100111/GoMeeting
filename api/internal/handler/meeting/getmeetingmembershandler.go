package meeting

import (
	"GoMeeting/pkg/ctxdata"
	"net/http"

	"GoMeeting/api/internal/logic/meeting"
	"GoMeeting/api/internal/svc"
	"GoMeeting/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMeetingMembersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMeetingMembersReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		//获取Ctx中的用户id
		req.UserId = r.Context().Value(ctxdata.JwtUserId).(uint64)
		l := meeting.NewGetMeetingMembersLogic(r.Context(), svcCtx)
		resp, err := l.GetMeetingMembers(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
