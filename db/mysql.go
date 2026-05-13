package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
func InstallMysqldump() error {
	wait := exec.RunShell("apk", []string{"add", "--no-cache", "mariadb-client"}, nil, "", false)
	result, code := wait.WaitToList()
	if code != 0 {
		return fmt.Errorf("安装 mariadb-client 失败（exit=%d）：%s", code, result.ToString(","))
	}
	return nil
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
		if err := InstallMysqldump(); err != nil {
			return 0, err
		}
		if !IsMysqldumpInstalled() {
			return 0, fmt.Errorf("mysqldump 未安装，且自动安装后仍未找到（当前镜像可能非 alpine，需手动安装 mariadb-client/mysql-client）")
		}
	}

	// pipefail：管道中任一命令失败都反映到退出码，避免 mysqldump 失败但 gzip 成功时被误判为 0
	// 2>errFile：把 mysqldump 的 stderr 单独写到文件，便于失败时取到真实错误原因
	// （否则 mysqldump 的 stderr 会跟 sh 的 stderr 混在一起，且 bufio reader 在快速失败时可能丢数据）
	errFile := fileName + ".err"
	cmd := fmt.Sprintf("set -o pipefail; mysqldump -h %s -P %d -u%s -p%s %s 2>%s | gzip > %s", host, port, username, password, database, errFile, fileName)
	args := []string{"-c", cmd}
	wait := exec.RunShell("sh", args, nil, "", false)
	result, code := wait.WaitToList()

	// 读取 mysqldump 的 stderr 内容，命令结束后不管成功失败都要清理
	var dumpErr string
	if data, readErr := os.ReadFile(errFile); readErr == nil {
		dumpErr = strings.TrimSpace(string(data))
	}
	file.Delete(errFile)

	buildErrMsg := func() string {
		msg := dumpErr
		if shellOut := result.ToString(","); shellOut != "" {
			if msg != "" {
				msg += " | "
			}
			msg += "shell: " + shellOut
		}
		return msg
	}

	// 备份失败时删除备份文件
	if code != 0 {
		file.Delete(fileName)
		return 0, fmt.Errorf("备份%s数据库失败（exit=%d）：%s", database, code, buildErrMsg())
	}
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return 0, fmt.Errorf("获取备份文件信息:%s,失败： %s", fileName, err.Error())
	}
	// 兜底：某些 sh 实现对 pipefail 支持不佳，空文件视为失败
	if fileInfo.Size() == 0 {
		file.Delete(fileName)
		return 0, fmt.Errorf("备份%s数据库失败：生成文件为空，%s", database, buildErrMsg())
	}
	return fileInfo.Size() / 1024, nil
}

// 恢复数据库
func RecoverMysql(host string, port int, username, password, database string, fileName string) error {
	// 安装 mysqldump
	if !IsMysqldumpInstalled() {
		if err := InstallMysqldump(); err != nil {
			return err
		}
		if !IsMysqldumpInstalled() {
			return fmt.Errorf("mysqldump 未安装，且自动安装后仍未找到（当前镜像可能非 alpine，需手动安装 mariadb-client/mysql-client）")
		}
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
