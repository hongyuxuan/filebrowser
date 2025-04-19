package static

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

//go:embed dist/*
var staticFiles embed.FS

// Handler 返回处理静态资源的 http.HandlerFunc
func Handler() http.HandlerFunc {
	// 剥离 dist 目录前缀
	distFS, _ := fs.Sub(staticFiles, "dist")
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取请求路径
		path := r.URL.Path

		// 如果请求的是静态文件（如 /assets/index-C2GhvC1w.js），直接返回文件
		if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/images/") || strings.HasPrefix(path, "/login/") {
			file, err := distFS.Open(strings.TrimPrefix(path, "/"))
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			defer file.Close()

			// 获取文件信息（包括修改时间）
			fileInfo, err := file.Stat()
			if err != nil {
				http.Error(w, "Failed to get file info", http.StatusInternalServerError)
				return
			}
			modTime := fileInfo.ModTime()

			// 将文件内容读取到内存中
			fileData, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, "Failed to read file", http.StatusInternalServerError)
				return
			}

			// 将文件数据包装为 io.ReadSeeker
			fileReader := bytes.NewReader(fileData)

			// 根据文件扩展名设置 MIME 类型
			ext := filepath.Ext(path)
			switch ext {
			case ".js":
				w.Header().Set("Content-Type", "application/javascript")
			case ".css":
				w.Header().Set("Content-Type", "text/css")
			case ".html":
				w.Header().Set("Content-Type", "text/html")
			case ".json":
				w.Header().Set("Content-Type", "application/json")
			case ".png":
				w.Header().Set("Content-Type", "image/png")
			case ".jpg", ".jpeg":
				w.Header().Set("Content-Type", "image/jpeg")
			case ".svg":
				w.Header().Set("Content-Type", "image/svg+xml")
			default:
				w.Header().Set("Content-Type", "text/plain")
			}

			// 返回文件内容
			http.ServeContent(w, r, path, modTime, fileReader)
			return
		}

		// 如果不是静态文件请求，返回 index.html（支持 Vue Router History 模式）
		file, err := distFS.Open("index.html")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		// 将文件内容读取到内存中
		fileData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		// 将文件数据包装为 io.ReadSeeker
		fileReader := bytes.NewReader(fileData)

		// 设置 MIME 类型为 text/html
		w.Header().Set("Content-Type", "text/html")

		// 返回 index.html
		http.ServeContent(w, r, "index.html", time.Now(), fileReader)
	}
}
