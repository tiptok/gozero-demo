package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-demo/app/usercenter/cmd/api/internal/logic/user"
	"zero-demo/app/usercenter/cmd/api/internal/svc"
	"zero-demo/app/usercenter/cmd/api/internal/types"
)

func UserGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserGetReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := user.NewUserGetLogic(r.Context(), svcCtx)
		resp, err := l.UserGet(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
