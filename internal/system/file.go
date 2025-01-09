package system

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/qycnet/palworld-server-tool-main/internal/logger"
)

func GetExecDir() (string, error) {
	// 获取可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		// 如果获取可执行文件路径出错，返回空字符串和错误信息
		return "", err
	}
	// 返回可执行文件所在的目录路径
	return filepath.Dir(exePath), nil
}

func CheckIsDir(path string) (bool, error) {
	// 获取文件信息
	fileInfo, err := os.Stat(path)
	if err != nil {
		// 如果获取文件信息出错，则返回false和错误信息
		return false, err
	}
	// 返回文件是否为目录的信息和nil
	return fileInfo.IsDir(), nil
}

func CopyDir(srcDir, dstDir string) error {
	// 使用 filepath.Walk 遍历 srcDir 目录
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		// 如果遍历过程中出现错误，则直接返回错误
		if err != nil {
			return err
		}
		// 获取 srcDir 到当前遍历路径的相对路径
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		// 拼接目标路径
		dstPath := filepath.Join(dstDir, relPath)
		// 如果当前遍历的文件是一个目录
		if info.IsDir() {
			// 创建目标目录，并赋予最大权限
			return os.MkdirAll(dstPath, os.ModePerm)
		}
		// 如果当前遍历的文件不是目录，则复制文件
		return CopyFile(path, dstPath)
	})
}

func CopyFile(srcFile, destFile string) error {
	// 打开源文件
	input, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	// 关闭源文件句柄
	defer input.Close()

	// 创建目标文件
	output, err := os.Create(destFile)
	if err != nil {
		return err
	}
	// 关闭目标文件句柄
	defer output.Close()

	// 将源文件内容复制到目标文件
	_, err = io.Copy(output, input)
	return err
}

func ZipDir(srcDir, zipFilePath string) error {
	// 创建zip文件
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 创建zip写入器
	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	// 遍历源目录
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 创建zip头信息
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// 将路径分隔符替换为'/'
		header.Name = strings.ReplaceAll(relPath, string(os.PathSeparator), "/")

		// 如果是目录，则添加'/'
		if info.IsDir() {
			header.Name += "/"
		} else {
			// 设置压缩方法
			header.Method = zip.Deflate
		}

		// 创建zip头
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		// 如果是文件，则将其写入zip文件
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := io.Copy(writer, file); err != nil {
				file.Close()
				return err
			}
		}
		return nil
	})
	return err
}

func UnzipDir(zipFile, destDir string) error {
	// 打开zip文件
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		// 构建目标路径
		path := filepath.Join(destDir, file.Name)

		// 如果是目录
		if file.FileInfo().IsDir() {
			// 创建目录
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
		} else {
			// 创建文件所在的目录
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}

			// 打开目标文件
			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			// 打开zip文件中的文件
			rc, err := file.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			// 将zip文件中的内容复制到目标文件
			if _, err := io.Copy(outFile, rc); err != nil {
				return err
			}
		}
	}
	return nil
}

func CleanAndCreateDir(dirPath string) error {
	// 检查dirPath是否存在
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		// 如果dirPath存在，则删除dirPath目录及其所有内容
		if err := os.RemoveAll(dirPath); err != nil {
			return err
		}
	}
	// 创建dirPath目录，并设置权限为0755
	return os.MkdirAll(dirPath, 0755)
}

// CheckAndCreateDir 检查指定路径的文件夹是否存在，如果不存在则创建它。
func CheckAndCreateDir(path string) error {
	// 检查目录是否存在
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// 如果目录不存在，则创建目录
		err := os.MkdirAll(path, 0755)
		if err != nil {
			// 如果创建目录失败，返回错误
			return err
		}
	} else if err != nil {
		// 如果检查目录存在性时出错，返回错误
		return err
	}
	return nil
}

func GetSavDir(path string) (string, error) {
	// 声明一个字符串变量levelFilePath
	var levelFilePath string
	// 调用GetLevelSavFilePath函数，传入path参数，并将返回值赋给levelFilePath和err
	levelFilePath, err := GetLevelSavFilePath(path)
	// 如果err不为nil，则返回空字符串和err
	if err != nil {
		return "", err
	}
	// 返回levelFilePath的目录路径和nil
	// 返回filepath.Dir(levelFilePath)，即levelFilePath的目录路径
	return filepath.Dir(levelFilePath), nil
}

func GetLevelSavFilePath(path string) (string, error) {
	var foundPath string
	// 遍历指定路径下的所有文件和文件夹
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 如果发生错误，则返回该错误
			return err
		}
		// 如果文件名是"Level.sav"
		if info.Name() == "Level.sav" {
			// 记录文件路径并返回一个表示文件已找到的错误
			foundPath = path
			return errors.New("file found")
		}
		return nil
	})

	// 如果没有找到文件路径
	if foundPath == "" {
		// 如果存在错误且错误不是表示文件已找到的错误
		if err != nil && !errors.Is(err, errors.New("file found")) {
			// 返回错误
			return "", err
		}
		// 如果文件未找到，则返回相应的错误
		return "", errors.New("file Level.sav not found")
	}

	// 返回找到的文件路径和nil错误
	return foundPath, nil
}

// LimitCacheZipFiles keeps only the latest `n` zip archives in the cache directory
func LimitCacheZipFiles(cacheDir string, n int) {
	// 读取缓存目录中的文件
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		// 如果读取目录时发生错误，记录错误并返回
		logger.Errorf("Error reading cache directory: %v\n", err)
		return
	}

	// 用于存储所有zip文件的切片
	zipFiles := []os.DirEntry{}
	// 遍历目录中的文件
	for _, file := range files {
		// 如果文件扩展名为.zip，则将其添加到zipFiles切片中
		if filepath.Ext(file.Name()) == ".zip" {
			zipFiles = append(zipFiles, file)
		}
	}

	// 如果zip文件的数量小于等于n，则直接返回
	if len(zipFiles) <= n {
		return
	}

	// 根据文件的创建时间对zip文件进行排序
	sort.Slice(zipFiles, func(i, j int) bool {
		createTimeI := GetEntryCreateTime(zipFiles[i])
		createTimeJ := GetEntryCreateTime(zipFiles[j])
		return createTimeI.After(createTimeJ)
	})

	// 删除多余的zip文件
	for i := n; i < len(zipFiles); i++ {
		// 构建要删除文件的路径
		path := filepath.Join(cacheDir, zipFiles[i].Name())
		// 删除文件
		err := os.Remove(path)
		if err != nil {
			// 如果删除文件时发生错误，记录错误
			logger.Errorf("Failed to delete excess zip file: %v\n", err)
		}
	}
}

type dirInfo struct {
	path       string
	createTime time.Time
}

// LimitCacheDir keeps only the latest `n` directories in the cache directory
func LimitCacheDir(cacheDirPrefix string, n int) error {
	// 获取临时目录
	tempDir := os.TempDir()

	// 读取临时目录中的条目
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		// 如果读取临时目录出错，记录错误并返回
		logger.Errorf("LimitCacheDir: error reading temp directory: %v\n", err)
		return err
	}

	// 存储符合条件的目录信息
	var dirs []dirInfo
	for _, entry := range entries {
		// 如果条目是目录且名称以cacheDirPrefix为前缀
		if entry.IsDir() && strings.HasPrefix(filepath.Base(entry.Name()), cacheDirPrefix) {
			// 获取目录路径
			dirPath := filepath.Join(tempDir, entry.Name())
			// 获取条目的创建时间
			createTime := GetEntryCreateTime(entry)
			// 将目录信息添加到dirs中
			dirs = append(dirs, dirInfo{path: dirPath, createTime: createTime})
		}
	}

	// 根据创建时间对dirs进行降序排序
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].createTime.After(dirs[j].createTime)
	})

	// 如果dirs长度大于n
	if len(dirs) > n {
		// 遍历dirs中第n个元素之后的所有元素
		for _, dir := range dirs[n:] {
			// 删除目录
			err := os.RemoveAll(dir.path)
			if err != nil {
				// 如果删除目录出错，记录错误并返回
				logger.Errorf("LimitCacheDir: error removing directory: %v\n", err)
				return err
			}
		}
	}

	return nil
}

func UnTarGzDir(tarStream io.Reader, destDir string) error {
	// 创建一个gzip解压缩器
	gzr, err := gzip.NewReader(tarStream)
	if err != nil {
		return err
	}
	defer gzr.Close()

	// 创建一个tar读取器
	tr := tar.NewReader(gzr)

	for {
		// 读取下一个tar头信息
		header, err := tr.Next()

		// 如果到达文件末尾，则跳出循环
		if err == io.EOF {
			break
		}
		// 如果读取过程中出现错误，则返回错误
		if err != nil {
			return err
		}

		// 生成目标路径
		target := filepath.Join(destDir, header.Name)

		// 根据tar头的类型标志进行不同的处理
		switch header.Typeflag {
		case tar.TypeDir:
			// 如果是目录，则创建目录
			if err := os.MkdirAll(target, os.ModePerm); err != nil {
				return err
			}
		case tar.TypeReg:
			// 如果是常规文件，则创建文件目录
			targetDir := filepath.Dir(target)
			if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
				return err
			}
			// 创建并打开文件
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				return err
			}
			// 复制数据到文件中
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		}
	}

	return nil
}
