package file

import (
	"io/fs"
	"os"
	"path/filepath"
)

// GetFiles 读取指定目录下的文件
// path：目录路径
// searchPattern：匹配文件名要包含的名称，搜索全部，传入""即可
// searchSubDir：是否要搜索子目录
func GetFiles(path string, searchPattern string, searchSubDir bool) []string {
	var files []string
	filepath.WalkDir(path, func(filePath string, dirInfo fs.DirEntry, err error) error {
		if path == filePath {
			return nil
		}
		// 目录不需要判断，filepath.Walk执行就包含递归了
		if !dirInfo.IsDir() {
			match := true
			if searchPattern != "" {
				match, _ = filepath.Match(filepath.Join(filepath.Dir(filePath), searchPattern), filePath)
			}
			if match {
				files = append(files, filePath)
			}
		} else if dirInfo.IsDir() && !searchSubDir {
			return fs.SkipDir
		}
		return nil
	})
	return files
}

// GetFiles 读取指定目录下的文件
// path：目录路径
// searchPattern：匹配文件名要包含的名称，搜索全部，传入""即可
// searchSubDir：是否要搜索子目录
func GetDirs(path string, searchPattern string, searchSubDir bool) []string {
	var dirs []string
	filepath.WalkDir(path, func(filePath string, dirInfo fs.DirEntry, err error) error {
		if path == filePath {
			return nil
		}
		if dirInfo.IsDir() {
			match := true
			if searchPattern != "" {
				match, _ = filepath.Match(filepath.Join(filepath.Dir(filePath), searchPattern), filePath)
			}
			if match {
				dirs = append(dirs, filePath)
			}
			// 目录不需要判断，filepath.Walk执行就包含递归了
			if !searchSubDir {
				return fs.SkipDir
			}
		}
		return nil
	})
	return dirs
}

// ClearFile 清空目录下的所有文件（但不删除path目录本身）
// path：目录路径
func ClearFile(path string) {
	_ = filepath.WalkDir(path, func(filePath string, dirInfo fs.DirEntry, err error) error {
		if path == filePath {
			return nil
		}
		_ = os.RemoveAll(filePath)
		return nil
	})
}

// IsExists 判断路径是否存在
// path：目录路径
func IsExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
