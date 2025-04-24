package s3

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/filebrowser/common/errorx"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletebucketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletebucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletebucketLogic {
	return &DeletebucketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletebucketLogic) Deletebucket(req *types.CreateBucketRequest) (resp *types.Response, err error) {
	if _, ok := l.svcCtx.S3Conn[req.Name]; !ok {
		err = errorx.NewDefaultError("Cannot find s3_endpoint with name: %s", req.Name)
		l.Logger.Error(err)
		return
	}
	exists, err := l.svcCtx.S3Conn[req.Name].Client.BucketExists(context.Background(), req.BucketName)
	if err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	if !exists {
		return nil, errorx.NewDefaultError("Bucket \"%s\" doesn't exists", req.BucketName)
	}
	if err = l.svcCtx.S3Conn[req.Name].Client.RemoveBucket(context.Background(), req.BucketName); err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	l.Logger.Infof("Remove bucket \"%s\" success", req.BucketName)
	return &types.Response{
		Code:    http.StatusOK,
		Message: "删除桶成功",
	}, nil
}
