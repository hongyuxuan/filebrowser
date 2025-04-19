package s3

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListconnLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListconnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListconnLogic {
	return &ListconnLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListconnLogic) Listconn() (resp *types.Response, err error) {
	return &types.Response{
		Code: http.StatusOK,
		Data: l.svcCtx.S3Conn,
	}, nil
}
