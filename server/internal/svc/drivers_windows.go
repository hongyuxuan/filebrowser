//go:build windows

package svc

import (
	"syscall"
	"unsafe"
)

func GetDrives() []string {
	var drives []string
	// 获取驱动器掩码
	kernel32, _ := syscall.LoadLibrary("kernel32.dll")
	getLogicalDrives, _ := syscall.GetProcAddress(kernel32, "GetLogicalDrives")
	ret, _, _ := syscall.Syscall(uintptr(getLogicalDrives), 0, 0, 0, 0)
	// 遍历驱动器
	for i := 0; i < 26; i++ {
		if ret&(1<<uint(i)) != 0 {
			drive := string('A'+i) + ":\\"
			drives = append(drives, drive)
		}
	}
	return drives
}

func GetDriveType(drive string) string {
	kernel32, _ := syscall.LoadLibrary("kernel32.dll")
	defer syscall.FreeLibrary(kernel32)
	getDriveType, _ := syscall.GetProcAddress(kernel32, "GetDriveTypeW")

	drivePtr, _ := syscall.UTF16PtrFromString(drive)
	ret, _, _ := syscall.Syscall(uintptr(getDriveType), 1, uintptr(unsafe.Pointer(drivePtr)), 0, 0)

	switch ret {
	case 0:
		return "无法确定"
	case 1:
		return "根目录不存在"
	case 2:
		return "可移动驱动器"
	case 3:
		return "固定驱动器"
	case 4:
		return "远程驱动器"
	case 5:
		return "CD-ROM"
	case 6:
		return "RAM磁盘"
	default:
		return "无效驱动器"
	}
}

func GetDiskFreeSpaceEx(drive string) (freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes int64, err error) {
	kernel32, _ := syscall.LoadLibrary("kernel32.dll")
	defer syscall.FreeLibrary(kernel32)
	getDiskFreeSpaceEx, _ := syscall.GetProcAddress(kernel32, "GetDiskFreeSpaceExW")

	drivePtr, _ := syscall.UTF16PtrFromString(drive)
	ret, _, _ := syscall.Syscall6(uintptr(getDiskFreeSpaceEx), 4,
		uintptr(unsafe.Pointer(drivePtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)),
		0, 0)

	if ret == 0 {
		err = syscall.GetLastError()
		return
	}
	return
}

func GetVolumeInformation(drive string) (volumeName, fileSystemName string, err error) {
	kernel32, _ := syscall.LoadLibrary("kernel32.dll")
	defer syscall.FreeLibrary(kernel32)
	getVolumeInformation, _ := syscall.GetProcAddress(kernel32, "GetVolumeInformationW")

	drivePtr, _ := syscall.UTF16PtrFromString(drive)
	volumeNameBuffer := make([]uint16, 256)
	fileSystemNameBuffer := make([]uint16, 256)

	ret, _, _ := syscall.Syscall9(uintptr(getVolumeInformation), 8,
		uintptr(unsafe.Pointer(drivePtr)),             // 驱动器路径
		uintptr(unsafe.Pointer(&volumeNameBuffer[0])), // 卷标缓冲区
		uintptr(len(volumeNameBuffer)),                // 卷标缓冲区大小
		0,                                             // 卷序列号（不需要）
		0,                                             // 最大文件长度（不需要）
		0,                                             // 文件系统标志（不需要）
		uintptr(unsafe.Pointer(&fileSystemNameBuffer[0])), // 文件系统名称缓冲区
		uintptr(len(fileSystemNameBuffer)),                // 文件系统名称缓冲区大小
		0)                                                 // 保留参数

	if ret == 0 {
		err = syscall.GetLastError()
		return
	}

	volumeName = syscall.UTF16ToString(volumeNameBuffer)
	if volumeName == "" {
		volumeName = "本地磁盘"
	}
	fileSystemName = syscall.UTF16ToString(fileSystemNameBuffer)
	return
}
