package logic

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/hongyuxuan/filebrowser/common/errorx"
	commontypes "github.com/hongyuxuan/filebrowser/common/types"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

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

func (l *PreviewLogic) Preview(req *types.DownloadRequest, w http.ResponseWriter) {
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

	fileInfo, err := os.Stat(req.File)
	if err != nil {
		httpx.Error(w, err)
		return
	}
	if fileInfo.Size() > sizeLimit { // > sizeLimit 的文件没必要提供预览功能，可直接走下载查看
		httpx.Error(w, errorx.NewDefaultError("大于 %s 的文件不提供预览，请下载文件后查看", humanize.Bytes(uint64(sizeLimit))))
		return
	}
	data, err := os.ReadFile(req.File)
	if err != nil {
		httpx.Error(w, err)
		return
	}

	// 设置响应头
	contentType := http.DetectContentType(data)
	logx.Infof("File=%s contentType=%s", req.File, contentType)
	if filepath.Ext(req.File) == ".svg" ||
		strings.Contains(contentType, "image/") ||
		strings.Contains(contentType, "text/plain") ||
		strings.Contains(contentType, "application/pdf") {
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Disposition", "inline; filename="+filepath.Base(req.File))
		// 写入文件内容
		if _, err = w.Write(data); err != nil {
			httpx.Error(w, err)
			return
		}
	} else {
		httpx.Error(w, errorx.NewDefaultError("该文件暂不支持预览"))
	}
}
