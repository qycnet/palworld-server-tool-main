package source

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/qycnet/palworld-server-tool-main/internal/logger"
	"github.com/qycnet/palworld-server-tool-main/internal/system"
)

func CopyFromLocal(src, way string) (string, error) {

	// 检查src是否为目录
	isDir, err := system.CheckIsDir(src)
	if err != nil {
		logger.Errorf("error checking if %s is a directory: %v\n", src, err)
	}

	// 获得Level.sav路径
	var levelPath string
	if isDir {
		// 如果src是目录，获取Level.sav文件路径
		levelPath, err = system.GetLevelSavFilePath(src)
		if err != nil {
			return "", errors.New("error finding Level.sav: \n" + err.Error())
		}
	} else {
		// 如果src不是目录，检查文件名是否为Level.sav
		if filepath.Base(src) == "Level.sav" {
			levelPath = src
		} else {
			return "", errors.New("specified file is not Level.sav and source is not a directory")
		}
	}

	// 获取Level.sav所在的目录
	savDir := filepath.Dir(levelPath)

	// 创建临时目录
	randId := uuid.New().String()
	tempDir := filepath.Join(os.TempDir(), "palworldsav-"+way+"-"+randId)
	if err = os.MkdirAll(tempDir, fs.ModePerm); err != nil {
		return "", err
	}

	// 拷贝文件
	files, err := filepath.Glob(filepath.Join(savDir, "*.sav"))
	if err != nil {
		return "", err
	}
	for _, file := range files {
		// 将文件拷贝到临时目录
		dist := filepath.Join(tempDir, filepath.Base(file))
		if err = system.CopyFile(file, dist); err != nil {
			return "", err
		}
	}

	// 拷贝Players目录
	playerDir := filepath.Join(savDir, "Players")
	distPlayerDir := filepath.Join(tempDir, "Players")
	if err = system.CopyDir(playerDir, distPlayerDir); err != nil {
		return "", err
	}

	// 返回临时目录中Level.sav的路径
	distLevelPath := filepath.Join(tempDir, "Level.sav")
	return distLevelPath, nil
}
