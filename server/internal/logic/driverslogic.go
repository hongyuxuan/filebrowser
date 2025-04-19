package logic

import (
	"context"
	"fmt"
	"net/http"
	"runtime"

	"github.com/hongyuxuan/filebrowser/internal/svc"
	"github.com/hongyuxuan/filebrowser/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DriversLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDriversLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DriversLogic {
	return &DriversLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DriversLogic) Drivers() (resp *types.Response, err error) {
	var drivers []types.Driver
	if runtime.GOOS == "windows" {
		for _, drive := range svc.GetDrives() {
			driveType := svc.GetDriveType(drive)
			// 获取磁盘空间信息
			var freeBytesAvailable, totalNumberOfBytes int64
			if freeBytesAvailable, totalNumberOfBytes, _, err = svc.GetDiskFreeSpaceEx(drive); err != nil {
				return
			}
			var volumeName, fileSystemName string
			if volumeName, fileSystemName, err = svc.GetVolumeInformation(drive); err != nil {
				return
			}
			usedBytes := totalNumberOfBytes - freeBytesAvailable
			drivers = append(drivers, types.Driver{
				DriverName:     drive,
				DriverType:     driveType,
				VolumeName:     volumeName,
				FileSystemName: fileSystemName,
				TotalSpace:     totalNumberOfBytes,
				TotalSpaceGB:   fmt.Sprintf("%.f GB", float64(totalNumberOfBytes)/1024/1024/1024),
				UsedSpace:      usedBytes,
				UsedSpaceGB:    fmt.Sprintf("%.f GB", float64(usedBytes)/1024/1024/1024),
				FreeSpaceGB:    fmt.Sprintf("%.f GB", float64(freeBytesAvailable)/1024/1024/1024),
			})
		}
	} else {
		drivers = []types.Driver{
			{
				DriverName: "/",
			},
		}
	}
	return &types.Response{
		Code: http.StatusOK,
		Data: drivers,
	}, nil
}
