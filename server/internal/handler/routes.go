// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	auth "github.com/hongyuxuan/filebrowser/internal/handler/auth"
	db "github.com/hongyuxuan/filebrowser/internal/handler/db"
	s3 "github.com/hongyuxuan/filebrowser/internal/handler/s3"
	"github.com/hongyuxuan/filebrowser/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Validateuser},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/version",
					Handler: versionHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/drivers",
					Handler: driversHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/listfile",
					Handler: listfileHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/rootpath",
					Handler: rootpathHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/upload",
					Handler: uploadHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/mkdir",
					Handler: mkdirHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/preview",
					Handler: previewHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/delete",
					Handler: deleteHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/filebrowser"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/download",
				Handler: downloadHandler(serverCtx),
			},
		},
		rest.WithPrefix("/filebrowser"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/user/info",
				Handler: auth.UserinfoHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/chpasswd",
				Handler: auth.ChpasswdHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/adduser",
				Handler: auth.AdduserHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/filebrowser/auth"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: auth.LoginHandler(serverCtx),
			},
		},
		rest.WithPrefix("/filebrowser/auth"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Validateuser},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/:tablename",
					Handler: db.ListdataHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/:tablename/:id",
					Handler: db.GetdataHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/:tablename",
					Handler: db.CreatedataHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/:tablename/:id",
					Handler: db.UpdatedataHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/:tablename/:id",
					Handler: db.DeletedataHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/filebrowser/db"),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Validateuser},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/listconnections",
					Handler: s3.ListconnHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/listbuckets",
					Handler: s3.ListbucketHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/listobjects",
					Handler: s3.ListobjectHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/mkdir",
					Handler: s3.MkdirHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/upload",
					Handler: s3.UploadHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/preview",
					Handler: s3.PreviewHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/delete",
					Handler: s3.DeleteHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/createbucket",
					Handler: s3.CreatebucketHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/deletebucket",
					Handler: s3.DeletebucketHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/filebrowser/s3"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/downloadobject",
				Handler: s3.DownloadobjectHandler(serverCtx),
			},
		},
		rest.WithPrefix("/filebrowser/s3"),
	)
}
