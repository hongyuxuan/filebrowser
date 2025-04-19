package s3

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/hongyuxuan/filebrowser/common/errorx"
	"github.com/hongyuxuan/filebrowser/common/utils"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"
	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DownloadobjectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListObjectRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		bucket, path := utils.GetS3BucketAndPath(req.Path)

		// 从 MinIO/S3 获取文件对象
		if _, ok := svcCtx.S3Conn[req.Name]; !ok {
			httpx.Error(w, errorx.NewDefaultError("Cannot find s3_endpoint with name: %s", req.Name))
			return
		}
		object, err := svcCtx.S3Conn[req.Name].Client.GetObject(context.Background(), bucket, strings.TrimPrefix(path, "/"), minio.GetObjectOptions{})
		if err != nil {
			httpx.Error(w, errorx.NewDefaultError("Failed to get object from s3: %v", err.Error()))
			return
		}
		defer object.Close()

		// 获取文件信息（用于设置 Content-Disposition）
		objInfo, err := object.Stat()
		if err != nil {
			httpx.Error(w, errorx.NewDefaultError("Failed to stat object: %v", err.Error()))
			return
		}

		// 设置响应头（强制下载）
		filename := filepath.Base(objInfo.Key)
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", objInfo.ContentType)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", objInfo.Size))

		// 流式传输文件内容
		_, err = io.Copy(w, object)
		if err != nil {
			httpx.Error(w, errorx.NewDefaultError("Failed to stream object: %v", err.Error()))
			return
		}
	}
}
