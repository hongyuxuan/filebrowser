//go:build linux

package svc

func GetDrives() []string {
	return []string{}
}

func GetDriveType(drive string) string {
	return ""
}

func GetDiskFreeSpaceEx(drive string) (freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes int64, err error) {
	return 0, 0, 0, nil
}

func GetVolumeInformation(drive string) (volumeName, fileSystemName string, err error) {
	return "", "", nil
}
