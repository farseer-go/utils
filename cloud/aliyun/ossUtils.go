package aliyun

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/utils/file"
)

// 阿里云OSS存储配置
type OSSConfig struct {
	AccessKeyID     string // AccessKeyID
	AccessKeySecret string // AccessKeySecret
	Endpoint        string // 填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	Region          string // 填写Bucket所在地域，以华东1（杭州）为例，填写为cn-hangzhou。其它Region请按实际情况填写。
	BucketName      string // BucketName
}

// 获取OSS客户端
func (receiver *OSSConfig) GetOssClient() (*oss.Client, string, error) {
	// 从环境变量中获取访问凭证。运行本代码示例之前，请确保已设置环境变量OSS_ACCESS_KEY_ID和OSS_ACCESS_KEY_SECRET。
	cfg := oss.LoadDefaultConfig().WithCredentialsProvider(credentials.NewStaticCredentialsProvider(receiver.AccessKeyID, receiver.AccessKeySecret))
	if receiver.Region != "" {
		cfg.WithRegion(receiver.Region)
	}
	if receiver.Endpoint != "" {
		cfg.WithEndpoint(receiver.Endpoint)
	}

	return oss.NewClient(cfg), receiver.BucketName, nil
}

// 上传备份文件到OSS
func (receiver *OSSConfig) UploadOSS(backupRoot string, fileNames []string) error {
	client, bucketName, err := receiver.GetOssClient()
	if err != nil || client == nil {
		return err
	}
	// 批量上传
	for _, fileName := range fileNames {
		f, err := os.Open(backupRoot + fileName)
		if err != nil {
			return fmt.Errorf("打开上传文件：%s 时，发生错误：%v", fileName, err)
		}
		defer f.Close()

		result, err := client.PutObject(context.TODO(), &oss.PutObjectRequest{
			Bucket: oss.Ptr(bucketName),
			Key:    oss.Ptr(fileName),
			Body:   f,
			Metadata: map[string]string{
				"x-oss-meta-local-time": time.Now().Format("2006-01-02 15:04:05"),
			},
		})

		if err != nil {
			return fmt.Errorf("上传文件：%s 时，发生错误：%v", fileName, err)
		}

		if result != nil && err == nil {
			//flog.Infof("OSS上传文件：%s 成功, ETag :%v\n", fileName, result.ETag)
			// 上传成功后，删除本地文件
			file.Delete(backupRoot + fileName)
		}
	}
	return nil
}

// 删除文件
func (receiver *OSSConfig) DeleteFile(fileName string) error {
	// 删除OSS文件
	client, bucketName, err := receiver.GetOssClient()
	if err != nil || client == nil {
		return err
	}

	// 执行删除对象的操作并处理结果
	_, err = client.DeleteObject(context.TODO(), &oss.DeleteObjectRequest{
		Bucket: oss.Ptr(bucketName), // 存储空间名称
		Key:    oss.Ptr(fileName),   // 对象名称
	})
	if err != nil {
		return fmt.Errorf("OSS删除文件：%s 时，发生错误：%v", fileName, err)
	}
	// 打印删除对象的结果
	//flog.Infof("OSS删除文件：%s 成功, ETag :%v\n", fileName, result)
	return nil
}

// 下载oss文件
func (receiver *OSSConfig) DownloadFile(backupRoot string, fileName string) error {
	client, bucketName, err := receiver.GetOssClient()
	if err != nil || client == nil {
		return err
	}

	// 创建获取对象的请求
	request := &oss.GetObjectRequest{
		Bucket: oss.Ptr(bucketName), // 存储空间名称
		Key:    oss.Ptr(fileName),   // 对象名称
	}

	// 执行获取对象的操作并处理结果
	result, err := client.GetObject(context.TODO(), request)
	if err != nil {
		return fmt.Errorf("无法获取文件：%s %v", fileName, err)
	}
	defer result.Body.Close() // 确保在函数结束时关闭响应体

	// 一次性读取整个文件内容
	data, err := io.ReadAll(result.Body)
	if err != nil {
		return fmt.Errorf("无法读取文件：%s %v", fileName, err)
	}

	// 将内容写入到文件
	path := filepath.Dir(backupRoot + fileName)
	if !file.IsExists(path) {
		file.CreateDir766(path)
	}
	err = os.WriteFile(backupRoot+fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("无法保存文件：%s %v", fileName, err)
	}
	return nil
}

// 备份历史数据
type FileObject struct {
	BackupId string            // 备份计划的ID
	FileName string            // 文件名
	CreateAt dateTime.DateTime // 备份时间
	Size     int64             // 备份文件大小（KB）
}

// 获取文件
func (receiver *OSSConfig) GetFileList(filePath string) (collections.List[FileObject], error) {
	lst := collections.NewList[FileObject]()
	client, bucketName, err := receiver.GetOssClient()
	if err != nil || client == nil {
		return lst, err
	}

	// 执行列举所有文件的操作
	lsRes, err := client.ListObjectsV2(context.TODO(), &oss.ListObjectsV2Request{
		Bucket:            oss.Ptr(bucketName),
		ContinuationToken: oss.Ptr(""),
		Prefix:            oss.Ptr(filePath), // 列举指定目录下的所有对象
		MaxKeys:           144,
	})

	if err != nil {
		return lst, fmt.Errorf("OSS读取文件列表时，发生错误：%v", err)
	}

	for _, object := range lsRes.Contents {
		lst.Add(FileObject{
			BackupId: bucketName,
			FileName: *object.Key,
			CreateAt: dateTime.New(*object.LastModified),
			Size:     object.Size / 1024,
		})
	}
	return lst, nil
}
