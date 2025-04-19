package logic

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RootpathLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRootpathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RootpathLogic {
	return &RootpathLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RootpathLogic) Rootpath() (resp *types.Response, err error) {
	return &types.Response{
		Code: http.StatusOK,
		Data: l.svcCtx.Config.RoutePath,
	}, nil
}
