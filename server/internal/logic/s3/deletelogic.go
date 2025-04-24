package s3

import (
	"context"
	"net/http"
	"strings"

	"github.com/hongyuxuan/filebrowser/common/errorx"
	"github.com/hongyuxuan/filebrowser/common/utils"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"
	"github.com/minio/minio-go/v7"

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

func (l *DeleteLogic) Delete(req *types.ListObjectRequest) (resp *types.Response, err error) {
	bucket, path := utils.GetS3BucketAndPath(req.Path)
	if _, ok := l.svcCtx.S3Conn[req.Name]; !ok {
		err = errorx.NewDefaultError("Cannot find s3_endpoint with name: %s", req.Name)
		l.Logger.Error(err)
		return
	}
	opts := minio.RemoveObjectOptions{}
	if _, _, _, _, err = l.svcCtx.S3Conn[req.Name].Client.GetObjectLockConfig(context.Background(), bucket); err == nil {
		opts.GovernanceBypass = true
	}
	if strings.HasSuffix(path, "/") { // 删除目录
		objectsCh := l.svcCtx.S3Conn[req.Name].Client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{
			Prefix:    strings.TrimPrefix(path, "/"),
			Recursive: true,
		})
		for object := range objectsCh {
			if object.Err != nil {
				l.Logger.Error(err)
				continue
			}
			if err = l.svcCtx.S3Conn[req.Name].Client.RemoveObject(context.Background(), bucket, object.Key, opts); err != nil {
				l.Logger.Error(err)
				return
			}
			l.Logger.Infof("Delete s3 file=%s success", object.Key)
		}
		l.Logger.Infof("Delete s3 dir=%s success", req.Path)
	} else {
		if err = l.svcCtx.S3Conn[req.Name].Client.RemoveObject(context.Background(), bucket, strings.TrimPrefix(path, "/"), opts); err != nil {
			l.Logger.Error(err)
			return
		}
		l.Logger.Infof("Delete s3 file=%s success", req.Path)
	}
	return &types.Response{
		Code:    http.StatusOK,
		Message: "删除成功",
	}, nil
}
