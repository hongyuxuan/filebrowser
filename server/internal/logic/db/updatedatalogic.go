package db

import (
	"context"
	"encoding/json"
	"net/http"

	commontypes "github.com/hongyuxuan/filebrowser/common/types"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatedataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatedataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatedataLogic {
	return &UpdatedataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatedataLogic) Updatedata(tablename, id string, data map[string]interface{}) (resp *types.Response, err error) {
	if err = l.svcCtx.Database.Table(tablename).Where("id = ?", id).Updates(data).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	if tablename == "s3_repository" {
		var s3repo commontypes.S3Repository
		b, _ := json.Marshal(data)
		json.Unmarshal(b, &s3repo)
		if _, ok := l.svcCtx.S3Conn[s3repo.Name]; ok {
			var client *minio.Client
			if client, err = minio.New(s3repo.S3Endpoint, &minio.Options{
				Creds:  credentials.NewStaticV4(s3repo.S3AccessKey, s3repo.S3SecretKey, ""),
				Secure: false,
			}); err != nil {
				l.Logger.Errorf("Failed to connect to s3_endpoint: %s: %v", s3repo.S3Endpoint, err)
				return
			}
			l.svcCtx.S3Conn[s3repo.Name] = commontypes.S3Conn{
				S3Endpoint: s3repo.S3Endpoint,
				Client:     client,
			}
		}
	}
	return &types.Response{
		Code:    http.StatusOK,
		Message: "更新成功",
	}, nil
}
