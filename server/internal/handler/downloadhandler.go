package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hongyuxuan/filebrowser/common/errorx"
	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func downloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DownloadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		file, err := os.Open(req.File)
		if err != nil {
			httpx.Error(w, errorx.NewDefaultError(err.Error()))
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			httpx.Error(w, errorx.NewDefaultError(err.Error()))
			return
		}
		if fileInfo.IsDir() {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, fmt.Sprintf("%s is a directory", req.File), nil))
			return
		}

		// 设置响应头
		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(req.File))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Accept-Ranges", "bytes")

		// 如果没有Range头，直接发送整个文件
		var rangeHeader string
		if rangeHeader = r.Header.Get("Range"); rangeHeader == "" {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
			http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
			return
		}

		// 处理分块下载
		const chunkSize = 1024 * 1024 // 1MB chunks
		totalSize := fileInfo.Size()

		// 解析Range头
		ranges, err := parseRange(rangeHeader, totalSize)
		if err != nil {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		// 如果是多范围请求，这里简化处理只支持单范围
		if len(ranges) != 1 {
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		start := ranges[0].start
		end := ranges[0].end

		// 确保结束位置不超过文件大小
		if end >= totalSize {
			end = totalSize - 1
		}

		// 如果请求的范围大于chunkSize，则限制为chunkSize
		if end-start+1 > chunkSize {
			end = start + chunkSize - 1
			if end >= totalSize {
				end = totalSize - 1
			}
		}

		// 设置部分内容的响应头
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, totalSize))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
		w.WriteHeader(http.StatusPartialContent)

		// 定位到开始位置
		if _, err := file.Seek(start, io.SeekStart); err != nil {
			httpx.Error(w, err)
			return
		}

		// 使用io.CopyN复制指定大小的数据
		if _, err := io.CopyN(w, file, end-start+1); err != nil {
			// 客户端可能提前关闭连接，不需要记录为错误
			if err != io.EOF {
				logx.Errorf("Error sending file chunk: %v", err)
			}
			return
		}
	}
}

type httpRange struct {
	start, end int64
}

func parseRange(s string, size int64) ([]httpRange, error) {
	if s == "" {
		return nil, nil
	}
	const b = "bytes="
	if !strings.HasPrefix(s, b) {
		return nil, errors.New("invalid range")
	}
	var ranges []httpRange
	for _, ra := range strings.Split(s[len(b):], ",") {
		ra = strings.TrimSpace(ra)
		if ra == "" {
			continue
		}
		i := strings.Index(ra, "-")
		if i < 0 {
			return nil, errors.New("invalid range")
		}
		start, end := strings.TrimSpace(ra[:i]), strings.TrimSpace(ra[i+1:])
		var r httpRange
		if start == "" {
			// 如 "-500" 表示最后500字节
			i, err := strconv.ParseInt(end, 10, 64)
			if err != nil {
				return nil, errors.New("invalid range")
			}
			if i > size {
				i = size
			}
			r.start = size - i
			r.end = size - 1
		} else {
			i, err := strconv.ParseInt(start, 10, 64)
			if err != nil || i >= size || i < 0 {
				return nil, errors.New("invalid range")
			}
			r.start = i
			if end == "" {
				// 如 "1000-" 表示从1000字节到文件结束
				r.end = size - 1
			} else {
				i, err := strconv.ParseInt(end, 10, 64)
				if err != nil || r.start > i {
					return nil, errors.New("invalid range")
				}
				if i >= size {
					i = size - 1
				}
				r.end = i
			}
		}
		ranges = append(ranges, r)
	}
	if len(ranges) == 0 {
		return nil, errors.New("invalid range")
	}
	return ranges, nil
}
