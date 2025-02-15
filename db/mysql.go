package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/utils/exec"
	"github.com/farseer-go/utils/file"
)

// 检查 mysqldump 是否已安装
func IsMysqldumpInstalled() bool {
	code, result := exec.RunShellCommand("mysqldump --version", nil, "", false)
	if code != 0 || len(result) == 0 {
		return false
	}
	// 检查输出中是否包含 "mysqldump" 关键字
	return strings.Contains(result[0], "mysqldump")
}

// 安装 mysqldump
func InstallMysqldump() {
	exec.RunShellCommand("apk add --no-cache mariadb-client", nil, "", false)
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

	mysqldumpCmd := fmt.Sprintf("mysqldump -h %s -P %d -u%s -p%s %s | gzip > %s", host, port, username, password, database, fileName)
	code, result := exec.RunShellCommand(mysqldumpCmd, nil, "", false)
	// 备份失败时删除备份文件
	if code != 0 {
		file.Delete(fileName)
		return 0, fmt.Errorf("备份%s数据库失败：%s", database, collections.NewList(result...).ToString(","))
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

	var cmd string
	switch fileExt {
	case ".gz":
		cmd = fmt.Sprintf("gzip -dc %s | mysql -h %s -P %d -u%s -p%s %s", filepath.Base(fileName), host, port, username, password, database)
	case ".sql":
		cmd = fmt.Sprintf("mysql -h %s -P %d -u%s -p%s %s < %s", host, port, username, password, database, fileName)
	default:
		return fmt.Errorf("未知的扩展名：%s", fileExt)
	}

	code, result := exec.RunShellCommand(cmd, nil, path, false)
	if code != 0 {
		return fmt.Errorf("还原SQL文件：%s 时失败：%s", cmd, collections.NewList(result...).ToString(","))
	}
	return nil
}
