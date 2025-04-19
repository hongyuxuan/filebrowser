package types

import "time"

type Driver struct {
	DriverName     string
	DriverType     string
	VolumeName     string
	FileSystemName string
	TotalSpace     int64
	TotalSpaceGB   string
	UsedSpace      int64
	UsedSpaceGB    string
	FreeSpaceGB    string
}

type FilePath struct {
	Path       string
	Name       string
	IsDir      bool
	Size       int64
	LastUpdate time.Time
}
