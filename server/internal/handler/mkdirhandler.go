package handler

import (
	"net/http"

	"github.com/hongyuxuan/filebrowser/common/errorx"
	"github.com/hongyuxuan/filebrowser/internal/logic"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func mkdirHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		l := logic.NewMkdirLogic(r.Context(), svcCtx)
		resp, err := l.Mkdir(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
