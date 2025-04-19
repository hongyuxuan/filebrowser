package handler

import (
	"net/http"

	"github.com/hongyuxuan/filebrowser/internal/logic"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func driversHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewDriversLogic(r.Context(), svcCtx)
		resp, err := l.Drivers()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
