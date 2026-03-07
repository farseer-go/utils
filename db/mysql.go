package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/utils/exec"
	"github.com/farseer-go/utils/file"
)

// 检查 mysqldump 是否已安装
func IsMysqldumpInstalled() bool {
	wait := exec.RunShell("mysqldump", []string{"--version"}, nil, "", false)
	result, code := wait.WaitToList()
	if code != 0 || result.Count() == 0 {
		return false
	}
	// 检查输出中是否包含 "mysqldump" 关键字
	return result.ContainsAny("mysqldump")
}

// 安装 mysqldump
func InstallMysqldump() {
	exec.RunShell("apk", []string{"add", "--no-cache", "mariadb-client"}, nil, "", false)
}

// 备份历史数据
type FileObject struct {
	Database string            // 数据库
	FileName string            // 文件名
	CreateAt dateTime.DateTime // 备份时间
	Size     int64             // 备份文件大小（KB）
}

// 备份MYSQL
func BackupMysql(host string, port int, username, password, database string, fileName string) (int64, error) {
	// 安装 mysqldump
	if !IsMysqldumpInstalled() {
		InstallMysqldump()
	}

	cmd := fmt.Sprintf("mysqldump -h %s -P %d -u%s -p%s %s | gzip > %s", host, port, username, password, database, fileName)
	args := []string{"-c", cmd}
	wait := exec.RunShell("sh", args, nil, "", false)
	result, code := wait.WaitToList()
	// 备份失败时删除备份文件
	if code != 0 {
		file.Delete(fileName)
		return 0, fmt.Errorf("备份%s数据库失败：%s", database, result.ToString(","))
	}
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return 0, fmt.Errorf("获取备份文件信息:%s,失败： %s", fileName, err.Error())
	}
	return fileInfo.Size() / 1024, nil
}

// 恢复数据库
func RecoverMysql(host string, port int, username, password, database string, fileName string) error {
	// 安装 mysqldump
	if !IsMysqldumpInstalled() {
		InstallMysqldump()
	}

	path := filepath.Dir(fileName)
	fileExt := filepath.Ext(fileName)

	var args []string
	switch fileExt {
	case ".gz":
		// 管道操作需要通过 shell 执行
		cmd := fmt.Sprintf("gzip -dc %s | mysql -h %s -P %d -u%s -p%s %s",
			filepath.Base(fileName), host, port, username, password, database)
		args = []string{"-c", cmd}
	case ".sql":
		// 重定向操作需要通过 shell 执行
		cmd := fmt.Sprintf("mysql -h %s -P %d -u%s -p%s %s < %s",
			host, port, username, password, database, fileName)
		args = []string{"-c", cmd}
	default:
		return fmt.Errorf("未知的扩展名：%s", fileExt)
	}

	wait := exec.RunShell("sh", args, nil, path, false)
	result, code := wait.WaitToList()
	if code != 0 {
		return fmt.Errorf("还原SQL文件失败：%s", result.ToString(","))
	}
	return nil
}
