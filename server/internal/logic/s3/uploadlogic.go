package s3

import (
	"context"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/hongyuxuan/filebrowser/common/errorx"
	"github.com/hongyuxuan/filebrowser/common/utils"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"
	"github.com/minio/minio-go/v7"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload(req *types.ListObjectRequest, file multipart.File, header *multipart.FileHeader) (resp *types.Response, err error) {
	if _, ok := l.svcCtx.S3Conn[req.Name]; !ok {
		err = errorx.NewDefaultError("Cannot find s3_endpoint with name: %s", req.Name)
		l.Logger.Error(err)
		return
	}
	// 确保目录路径以 "/" 结尾
	if !strings.HasSuffix(req.Path, "/") {
		req.Path += "/"
	}
	bucket, path := utils.GetS3BucketAndPath(req.Path)
	filename := strings.TrimPrefix(path, "/") + header.Filename
	l.Logger.Debugf("Upload object=%s to bucket=%s", filename, bucket)
	if _, err = l.svcCtx.S3Conn[req.Name].Client.PutObject(
		context.Background(),
		bucket,
		filename,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: "application/octet-stream", // 可选
		},
	); err != nil {
		return
	}
	return &types.Response{
		Code:    http.StatusOK,
		Message: "上传成功",
	}, nil
}
