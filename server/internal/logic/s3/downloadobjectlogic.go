package s3

import (
	"context"

	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadobjectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadobjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadobjectLogic {
	return &DownloadobjectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadobjectLogic) Downloadobject(req *types.ListObjectRequest) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
