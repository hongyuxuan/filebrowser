package db

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/filebrowser/common/types"
	"github.com/hongyuxuan/filebrowser/common/utils"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListdataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListdataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListdataLogic {
	return &ListdataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListdataLogic) Listdata(req *commontypes.GetDataReq) (resp *types.Response, err error) {
	_, role := utils.GetPayload(l.ctx)
	data := []map[string]interface{}{}
	tx := l.svcCtx.Database.Table(req.Tablename)
	var count int64
	utils.SetTx(tx, req, &count, role, nil)
	if err = tx.Find(&data).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: commontypes.ListResult{
			Total:   int(count),
			Results: data,
		},
	}
	return
}
