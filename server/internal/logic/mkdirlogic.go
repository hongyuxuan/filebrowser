package logic

import (
	"context"
	"net/http"
	"os"

	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MkdirLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMkdirLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MkdirLogic {
	return &MkdirLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MkdirLogic) Mkdir(req *types.UploadRequest) (resp *types.Response, err error) {
	if err = os.MkdirAll(req.Path, 0755); err != nil {
		return
	}
	// 设置属主和属组
	if req.UID != 0 && req.GID != 0 {
		os.Chown(req.Path, req.UID, req.GID)
	}
	return &types.Response{
		Code:    http.StatusOK,
		Message: "创建目录成功",
	}, nil
}
