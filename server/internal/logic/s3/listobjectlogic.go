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

type ListobjectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListobjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListobjectLogic {
	return &ListobjectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListobjectLogic) Listobject(req *types.ListObjectRequest) (resp *types.Response, err error) {
	// 确保路径以 "/" 结尾（如果是目录）
	if req.Path != "" && !strings.HasSuffix(req.Path, "/") {
		req.Path += "/"
	}
	bucket, path := utils.GetS3BucketAndPath(req.Path)
	// 列出对象
	if _, ok := l.svcCtx.S3Conn[req.Name]; !ok {
		err = errorx.NewDefaultError("Cannot find s3_endpoint with name: %s", req.Name)
		l.Logger.Error(err)
		return
	}
	l.Logger.Debugf("req.Path=%s bucket=%s prefix=%s", req.Path, bucket, path)
	objectCh := l.svcCtx.S3Conn[req.Name].Client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{
		Prefix:    strings.TrimPrefix(path, "/"),
		Recursive: false, // 不递归，仅当前目录
	})
	files := []types.FilePath{}
	for object := range objectCh {
		if object.Err != nil {
			l.Logger.Error(err)
			return nil, object.Err
		}
		l.Logger.Debugf("req.Path=%s object.Key=%s", req.Path, object.Key)
		objectKey := strings.TrimPrefix(object.Key, "/")
		files = append(files, types.FilePath{
			Name:       strings.TrimSuffix(strings.TrimPrefix(objectKey, strings.TrimPrefix(path, "/")), "/"),
			Path:       "/" + bucket + "/" + strings.TrimSuffix(objectKey, "/"),
			IsDir:      strings.HasSuffix(object.Key, "/"),
			LastUpdate: object.LastModified,
			Size:       object.Size,
		})
	}

	return &types.Response{
		Code: http.StatusOK,
		Data: files,
	}, nil
}
