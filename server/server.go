package main

import (
	"fmt"
	"net/http"

	"github.com/hongyuxuan/filebrowser/common/errorx"
	"github.com/hongyuxuan/filebrowser/common/utils"
	"github.com/hongyuxuan/filebrowser/internal/config"
	"github.com/hongyuxuan/filebrowser/internal/handler"
	"github.com/hongyuxuan/filebrowser/internal/static"
	"github.com/hongyuxuan/filebrowser/internal/svc"

	"github.com/alecthomas/kingpin/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var (
	configFile = kingpin.Flag("config", "config file").Short('f').Default("").String()
	port       = kingpin.Flag("port", "Listen port").Short('t').Default("9121").Int()
	logLevel   = kingpin.Flag("log.level", "Log level.").Default("").String()
	rootPath   = kingpin.Flag("root", "Root path of directory").Short('r').Default("").String()
	dbfile     = kingpin.Flag("db", "SQLite database file.").Short('d').Default("filebrowser.db").String()

	/* print app version */
	AppVersion = "unknown"
	GoVersion  = "unknown"
	BuildTime  = "unknown"
	OsArch     = "unknown"
	Author     = "unknown"
)

func main() {
	// Parse flags
	kingpin.Version(printVersion())
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	c := config.NewConfig(configFile, logLevel, rootPath, dbfile, port)

	server := rest.MustNewServer(c.RestConf, rest.WithUnauthorizedCallback(func(w http.ResponseWriter, r *http.Request, err error) {
		httpx.Error(w, errorx.NewError(http.StatusUnauthorized, "jwttoken is invalid or expired", nil))
	}))
	defer server.Stop()

	// 注册静态资源处理器
	staticHandler := static.Handler()
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/",
		Handler: staticHandler,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/assets/:file", // 处理 /assets/ 路径
		Handler: staticHandler,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/images/:file", // 处理 /images/ 路径
		Handler: staticHandler,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/login/:file", // 处理 /login/ 路径
		Handler: staticHandler,
	})
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/:path",
		Handler: staticHandler,
	})

	ctx := svc.NewServiceContext(c)
	ctx.SetVersion(AppVersion)

	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *errorx.BrowserError:
			return e.Code, e.GetData()
		default:
			return http.StatusInternalServerError, errorx.HttpErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
	})

	// initdb
	utils.InitDB(ctx.Database)

	// init s3
	ctx.SetS3()

	logx.Infof("Starting server at %s:%d...", c.Host, c.Port)
	server.Start()
}

func printVersion() string {
	return fmt.Sprintf("App version: %s\nGo version:  %s\nBuild Time:  %s\nOS/Arch:     %s\nAuthor:      %s\n", AppVersion, GoVersion, BuildTime, OsArch, Author)
}
