package db

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/filebrowser/common/types"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletedataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletedataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletedataLogic {
	return &DeletedataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletedataLogic) Deletedata(req *types.DataByIdReq) (resp *types.Response, err error) {
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "删除成功",
	}

	if req.Tablename == "s3_repository" {
		var s3repo commontypes.S3Repository
		if err = l.svcCtx.Database.First(&s3repo, req.Id).Error; err != nil {
			l.Logger.Error(err)
			return
		}
		if err = l.svcCtx.Database.Delete(&s3repo).Error; err != nil {
			l.Logger.Error(err)
			return
		}
		if _, ok := l.svcCtx.S3Conn[s3repo.Name]; ok {
			delete(l.svcCtx.S3Conn, s3repo.Name)
			l.Logger.Infof("Removed s3_endpoint: %s from S3Conn", s3repo.S3Endpoint)
		}
	} else {
		data := map[string]interface{}{}
		if err = l.svcCtx.Database.Table(req.Tablename).Where("id = ?", req.Id).Delete(&data).Error; err != nil {
			l.Logger.Error(err)
			return
		}
	}
	return
}
