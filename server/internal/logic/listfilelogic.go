package logic

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListfileLogic {
	return &ListfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListfileLogic) Listfile(req *types.ListFileRequest) (resp *types.Response, err error) {
	if req.Path == "" {
		req.Path = l.svcCtx.Config.RoutePath
	}
	// 计算根目录的深度
	rootParts := strings.Split(filepath.ToSlash(req.Path), "/")
	files := []types.FilePath{}

	// 过滤掉空字符串（例如 Windows 的根目录 "D:/" 会分割为 ["D:", ""]）
	var filteredRootParts []string
	for _, part := range rootParts {
		if part != "" {
			filteredRootParts = append(filteredRootParts, part)
		}
	}
	rootDepth := len(filteredRootParts)

	err = filepath.Walk(req.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) { // 如果是权限错误，跳过该目录
				return nil
			}
			return err
		}
		// 计算当前路径的深度
		pathParts := strings.Split(filepath.ToSlash(path), "/")
		// 过滤掉空字符串
		var filteredPathParts []string
		for _, part := range pathParts {
			if part != "" {
				filteredPathParts = append(filteredPathParts, part)
			}
		}
		currentDepth := len(filteredPathParts) - rootDepth
		if currentDepth > 1 {
			return filepath.SkipDir
		}
		if path != req.Path {
			files = append(files, types.FilePath{
				Path:       path,
				Name:       info.Name(),
				IsDir:      info.IsDir(),
				LastUpdate: info.ModTime(),
				Size:       info.Size(),
			})
		}
		return nil
	})
	return &types.Response{
		Code: http.StatusOK,
		Data: files,
	}, err
}
