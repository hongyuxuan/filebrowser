package logic

import (
	"context"
	"net/http"
	"os"

	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DownloadRequest) (resp *types.Response, err error) {
	fileInfo, err := os.Stat(req.File)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		if err = os.RemoveAll(req.File); err != nil {
			return
		}
	} else {
		if err = os.Remove(req.File); err != nil {
			return
		}
	}
	return &types.Response{
		Code:    http.StatusOK,
		Message: "删除成功",
	}, nil
}
