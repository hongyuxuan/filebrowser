package handler

import (
	"net/http"

	"github.com/hongyuxuan/filebrowser/internal/logic"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func rootpathHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewRootpathLogic(r.Context(), svcCtx)
		resp, err := l.Rootpath()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
