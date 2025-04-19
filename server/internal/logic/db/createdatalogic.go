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

type CreatedataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatedataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatedataLogic {
	return &CreatedataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatedataLogic) Createdata(tablename string, data map[string]interface{}) (resp *types.Response, err error) {
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "新增成功",
	}
	if err = l.svcCtx.Database.Table(tablename).Create(data).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	if tablename == "s3_repository" {
		var s3repo commontypes.S3Repository
		b, _ := json.Marshal(data)
		json.Unmarshal(b, &s3repo)
		var client *minio.Client
		if client, err = minio.New(s3repo.S3Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(s3repo.S3AccessKey, s3repo.S3SecretKey, ""),
			Region: s3repo.S3Region,
			Secure: s3repo.UseSecure,
		}); err != nil {
			l.Logger.Errorf("Failed to connect to s3_endpoint: %s: %v", s3repo.S3Endpoint, err)
			return
		}
		l.Logger.Infof("Successfully connect to s3_endpoint: %s", s3repo.S3Endpoint)
		l.svcCtx.S3Conn[s3repo.Name] = commontypes.S3Conn{
			S3Endpoint: s3repo.S3Endpoint,
			Client:     client,
		}
	}
	return
}
