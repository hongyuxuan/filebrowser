package handler

import (
	"net/http"

	"github.com/hongyuxuan/filebrowser/common/errorx"
	"github.com/hongyuxuan/filebrowser/internal/logic"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func uploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		// 解析表单，限制上传文件大小
		r.ParseMultipartForm(1 << 30) // 限制1GB

		file, handler, err := r.FormFile("file")
		if err != nil {
			httpx.Error(w, err)
			return
		}
		defer file.Close()

		l := logic.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.Upload(&req, file, handler)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
