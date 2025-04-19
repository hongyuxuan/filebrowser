package auth

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/filebrowser/common/utils"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChpasswdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChpasswdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChpasswdLogic {
	return &ChpasswdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChpasswdLogic) Chpasswd(req *types.ChpasswdReq) (resp *types.Response, err error) {
	if err = utils.ModifyPassword(req.Username, req.OldPassword, req.NewPassword, l.svcCtx.Database); err != nil {
		l.Logger.Error(err)
		return
	}
	l.Logger.Infof("User \"%s\" changed password success")
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "修改密码成功",
	}
	return
}
