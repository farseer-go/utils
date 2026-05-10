package db

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/stopwatch"
	"github.com/farseer-go/utils/exec"
	"github.com/farseer-go/utils/file"
)

// 检查 clickhouse-client 是否已安装
func IsClickhouseClientInstalled() bool {
	wait := exec.RunShell("clickhouse-client", []string{"--version"}, nil, "", false)
	result, code := wait.WaitToList()
	if code != 0 || result.Count() == 0 {
		return false
	}
	return result.ContainsAny("clickhouse")
}

func InstallClickhouseClient() {
	exec.RunShell("apk", []string{"add", "--no-cache", "clickhouse-client"}, nil, "", false)
}

func buildClickhouseClientArgs(host string, port int, username, password, database string) string {
	args := fmt.Sprintf("--host=%s --port=%d --database=%s", host, port, database)
	if username != "" {
		args += fmt.Sprintf(" --user=%s", username)
	}
	if password != "" {
		args += fmt.Sprintf(" --password=%s", password)
	}
	return args
}

// BackupClickhouse 备份Clickhouse数据库（Native格式）
// tables: 表列表
// ddlStatements: 每张表对应的DDL（DROP + CREATE），key=tableName, value=ddl内容
// 返回最终tar.gz文件大小(KB)
func BackupClickhouse(host string, port int, username, password, database string, tables []string, ddlStatements map[string]string, fileName string) (int64, error) {
	if !IsClickhouseClientInstalled() {
		InstallClickhouseClient()
	}

	filePath := filepath.Dir(fileName)
	file.CreateDir766(filePath)

	// 写入schema文件
	schemaFile := filepath.Join(filePath, database+".schema.sql")
	file.Delete(schemaFile)
	file.WriteString(schemaFile, "")
	for _, tableName := range tables {
		if ddl, ok := ddlStatements[tableName]; ok {
			file.AppendLine(schemaFile, ddl)
		}
	}

	connArgs := buildClickhouseClientArgs(host, port, username, password, database)

	// 导出每张表的Native数据
	for _, tableName := range tables {
		sw := stopwatch.StartNew()
		query := fmt.Sprintf("SELECT * FROM %s.%s FORMAT Native", database, tableName)
		dataFile := filepath.Join(filePath, tableName+".native")
		cmd := fmt.Sprintf("clickhouse-client %s --query=%q > %s", connArgs, query, dataFile)
		wait := exec.RunShell("sh", []string{"-c", cmd}, nil, filePath, false)
		result, exitCode := wait.WaitToList()
		if exitCode != 0 {
			return 0, fmt.Errorf("导出%s.%s失败：%s", database, tableName, result.ToString(","))
		}
		flog.Infof("导出clickhouse %s.%s 使用了：%s", database, tableName, sw.GetMillisecondsText())
	}

	// 打包为tar.gz
	var tarFiles []string
	tarFiles = append(tarFiles, filepath.Base(schemaFile))
	for _, tableName := range tables {
		tarFiles = append(tarFiles, tableName+".native")
	}
	tarCmd := fmt.Sprintf("tar -czf %s %s", fileName, strings.Join(tarFiles, " "))
	wait := exec.RunShell("sh", []string{"-c", tarCmd}, nil, filePath, false)
	result, exitCode := wait.WaitToList()
	if exitCode != 0 {
		return 0, fmt.Errorf("压缩备份文件失败：%s", result.ToString(","))
	}

	// 清理临时文件
	file.Delete(schemaFile)
	for _, tableName := range tables {
		file.Delete(filepath.Join(filePath, tableName+".native"))
	}

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return 0, fmt.Errorf("获取备份文件信息:%s,失败： %s", fileName, err.Error())
	}
	return fileInfo.Size() / 1024, nil
}

// RecoverClickhouse 恢复Clickhouse数据库（Native格式）
func RecoverClickhouse(host string, port int, username, password, database string, fileName string) error {
	if !IsClickhouseClientInstalled() {
		InstallClickhouseClient()
	}

	// 解压tar.gz到临时目录
	extractDir := fileName + "_extract"
	file.CreateDir766(extractDir)
	defer os.RemoveAll(extractDir)

	tarCmd := fmt.Sprintf("tar -xzf %s -C %s", fileName, extractDir)
	wait := exec.RunShell("sh", []string{"-c", tarCmd}, nil, filepath.Dir(fileName), false)
	result, exitCode := wait.WaitToList()
	if exitCode != 0 {
		return fmt.Errorf("解压备份文件%s失败：%s", fileName, result.ToString(","))
	}

	// 找到schema文件并执行DDL
	schemaFiles, _ := filepath.Glob(extractDir + "/*.schema.sql")
	if len(schemaFiles) == 0 {
		return fmt.Errorf("备份文件中未找到schema文件")
	}

	connArgs := buildClickhouseClientArgs(host, port, username, password, database)

	fSql, err := os.Open(schemaFiles[0])
	if err != nil {
		return fmt.Errorf("打开schema文件失败: %v", err)
	}
	defer fSql.Close()

	scanner := bufio.NewScanner(fSql)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, cap(buf))

	var sqlBuilder strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "--") {
			continue
		}
		sqlBuilder.WriteString(line + "\n")
		if strings.HasSuffix(line, ";") {
			sqlStatement := sqlBuilder.String()
			sqlBuilder.Reset()
			sw := stopwatch.StartNew()
			cmd := fmt.Sprintf("clickhouse-client %s --query=%q", connArgs, sqlStatement)
			wait := exec.RunShell("sh", []string{"-c", cmd}, nil, extractDir, false)
			result, exitCode := wait.WaitToList()
			if exitCode != 0 {
				return fmt.Errorf("执行DDL失败: %s\nSQL: %s", result.ToString(","), sqlStatement)
			}
			flog.Infof("还原%s DDL执行 使用了：%s", database, sw.GetMillisecondsText())
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取schema文件失败: %v", err)
	}

	// 导入每张表的Native数据
	nativeFiles, _ := filepath.Glob(extractDir + "/*.native")
	for _, nativeFile := range nativeFiles {
		tableName := strings.TrimSuffix(filepath.Base(nativeFile), ".native")
		sw := stopwatch.StartNew()

		fileInfo, _ := os.Stat(nativeFile)
		if fileInfo != nil && fileInfo.Size() == 0 {
			flog.Infof("还原%s.%s 跳过空表", database, tableName)
			continue
		}

		query := fmt.Sprintf("INSERT INTO %s.%s FORMAT Native", database, tableName)
		cmd := fmt.Sprintf("clickhouse-client %s --query=%q < %s", connArgs, query, nativeFile)
		wait := exec.RunShell("sh", []string{"-c", cmd}, nil, extractDir, false)
		result, exitCode := wait.WaitToList()
		if exitCode != 0 {
			return fmt.Errorf("还原%s.%s失败：%s", database, tableName, result.ToString(","))
		}
		flog.Infof("还原%s.%s 使用了：%s", database, tableName, sw.GetMillisecondsText())
	}
	return nil
}
