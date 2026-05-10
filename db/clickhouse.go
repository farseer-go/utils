package db

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/stopwatch"
	"github.com/farseer-go/utils/exec"
	"github.com/farseer-go/utils/file"
)

// BackupClickhouse 备份Clickhouse数据库（CSV格式）
// tables: 表列表
// ddlStatements: 每张表对应的DDL（DROP + CREATE），key=tableName, value=ddl内容
// 返回最终tar.gz文件大小(KB)
func BackupClickhouse(db *sql.DB, database string, tables []string, ddlStatements map[string]string, fileName string) (int64, error) {
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

	// 导出每张表的CSV数据
	for _, tableName := range tables {
		sw := stopwatch.StartNew()
		dataFile := filepath.Join(filePath, tableName+".csv")

		rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s.%s", database, tableName))
		if err != nil {
			return 0, fmt.Errorf("导出%s.%s失败：%v", database, tableName, err)
		}

		f, err := os.Create(dataFile)
		if err != nil {
			rows.Close()
			return 0, fmt.Errorf("创建文件%s失败：%v", dataFile, err)
		}

		columns, _ := rows.Columns()
		f.WriteString(strings.Join(columns, "\t") + "\n")

		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		for rows.Next() {
			if err := rows.Scan(valuePtrs...); err != nil {
				f.Close()
				rows.Close()
				return 0, fmt.Errorf("导出%s.%s scan失败：%v", database, tableName, err)
			}
			var line []string
			for _, v := range values {
				line = append(line, fmt.Sprintf("%v", v))
			}
			f.WriteString(strings.Join(line, "\t") + "\n")
		}
		f.Close()
		rows.Close()

		flog.Infof("导出clickhouse %s.%s 使用了：%s", database, tableName, sw.GetMillisecondsText())
	}

	// 打包为tar.gz
	var tarFiles []string
	tarFiles = append(tarFiles, filepath.Base(schemaFile))
	for _, tableName := range tables {
		tarFiles = append(tarFiles, tableName+".csv")
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
		file.Delete(filepath.Join(filePath, tableName+".csv"))
	}

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return 0, fmt.Errorf("获取备份文件信息:%s,失败： %s", fileName, err.Error())
	}
	return fileInfo.Size() / 1024, nil
}

// RecoverClickhouse 恢复Clickhouse数据库（CSV格式）
func RecoverClickhouse(db *sql.DB, database string, fileName string) error {
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
			_, err := db.Exec(sqlStatement)
			if err != nil {
				return fmt.Errorf("执行DDL失败: %v\nSQL: %s", err, sqlStatement)
			}
			flog.Infof("还原%s DDL执行 使用了：%s", database, sw.GetMillisecondsText())
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取schema文件失败: %v", err)
	}

	// 导入每张表的CSV数据
	csvFiles, _ := filepath.Glob(extractDir + "/*.csv")
	for _, csvFile := range csvFiles {
		tableName := strings.TrimSuffix(filepath.Base(csvFile), ".csv")
		sw := stopwatch.StartNew()

		fileInfo, _ := os.Stat(csvFile)
		if fileInfo != nil && fileInfo.Size() == 0 {
			flog.Infof("还原%s.%s 跳过空表", database, tableName)
			continue
		}

		f, err := os.Open(csvFile)
		if err != nil {
			return fmt.Errorf("打开%s失败: %v", csvFile, err)
		}

		lineScanner := bufio.NewScanner(f)
		lineScanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

		// 读取列头
		if !lineScanner.Scan() {
			f.Close()
			continue
		}
		columns := strings.Split(lineScanner.Text(), "\t")
		placeholders := make([]string, len(columns))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		insertSQL := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s)",
			database, tableName, strings.Join(columns, ","), strings.Join(placeholders, ","))

		tx, err := db.Begin()
		if err != nil {
			f.Close()
			return fmt.Errorf("开启事务失败: %v", err)
		}

		stmt, err := tx.Prepare(insertSQL)
		if err != nil {
			tx.Rollback()
			f.Close()
			return fmt.Errorf("准备INSERT语句失败: %v\nSQL: %s", err, insertSQL)
		}

		rowCount := 0
		for lineScanner.Scan() {
			line := lineScanner.Text()
			fields := strings.Split(line, "\t")
			args := make([]any, len(fields))
			for i, v := range fields {
				args[i] = v
			}
			_, err := stmt.Exec(args...)
			if err != nil {
				stmt.Close()
				tx.Rollback()
				f.Close()
				return fmt.Errorf("还原%s.%s插入数据失败：%v", database, tableName, err)
			}
			rowCount++
		}

		stmt.Close()
		if err := tx.Commit(); err != nil {
			f.Close()
			return fmt.Errorf("还原%s.%s提交事务失败：%v", database, tableName, err)
		}
		f.Close()

		flog.Infof("还原%s.%s (%d行) 使用了：%s", database, tableName, rowCount, sw.GetMillisecondsText())
	}
	return nil
}
