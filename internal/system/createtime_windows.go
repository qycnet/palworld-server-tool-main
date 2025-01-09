//go:build windows

package system

import (
	"os"
	"syscall"
	"time"
)

func GetFileCreateTime(info os.FileInfo) time.Time {
	// 尝试将 FileInfo 的 Sys() 返回值转换为 *syscall.Win32FileAttributeData 类型
	stat, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		// 如果转换失败，则返回零时间
		return time.Time{}
	}
	// 将 CreationTime 转换为 Unix 时间戳并返回
	return time.Unix(0, stat.CreationTime.Nanoseconds())
}

func GetEntryCreateTime(info os.DirEntry) time.Time {
	// 获取文件信息
	winFileInfo, err := info.Info()
	if err != nil {
		// 如果获取文件信息出错，则返回零时间
		return time.Time{}
	}

	// 尝试将文件信息的系统部分转换为Win32FileAttributeData类型
	stat, ok := winFileInfo.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		// 如果转换失败，则返回零时间
		return time.Time{}
	}

	// 将文件的创建时间转换为Unix时间并返回
	return time.Unix(0, stat.CreationTime.Nanoseconds())
}
