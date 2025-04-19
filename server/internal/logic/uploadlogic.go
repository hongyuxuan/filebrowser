package logic

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload(req *types.UploadRequest, file multipart.File, header *multipart.FileHeader) (resp *types.Response, err error) {
	filePath := path.Join(req.Path, header.Filename)

	// 创建目标文件
	out, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// 拷贝文件内容
	_, err = io.Copy(out, file)
	if err != nil {
		return nil, err
	}

	// 设置属主和属组
	if req.UID != 0 && req.GID != 0 {
		os.Chown(filePath, req.UID, req.GID)
	}

	// 返回结果
	return &types.Response{
		Code:    http.StatusOK,
		Message: "上传成功",
	}, nil
}
