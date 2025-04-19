package handler

import (
	"net/http"

	"github.com/hongyuxuan/filebrowser/internal/logic"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/hongyuxuan/filebrowser/common/errorx"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func listfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListFileRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		l := logic.NewListfileLogic(r.Context(), svcCtx)
		resp, err := l.Listfile(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
