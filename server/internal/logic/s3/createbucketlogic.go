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

type CreatebucketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatebucketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatebucketLogic {
	return &CreatebucketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatebucketLogic) Createbucket(req *types.CreateBucketRequest) (resp *types.Response, err error) {
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
	if exists {
		return nil, errorx.NewDefaultError("Bucket \"%s\" has already exists", req.BucketName)
	}
	if err = l.svcCtx.S3Conn[req.Name].Client.MakeBucket(context.Background(), req.BucketName, minio.MakeBucketOptions{
		Region:        l.svcCtx.S3Conn[req.Name].S3Region,
		ObjectLocking: req.ObjectLocking,
	}); err != nil {
		l.Logger.Error(err)
		return
	}
	l.Logger.Infof("Create bucket \"%s\" success", req.BucketName)
	if req.Versioning {
		if err := l.svcCtx.S3Conn[req.Name].Client.EnableVersioning(context.Background(), req.BucketName); err != nil {
			l.Logger.Error(err)
		} else {
			l.Logger.Infof("Bucket \"%s\" enable versioning success", req.BucketName)
		}
	}
	return &types.Response{
		Code:    http.StatusOK,
		Message: "创建桶成功",
	}, nil
}
