package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-demo/app/usercenter/cmd/api/internal/logic/user"
	"zero-demo/app/usercenter/cmd/api/internal/svc"
	"zero-demo/app/usercenter/cmd/api/internal/types"
)

func UserUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := user.NewUserUpdateLogic(r.Context(), svcCtx)
		resp, err := l.UserUpdate(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
