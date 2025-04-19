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

type MkdirLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMkdirLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MkdirLogic {
	return &MkdirLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MkdirLogic) Mkdir(req *types.ListObjectRequest) (resp *types.Response, err error) {
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
	// 上传一个空的.keep文件，为了占住目录
	filename := strings.TrimPrefix(path, "/") + ".keep"
	l.Logger.Debugf("Upload object=%s to bucket=%s", filename, bucket)
	keep := "I am a keep file"
	if _, err = l.svcCtx.S3Conn[req.Name].Client.PutObject(
		context.Background(),
		bucket,
		filename,
		strings.NewReader(keep), // 空内容
		int64(len(keep)),        // 大小为 0
		minio.PutObjectOptions{ContentType: "application/x-directory"}, // 可选标记
	); err != nil {
		return
	}
	return &types.Response{
		Code:    http.StatusOK,
		Message: "创建目录成功",
	}, nil
}
