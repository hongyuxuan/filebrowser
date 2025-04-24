package svc

import (
	"time"

	"github.com/golang-jwt/jwt"
	commontypes "github.com/hongyuxuan/filebrowser/common/types"
	"github.com/hongyuxuan/filebrowser/common/utils"
	"github.com/hongyuxuan/filebrowser/internal/config"
	"github.com/hongyuxuan/filebrowser/internal/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	Database     *gorm.DB
	S3Conn       map[string]commontypes.S3Conn
	Validateuser rest.Middleware
	Version      string
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		Database:     utils.NewSQLite(c.Database, c.Log.Level),
		S3Conn:       make(map[string]commontypes.S3Conn),
		Validateuser: middleware.NewValidateuserMiddleware().Handle,
	}
}

func (s *ServiceContext) GetJwtToken(user commontypes.User, expireTime *int64) (accessToken string, now int64, err error) {
	payloads := map[string]interface{}{
		"username": user.Username,
		"role":     user.Role,
	}
	if expireTime == nil {
		expireTime = &s.Config.Auth.AccessExpire
	}
	now = time.Now().Unix()
	if accessToken, err = s.generateToken(now, payloads, *expireTime); err != nil {
		return
	}
	return
}

func (s *ServiceContext) SetS3() {
	var s3repos []commontypes.S3Repository
	if err := s.Database.Find(&s3repos).Error; err != nil {
		logx.Errorf("Failed to get s3_endpoint from database: %v", err)
		return
	}
	for _, s3 := range s3repos {
		client, err := minio.New(s3.S3Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(s3.S3AccessKey, s3.S3SecretKey, ""),
			Region: s3.S3Region,
			Secure: s3.UseSecure,
		})
		if err != nil {
			logx.Errorf("Failed to connect to s3_endpoint: %s: %v", s3.S3Endpoint, err)
			continue
		}
		logx.Infof("Successfully connect to s3_endpoint: %s", s3.S3Endpoint)
		s.S3Conn[s3.Name] = commontypes.S3Conn{
			S3Endpoint: s3.S3Endpoint,
			S3Region:   s3.S3Region,
			Client:     client,
		}
	}
}

func (s *ServiceContext) SetVersion(version string) {
	s.Version = version
}

func (s *ServiceContext) generateToken(iat int64, payloads map[string]interface{}, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["payloads"] = payloads
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(s.Config.Auth.AccessSecret))
}
