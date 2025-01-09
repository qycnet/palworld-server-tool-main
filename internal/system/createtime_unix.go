//go:build !windows

package system

import (
	"os"
	"time"
)

func GetFileCreateTime(info os.FileInfo) time.Time {
	// 返回文件的修改时间
	return info.ModTime()
}

func GetEntryCreateTime(info os.DirEntry) time.Time {
	// 获取文件的Unix文件信息
	unixFileInfo, err := info.Info()
	if err != nil {
		// 如果获取信息失败，则返回一个空的time.Time对象
		return time.Time{}
	}
	// 返回文件的修改时间
	return unixFileInfo.ModTime()
}
