package auth

import (
	"net/http"

	"github.com/hongyuxuan/filebrowser/internal/logic/auth"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserinfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewUserinfoLogic(r.Context(), svcCtx)
		resp, err := l.Userinfo()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
