package db

import (
	"fmt"
	"os"
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
