package s3

import (
	"net/http"

	"github.com/hongyuxuan/filebrowser/internal/logic/s3"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListconnHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := s3.NewListconnLogic(r.Context(), svcCtx)
		resp, err := l.Listconn()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
