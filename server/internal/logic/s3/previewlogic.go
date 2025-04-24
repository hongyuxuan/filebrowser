package s3

import (
	"context"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/hongyuxuan/filebrowser/common/errorx"
	commontypes "github.com/hongyuxuan/filebrowser/common/types"
	"github.com/hongyuxuan/filebrowser/common/utils"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"
	"github.com/minio/minio-go/v7"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type PreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreviewLogic {
	return &PreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreviewLogic) Preview(req *types.ListObjectRequest, w http.ResponseWriter) {
	// get settings
	var settings commontypes.Settings
	if err := l.svcCtx.Database.First(&settings, "setting_key = ?", "preview_file_size").Error; err != nil {
		httpx.Error(w, err)
		return
	}
	sizeLimit, err := strconv.ParseInt(settings.SettingValue, 10, 64)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	bucket, path := utils.GetS3BucketAndPath(req.Path)

	// 从 MinIO/S3 获取文件对象
	if _, ok := l.svcCtx.S3Conn[req.Name]; !ok {
		httpx.Error(w, errorx.NewDefaultError("Cannot find s3_endpoint with name: %s", req.Name))
		return
	}
	object, err := l.svcCtx.S3Conn[req.Name].Client.GetObject(context.Background(), bucket, strings.TrimPrefix(path, "/"), minio.GetObjectOptions{})
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
	if objInfo.Size > sizeLimit { // > sizeLimit 的文件没必要提供预览功能，可直接走下载查看
		httpx.Error(w, errorx.NewDefaultError("大于 %s 的文件不提供预览，请下载文件后查看", humanize.Bytes(uint64(sizeLimit))))
		return
	}

	data, err := io.ReadAll(object)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	// 设置响应头
	contentType := http.DetectContentType(data)
	logx.Infof("File=%s contentType=%s", req.Path, contentType)
	if filepath.Ext(req.Path) == ".svg" ||
		strings.Contains(contentType, "image/") ||
		strings.Contains(contentType, "text/plain") ||
		strings.Contains(contentType, "application/pdf") {
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition", "inline; filename="+filepath.Base(req.Path))
		// 写入文件内容
		if _, err = w.Write(data); err != nil {
			httpx.Error(w, err)
			return
		}
	} else {
		httpx.Error(w, errorx.NewDefaultError("该文件暂不支持预览"))
	}
}
