package db

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hongyuxuan/filebrowser/common/errorx"
	"github.com/hongyuxuan/filebrowser/internal/logic/db"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdatedataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 此处未使用go-zero提供的httpx.Parse解析，改为手动解析
		var data map[string]interface{}
		r.Body = http.MaxBytesReader(w, r.Body, 10485760) // 10MB限制
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			httpx.Error(w, err)
			return
		}
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 5 {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, "Path params error", nil))
		}
		tablename := pathParts[3]
		id := pathParts[4]

		l := db.NewUpdatedataLogic(r.Context(), svcCtx)
		resp, err := l.Updatedata(tablename, id, data)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
