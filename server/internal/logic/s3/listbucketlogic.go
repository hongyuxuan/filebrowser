package s3

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/filebrowser/common/errorx"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"
	"github.com/minio/minio-go/v7"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListbucketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListbucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListbucketLogic {
	return &ListbucketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListbucketLogic) Listbucket(req *types.ListBucketRequest) (resp *types.Response, err error) {
	var buckets []minio.BucketInfo
	if _, ok := l.svcCtx.S3Conn[req.Name]; !ok {
		err = errorx.NewDefaultError("Cannot find s3_endpoint with name: %s", req.Name)
		l.Logger.Error(err)
		return
	}
	if buckets, err = l.svcCtx.S3Conn[req.Name].Client.ListBuckets(context.Background()); err != nil {
		l.Logger.Error(err)
		return
	}
	return &types.Response{
		Code: http.StatusOK,
		Data: buckets,
	}, nil
}
